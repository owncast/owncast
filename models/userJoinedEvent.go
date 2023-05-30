package models

import "time"

// UserJoinedEvent represents an event when a user joins the chat.
type UserJoinedEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Username  string    `json:"username"`
	Type      EventType `json:"type"`
	ID        string    `json:"id"`
}
