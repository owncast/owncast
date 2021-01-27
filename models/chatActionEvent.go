package models

import "time"

// ChatActionEvent represents a generic action that took place by a chat user.
type ChatActionEvent struct {
	Username  string    `json:"username"`
	Type      EventType `json:"type"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message,omitempty"`
}
