package webhooks

import (
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/models"
)

// SendChatEvent will send a chat event to webhook destinations.
func (w *LiveWebhookManager) SendChatEvent(chatEvent *events.UserMessageEvent) {
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

	w.SendEventToWebhooks(webhookEvent)
}

// SendChatEventUsernameChanged will send a username changed event to webhook destinations.
func (w *LiveWebhookManager) SendChatEventUsernameChanged(event events.NameChangeEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserNameChanged,
		EventData: event,
	}

	w.SendEventToWebhooks(webhookEvent)
}

// SendChatEventUserJoined sends a webhook notifying that a user has joined.
func (w *LiveWebhookManager) SendChatEventUserJoined(event events.UserJoinedEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserJoined,
		EventData: event,
	}

	w.SendEventToWebhooks(webhookEvent)
}

// SendChatEventUserParted sends a webhook notifying that a user has parted.
func SendChatEventUserParted(event events.UserPartEvent) {
	webhookEvent := WebhookEvent{
		Type:      events.UserParted,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}

// SendChatEventSetMessageVisibility sends a webhook notifying that the visibility of one or more
// messages has changed.
func (w *LiveWebhookManager) SendChatEventSetMessageVisibility(event events.SetMessageVisibilityEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.VisibiltyToggled,
		EventData: event,
	}

	w.SendEventToWebhooks(webhookEvent)
}
