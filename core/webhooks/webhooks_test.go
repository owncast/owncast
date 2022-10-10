package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	jsonpatch "gopkg.in/evanphx/json-patch.v5"
)

func TestMain(m *testing.M) {
	dbFile, err := os.CreateTemp(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}
	dbFile.Close()
	defer os.Remove(dbFile.Name())

	if err := data.SetupPersistence(dbFile.Name()); err != nil {
		panic(err)
	}

	InitWorkerPool()
	defer close(queue)

	m.Run()
}

// Because the other tests use `sendEventToWebhooks` with a `WaitGroup` to know when the test completes,
// this test ensures that `SendToWebhooks` without a `WaitGroup` doesn't panic.
func TestPublicSend(t *testing.T) {
	// Send enough events to be sure at least one worker delivers a second event.
	const eventsCount = webhookWorkerPoolSize + 1

	var wg sync.WaitGroup
	wg.Add(eventsCount)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Done()
	}))
	defer svr.Close()

	hook, err := data.InsertWebhook(svr.URL, []models.EventType{models.MessageSent})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := data.DeleteWebhook(hook); err != nil {
			t.Error(err)
		}
	}()

	for i := 0; i < eventsCount; i++ {
		wh := WebhookEvent{
			EventData: struct{}{},
			Type:      models.MessageSent,
		}
		SendEventToWebhooks(wh)
	}

	wg.Wait()
}

// Make sure that events are only sent to interested endpoints.
func TestRouting(t *testing.T) {
	eventTypes := []models.EventType{models.PING, models.PONG}

	calls := map[models.EventType]int{}
	var lock sync.Mutex
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) < 1 || r.URL.Path[0] != '/' {
			t.Fatalf("Got unexpected path %v", r.URL.Path)
		}
		pathType := r.URL.Path[1:]
		var body WebhookEvent
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Type != pathType {
			t.Fatalf("Got %v payload on %v endpoint", body.Type, pathType)
		}
		lock.Lock()
		defer lock.Unlock()
		calls[pathType] += 1
	}))
	defer svr.Close()

	for _, eventType := range eventTypes {
		hook, err := data.InsertWebhook(svr.URL+"/"+eventType, []models.EventType{eventType})
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := data.DeleteWebhook(hook); err != nil {
				t.Error(err)
			}
		}()
	}

	var wg sync.WaitGroup

	for _, eventType := range eventTypes {
		wh := WebhookEvent{
			EventData: struct{}{},
			Type:      eventType,
		}
		sendEventToWebhooks(wh, &wg)
	}

	wg.Wait()

	for _, eventType := range eventTypes {
		if calls[eventType] != 1 {
			t.Errorf("Expected %v to be called exactly once but it was called %v times", eventType, calls[eventType])
		}
	}
}

// Make sure that events are sent to all interested endpoints.
func TestMultiple(t *testing.T) {
	const times = 2

	var calls uint32
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(&calls, 1)
	}))
	defer svr.Close()

	for i := 0; i < times; i++ {
		hook, err := data.InsertWebhook(fmt.Sprintf("%v/%v", svr.URL, i), []models.EventType{models.MessageSent})
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := data.DeleteWebhook(hook); err != nil {
				t.Error(err)
			}
		}()
	}

	var wg sync.WaitGroup

	wh := WebhookEvent{
		EventData: struct{}{},
		Type:      models.MessageSent,
	}
	sendEventToWebhooks(wh, &wg)

	wg.Wait()

	if atomic.LoadUint32(&calls) != times {
		t.Errorf("Expected event to be sent exactly %v times but it was sent %v times", times, atomic.LoadUint32(&calls))
	}
}

// Make sure when a webhook is used its last used timestamp is updated.
func TestTimestamps(t *testing.T) {
	const tolerance = time.Second
	start := time.Now()
	eventTypes := []models.EventType{models.PING, models.PONG}
	handlerIds := []int{0, 0}
	handlers := []*models.Webhook{nil, nil}
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer svr.Close()

	for i, eventType := range eventTypes {
		hook, err := data.InsertWebhook(svr.URL+"/"+eventType, []models.EventType{eventType})
		if err != nil {
			t.Fatal(err)
		}
		handlerIds[i] = hook
		defer func() {
			if err := data.DeleteWebhook(hook); err != nil {
				t.Error(err)
			}
		}()
	}

	var wg sync.WaitGroup

	wh := WebhookEvent{
		EventData: struct{}{},
		Type:      eventTypes[0],
	}
	sendEventToWebhooks(wh, &wg)

	wg.Wait()

	hooks, err := data.GetWebhooks()
	if err != nil {
		t.Fatal(err)
	}

	for h, hook := range hooks {
		for i, handlerId := range handlerIds {
			if hook.ID == handlerId {
				handlers[i] = &hooks[h]
			}
		}
	}

	if handlers[0] == nil {
		t.Fatal("First handler was not found in registered handlers")
	}
	if handlers[1] == nil {
		t.Fatal("Second handler was not found in registered handlers")
	}

	end := time.Now()

	if handlers[0].Timestamp.Add(tolerance).Before(start) {
		t.Errorf("First handler timestamp %v should not be before start of test %v", handlers[0].Timestamp, start)
	}

	if handlers[0].Timestamp.Add(tolerance).Before(handlers[1].Timestamp) {
		t.Errorf("Second handler timestamp %v should not be before first handler timestamp %v", handlers[1].Timestamp, handlers[0].Timestamp)
	}

	if end.Add(tolerance).Before(handlers[1].Timestamp) {
		t.Errorf("End of test %v should not be before second handler timestamp %v", end, handlers[1].Timestamp)
	}

	if handlers[0].LastUsed == nil {
		t.Error("First handler last used should have been set")
	} else if handlers[0].LastUsed.Add(tolerance).Before(handlers[1].Timestamp) {
		t.Errorf("First handler last used %v should not be before second handler timestamp %v", handlers[0].LastUsed, handlers[1].Timestamp)
	} else if end.Add(tolerance).Before(*handlers[0].LastUsed) {
		t.Errorf("End of test %v should not be before first handler last used %v", end, handlers[0].LastUsed)
	}

	if handlers[1].LastUsed != nil {
		t.Error("Second handler last used should not have been set")
	}
}

// Make sure up to the expected number of events can be fired in parallel.
func TestParallel(t *testing.T) {
	var calls uint32

	var wgStart sync.WaitGroup
	finished := make(chan int)
	wgStart.Add(webhookWorkerPoolSize)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		myId := atomic.AddUint32(&calls, 1)

		// We made it to the pool size + 1 event, so we're done with the test.
		if myId == webhookWorkerPoolSize+1 {
			close(finished)
			return
		}

		// Wait until all the handlers are started.
		wgStart.Done()
		wgStart.Wait()

		// The first handler just returns so the pool size + 1 event can be handled.
		if myId != 1 {
			// The other handlers will wait for pool size + 1.
			_ = <-finished
		}
	}))
	defer svr.Close()

	hook, err := data.InsertWebhook(svr.URL, []models.EventType{models.MessageSent})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := data.DeleteWebhook(hook); err != nil {
			t.Error(err)
		}
	}()

	var wgMessages sync.WaitGroup

	for i := 0; i < webhookWorkerPoolSize+1; i++ {
		wh := WebhookEvent{
			EventData: struct{}{},
			Type:      models.MessageSent,
		}
		sendEventToWebhooks(wh, &wgMessages)
	}

	wgMessages.Wait()
}

// Send an event, capture it, and verify that it has the expected payload.
func checkPayload(t *testing.T, eventType models.EventType, send func(), expectedJson string) {
	eventChannel := make(chan WebhookEvent)

	// Set up a server.
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := WebhookEvent{}
		json.NewDecoder(r.Body).Decode(&data)
		eventChannel <- data
	}))
	defer svr.Close()

	// Subscribe to the webhook.
	hook, err := data.InsertWebhook(svr.URL, []models.EventType{eventType})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := data.DeleteWebhook(hook); err != nil {
			t.Error(err)
		}
	}()

	// Send and capture the event.
	send()
	event := <-eventChannel

	if event.Type != eventType {
		t.Errorf("Got event type %v but expected %v", event.Type, eventType)
	}

	// Compare.
	payloadJson, err := json.MarshalIndent(event.EventData, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Actual payload:\n%s", payloadJson)

	if !jsonpatch.Equal(payloadJson, []byte(expectedJson)) {
		diff, err := jsonpatch.CreateMergePatch(payloadJson, []byte(expectedJson))
		if err != nil {
			t.Fatal(err)
		}
		var out bytes.Buffer
		if err := json.Indent(&out, diff, "", "  "); err != nil {
			t.Fatal(err)
		}
		t.Errorf("Expected difference from actual payload:\n%s", out.Bytes())
	}
}
