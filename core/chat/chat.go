package chat

import (
	"net/http"
	"sort"

	"github.com/owncast/owncast/core/chat/events"
	log "github.com/sirupsen/logrus"
)

func Start() error {
	setupPersistence()

	_server = NewChat()

	go _server.Run()

	log.Traceln("Chat server started with max connection count of", _server.maxClientCount)

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

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].ConnectedAt.Before(clients[j].ConnectedAt)
	})

	return clients
}

func SendSystemMessage(text string, ephemeral bool) error {
	message := events.SystemMessageEvent{
		MessageEvent: events.MessageEvent{
			Body: text,
		},
	}
	message.SetDefaults()
	message.RenderBody()

	if err := Broadcast(&message); err != nil {
		log.Errorln("error sending system message", err)
	}

	if !ephemeral {
		saveEvent(message.Id, "system", message.Body, message.GetMessageType(), nil, message.Timestamp)
	}

	return nil
}

func SendSystemAction(text string, ephemeral bool) error {
	message := events.ActionEvent{
		MessageEvent: events.MessageEvent{
			Body: text,
		},
	}

	message.SetDefaults()
	message.RenderBody()

	if err := Broadcast(&message); err != nil {
		log.Errorln("error sending system chat action")
	}

	if !ephemeral {
		saveEvent(message.Id, "action", message.Body, message.GetMessageType(), nil, message.Timestamp)
	}

	return nil
}

func Broadcast(event events.OutboundEvent) error {
	return _server.Broadcast(event.GetBroadcastPayload())
}

func HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	_server.HandleClientConnection(w, r)
}

// DisconnectUser will forcefully disconnect all clients belonging to a user by ID.
func DisconnectUser(userID string) {
	_server.DisconnectUser(userID)
}
