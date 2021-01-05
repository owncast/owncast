package webhooks

import "github.com/owncast/owncast/models"

func SendStreamStatusEvent(eventType models.EventType) {
	SendEventToWebhooks(WebhookEvent{Type: eventType})
}
