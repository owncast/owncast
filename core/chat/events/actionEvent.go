package events

// ActionEvent represents an action that took place, not a chat message.
type ActionEvent struct {
	Event
	MessageEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *ActionEvent) GetBroadcastPayload(serverName string) EventPayload {
	return EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"type":      e.GetMessageType(),
		"user": EventPayload{
			"displayName": serverName,
		},
	}
}

// GetMessageType will return the type of message.
func (e *ActionEvent) GetMessageType() EventType {
	return ChatActionSent
}
