package models

//ChatListener represents the listener for the chat server
type ChatListener interface {
	ClientAdded(clientID string)
	ClientRemoved(clientID string)
	MessageSent(message ChatMessage)
}
