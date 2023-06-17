package webhooks

import (
	"time"

	"github.com/owncast/owncast/models"
	"github.com/teris-io/shortid"
)

// SendStreamStatusEvent will send all webhook destinations the current stream status.
func (w *LiveWebhookManager) SendStreamStatusEvent(eventType models.EventType) {
	w.sendStreamStatusEvent(eventType, shortid.MustGenerate(), time.Now())
}

func (w *LiveWebhookManager) sendStreamStatusEvent(eventType models.EventType, id string, timestamp time.Time) {
	w.SendEventToWebhooks(WebhookEvent{
		Type: eventType,
		EventData: map[string]interface{}{
			"id":          id,
			"name":        data.GetServerName(),
			"summary":     data.GetServerSummary(),
			"streamTitle": data.GetStreamTitle(),
			"status":      w.getStatus(),
			"timestamp":   timestamp,
		},
	})
}
