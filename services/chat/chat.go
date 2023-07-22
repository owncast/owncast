package chat

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage/chatrepository"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/userrepository"
	log "github.com/sirupsen/logrus"
)

type Chat struct {
	getStatus        func() *models.Status
	server           *Server
	configRepository *configrepository.SqlConfigRepository
	emojis           *emojis
}

func New() *Chat {
	return &Chat{
		configRepository: configrepository.Get(),
		emojis:           newEmojis(),
	}
}

var temporaryGlobalInstance *Chat

// GetConfig returns the temporary global instance.
// Remove this after dependency injection is implemented.
func Get() *Chat {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}

	return temporaryGlobalInstance
}

// Start begins the chat server.
func (c *Chat) Start(getStatusFunc func() *models.Status) error {
	c.getStatus = getStatusFunc
	c.server = NewChat()

	go c.server.Run()

	log.Traceln("Chat server started with max connection count of", c.server.maxSocketConnectionLimit)

	return nil
}

// FindClientByID will return a single connected client by ID.
func (c *Chat) FindClientByID(clientID uint) (*Client, bool) {
	client, found := c.server.clients[clientID]
	return client, found
}

// GetClients will return all the current chat clients connected.
func (c *Chat) GetClients() []*Client {
	clients := []*Client{}

	if c.server == nil {
		return clients
	}

	// Convert the keyed map to a slice.
	for _, client := range c.server.clients {
		clients = append(clients, client)
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].ConnectedAt.Before(clients[j].ConnectedAt)
	})

	return clients
}

// SendSystemMessage will send a message string as a system message to all clients.
func (c *Chat) SendSystemMessage(text string, ephemeral bool) error {
	message := SystemMessageEvent{
		MessageEvent: MessageEvent{
			Body: text,
		},
	}
	message.SetDefaults()
	message.RenderBody()
	message.DisplayName = c.configRepository.GetServerName()

	if err := c.Broadcast(&message); err != nil {
		log.Errorln("error sending system message", err)
	}

	if !ephemeral {
		cr := chatrepository.Get()
		cr.SaveEvent(message.ID, nil, message.Body, message.GetMessageType(), nil, message.Timestamp, nil, nil, nil, nil)
	}

	return nil
}

// SendFediverseAction will send a message indicating some Fediverse engagement took place.
func (c *Chat) SendFediverseAction(eventType string, userAccountName string, image *string, body string, link string) error {
	message := FediverseEngagementEvent{
		Event: models.Event{
			Type: eventType,
		},
		MessageEvent: MessageEvent{
			Body: body,
		},
		UserAccountName: userAccountName,
		Image:           image,
		Link:            link,
	}

	message.SetDefaults()
	message.RenderBody()

	if err := c.Broadcast(&message); err != nil {
		log.Errorln("error sending system message", err)
		return err
	}

	cr := chatrepository.Get()
	cr.SaveFederatedAction(message)

	return nil
}

// SendSystemAction will send a system action string as an action event to all clients.
func (c *Chat) SendSystemAction(text string, ephemeral bool) error {
	message := ActionEvent{
		MessageEvent: MessageEvent{
			Body: text,
		},
	}

	message.SetDefaults()
	message.RenderBody()

	if err := c.Broadcast(&message); err != nil {
		log.Errorln("error sending system chat action")
	}

	if !ephemeral {
		cr := chatrepository.Get()
		cr.SaveEvent(message.ID, nil, message.Body, message.GetMessageType(), nil, message.Timestamp, nil, nil, nil, nil)
	}

	return nil
}

// SendAllWelcomeMessage will send the chat message to all connected clients.
func (c *Chat) SendAllWelcomeMessage() {
	c.server.sendAllWelcomeMessage()
}

// SendSystemMessageToClient will send a single message to a single connected chat client.
func (c *Chat) SendSystemMessageToClient(clientID uint, text string) {
	if client, foundClient := c.FindClientByID(clientID); foundClient {
		c.server.sendSystemMessageToClient(client, text)
	}
}

// Broadcast will send all connected clients the outbound object provided.
func (c *Chat) Broadcast(event OutboundEvent) error {
	return c.server.Broadcast(event.GetBroadcastPayload())
}

// HandleClientConnection handles a single inbound websocket connection.
func (c *Chat) HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	c.server.HandleClientConnection(w, r)
}

// DisconnectClients will forcefully disconnect all clients belonging to a user by ID.
func (c *Chat) DisconnectClients(clients []*Client) {
	c.server.DisconnectClients(clients)
}

func (c *Chat) GetClientsForUser(userID string) ([]*Client, error) {
	return c.server.GetClientsForUser(userID)
}

// SendConnectedClientInfoToUser will find all the connected clients assigned to a user
// and re-send each the connected client info.
func (c *Chat) SendConnectedClientInfoToUser(userID string) error {
	clients, err := c.server.GetClientsForUser(userID)
	if err != nil {
		return err
	}

	userRepository := userrepository.Get()

	// Get an updated reference to the user.
	user := userRepository.GetUserByID(userID)
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if err != nil {
		return err
	}

	for _, client := range clients {
		// Update the client's reference to its user.
		client.User = user
		// Send the update to the client.
		client.sendConnectedClientInfo()
	}

	return nil
}

// SendActionToUser will send system action text to all connected clients
// assigned to a user ID.
func (c *Chat) SendActionToUser(userID string, text string) error {
	clients, err := c.server.GetClientsForUser(userID)
	if err != nil {
		return err
	}

	for _, client := range clients {
		c.server.sendActionToClient(client, text)
	}

	return nil
}
