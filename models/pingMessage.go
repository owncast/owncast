package models

// PingMessage represents a ping message between the client and server.
type PingMessage struct {
	MessageType ChatMessageType `json:"type"`
}
