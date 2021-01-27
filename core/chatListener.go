package core

import (
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/models"
)

// ChatListenerImpl the implementation of the chat client.
type ChatListenerImpl struct{}

// ClientAdded is for when a client is added the system.
func (cl ChatListenerImpl) ClientAdded(client models.Client) {
	SetClientActive(client)
}

// ClientRemoved is for when a client disconnects/is removed.
func (cl ChatListenerImpl) ClientRemoved(clientID string) {
	RemoveClient(clientID)
}

// MessageSent is for when a message is sent.
func (cl ChatListenerImpl) MessageSent(message models.ChatEvent) {
}

// SendMessageToChat sends a message to the chat server.
func SendMessageToChat(message models.ChatEvent) error {
	chat.SendMessage(message)

	return nil
}

// GetAllChatMessages gets all of the chat messages.
func GetAllChatMessages(filtered bool) []models.ChatEvent {
	return chat.GetMessages(filtered)
}
