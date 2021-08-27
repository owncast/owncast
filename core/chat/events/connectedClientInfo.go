package events

import "github.com/owncast/owncast/core/user"

type ConnectedClientInfo struct {
	Event
	User *user.User `json:"user"`
}
