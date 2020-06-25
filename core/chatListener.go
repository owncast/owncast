package core

import (
	"errors"

	"github.com/gabek/owncast/core/chat"
	"github.com/gabek/owncast/models"
)

//ChatListenerImpl the implementation of the chat client
type ChatListenerImpl struct{}

//ClientAdded is for when a client is added the system
func (cl ChatListenerImpl) ClientAdded(clientID string) {
	SetClientActive(clientID)
}

//ClientRemoved is for when a client disconnects/is removed
func (cl ChatListenerImpl) ClientRemoved(clientID string) {
	RemoveClient(clientID)
}

//MessageSent is for when a message is sent
func (cl ChatListenerImpl) MessageSent(message models.ChatMessage) {
}

//SendMessageToChat sends a message to the chat server
func SendMessageToChat(message models.ChatMessage) error {
	if !message.Valid() {
		return errors.New("invalid chat message; id, author, and body are required")
	}

	chat.SendMessage(message)

	return nil
}

//GetAllChatMessages gets all of the chat messages
func GetAllChatMessages() []models.ChatMessage {
	return chat.GetMessages()
}
