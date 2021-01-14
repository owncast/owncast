package models

import "time"

// UserJoinedEvent represents a user joining the chat.
type UserJoinedEvent struct {
	Username  string    `json:"username"`
	Type      EventType `json:"type"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
