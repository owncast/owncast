package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/user"

	"github.com/owncast/owncast/core/data"
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
		go func(payload WebhookEvent, webhook models.Webhook) {
			log.Debugf("Event %s sent to Webhook %s", payload.Type, webhook.URL)
			if err := sendWebhook(webhook.URL, payload); err != nil {
				log.Errorf("Event: %s failed to send to webhook: %s  Error: %s", payload.Type, webhook.URL, err)
			}
		}(payload, webhook)
	}
}

func sendWebhook(url string, payload WebhookEvent) error {
	jsonText, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonText))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := data.SetWebhookAsUsed(url); err != nil {
		log.Warnln(err)
	}

	return nil
}
