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
			BaseWebhookData: BaseWebhookData{
				Status:    getStatus(),
				ServerURL: getServerURL(),
			},
			User:      chatEvent.User,
			Body:      chatEvent.Body,
			ClientID:  chatEvent.ClientID,
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
		Type: models.UserNameChanged,
		EventData: &WebhookNameChangeEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    getStatus(),
				ServerURL: getServerURL(),
			},
			ID:        event.ID,
			Timestamp: event.Timestamp,
			User:      event.User,
			NewName:   event.NewName,
		},
	}

	SendEventToWebhooks(webhookEvent)
}

// SendChatEventUserJoined sends a webhook notifying that a user has joined.
func SendChatEventUserJoined(event events.UserJoinedEvent) {
	webhookEvent := WebhookEvent{
		Type: models.UserJoined,
		EventData: &WebhookUserJoinedEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    getStatus(),
				ServerURL: getServerURL(),
			},
			ID:        event.ID,
			Timestamp: event.Timestamp,
			User:      event.User,
		},
	}

	SendEventToWebhooks(webhookEvent)
}

// SendChatEventUserParted sends a webhook notifying that a user has parted.
func SendChatEventUserParted(event events.UserPartEvent) {
	webhookEvent := WebhookEvent{
		Type: events.UserParted,
		EventData: &WebhookUserPartEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    getStatus(),
				ServerURL: getServerURL(),
			},
			ID:        event.ID,
			Timestamp: event.Timestamp,
			User:      event.User,
		},
	}

	SendEventToWebhooks(webhookEvent)
}

// SendChatEventSetMessageVisibility sends a webhook notifying that the visibility of one or more
// messages has changed.
func SendChatEventSetMessageVisibility(event events.SetMessageVisibilityEvent) {
	webhookEvent := WebhookEvent{
		Type: models.VisibiltyToggled,
		EventData: &WebhookVisibilityToggleEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    getStatus(),
				ServerURL: getServerURL(),
			},
			ID:         event.ID,
			Timestamp:  event.Timestamp,
			User:       event.User,
			Visible:    event.Visible,
			MessageIDs: event.MessageIDs,
		},
	}

	SendEventToWebhooks(webhookEvent)
}
