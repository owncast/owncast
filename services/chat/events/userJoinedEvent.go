package events

// UserJoinedEvent is the event fired when a user joins chat.
type UserJoinedEvent struct {
	Event
	UserEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *UserJoinedEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		payloadKeyType:      UserJoined,
		"id":                e.ID,
		payloadKeyTimestamp: e.Timestamp,
		payloadKeyUser:      e.User,
	}
}
