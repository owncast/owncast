package models

import "github.com/owncast/owncast/utils"

// Follower is our internal representation of a single follower within Owncast.
type Follower struct {
	// ActorIRI is the IRI of the remote actor.
	ActorIRI string `json:"link"`
	// Inbox is the inbox URL of the remote follower
	Inbox string `json:"-"`
	// Name is the display name of the follower.
	Name string `json:"name"`
	// Username is the account username of the remote actor.
	Username string `json:"username"`
	// Image is the avatar image of the follower.
	Image string `json:"image"`
	// Timestamp is when this follow request was created.
	Timestamp utils.NullTime `json:"timestamp,omitempty"`
	// DisabledAt is when this follower was rejected or disabled.
	DisabledAt utils.NullTime `json:"disabledAt,omitempty"`
}
