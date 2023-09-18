package webhooks

import (
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/teris-io/shortid"
)

// SendStreamStatusEvent will send all webhook destinations the current stream status.
func SendStreamStatusEvent(eventType models.EventType, streamID string) {
	sendStreamStatusEvent(eventType, shortid.MustGenerate(), streamID, time.Now())
}

func sendStreamStatusEvent(eventType models.EventType, id, streamID string, timestamp time.Time) {
	SendEventToWebhooks(WebhookEvent{
		Type: eventType,
		EventData: map[string]interface{}{
			"id":          id,
			"name":        data.GetServerName(),
			"summary":     data.GetServerSummary(),
			"streamTitle": data.GetStreamTitle(),
			"status":      getStatus(),
			"timestamp":   timestamp,
			"streamID":    streamID,
		},
	})
}
