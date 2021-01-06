package models

// NameChangeEvent represents a user changing their name in chat.
type NameChangeEvent struct {
	OldName string    `json:"oldName"`
	NewName string    `json:"newName"`
	Image   string    `json:"image"`
	Type    EventType `json:"type"`
	ID      string    `json:"id"`
}
