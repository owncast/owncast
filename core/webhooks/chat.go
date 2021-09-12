package webhooks

import (
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/models"
)

// SendChatEvent will send a chat event to webhook destinations.
func SendChatEvent(chatEvent *events.UserMessageEvent) {
	webhookEvent := WebhookEvent{
		Type: chatEvent.GetMessageType(),
		EventData: &WebhookChatMessage{
			User:      chatEvent.User,
			Body:      chatEvent.Body,
			RawBody:   chatEvent.RawBody,
			ID:        chatEvent.ID,
			Visible:   chatEvent.HiddenAt == nil,
			Timestamp: &chatEvent.Timestamp,
		},
	}

	SendEventToWebhooks(webhookEvent)
}

// SendChatEventUsernameChanged will send a username changed event to webhook destinations.
func SendChatEventUsernameChanged(event events.NameChangeEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserNameChanged,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}

// SendChatEventUserJoined sends a webhook notifying that a user has joined.
func SendChatEventUserJoined(event events.UserJoinedEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserJoined,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}
