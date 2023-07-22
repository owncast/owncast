package models

import "time"

// UserEvent is an event with an associated user.
type UserEvent struct {
	User     *User      `json:"user"`
	HiddenAt *time.Time `json:"hiddenAt,omitempty"`
	ClientID uint       `json:"clientId,omitempty"`
}
