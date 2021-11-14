package webhooks

import (
	"time"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/models"
)

// WebhookEvent represents an event sent as a webhook.
type WebhookEvent struct {
	Type      models.EventType `json:"type"` // messageSent | userJoined | userNameChange
	EventData interface{}      `json:"eventData,omitempty"`
}

// WebhookChatMessage represents a single chat message sent as a webhook payload.
type WebhookChatMessage struct {
	User      *user.User `json:"user,omitempty"`
	ClientID  uint       `json:"clientId,omitempty"`
	Body      string     `json:"body,omitempty"`
	RawBody   string     `json:"rawBody,omitempty"`
	ID        string     `json:"id,omitempty"`
	Visible   bool       `json:"visible"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// SendEventToWebhooks will send a single webhook event to all webhook destinations.
func SendEventToWebhooks(payload WebhookEvent) {
	webhooks := data.GetWebhooksForEvent(payload.Type)

	for _, webhook := range webhooks {
		go addToQueue(webhook, payload)
	}
}
