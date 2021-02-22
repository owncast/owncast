package models

// ChatListener represents the listener for the chat server.
type ChatListener interface {
	ClientAdded(client Client)
	ClientRemoved(clientID string)
	MessageSent(message ChatEvent)
	IsStreamConnected() bool
}
