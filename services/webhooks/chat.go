package webhooks

import (
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat/events"
)

// SendChatEvent dispatches a chat-message event to webhook destinations.
func (s *Service) SendChatEvent(chatEvent *events.UserMessageEvent) {
	webhookEvent := WebhookEvent{
		Type: chatEvent.GetMessageType(),
		EventData: &WebhookChatMessage{
			BaseWebhookData: BaseWebhookData{
				Status:    s.getStatus(),
				ServerURL: s.serverURL(),
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

	s.SendEventToWebhooks(webhookEvent)
}

// SendChatEventUsernameChanged dispatches a username-changed event.
func (s *Service) SendChatEventUsernameChanged(event events.NameChangeEvent) {
	webhookEvent := WebhookEvent{
		Type: models.UserNameChanged,
		EventData: &WebhookNameChangeEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    s.getStatus(),
				ServerURL: s.serverURL(),
			},
			ID:        event.ID,
			Timestamp: event.Timestamp,
			User:      event.User,
			NewName:   event.NewName,
		},
	}

	s.SendEventToWebhooks(webhookEvent)
}

// SendChatEventUserJoined dispatches a user-joined event.
func (s *Service) SendChatEventUserJoined(event events.UserJoinedEvent) {
	webhookEvent := WebhookEvent{
		Type: models.UserJoined,
		EventData: &WebhookUserJoinedEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    s.getStatus(),
				ServerURL: s.serverURL(),
			},
			ID:        event.ID,
			Timestamp: event.Timestamp,
			User:      event.User,
		},
	}

	s.SendEventToWebhooks(webhookEvent)
}

// SendChatEventUserParted dispatches a user-parted event.
func (s *Service) SendChatEventUserParted(event events.UserPartEvent) {
	webhookEvent := WebhookEvent{
		Type: events.UserParted,
		EventData: &WebhookUserPartEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    s.getStatus(),
				ServerURL: s.serverURL(),
			},
			ID:        event.ID,
			Timestamp: event.Timestamp,
			User:      event.User,
		},
	}

	s.SendEventToWebhooks(webhookEvent)
}

// SendChatEventSetMessageVisibility dispatches a message-visibility-changed event.
func (s *Service) SendChatEventSetMessageVisibility(event events.SetMessageVisibilityEvent) {
	webhookEvent := WebhookEvent{
		Type: models.VisibiltyToggled,
		EventData: &WebhookVisibilityToggleEventData{
			BaseWebhookData: BaseWebhookData{
				Status:    s.getStatus(),
				ServerURL: s.serverURL(),
			},
			ID:         event.ID,
			Timestamp:  event.Timestamp,
			User:       event.User,
			Visible:    event.Visible,
			MessageIDs: event.MessageIDs,
		},
	}

	s.SendEventToWebhooks(webhookEvent)
}
