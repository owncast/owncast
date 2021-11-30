package events

import "github.com/owncast/owncast/core/data"

// FediverseEngagementEvent is a message displayed in chat on representing an action on the Fediverse.
type FediverseEngagementEvent struct {
	Event
	MessageEvent
	Image    string `json:"image"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Link     string `json:"link"`
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *FediverseEngagementEvent) GetBroadcastPayload() EventPayload {
	return EventPayload{
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"body":      e.Body,
		"title":     e.Title,
		"subtitle":  e.Subtitle,
		"image":     e.Image,
		"type":      e.Event.Type,
		"user": EventPayload{
			"displayName": data.GetServerName(),
		},
	}
}

// GetMessageType will return the event type for this message.
func (e *FediverseEngagementEvent) GetMessageType() EventType {
	return e.Event.Type
}
