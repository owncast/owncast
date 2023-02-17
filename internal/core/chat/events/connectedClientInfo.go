package events

import (
	"github.com/owncast/owncast/internal/core/user"
)

// ConnectedClientInfo represents the information about a connected client.
type ConnectedClientInfo struct {
	Event
	User *user.User `json:"user"`
}
