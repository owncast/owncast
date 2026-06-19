package chat

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/webhooks"
)

// fakeJoinPartConfig overrides only GetChatJoinPartMessagesEnabled; every other
// ConfigRepository method is inherited from the embedded (nil) interface and is
// never called by sendUserJoinedMessage.
type fakeJoinPartConfig struct {
	configrepository.ConfigRepository
	enabled bool
}

func (f fakeJoinPartConfig) GetChatJoinPartMessagesEnabled() bool { return f.enabled }

// TestSendUserJoinedMessage_WebhookFiresRegardlessOfJoinPartSetting guards the
// fix for #4950: the USER_JOINED webhook must fire on a genuine join even when
// the admin "show join/part messages" setting is off (it only governs the
// visible chat broadcast), mirroring how the PART path already behaves.
func TestSendUserJoinedMessage_WebhookFiresRegardlessOfJoinPartSetting(t *testing.T) {
	for _, joinPartEnabled := range []bool{true, false} {
		joinPartEnabled := joinPartEnabled
		t.Run(map[bool]string{true: "messages-enabled", false: "messages-disabled"}[joinPartEnabled], func(t *testing.T) {
			// A real datastore-backed webhooks service so dispatch actually runs.
			dbFile, err := os.CreateTemp(t.TempDir(), "owncast-test-*.db")
			if err != nil {
				t.Fatal(err)
			}
			dbFile.Close()

			ds, err := datastore.SetupPersistence(dbFile.Name(), t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			// Release the DB file handle before t.TempDir's RemoveAll runs;
			// Windows cannot delete a file that is still open (registered after
			// t.TempDir so it runs first, LIFO).
			t.Cleanup(func() { _ = ds.DB.Close() })

			realConfig := configrepository.New(ds)
			realConfig.SetServerURL("http://localhost:8080")
			webhookRepo := webhookrepository.New(ds)

			whSvc := webhooks.New(webhooks.Deps{
				GetStatus:         func() models.Status { return models.Status{Online: true} },
				ConfigRepository:  realConfig,
				WebhookRepository: webhookRepo,
			})
			whSvc.Start()

			received := make(chan struct{}, 1)
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				received <- struct{}{}
			}))
			defer svr.Close()

			if _, err := webhookRepo.InsertWebhook(svr.URL, []models.EventType{models.UserJoined}); err != nil {
				t.Fatal(err)
			}

			chatSvc := New(Deps{
				Webhooks:         whSvc,
				ConfigRepository: fakeJoinPartConfig{enabled: joinPartEnabled},
				GetStatus:        func() models.Status { return models.Status{Online: true} },
			})

			// An already-connected observer so we can detect whether the join was
			// broadcast to chat (Broadcast delivers to registered clients).
			observer := &Client{
				User: &models.User{ID: "observer", DisplayName: "Observer"},
				Id:   99,
				send: make(chan []byte, 8),
			}
			chatSvc.clients[observer.Id] = observer

			joining := &Client{
				User: &models.User{ID: "u1", DisplayName: "Tester"},
				Id:   1,
				send: make(chan []byte, 8),
			}

			chatSvc.sendUserJoinedMessage(joining)

			// The webhook must always fire for a real join (the #4950 bug: it did not when messages were disabled).
			select {
			case <-received:
			case <-time.After(3 * time.Second):
				t.Fatalf("USER_JOINED webhook did not fire (joinPartEnabled=%v)", joinPartEnabled)
			}

			// The visible chat broadcast is gated on the display setting.
			gotBroadcast := len(observer.send) > 0
			if gotBroadcast != joinPartEnabled {
				t.Errorf("broadcast sent=%v, want %v (joinPartEnabled=%v)", gotBroadcast, joinPartEnabled, joinPartEnabled)
			}
		})
	}
}
