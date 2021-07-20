package events

// UserJoinedEvent is the event fired when a user joins chat.
type UserJoinedEvent struct {
	Event
	UserEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *UserJoinedEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"type":      UserJoined,
		"id":        e.Id,
		"timestamp": e.Timestamp,
		"user":      e.User,
	}
}
