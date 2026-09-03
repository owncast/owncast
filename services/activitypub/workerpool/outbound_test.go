package workerpool

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"

	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/persistence/configrepository"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/datastore"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("OWNCAST_ALLOW_INTERNAL_FEDERATION", "true")
	_ = os.Setenv("OWNCAST_INSECURE_SKIP_VERIFY", "true")
	os.Exit(m.Run())
}

func TestDeliveryRetriesAfterTransientFailure(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	service, database := newTestService(t)
	service.retryDelay = func(int64) time.Duration { return 20 * time.Millisecond }
	service.Start()
	defer service.Stop()

	destination, _ := url.Parse(server.URL)
	actor, _ := url.Parse("https://sender.example/federation/user/streamer")
	if err := service.Enqueue(Delivery{
		Inbox:        destination,
		Payload:      []byte(`{"type":"Accept"}`),
		ActorIRI:     actor,
		ActivityType: "Accept",
	}); err != nil {
		t.Fatalf("enqueue delivery: %v", err)
	}

	waitFor(t, 20*time.Second, func() bool { return attempts.Load() == 3 })
	if count := deliveryCount(t, database.db); count != 0 {
		t.Fatalf("delivery queue count = %d, want 0 after success", count)
	}
}

func TestCircuitBreakerBackoffLogsAtDebug(t *testing.T) {
	originalLevel := log.GetLevel()
	log.SetLevel(log.DebugLevel)
	defer log.SetLevel(originalLevel)

	hook := logtest.NewGlobal()
	defer hook.Reset()

	service, _ := newTestService(t)
	const domain = "failing.example.com"
	for range circuitBreakerFailureThreshold {
		service.recordDomainFailure(domain)
	}

	var backoffLogs int
	for _, entry := range hook.AllEntries() {
		if strings.Contains(entry.Message, "backing off") {
			backoffLogs++
			if entry.Level != log.DebugLevel {
				t.Errorf("backoff log %q level = %s, want debug", entry.Message, entry.Level)
			}
		}
	}
	if backoffLogs != 1 {
		t.Errorf("backoff log count = %d, want 1", backoffLogs)
	}
}

func TestDeliverySurvivesWorkerRestart(t *testing.T) {
	var available atomic.Bool
	var attempts atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		if !available.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	first, database := newTestService(t)
	first.retryDelay = func(int64) time.Duration { return 200 * time.Millisecond }
	first.Start()
	destination, _ := url.Parse(server.URL)
	actor, _ := url.Parse("https://sender.example/federation/user/streamer")
	if err := first.Enqueue(Delivery{
		Inbox:        destination,
		Payload:      []byte(`{"type":"Accept"}`),
		ActorIRI:     actor,
		ActivityType: "Accept",
	}); err != nil {
		t.Fatalf("enqueue delivery: %v", err)
	}
	waitFor(t, 5*time.Second, func() bool { return attempts.Load() == 1 })
	first.Stop()

	available.Store(true)
	restarted := New(Deps{WorkerPoolSize: 1, Datastore: database.datastore, Signer: database.signer})
	restarted.Start()
	defer restarted.Stop()

	waitFor(t, 2*time.Second, func() bool { return attempts.Load() == 2 })
	if count := deliveryCount(t, database.db); count != 0 {
		t.Fatalf("delivery queue count = %d, want 0 after restarted worker succeeds", count)
	}
}

func TestCoalescedDeliveryKeepsNewestPayload(t *testing.T) {
	service, database := newTestService(t)
	destination, _ := url.Parse("https://receiver.example/inbox")
	actor, _ := url.Parse("https://sender.example/federation/user/streamer")

	for _, payload := range []string{"old", "new"} {
		if err := service.Enqueue(Delivery{
			Inbox:        destination,
			Payload:      []byte(payload),
			ActorIRI:     actor,
			ActivityType: "Offer",
			CoalesceKey:  "stream-status",
		}); err != nil {
			t.Fatalf("enqueue %q: %v", payload, err)
		}
	}

	var payload string
	var count int
	if err := database.db.QueryRow(`SELECT CAST(payload AS TEXT), count(*) FROM ap_delivery_queue`).Scan(&payload, &count); err != nil {
		t.Fatalf("read coalesced delivery: %v", err)
	}
	if payload != "new" || count != 1 {
		t.Fatalf("coalesced delivery = (%q, %d rows), want newest payload in one row", payload, count)
	}
}

func TestCoalescedDeliveryRejectsOlderVersion(t *testing.T) {
	service, database := newTestService(t)
	destination, _ := url.Parse("https://receiver.example/inbox")
	actor, _ := url.Parse("https://sender.example/federation/user/streamer")

	for _, delivery := range []Delivery{
		{Inbox: destination, Payload: []byte("new"), ActorIRI: actor, ActivityType: "Update", CoalesceKey: "event", CoalesceVersion: 2},
		{Inbox: destination, Payload: []byte("old"), ActorIRI: actor, ActivityType: "Update", CoalesceKey: "event", CoalesceVersion: 1},
	} {
		if err := service.Enqueue(delivery); err != nil {
			t.Fatal(err)
		}
	}

	var payload string
	if err := database.db.QueryRow(`SELECT CAST(payload AS TEXT) FROM ap_delivery_queue`).Scan(&payload); err != nil {
		t.Fatal(err)
	}
	if payload != "new" {
		t.Fatalf("coalesced payload = %q, want newer version", payload)
	}
}

func TestOrderedDeliveriesClaimInSequence(t *testing.T) {
	service, database := newTestService(t)
	destination, _ := url.Parse("https://receiver.example/inbox")
	actor, _ := url.Parse("https://sender.example/federation/user/streamer")
	for _, delivery := range []Delivery{
		{Inbox: destination, Payload: []byte("create"), ActorIRI: actor, ActivityType: "Create", CoalesceKey: "create", OrderingKey: "event", BlocksFollowing: true},
		{Inbox: destination, Payload: []byte("reminder"), ActorIRI: actor, ActivityType: "Create", CoalesceKey: "reminder", OrderingKey: "event"},
		{Inbox: destination, Payload: []byte("update"), ActorIRI: actor, ActivityType: "Update", CoalesceKey: "state", OrderingKey: "event"},
	} {
		if err := service.Enqueue(delivery); err != nil {
			t.Fatal(err)
		}
	}

	queries := database.datastore.GetQueries()
	now := time.Now()
	first, err := queries.ClaimActivityPubDelivery(context.Background(), db.ClaimActivityPubDeliveryParams{
		ClaimedUntil: sql.NullTime{Time: now.Add(time.Minute), Valid: true},
		Now:          now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(first.Payload) != "create" {
		t.Fatalf("first claimed payload = %q", first.Payload)
	}
	if _, err := queries.ClaimActivityPubDelivery(context.Background(), db.ClaimActivityPubDeliveryParams{
		ClaimedUntil: sql.NullTime{Time: now.Add(time.Minute), Valid: true},
		Now:          now,
	}); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("claim with pending predecessor error = %v, want sql.ErrNoRows", err)
	}
	if _, err := queries.CompleteActivityPubDelivery(context.Background(), db.CompleteActivityPubDeliveryParams{
		ID: first.ID, Revision: first.Revision,
	}); err != nil {
		t.Fatal(err)
	}
	second, err := queries.ClaimActivityPubDelivery(context.Background(), db.ClaimActivityPubDeliveryParams{
		ClaimedUntil: sql.NullTime{Time: now.Add(time.Minute), Valid: true},
		Now:          now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(second.Payload) != "reminder" {
		t.Fatalf("second claimed payload = %q", second.Payload)
	}
	third, err := queries.ClaimActivityPubDelivery(context.Background(), db.ClaimActivityPubDeliveryParams{
		ClaimedUntil: sql.NullTime{Time: now.Add(time.Minute), Valid: true},
		Now:          now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(third.Payload) != "update" {
		t.Fatalf("third claimed payload = %q", third.Payload)
	}
}

func TestRemovedClaimedDeliveryIsNotSent(t *testing.T) {
	var attempts atomic.Int64
	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		attempts.Add(1)
	}))
	defer server.Close()

	service, database := newTestService(t)
	destination, _ := url.Parse(server.URL)
	actor, _ := url.Parse("https://sender.example/federation/user/streamer")
	if err := service.Enqueue(Delivery{
		Inbox: destination, Payload: []byte(`{"type":"Create"}`), ActorIRI: actor,
		ActivityType: "Create", CoalesceKey: "reminder", OrderingKey: "event",
	}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	job, err := database.datastore.GetQueries().ClaimActivityPubDelivery(context.Background(), db.ClaimActivityPubDeliveryParams{
		ClaimedUntil: sql.NullTime{Time: now.Add(time.Minute), Valid: true},
		Now:          now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := database.db.Exec(`DELETE FROM ap_delivery_queue WHERE id = ?`, job.ID); err != nil {
		t.Fatal(err)
	}

	service.httpClient = server.Client()
	service.deliver(job)
	if attempts.Load() != 0 {
		t.Fatal("removed claimed delivery was sent")
	}
}

type testDatabase struct {
	db        *sql.DB
	datastore *datastore.Datastore
	signer    *apcrypto.Signer
}

func newTestService(t *testing.T) (*Service, *testDatabase) {
	t.Helper()
	directory := t.TempDir()
	ds, err := datastore.SetupPersistence(filepath.Join(directory, "owncast.db"), directory)
	if err != nil {
		t.Fatalf("set up persistence: %v", err)
	}
	t.Cleanup(func() { _ = ds.DB.Close() })

	configRepo := configrepository.New(ds)
	privateKey, publicKey, err := apcrypto.GenerateKeys()
	if err != nil {
		t.Fatalf("generate signing keys: %v", err)
	}
	if err := configRepo.SetPrivateKey(string(privateKey)); err != nil {
		t.Fatalf("store private key: %v", err)
	}
	if err := configRepo.SetPublicKey(string(publicKey)); err != nil {
		t.Fatalf("store public key: %v", err)
	}
	signer := apcrypto.New(apcrypto.Deps{ConfigRepository: configRepo})
	testDB := &testDatabase{db: ds.DB, datastore: ds, signer: signer}
	return New(Deps{WorkerPoolSize: 1, Datastore: ds, Signer: signer}), testDB
}

func deliveryCount(t *testing.T, database *sql.DB) int {
	t.Helper()
	var count int
	if err := database.QueryRow(`SELECT count(*) FROM ap_delivery_queue`).Scan(&count); err != nil {
		t.Fatalf("count deliveries: %v", err)
	}
	return count
}

func waitFor(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("condition was not met before timeout")
}
