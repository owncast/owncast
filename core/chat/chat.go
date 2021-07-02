package chat

import (
	"fmt"
	"net/http"

	"github.com/owncast/owncast/core/chat/events"
)

func Start() error {
	setupPersistence()

	_server = NewChat()

	go _server.Run()

	fmt.Println("Chat server started with max connection count of", _server.maxClientCount)

	return nil
}
func GetClient(id uint) *ChatClient {
	return _server.clients[id]
}

func GetClients() []*ChatClient {
	clients := []*ChatClient{}

	for _, client := range _server.clients {
		clients = append(clients, client)
	}

	return clients
}

func SendSystemMessage(text string, ephemeral bool) error {
	return nil
}

func SendSystemAction(text string, ephemeral bool) error {
	return nil
}

func Broadcast(event events.OutboundEvent) error {
	// TODO: Save event to database
	return _server.Broadcast(event.GetBroadcastPayload())
}

func HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	_server.HandleClientConnection(w, r)
}

// DisconnectUser will forcefully disconnect all clients belonging to a user by ID.
func DisconnectUser(userID string) {
	_server.DisconnectUser(userID)
}

// // GetMessages gets all of the messages.
// func GetMessages() []events.Event {
// 	if _server == nil {
// 		return []events.Event{}
// 	}

// 	return getChatHistory()
// }

// func GetModerationChatMessages() []events.Event {
// 	return getChatModerationHistory()
// }
