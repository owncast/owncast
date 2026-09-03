package webhooks

import (
	"testing"
	"time"

	"github.com/owncast/owncast/models"
)

func TestSendScheduledEvent(t *testing.T) {
	event := models.ScheduledEvent{
		ID:              "event-1",
		Name:            "Weekly stream",
		Description:     "A scheduled stream",
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
