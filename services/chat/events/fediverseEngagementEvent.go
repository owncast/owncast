package events

// FediverseEngagementEvent is a message displayed in chat on representing an action on the Fediverse.
type FediverseEngagementEvent struct {
	Event
	MessageEvent
	Image           *string `json:"image"`
	Link            string  `json:"link"`
	UserAccountName string  `json:"title"`
	// ServerName is the configured server display name baked in at
	// construction time. Populated by chat.Service when the event is
	// created so GetBroadcastPayload no longer needs to read from a
	// package-level configRepository. Excluded from JSON because it's a
	// transport detail, not a persisted field.
	ServerName string `json:"-"`
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *FediverseEngagementEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"id":                e.ID,
		payloadKeyTimestamp: e.Timestamp,
		payloadKeyBody:      e.Body,
		"image":             e.Image,
		payloadKeyType:      e.Type,
		"title":             e.UserAccountName,
		"link":              e.Link,
		payloadKeyUser: EventPayload{
			"displayName": e.ServerName,
		},
	}
}

// GetMessageType will return the event type for this message.
func (e *FediverseEngagementEvent) GetMessageType() EventType {
	return e.Type
}
