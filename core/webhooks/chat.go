package webhooks

import (
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/models"
)

func SendChatEvent(chatEvent *events.UserMessageEvent) {
	webhookEvent := WebhookEvent{
		Type: chatEvent.GetMessageType(),
		EventData: &WebhookChatMessage{
			User:      chatEvent.User,
			Body:      chatEvent.Body,
			RawBody:   chatEvent.Body,
			ID:        chatEvent.Id,
			Visible:   chatEvent.HiddenAt == nil,
			Timestamp: &chatEvent.Timestamp,
		},
	}

	SendEventToWebhooks(webhookEvent)
}

func SendChatEventUsernameChanged(event events.NameChangeEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserNameChanged,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}

func SendChatEventUserJoined(event events.UserJoinedEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserNameChanged,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}
