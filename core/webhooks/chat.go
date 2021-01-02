package webhooks

import (
	"github.com/owncast/owncast/models"
)

// TODO: relocate to shared location that we and chat package can use
const (
	CHAT             = "CHAT"
	NAMECHANGE       = "NAME_CHANGE"
	PING             = "PING"
	PONG             = "PONG"
	VISIBILITYUPDATE = "VISIBILITY-UPDATE"
)

func SendChatEvent(chatEvent models.ChatEvent) {
	webhookEvent := WebhookEvent{}

	// TODO: handle errors here instead of returning
	switch chatEvent.MessageType {
	case CHAT:
		webhookEvent.Type = MessageSent
	case NAMECHANGE:
		webhookEvent.Type = UserNameChanged
	case VISIBILITYUPDATE:
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
