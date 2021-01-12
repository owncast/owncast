package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

type WebhookEvent struct {
	Type      models.EventType `json:"type"` // messageSent | userJoined | userNameChange
	EventData interface{}      `json:"eventData,omitempty"`
}

type WebhookChatMessage struct {
	Author    string     `json:"author,omitempty"`
	Body      string     `json:"body,omitempty"`
	RawBody   string     `json:"rawBody,omitempty"`
	ID        string     `json:"id,omitempty"`
	Visible   bool       `json:"visible"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

func SendEventToWebhooks(payload WebhookEvent) {
	webhooks := data.GetWebhooksForEvent(payload.Type)

	for _, webhook := range webhooks {
		log.Debugf("Event %s sent to Webhook %s", payload.Type, webhook.Url)
		if err := sendWebhook(webhook.Url, payload); err != nil {
			log.Errorf("Event: %s failed to send to webhook: %s  Error: %s", payload.Type, webhook.Url, err)
		}
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

	data.SetWebhookAsUsed(url)

	return nil
}
