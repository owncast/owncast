package webhooks

import (
	"github.com/owncast/owncast/models"
)

func SendChatEvent(chatEvent models.ChatEvent) {
	webhookEvent := WebhookEvent{
		Type: chatEvent.MessageType,
		EventData: &WebhookChatMessage{
			Author:    chatEvent.Author,
			Body:      chatEvent.Body,
			RawBody:   chatEvent.RawBody,
			ID:        chatEvent.ID,
			Visible:   chatEvent.Visible,
			Timestamp: &chatEvent.Timestamp,
		},
	}

	SendEventToWebhooks(webhookEvent)
}

func SendChatEventUsernameChanged(event models.NameChangeEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserNameChanged,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}

func SendChatEventUserJoined(event models.UserJoinedEvent) {
	webhookEvent := WebhookEvent{
		Type:      models.UserNameChanged,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}
