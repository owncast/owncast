package chat

import (
	"errors"
	"net/http"
	"sort"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

var getStatus func() models.Status

// Start begins the chat server.
func Start(getStatusFunc func() models.Status) error {
	setupPersistence()

	getStatus = getStatusFunc
	_server = NewChat()

	go _server.Run()

	log.Traceln("Chat server started with max connection count of", _server.maxSocketConnectionLimit)

	return nil
}

// GetClientsForUser will return chat connections that are owned by a specific user.
func GetClientsForUser(userID string) ([]*Client, error) {
	clients := map[string][]*Client{}
	_server.clients.Range(func(k, v interface{}) bool {
		client := v.(*Client)
		clients[client.User.ID] = append(clients[client.User.ID], client)
		return true
	})

	if _, exists := clients[userID]; !exists {
		return nil, errors.New("no connections for user found")
	}

	return clients[userID], nil
}

// FindClientByID will return a single connected client by ID.
func FindClientByID(clientID uint) (*Client, bool) {
	client, found := _server.clients.Load(clientID)
	return client.(*Client), found
}

// GetClients will return all the current chat clients connected.
func GetClients() []*Client {
	clients := []*Client{}

	// Convert the keyed map to a slice.
	_server.clients.Range(func(k, v interface{}) bool {
		client := v.(*Client)
		clients = append(clients, client)
		return true
	})

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].ConnectedAt.Before(clients[j].ConnectedAt)
	})

	return clients
}

// SendSystemMessage will send a message string as a system message to all clients.
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
		saveEvent(message.ID, "system", message.Body, message.GetMessageType(), nil, message.Timestamp)
	}

	return nil
}

// SendSystemAction will send a system action string as an action event to all clients.
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
		saveEvent(message.ID, "action", message.Body, message.GetMessageType(), nil, message.Timestamp)
	}

	return nil
}

// SendAllWelcomeMessage will send the chat message to all connected clients.
func SendAllWelcomeMessage() {
	_server.sendAllWelcomeMessage()
}

// SendSystemMessageToClient will send a single message to a single connected chat client.
func SendSystemMessageToClient(clientID uint, text string) {
	if client, foundClient := FindClientByID(clientID); foundClient {
		_server.sendSystemMessageToClient(client, text)
	}
}

// Broadcast will send all connected clients the outbound object provided.
func Broadcast(event events.OutboundEvent) error {
	return _server.Broadcast(event.GetBroadcastPayload())
}

// HandleClientConnection handles a single inbound websocket connection.
func HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	_server.HandleClientConnection(w, r)
}

// DisconnectUser will forcefully disconnect all clients belonging to a user by ID.
func DisconnectUser(userID string) {
	_server.DisconnectUser(userID)
}
