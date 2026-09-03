package apmodels

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
)

func TestScheduledEventFEP8a8eRepresentation(t *testing.T) {
	createdAt := time.Date(2030, time.September, 8, 12, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	event, err := testBuilder.MakeScheduledEvent(models.ScheduledEvent{
		ID:              "weekly-show",
		Name:            "Weekly show",
		Description:     "A scheduled example stream.",
		StartTime:       time.Date(2030, time.September, 10, 19, 0, 0, 0, time.UTC),
		DurationMinutes: 90,
		Timezone:        "America/New_York",
		Status:          models.ScheduledEventStatusScheduled,
		CreatedAt:       &createdAt,
		UpdatedAt:       &updatedAt,
	})
	if err != nil {
		t.Fatal(err)
	}

	data, err := SerializeEvent(event)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}

	assertEventValue(t, got, "type", "Event")
	assertEventValue(t, got, "id", "https://my.cool.site.biz/federation/event/weekly-show")
	assertEventValue(t, got, "name", "Weekly show")
	assertEventValue(t, got, "content", "A scheduled example stream.")
	assertEventValue(t, got, "startTime", "2030-09-10T19:00:00Z")
	assertEventValue(t, got, "endTime", "2030-09-10T20:30:00Z")
	assertEventValue(t, got, "duration", "PT1H30M")
	assertEventValue(t, got, "timezone", "America/New_York")
	assertEventValue(t, got, "eventStatus", "EventScheduled")
	assertEventValue(t, got, "joinMode", "none")
	assertEventValue(t, got, "url", "https://my.cool.site.biz/schedule/weekly-show")
	assertEventValue(t, got, "published", "2030-09-08T12:00:00Z")
	assertEventValue(t, got, "updated", "2030-09-08T13:00:00Z")
	assertEventRecipients(t, got, "to", "https://www.w3.org/ns/activitystreams#Public")
	assertEventRecipients(t, got, "cc", "https://my.cool.site.biz/federation/user/streamer/followers")

	contexts, ok := got["@context"].([]interface{})
	if !ok || len(contexts) != 3 || contexts[1] != eventContext || contexts[2] != schemaContext {
		t.Fatalf("@context = %#v", got["@context"])
	}
	organizers, ok := got["organizers"].(map[string]interface{})
	if !ok || organizers["type"] != "OrganizersCollection" || organizers["totalItems"] != float64(1) {
		t.Fatalf("organizers = %#v", got["organizers"])
	}
	location, ok := got["location"].(map[string]interface{})
	if !ok || location["type"] != "VirtualLocation" || location["url"] != "https://my.cool.site.biz/schedule/weekly-show" {
		t.Fatalf("location = %#v", got["location"])
	}

	attachments, ok := got["attachment"].([]interface{})
	if !ok || len(attachments) != 1 {
		t.Fatalf("attachment = %#v", got["attachment"])
	}
	attachment, ok := attachments[0].(map[string]interface{})
	if !ok || attachment["type"] != "Image" ||
		attachment["url"] != "https://my.cool.site.biz/logo/external" ||
		attachment["mediaType"] != "image/svg+xml" {
		t.Fatalf("attachment = %#v", attachments[0])
	}
}

func assertEventRecipients(t *testing.T, event map[string]interface{}, key, want string) {
	t.Helper()
	recipients, ok := event[key].([]interface{})
	if !ok || len(recipients) != 1 || recipients[0] != want {
		t.Errorf("%s = %#v, want [%q]", key, event[key], want)
	}
}

func assertEventValue(t *testing.T, event map[string]interface{}, key string, want interface{}) {
	t.Helper()
	if got := event[key]; got != want {
		t.Errorf("%s = %#v, want %#v", key, got, want)
	}
}
