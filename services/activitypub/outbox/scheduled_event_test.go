package outbox

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"code.superseriousbusiness.org/activity/streams"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	appersistence "github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	apresolvers "github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/activitypub/workerpool"
	"github.com/owncast/owncast/services/datastore"
)

type scheduledEventFollowers struct {
	followersrepository.FollowersRepository
	inboxes []string
}

func (f *scheduledEventFollowers) GetUniqueDeliveryInboxes() ([]string, error) {
	return f.inboxes, nil
}

func TestSendScheduledEventQueuesCreateForEveryFollowerInbox(t *testing.T) {
	service, persistence, ds := newScheduledEventService(t, []string{"https://example.com/inbox", "https://www.example.com/inbox"})
	createdAt := time.Date(2030, time.September, 8, 12, 0, 0, 0, time.UTC)
	if err := service.SendScheduledEvent(models.ScheduledEvent{
		ID:              "weekly-show",
		Name:            "Weekly show",
		Description:     "A scheduled example stream.",
		StartTime:       createdAt.Add(48 * time.Hour),
		DurationMinutes: 90,
		Timezone:        "America/New_York",
		CreatedAt:       &createdAt,
	}); err != nil {
		t.Fatal(err)
	}
	if err := service.SendScheduledEvent(models.ScheduledEvent{
		ID:              "weekly-show",
		Name:            "Weekly show",
		Description:     "A scheduled example stream.",
		StartTime:       createdAt.Add(48 * time.Hour),
		DurationMinutes: 90,
		Timezone:        "America/New_York",
		CreatedAt:       &createdAt,
	}); err != nil {
		t.Fatal(err)
	}

	rows, err := ds.DB.Query(`SELECT inbox, payload FROM ap_delivery_queue ORDER BY inbox`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var deliveries int
	for rows.Next() {
		var inbox string
		var payload []byte
		if err := rows.Scan(&inbox, &payload); err != nil {
			t.Fatal(err)
		}
		activity := decodePayload(t, payload)
		if activity["type"] != "Create" || activity["id"] != "https://owncast.example.com/federation/event/weekly-show/activity/create" {
			t.Errorf("%s activity = %#v", inbox, activity)
		}
		object, ok := activity["object"].(map[string]interface{})
		if !ok || object["type"] != "Event" || object["id"] != "https://owncast.example.com/federation/event/weekly-show" {
			t.Errorf("%s object = %#v", inbox, activity["object"])
		}
		attachment, ok := object["attachment"].([]interface{})
		if !ok || len(attachment) != 1 {
			t.Errorf("%s attachment = %#v", inbox, object["attachment"])
		}
		deliveries++
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if deliveries != 2 {
		t.Fatalf("queued %d deliveries, want one for each follower inbox", deliveries)
	}

	stored, _, _, err := persistence.GetObjectByIRI("https://owncast.example.com/federation/event/weekly-show")
	if err != nil {
		t.Fatal(err)
	}
	event := decodePayload(t, []byte(stored))
	if event["type"] != "Event" || event["name"] != "Weekly show" {
		t.Fatalf("stored event = %#v", event)
	}
	storedCreate, _, _, err := persistence.GetObjectByIRI("https://owncast.example.com/federation/event/weekly-show/activity/create")
	if err != nil {
		t.Fatal(err)
	}
	if create := decodePayload(t, []byte(storedCreate)); create["type"] != "Create" {
		t.Fatalf("stored activity = %#v", create)
	}

	outboxPersistence := appersistence.New(ds, apresolvers.New(apresolvers.Deps{}))
	collection, err := outboxPersistence.GetOutbox(50, 0)
	if err != nil {
		t.Fatal(err)
	}
	serialized, err := streams.Serialize(collection)
	if err != nil {
		t.Fatal(err)
	}
	outboxPayload, err := json.Marshal(serialized)
	if err != nil {
		t.Fatal(err)
	}
	outbox := decodePayload(t, outboxPayload)
	create, ok := outbox["orderedItems"].(map[string]interface{})
	if !ok || create["type"] != "Create" {
		t.Fatalf("outbox item = %#v", outbox["orderedItems"])
	}
}

func TestScheduledEventLifecyclePublishesCompleteUpdatesAndTombstone(t *testing.T) {
	service, persistence, ds := newScheduledEventService(t, []string{"https://example.com/inbox"})
	updatedAt := time.Date(2030, time.September, 8, 13, 0, 0, 0, time.UTC)
	event := models.ScheduledEvent{
		ID:                "weekly-show",
		Name:              "Updated weekly show",
		Description:       "",
		StartTime:         updatedAt.Add(24 * time.Hour),
		DurationMinutes:   60,
		Timezone:          "UTC",
		Status:            models.ScheduledEventStatusCancelled,
		UpdatedAt:         &updatedAt,
		FederationVersion: 2,
	}
	storeScheduledEvent(t, ds, event)
	original := event
	original.Name = "Weekly show"
	original.Status = models.ScheduledEventStatusScheduled
	if err := service.SendScheduledEvent(original); err != nil {
		t.Fatal(err)
	}
	if err := service.SendScheduledEventReminder(original, 1, "Starting soon"); err != nil {
		t.Fatal(err)
	}
	if err := service.SendScheduledEventUpdate(event); err != nil {
		t.Fatal(err)
	}
	if err := service.SendScheduledEventUpdate(event); err != nil {
		t.Fatal(err)
	}
	var queuedUpdates int
	if err := ds.DB.QueryRow(`SELECT COUNT(*) FROM ap_delivery_queue WHERE activity_type = 'Update'`).Scan(&queuedUpdates); err != nil {
		t.Fatal(err)
	}
	if queuedUpdates != 1 {
		t.Fatalf("queued identical updates = %d, want 1", queuedUpdates)
	}
	var queuedActivities int
	if err := ds.DB.QueryRow(`SELECT COUNT(*) FROM ap_delivery_queue`).Scan(&queuedActivities); err != nil {
		t.Fatal(err)
	}
	if queuedActivities != 2 {
		t.Fatalf("queued lifecycle activities = %d, want Create followed by latest state", queuedActivities)
	}
	var queuedOrder string
	if err := ds.DB.QueryRow(`SELECT group_concat(activity_type, ',') FROM (SELECT activity_type FROM ap_delivery_queue ORDER BY id)`).Scan(&queuedOrder); err != nil {
		t.Fatal(err)
	}
	if queuedOrder != "Create,Update" {
		t.Fatalf("queued lifecycle order = %q", queuedOrder)
	}
	var queuedReminders int
	if err := ds.DB.QueryRow(`SELECT COUNT(*) FROM ap_delivery_queue WHERE coalesce_key LIKE '%:reminder:%'`).Scan(&queuedReminders); err != nil {
		t.Fatal(err)
	}
	if queuedReminders != 0 {
		t.Fatalf("queued reminders after cancellation = %d", queuedReminders)
	}

	update := queuedActivity(t, ds, "Update")
	object, ok := update["object"].(map[string]interface{})
	if !ok || object["id"] != "https://owncast.example.com/federation/event/weekly-show" || object["eventStatus"] != "EventCancelled" || object["content"] != "" {
		t.Fatalf("update object = %#v", update["object"])
	}
	if _, ok := object["attachment"].([]interface{}); !ok {
		t.Fatalf("update attachment is not an array: %#v", object["attachment"])
	}

	if err := service.SendScheduledEventDelete(event); err != nil {
		t.Fatal(err)
	}
	deleted := queuedActivity(t, ds, "Delete")
	tombstone, ok := deleted["object"].(map[string]interface{})
	if !ok || tombstone["type"] != "Tombstone" || tombstone["formerType"] != "Event" {
		t.Fatalf("delete object = %#v", deleted["object"])
	}
	stale := event
	stale.Name = "stale update"
	stale.FederationVersion--
	if err := service.SendScheduledEventUpdate(stale); err != nil {
		t.Fatal(err)
	}
	if err := service.SendScheduledEventReminder(stale, 2, "Stale reminder"); err != nil {
		t.Fatal(err)
	}
	if err := ds.DB.QueryRow(`SELECT COUNT(*) FROM ap_delivery_queue WHERE coalesce_key LIKE '%:reminder:%'`).Scan(&queuedReminders); err != nil {
		t.Fatal(err)
	}
	if queuedReminders != 0 {
		t.Fatalf("queued stale reminders after deletion = %d", queuedReminders)
	}
	stored, _, _, err := persistence.GetObjectByIRI("https://owncast.example.com/federation/event/weekly-show")
	if err != nil {
		t.Fatal(err)
	}
	if direct := decodePayload(t, []byte(stored)); direct["type"] != "Tombstone" {
		t.Fatalf("direct object = %#v", direct)
	}
	var queuedStateType string
	if err := ds.DB.QueryRow(`SELECT activity_type FROM ap_delivery_queue WHERE coalesce_key = ?`, "scheduled-event:weekly-show:state").Scan(&queuedStateType); err != nil {
		t.Fatal(err)
	}
	if queuedStateType != "Delete" {
		t.Fatalf("queued state after stale update = %q", queuedStateType)
	}
}

func TestScheduledEventDeleteRollsBackWhenQueueingFails(t *testing.T) {
	service, persistence, ds := newScheduledEventService(t, []string{"https://example.com/inbox"})
	event := models.ScheduledEvent{ID: "queue-failure", Name: "Event", FederationVersion: 1}
	if err := service.SendScheduledEvent(event); err != nil {
		t.Fatal(err)
	}
	if _, err := ds.DB.Exec(`CREATE TRIGGER fail_scheduled_event_delivery BEFORE INSERT ON ap_delivery_queue BEGIN SELECT RAISE(ABORT, 'queue failed'); END`); err != nil {
		t.Fatal(err)
	}

	if err := service.SendScheduledEventDelete(event); err == nil {
		t.Fatal("SendScheduledEventDelete() succeeded with a failing delivery queue")
	}
	stored, _, _, err := persistence.GetObjectByIRI("https://owncast.example.com/federation/event/queue-failure")
	if err != nil {
		t.Fatal(err)
	}
	if object := decodePayload(t, []byte(stored)); object["type"] != "Event" {
		t.Fatalf("canonical object after rollback = %#v", object)
	}
}

func TestScheduledEventReminderUsesStableNote(t *testing.T) {
	service, persistence, ds := newScheduledEventService(t, []string{"https://example.com/inbox"})
	event := models.ScheduledEvent{
		ID:                "weekly-show",
		Name:              "Weekly show",
		StartTime:         time.Now().Add(time.Hour),
		DurationMinutes:   60,
		Timezone:          "UTC",
		Status:            models.ScheduledEventStatusScheduled,
		FederationVersion: 1,
	}
	storeScheduledEvent(t, ds, event)
	if err := service.SendScheduledEvent(event); err != nil {
		t.Fatal(err)
	}
	if err := service.SendScheduledEventReminder(event, 1, "Starting soon"); err != nil {
		t.Fatal(err)
	}
	activity := queuedActivity(t, ds, "Create")
	if activity["id"] != "https://owncast.example.com/federation/event/weekly-show/reminder/1/activity" {
		t.Fatalf("reminder activity id = %#v", activity["id"])
	}
	note, ok := activity["object"].(map[string]interface{})
	if !ok || note["id"] != "https://owncast.example.com/federation/event/weekly-show/reminder/1" || note["inReplyTo"] != "https://owncast.example.com/federation/event/weekly-show" {
		t.Fatalf("reminder note = %#v", activity["object"])
	}
	if _, _, _, err := persistence.GetObjectByIRI("https://owncast.example.com/federation/event/weekly-show/reminder/1"); err != nil {
		t.Fatal(err)
	}
	collection, err := persistence.GetOutbox(50, 0)
	if err != nil {
		t.Fatal(err)
	}
	serialized, err := streams.Serialize(collection)
	if err != nil {
		t.Fatal(err)
	}
	items, ok := serialized["orderedItems"].([]interface{})
	if !ok {
		t.Fatalf("actor outbox items = %#v", serialized["orderedItems"])
	}
	found := false
	for _, item := range items {
		outboxActivity, ok := item.(map[string]interface{})
		if ok && outboxActivity["id"] == activity["id"] {
			found = true
		}
	}
	if !found {
		t.Fatalf("actor outbox omits delivered reminder id %#v", activity["id"])
	}
}

func storeScheduledEvent(t *testing.T, ds *datastore.Datastore, event models.ScheduledEvent) {
	t.Helper()
	if _, err := ds.DB.Exec(
		`INSERT INTO stream_events(id, name, start_time, duration_minutes, timezone, status, federated_at, federation_version) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?)`,
		event.ID, event.Name, event.StartTime, event.DurationMinutes, event.Timezone, models.ScheduledEventStatusScheduled, event.FederationVersion,
	); err != nil {
		t.Fatal(err)
	}
}

func newScheduledEventService(t *testing.T, inboxes []string) (*Service, *appersistence.Service, *datastore.Datastore) {
	t.Helper()
	directory := t.TempDir()
	ds, err := datastore.SetupPersistence(filepath.Join(directory, "owncast.db"), directory)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ds.DB.Close() })
	configRepo := configrepository.New(ds)
	for _, set := range []func() error{
		func() error { return configRepo.SetServerURL("https://owncast.example.com") },
		func() error { return configRepo.SetServerName("Test stream") },
		func() error { return configRepo.SetFederationUsername("streamer") },
		func() error { return configRepo.SetFederationEnabled(true) },
		func() error { return configRepo.SetFederationIsPrivate(false) },
	} {
		if err := set(); err != nil {
			t.Fatal(err)
		}
	}
	persistence := appersistence.New(ds, nil)
	service := New(Deps{
		Persistence:      persistence,
		Workerpool:       workerpool.New(workerpool.Deps{WorkerPoolSize: 1, Datastore: ds}),
		Followers:        &scheduledEventFollowers{inboxes: inboxes},
		ConfigRepository: configRepo,
		Builder:          apmodels.New(apmodels.Deps{ConfigRepository: configRepo}),
	})
	return service, persistence, ds
}

func queuedActivity(t *testing.T, ds *datastore.Datastore, activityType string) map[string]interface{} {
	t.Helper()
	var payload []byte
	if err := ds.DB.QueryRow(`SELECT payload FROM ap_delivery_queue WHERE activity_type = ? ORDER BY id DESC LIMIT 1`, activityType).Scan(&payload); err != nil {
		t.Fatal(err)
	}
	return decodePayload(t, payload)
}

func decodePayload(t *testing.T, payload []byte) map[string]interface{} {
	t.Helper()
	var decoded map[string]interface{}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatal(err)
	}
	return decoded
}
