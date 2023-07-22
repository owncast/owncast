package chat

import "github.com/owncast/owncast/models"

// UserPartEvent is the event fired when a user leaves chat.
type UserPartEvent struct {
	models.Event
	models.UserEvent
}

// GetBroadcastPayload will return the object to send to all chat users.
func (e *UserPartEvent) GetBroadcastPayload() models.EventPayload {
	return models.EventPayload{
		"type":      models.UserParted,
		"id":        e.ID,
		"timestamp": e.Timestamp,
		"user":      e.User,
	}
}
