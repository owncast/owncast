package handlers

import (
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/storage/chatrepository"
	"github.com/owncast/owncast/storage/configrepository"
)

type Handlers struct {
	configRepository *configrepository.SqlConfigRepository
	chatService      *chat.Chat
	chatRepository   *chatrepository.ChatRepository
}

// New creates a new instances of web server handlers.
func New() *Handlers {
	return &Handlers{
		configRepository: configrepository.Get(),
		chatService:      chat.Get(),
		chatRepository:   chatrepository.Get(),
	}
}
