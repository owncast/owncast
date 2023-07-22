package models

// ConnectedClientInfo represents the information about a connected client.
type ConnectedClientInfo struct {
	Event
	User *User `json:"user"`
}
