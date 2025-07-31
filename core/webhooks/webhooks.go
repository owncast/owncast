package webhooks

import (
	"sync"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"
)

// BaseWebhookData contains common fields shared across all webhook event data.
type BaseWebhookData struct {
	Status    models.Status `json:"status"`
	ServerURL string        `json:"serverURL,omitempty"`
}

// WebhookEvent represents an event sent as a webhook.
type WebhookEvent struct {
	EventData interface{}      `json:"eventData,omitempty"`
	Type      models.EventType `json:"type"` // messageSent | userJoined | userNameChange
}

// WebhookChatMessage represents a single chat message sent as a webhook payload.
type WebhookChatMessage struct {
	BaseWebhookData
	User      *models.User `json:"user,omitempty"`
	Timestamp *time.Time   `json:"timestamp,omitempty"`
	Body      string       `json:"body,omitempty"`
	RawBody   string       `json:"rawBody,omitempty"`
	ID        string       `json:"id,omitempty"`
	ClientID  uint         `json:"clientId,omitempty"`
	Visible   bool         `json:"visible"`
}

// WebhookUserJoinedEventData represents user joined event data sent as a webhook payload.
type WebhookUserJoinedEventData struct {
	BaseWebhookData
	ID        string       `json:"id"`
	Timestamp time.Time    `json:"timestamp"`
	User      *models.User `json:"user"`
}

// WebhookUserPartEventData represents user parted event data sent as a webhook payload.
type WebhookUserPartEventData struct {
	BaseWebhookData
	ID        string       `json:"id"`
	Timestamp time.Time    `json:"timestamp"`
	User      *models.User `json:"user"`
}

// WebhookNameChangeEventData represents name change event data sent as a webhook payload.
type WebhookNameChangeEventData struct {
	BaseWebhookData
	ID        string       `json:"id"`
	Timestamp time.Time    `json:"timestamp"`
	User      *models.User `json:"user"`
	NewName   string       `json:"newName"`
}

// WebhookVisibilityToggleEventData represents message visibility toggle event data sent as a webhook payload.
type WebhookVisibilityToggleEventData struct {
	BaseWebhookData
	ID         string       `json:"id"`
	Timestamp  time.Time    `json:"timestamp"`
	User       *models.User `json:"user"`
	Visible    bool         `json:"visible"`
	MessageIDs []string     `json:"ids"`
}

// SendEventToWebhooks will send a single webhook event to all webhook destinations.
func SendEventToWebhooks(payload WebhookEvent) {
	sendEventToWebhooks(payload, nil)
}

func sendEventToWebhooks(payload WebhookEvent, wg *sync.WaitGroup) {
	webhooksRepo := webhookrepository.Get()
	webhooks := webhooksRepo.GetWebhooksForEvent(payload.Type)

	for _, webhook := range webhooks {
		// Use wg to track the number of notifications to be sent.
		if wg != nil {
			wg.Add(1)
		}
		addToQueue(webhook, payload, wg)
	}
}

func getServerURL() string {
	configRepo := configrepository.Get()
	return configRepo.GetServerURL()
}
