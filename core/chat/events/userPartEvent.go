package events

// UserPartEvent is the event fired when a user leaves chat.
type UserPartEvent struct {
	Event
	UserEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *UserPartEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"type":      UserParted,
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"user":      e.User,
	}
}
