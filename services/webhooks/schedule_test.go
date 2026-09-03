package webhooks

import (
	"context"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"
	"github.com/owncast/owncast/services/dispatcher"
)

func TestSendScheduledEvent(t *testing.T) {
	event := models.ScheduledEvent{
		ID:              "event-1",
		Name:            "Weekly stream",
		Description:     "A scheduled stream",
		ReminderMessage: "Join us in ten minutes!",
		StartTime:       time.Date(2030, time.September, 8, 18, 0, 0, 0, time.UTC),
		DurationMinutes: 60,
		Timezone:        "America/Los_Angeles",
		Status:          models.ScheduledEventStatusScheduled,
	}
	timestamp := time.Date(2030, time.September, 8, 17, 50, 0, 0, time.UTC)

	checkPayload(t, models.ScheduledEvents, func() {
		testSvc.sendScheduledEvent(event, models.ScheduledEventWarning, timestamp)
	}, `{
		"action": "10-minute-warning",
		"event": {
			"description": "A scheduled stream",
			"reminderMessage": "Join us in ten minutes!",
			"durationMinutes": 60,
			"id": "event-1",
			"name": "Weekly stream",
			"startTime": "2030-09-08T18:00:00Z",
			"status": "scheduled",
			"timezone": "America/Los_Angeles"
		},
		"serverURL": "http://localhost:8080",
		"status": {
			"lastConnectTime": null,
			"lastDisconnectTime": null,
			"online": true,
			"overallMaxViewerCount": 420,
			"sessionMaxViewerCount": 69,
			"streamTitle": "my stream",
			"versionNumber": "1.2.3",
			"viewerCount": 5
		},
		"timestamp": "2030-09-08T17:50:00Z"
	}`)
}

func TestScheduledEventUsesDefaultReminderMessage(t *testing.T) {
	configRepository := configrepository.New(testDatastore)
	if err := configRepository.SetScheduleReminderMessage("Default event reminder"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := configRepository.SetScheduleReminderMessage(""); err != nil {
			t.Error(err)
		}
	})

	events := dispatcher.New()
	received := make(chan string, 1)
	events.AddListener(func(_ context.Context, event dispatcher.Event) {
		data := event.Payload.(WebhookEvent).EventData.(*WebhookScheduledEventData)
		received <- data.Event.ReminderMessage
	})
	service := New(Deps{
		GetStatus:         fakeGetStatus,
		ConfigRepository:  configRepository,
		WebhookRepository: webhookrepository.New(testDatastore),
		Events:            events,
	})

	service.sendScheduledEvent(models.ScheduledEvent{}, models.ScheduledEventWarning, time.Now())

	if got := <-received; got != "Default event reminder" {
		t.Errorf("reminder message = %q, want %q", got, "Default event reminder")
	}
}
