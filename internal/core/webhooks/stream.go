package webhooks

import (
	"time"

	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/internal/core/data"
	"github.com/owncast/owncast/internal/models"
)

// SendStreamStatusEvent will send all webhook destinations the current stream status.
func SendStreamStatusEvent(eventType models.EventType) {
	sendStreamStatusEvent(eventType, shortid.MustGenerate(), time.Now())
}

func sendStreamStatusEvent(eventType models.EventType, id string, timestamp time.Time) {
	SendEventToWebhooks(WebhookEvent{
		Type: eventType,
		EventData: map[string]interface{}{
			"id":          id,
			"name":        data.GetServerName(),
			"summary":     data.GetServerSummary(),
			"streamTitle": data.GetStreamTitle(),
			"timestamp":   timestamp,
		},
	})
}
