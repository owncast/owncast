package chat

import (
	"time"

	"github.com/owncast/owncast/models"
	"github.com/teris-io/shortid"
)

// UserMessageEvent is an inbound message from a user.
type UserMessageEvent struct {
	models.Event
	models.UserEvent
	MessageEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *UserMessageEvent) GetBroadcastPayload() models.EventPayload {
	return models.EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"user":      e.User,
		"type":      models.MessageSent,
		"visible":   e.HiddenAt == nil,
	}
}

// GetMessageType will return the event type for this message.
func (e *UserMessageEvent) GetMessageType() models.EventType {
	return models.MessageSent
}

// SetDefaults will set default properties of all inbound events.
func (e *UserMessageEvent) SetDefaults() {
	e.ID = shortid.MustGenerate()
	e.Timestamp = time.Now()
	e.RenderAndSanitizeMessageBody()
}
