package events

// SystemMessageEvent
type ActionEvent struct {
	Event
	MessageEvent
}

// SystemMessageEvent will return the object to send to all chat users.
func (e *ActionEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"id":        e.Id,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"type":      e.GetMessageType(),
	}
}

func (e *ActionEvent) GetMessageType() EventType {
	return Event_ChatActionSent
}
