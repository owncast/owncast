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
		"id":        e.Id,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"user":      e.User,
		"type":      Event_MessageSent,
		"visible":   e.HiddenAt == nil,
	}
}

func (e *UserMessageEvent) GetMessageType() EventType {
	return Event_MessageSent
}
