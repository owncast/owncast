// Package chat is the chat subsystem service: a websocket server that
// accepts client connections, broadcasts inbound messages, fans out
// system + action messages, and forwards chat events to webhooks.
// Construct via New(Deps) and call Start to launch the broadcast loop.
package chat

import (
	"errors"
	"net/http"
	"sort"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/authrepository"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/dispatcher"
	"github.com/owncast/owncast/services/webhooks"
)

// Deps lists the explicit dependencies for the chat Service.
type Deps struct {
	GetStatus             func() models.Status
	Webhooks              *webhooks.Service
	Datastore             *datastore.Datastore
	ConfigRepository      configrepository.ConfigRepository
	AuthRepository        authrepository.AuthRepository
	ChatMessageRepository chatmessagerepository.ChatMessageRepository
	UserRepository        userrepository.UserRepository
	Events                *dispatcher.Dispatcher
}

// New constructs a chat Service. Call Start to launch the broadcast
// loop and register the Prometheus counter.
func New(deps Deps) *Service {
	s := newServer(deps.Webhooks)
	s.getStatus = deps.GetStatus
	s.datastore = deps.Datastore
	s.configRepository = deps.ConfigRepository
	s.authRepository = deps.AuthRepository
	s.chatMessageRepository = deps.ChatMessageRepository
	s.userRepository = deps.UserRepository
	s.events = deps.Events
	return s
}

// SetGetStatus wires (or rewires) the stream-status callback. Exists
// because the composition root has a small cycle: chat needs
// stream.GetStatus; the stream service needs *chat.Service. main.go
// constructs chat first with a nil callback, then fills it in once
// streamSvc exists. Must be called before Start.
func (s *Service) SetGetStatus(fn func() models.Status) {
	s.getStatus = fn
}

// Start initializes persistence, launches the broadcast loop, and
// registers the Prometheus counter. Safe to call once.
func (s *Service) Start() error {
	s.setupPersistence()

	go s.Run()

	log.Traceln("Chat server started with max connection count of", s.maxSocketConnectionLimit)

	s.chatMessagesSentCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "total_chat_message_count",
		Help: "The number of chat messages incremented over time.",
		ConstLabels: map[string]string{
			"version": config.VersionNumber,
			"host":    s.configRepository.GetServerURL(),
		},
	})

	return nil
}

// GetClientsForUser will return chat connections that are owned by a specific user.
func (s *Service) GetClientsForUser(userID string) ([]*Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clients := map[string][]*Client{}

	for _, client := range s.clients {
		clients[client.User.ID] = append(clients[client.User.ID], client)
	}

	if _, exists := clients[userID]; !exists {
		return nil, errors.New("no connections for user found")
	}

	return clients[userID], nil
}

// FindClientByID will return a single connected client by ID.
func (s *Service) FindClientByID(clientID uint) (*Client, bool) {
	client, found := s.clients[clientID]
	return client, found
}

// GetClients will return all the current chat clients connected.
func (s *Service) GetClients() []*Client {
	clients := []*Client{}

	if s == nil {
		return clients
	}

	for _, client := range s.clients {
		clients = append(clients, client)
	}

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].ConnectedAt.Before(clients[j].ConnectedAt)
	})

	return clients
}

// SendSystemMessage sends a message string as a system message to all clients.
func (s *Service) SendSystemMessage(text string, ephemeral bool) error {
	message := events.SystemMessageEvent{
		MessageEvent: events.MessageEvent{
			Body: text,
		},
		ServerName: s.configRepository.GetServerName(),
	}
	message.SetDefaults()
	message.RenderBody()

	if err := s.BroadcastEvent(&message); err != nil {
		log.Errorln("error sending system message", err)
	}

	if !ephemeral {
		s.chatMessageRepository.SaveEvent(message.ID, nil, message.Body, message.GetMessageType(), nil, message.Timestamp, nil, nil, nil, nil)
	}

	return nil
}

// SendMessageAsBot broadcasts and persists a chat message authored by a bot
// user, running it through the same render/sanitize, broadcast, webhook, and
// persistence path as a message from a connected client.
//
// The user MUST have IsBot=true. This is deliberately the only chat-send API
// that accepts a user identity — there's no general "send as any user" path,
// so a plugin (the only thing that can reach this) cannot impersonate a real
// chatter. The plugin host's chatbot provisioner creates exactly one bot
// identity per plugin (with IsBot=true and DisplayName=plugin name) and
// passes it here when the plugin calls owncast.chat.send.
func (s *Service) SendMessageAsBot(user *models.User, text string) error {
	if user == nil {
		return errors.New("SendMessageAsBot requires a user")
	}
	if !user.IsBot {
		return errors.New("SendMessageAsBot: user must be a bot (IsBot=true) — impersonating a non-bot is not allowed")
	}

	event := events.UserMessageEvent{}
	event.Body = text
	event.SetDefaults() // generates ID + timestamp, renders and sanitizes the body
	// SetDefaults doesn't set Type (the websocket-inbound path picks it up
	// from the client JSON). We're synthesizing the event server-side, so set
	// it explicitly — without this the message is stored with an empty
	// eventType, and the web frontend can't categorize it and crashes when
	// rendering chat history.
	event.Type = events.MessageSent
	event.User = user

	if event.Empty() {
		return nil
	}

	if err := s.Broadcast(event.GetBroadcastPayload()); err != nil {
		return err
	}

	s.webhooks.SendChatEvent(&event)
	s.chatMessagesSentCounter.Inc()
	s.chatMessageRepository.SaveUserMessage(event)

	return nil
}

// SendFediverseAction sends a message indicating some Fediverse engagement took place.
func (s *Service) SendFediverseAction(eventType string, userAccountName string, image *string, body string, link string) error {
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
		ServerName:      s.configRepository.GetServerName(),
	}

	message.SetDefaults()

	if err := s.BroadcastEvent(&message); err != nil {
		log.Errorln("error sending system message", err)
		return err
	}

	s.chatMessageRepository.SaveFederatedAction(message)

	return nil
}

// SendSystemAction sends a system action string as an action event to all clients.
func (s *Service) SendSystemAction(text string, ephemeral bool) error {
	message := events.ActionEvent{
		MessageEvent: events.MessageEvent{
			Body: text,
		},
	}

	message.SetDefaults()
	message.RenderBody()

	if err := s.BroadcastEvent(&message); err != nil {
		log.Errorln("error sending system chat action")
	}

	if !ephemeral {
		s.chatMessageRepository.SaveEvent(message.ID, nil, message.Body, message.GetMessageType(), nil, message.Timestamp, nil, nil, nil, nil)
	}

	return nil
}

// SendAllWelcomeMessage sends the chat welcome message to all connected clients.
func (s *Service) SendAllWelcomeMessage() {
	s.sendAllWelcomeMessage()
}

// SendSystemMessageToClient sends a single message to a single connected chat client.
func (s *Service) SendSystemMessageToClient(clientID uint, text string) {
	if client, foundClient := s.FindClientByID(clientID); foundClient {
		s.sendSystemMessageToClient(client, text)
	}
}

// BroadcastEvent sends all connected clients the outbound object provided.
func (s *Service) BroadcastEvent(event events.OutboundEvent) error {
	return s.Broadcast(event.GetBroadcastPayload())
}

// HandleClientConnection handles a single inbound websocket connection.
// Wrapper around the lower-level method so router can bind cleanly.
func (s *Service) HandleWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	s.HandleClientConnection(w, r)
}
