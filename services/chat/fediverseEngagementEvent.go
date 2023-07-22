package chat

import "github.com/owncast/owncast/models"

// FediverseEngagementEvent is a message displayed in chat on representing an action on the Fediverse.
type FediverseEngagementEvent struct {
	models.Event
	MessageEvent
	Image           *string `json:"image"`
	Link            string  `json:"link"`
	UserAccountName string  `json:"title"`
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *FediverseEngagementEvent) GetBroadcastPayload() models.EventPayload {
	return models.EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"image":     e.Image,
		"type":      e.Event.Type,
		"title":     e.UserAccountName,
		"link":      e.Link,
		"user": models.EventPayload{
			"displayName": "Owncast",
		},
	}
}

// GetMessageType will return the event type for this message.
func (e *FediverseEngagementEvent) GetMessageType() models.EventType {
	return e.Event.Type
}
