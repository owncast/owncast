package models

import "github.com/owncast/owncast/utils"

// Follower represents a single remote follower of the instance.
type Follower struct {
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Image     string         `json:"image"`
	Link      string         `json:"link"`
	Inbox     string         `json:"-"`
	Timestamp utils.NullTime `json:"timestamp,omitempty"`
}
