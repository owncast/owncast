package chat

import "github.com/owncast/owncast/models"

// SystemMessageEvent is a message displayed in chat on behalf of the server.
type SystemMessageEvent struct {
	models.Event
	MessageEvent
	DisplayName string
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *SystemMessageEvent) GetBroadcastPayload() models.EventPayload {
	return models.EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"type":      models.SystemMessageSent,
		"user": models.EventPayload{
			"displayName": e.DisplayName,
		},
	}
}

// GetMessageType will return the event type for this message.
func (e *SystemMessageEvent) GetMessageType() models.EventType {
	return models.SystemMessageSent
}
