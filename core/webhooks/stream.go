package webhooks

import (
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/teris-io/shortid"
)

func SendStreamStatusEvent(eventType models.EventType) {
	SendEventToWebhooks(WebhookEvent{
		Type: eventType,
		EventData: map[string]interface{}{
			"id":          shortid.MustGenerate(),
			"name":        data.GetServerName(),
			"summary":     data.GetServerSummary(),
			"streamTitle": data.GetStreamTitle(),
			"timestamp":   time.Now(),
		},
	})
}
