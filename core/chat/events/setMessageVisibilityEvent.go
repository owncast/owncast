package events

// SetMessageVisibilityEvent is the event fired when one or more message
// visibilities are changed.
type SetMessageVisibilityEvent struct {
	Event
	UserMessageEvent
	MessageIDs []string
	Visible    bool
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *SetMessageVisibilityEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"type":      VisibiltyUpdate,
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"ids":       e.MessageIDs,
		"visible":   e.Visible,
	}
}
