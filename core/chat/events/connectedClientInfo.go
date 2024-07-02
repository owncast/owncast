package events

import "github.com/owncast/owncast/models"

// ConnectedClientInfo represents the information about a connected client.
type ConnectedClientInfo struct {
	User *models.User `json:"user"`
	Event
}
