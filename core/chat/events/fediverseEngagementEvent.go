package events

import "github.com/owncast/owncast/core/data"

// FediverseEngagementEvent is a message displayed in chat on representing an action on the Fediverse.
type FediverseEngagementEvent struct {
	Event
	MessageEvent
	Image           *string `json:"image"`
	Link            string  `json:"link"`
	UserAccountName string  `json:"title"`
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *FediverseEngagementEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"image":     e.Image,
		"type":      e.Event.Type,
		"title":     e.UserAccountName,
		"link":      e.Link,
		"user": EventPayload{
			"displayName": data.GetServerName(),
		},
	}
}

// GetMessageType will return the event type for this message.
func (e *FediverseEngagementEvent) GetMessageType() EventType {
	return e.Event.Type
}
