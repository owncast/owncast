package core

import (
	log "github.com/sirupsen/logrus"

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
	log.Printf("Message sent to all: %s", message.String())
}
