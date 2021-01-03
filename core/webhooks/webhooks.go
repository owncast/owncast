package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
)

type EventType string

const (
	MessageSent      EventType = "messageSent"
	UserJoined       EventType = "userJoined"
	UserNameChanged  EventType = "usernameChanged"
	VisibiltyToggled EventType = "visibilityToggle"
	StreamStarted    EventType = "streamStarted"
	StreamStopped    EventType = "streamStopped"
)

func (e EventType) toString() string {
	return fmt.Sprintf("%s", e)
}

type WebhookEvent struct {
	Type      EventType   `json:"type"` // messageSent | userJoined | userNameChange
	EventData interface{} `json:"eventData,omitempty"`
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
	for _, webhook := range config.Config.Webhooks {
		log.Debugf("Checking Webhook %s to send event: %s", webhook.Url, payload.Type)

		eventsAccepted := strings.Join(webhook.Events, ",")

		if strings.Contains(eventsAccepted, payload.Type.toString()) || eventsAccepted == "" {
			log.Debugf("Event sent to Webhook %s", webhook.Url)
			if err := sendWebhook(webhook.Url, payload); err != nil {
				log.Infof("Event: %s failed to send to webhook: %s  Error: %s", payload.Type, webhook.Url, err)
			}
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

	return nil
}
