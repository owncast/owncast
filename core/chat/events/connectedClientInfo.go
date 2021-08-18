package events

import "github.com/owncast/owncast/core/user"

type ConnecteClientInfo struct {
	Event
	User *user.User `json:"user"`
}
