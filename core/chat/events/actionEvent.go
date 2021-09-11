package events

type ActionEvent struct {
	Event
	MessageEvent
}

// ActionEvent will return the object to send to all chat users.
func (e *ActionEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"type":      e.GetMessageType(),
	}
}

func (e *ActionEvent) GetMessageType() EventType {
	return ChatActionSent
}
