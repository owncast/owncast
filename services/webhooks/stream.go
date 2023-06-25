package webhooks

import (
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/teris-io/shortid"
)

var configRepository = configrepository.Get()

// SendStreamStatusEvent will send all webhook destinations the current stream status.
func (w *LiveWebhookManager) SendStreamStatusEvent(eventType models.EventType) {
	w.sendStreamStatusEvent(eventType, shortid.MustGenerate(), time.Now())
}

func (w *LiveWebhookManager) sendStreamStatusEvent(eventType models.EventType, id string, timestamp time.Time) {
	w.SendEventToWebhooks(WebhookEvent{
		Type: eventType,
		EventData: map[string]interface{}{
			"id":          id,
			"name":        configRepository.GetServerName(),
			"summary":     configRepository.GetServerSummary(),
			"streamTitle": configRepository.GetStreamTitle(),
			"status":      w.getStatus(),
			"timestamp":   timestamp,
		},
	})
}
