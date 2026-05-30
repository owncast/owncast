package events

// SystemMessageEvent is a message displayed in chat on behalf of the server.
type SystemMessageEvent struct {
	Event
	MessageEvent
	// ServerName is the configured server display name baked in at
	// construction time. Populated by chat.Service when the event is
	// created so GetBroadcastPayload no longer needs to read from a
	// package-level configRepository. Excluded from JSON because it's a
	// transport detail, not a persisted field.
	ServerName string `json:"-"`
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *SystemMessageEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"id":                e.ID,
		payloadKeyTimestamp: e.Timestamp,
		payloadKeyBody:      e.Body,
		payloadKeyType:      SystemMessageSent,
		payloadKeyUser: EventPayload{
			"displayName": e.ServerName,
		},
	}
}

// GetMessageType will return the event type for this message.
func (e *SystemMessageEvent) GetMessageType() EventType {
	return SystemMessageSent
}
