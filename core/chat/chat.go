package chat

import (
	"errors"
	"net/http"
	"sort"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	getStatus               func() models.Status
	chatMessagesSentCounter prometheus.Gauge
)

// Start begins the chat server.
func Start(getStatusFunc func() models.Status) error {
	setupPersistence()

	getStatus = getStatusFunc
	_server = NewChat()

	go _server.Run()

	log.Traceln("Chat server started with max connection count of", _server.maxSocketConnectionLimit)

	chatMessagesSentCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_chat_message_count",
		Help: "The number of chat messages incremented over time.",
		ConstLabels: map[string]string{
			"version": config.VersionNumber,
			"host":    data.GetServerURL(),
		},
	})

	return nil
}

// GetClientsForUser will return chat connections that are owned by a specific user.
func GetClientsForUser(userID string) ([]*Client, error) {
	clients := map[string][]*Client{}

	for _, client := range _server.clients {
		clients[client.User.ID] = append(clients[client.User.ID], client)
	}

	if _, exists := clients[userID]; !exists {
		return nil, errors.New("no connections for user found")
	}

	return clients[userID], nil
}

// FindClientByID will return a single connected client by ID.
func FindClientByID(clientID uint) (*Client, bool) {
	client, found := _server.clients[clientID]
	return client, found
}

// GetClients will return all the current chat clients connected.
func GetClients() []*Client {
	clients := []*Client{}

	if _server == nil {
		return clients
	}

	// Convert the keyed map to a slice.
	for _, client := range _server.clients {
		clients = append(clients, client)
	}

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
		saveEvent(message.ID, nil, message.Body, message.GetMessageType(), nil, message.Timestamp, nil, nil, nil, nil)
	}

	return nil
}

// SendFediverseAction will send a message indicating some Fediverse engagement took place.
func SendFediverseAction(eventType string, userAccountName string, image *string, body string, link string) error {
	message := events.FediverseEngagementEvent{
		Event: events.Event{
			Type: eventType,
		},
		MessageEvent: events.MessageEvent{
			Body: body,
		},
		UserAccountName: userAccountName,
		Image:           image,
		Link:            link,
	}

	message.SetDefaults()
	message.RenderBody()

	if err := Broadcast(&message); err != nil {
		log.Errorln("error sending system message", err)
		return err
	}

	saveFederatedAction(message)

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
		saveEvent(message.ID, nil, message.Body, message.GetMessageType(), nil, message.Timestamp, nil, nil, nil, nil)
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
