package models

//NameChangeEvent represents a user changing their name in chat
type NameChangeEvent struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
	Type    string `json:"type"`
}
