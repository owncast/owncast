package outbox

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	appersistence "github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
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
	workers := workerpool.New(workerpool.Deps{WorkerPoolSize: 1, Datastore: ds})
	builder := apmodels.New(apmodels.Deps{ConfigRepository: configRepo})
	service := New(Deps{
		Persistence:      persistence,
		Workerpool:       workers,
		Followers:        &scheduledEventFollowers{inboxes: []string{"https://example.com/inbox", "https://www.example.com/inbox"}},
		ConfigRepository: configRepo,
		Builder:          builder,
	})

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
		var activity map[string]interface{}
		if err := json.Unmarshal(payload, &activity); err != nil {
			t.Fatal(err)
		}
		if activity["type"] != "Create" {
			t.Errorf("%s activity type = %#v", inbox, activity["type"])
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
	var event map[string]interface{}
	if err := json.Unmarshal([]byte(stored), &event); err != nil {
		t.Fatal(err)
	}
	if event["type"] != "Event" || event["name"] != "Weekly show" {
		t.Fatalf("stored event = %#v", event)
	}
}
