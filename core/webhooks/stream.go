package webhooks

import (
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

func SendStreamStatusEvent(eventType models.EventType) {
	SendEventToWebhooks(WebhookEvent{
		Type: eventType,
		EventData: map[string]interface{}{
			"name":        data.GetServerName(),
			"summary":     data.GetServerSummary(),
			"streamTitle": data.GetStreamTitle(),
		},
	})
}
