package webhooks

import (
	"github.com/owncast/owncast/models"
)

func SendChatEvent(chatEvent models.ChatEvent) {
	webhookEvent := WebhookEvent{}

	// TODO: handle errors here instead of returning
	switch chatEvent.MessageType {
	case models.MessageSent:
		webhookEvent.Type = MessageSent
	case models.UserNameChanged:
		webhookEvent.Type = UserNameChanged
	case models.VisibiltyToggled:
		webhookEvent.Type = VisibiltyToggled
	}

	webhookEvent.EventData = &WebhookChatMessage{
		Author:    chatEvent.Author,
		Body:      chatEvent.Body,
		RawBody:   chatEvent.RawBody,
		ID:        chatEvent.ID,
		Visible:   chatEvent.Visible,
		Timestamp: &chatEvent.Timestamp,
	}

	SendEventToWebhooks(webhookEvent)
}

func SendChatEventUsernameChanged(event models.NameChangeEvent) {
	webhookEvent := WebhookEvent{
		Type:      UserNameChanged,
		EventData: event,
	}

	SendEventToWebhooks(webhookEvent)
}
