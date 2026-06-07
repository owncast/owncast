// Package webhooks dispatches Owncast events (chat messages, stream
// status changes, fediverse engagement, …) to user-configured HTTP
// destinations via a bounded worker pool. Construct a *Service in
// main.go with the stream status callback and the followers repository,
// then call Start to launch the workers.
package webhooks

import (
	"context"
	"sync"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/dispatcher"
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

// WebhookFediverseEngagementFollowEventData represents Fediverse follow engagement
// event data sent as a webhook payload. Unlike stream-state events it does not
// embed BaseWebhookData; the live stream status is unrelated to a follow action.
// ServerURL is included so receivers can identify which Owncast instance sent
// the webhook.
type WebhookFediverseEngagementFollowEventData struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Image     string    `json:"image"`
	ServerURL string    `json:"serverURL,omitempty"`
}

// SendEventToWebhooks fans an event out to every configured webhook
// destination subscribed to that event type.
func (s *Service) SendEventToWebhooks(payload WebhookEvent) {
	s.sendEventToWebhooks(payload, nil)
}

// sendEventToWebhooks is the dispatch internal that also accepts an
// optional waitgroup so callers (tests, batched senders) can wait for
// all destinations to be drained.
func (s *Service) sendEventToWebhooks(payload WebhookEvent, wg *sync.WaitGroup) {
	// Publish to the shared dispatcher so in-process consumers (e.g. the
	// plugin host) see every event, even when no HTTP webhook destinations
	// are configured for it.
	if s.events != nil {
		s.events.Publish(context.Background(), dispatcher.Event{Type: payload.Type, Payload: payload})
	}

	webhooks := s.webhookRepository.GetWebhooksForEvent(payload.Type)

	for _, webhook := range webhooks {
		if wg != nil {
			wg.Add(1)
		}
		s.addToQueue(webhook, payload, wg)
	}
}

// serverURL is the helper used in event payloads.
func (s *Service) serverURL() string {
	return s.configRepository.GetServerURL()
}

// Followers exposes the followers repository the service was
// constructed with. Used by the fediverse-engagement event builder.
func (s *Service) Followers() followersrepository.FollowersRepository {
	return s.followers
}
