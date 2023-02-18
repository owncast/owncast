package webhooks

import (
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/models"
)

// SendChatEvent will send a chat event to webhook destinations.
func (s *Service) SendChatEvent(chatEvent *events.UserMessageEvent) {
	webhookEvent := WebhookEvent{
		Type: chatEvent.GetMessageType(),
		EventData: &WebhookChatMessage{
			User:      chatEvent.User,
			Body:      chatEvent.Body,
			ClientID:  chatEvent.ClientID,
			RawBody:   chatEvent.RawBody,
			ID:        chatEvent.ID,
			Visible:   chatEvent.HiddenAt == nil,
			Timestamp: &chatEvent.Timestamp,
		},
	}

	s.SendEventToWebhooks(webhookEvent)
}

// SendChatEventUsernameChanged will send a username changed event to webhook destinations.
func (s *Service) SendChatEventUsernameChanged(event events.NameChangeEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserNameChanged,
		EventData: event,
	}

	s.SendEventToWebhooks(webhookEvent)
}

// SendChatEventUserJoined sends a webhook notifying that a user has joined.
func (s *Service) SendChatEventUserJoined(event events.UserJoinedEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserJoined,
		EventData: event,
	}

	s.SendEventToWebhooks(webhookEvent)
}

// SendChatEventSetMessageVisibility sends a webhook notifying that the visibility of one or more
// messages has changed.
func (s *Service) SendChatEventSetMessageVisibility(event events.SetMessageVisibilityEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.VisibiltyToggled,
		EventData: event,
	}

	s.SendEventToWebhooks(webhookEvent)
}
