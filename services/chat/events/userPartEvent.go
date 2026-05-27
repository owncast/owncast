package events

// UserPartEvent is the event fired when a user leaves chat.
type UserPartEvent struct {
	Event
	UserEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *UserPartEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		payloadKeyType:      UserParted,
		"id":                e.ID,
		payloadKeyTimestamp: e.Timestamp,
		payloadKeyUser:      e.User,
	}
}
