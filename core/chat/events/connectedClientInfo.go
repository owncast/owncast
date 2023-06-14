package events

import "github.com/owncast/owncast/models"

// ConnectedClientInfo represents the information about a connected client.
type ConnectedClientInfo struct {
	Event
	User *models.User `json:"user"`
}
