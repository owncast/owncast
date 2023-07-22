package chat

import "github.com/owncast/owncast/models"

// ActionEvent represents an action that took place, not a chat message.
type ActionEvent struct {
	models.Event
	MessageEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *ActionEvent) GetBroadcastPayload() models.EventPayload {
	return models.EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"type":      e.GetMessageType(),
	}
}

// GetMessageType will return the type of message.
func (e *ActionEvent) GetMessageType() models.EventType {
	return models.ChatActionSent
}
