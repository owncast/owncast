package events

// UserDisabledEvent is the event fired when a user is banned/blocked and disconnected from chat.
type UserDisabledEvent struct {
	Event
	UserEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *UserDisabledEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"type":      Event_Error_User_Disabled,
		"id":        e.Id,
		"timestamp": e.Timestamp,
		"user":      e.User,
	}
}
