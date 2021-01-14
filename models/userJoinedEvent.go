package models

import "time"

// NameChangeEvent represents a user changing their name in chat.
type UserJoinedEvent struct {
	Username  string    `json:"username"`
	Type      EventType `json:"type"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
