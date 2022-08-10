package events

// NameChangeEvent is received when a user changes their chat display name.
type NameChangeEvent struct {
	Event
	UserEvent
	NewName string `json:"newName"`
}

// ColorChangeEvent is received when a user changes their chat display color.
type ColorChangeEvent struct {
	Event
	UserEvent
	NewColor int `json:"newColor"`
}

// NameChangeBroadcast represents a user changing their chat display name.
type NameChangeBroadcast struct {
	Event
	OutboundEvent
	UserEvent
	Oldname string `json:"oldName"`
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *NameChangeBroadcast) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"user":      e.User,
		"oldName":   e.Oldname,
		"type":      UserNameChanged,
	}
}
