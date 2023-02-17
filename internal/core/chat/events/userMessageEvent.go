package events

// UserMessageEvent is an inbound message from a user.
type UserMessageEvent struct {
	Event
	UserEvent
	MessageEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *UserMessageEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"user":      e.User,
		"type":      MessageSent,
		"visible":   e.HiddenAt == nil,
	}
}

// GetMessageType will return the event type for this message.
func (e *UserMessageEvent) GetMessageType() EventType {
	return MessageSent
}
