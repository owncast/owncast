package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/geoip"
	"github.com/owncast/owncast/utils"
)

var _server *Server

// Server represents an instance of the chat server.
type Server struct {
	clients map[uint]*Client

	// send outbound message payload to all clients
	outbound chan []byte

	// receive inbound message payload from all clients
	inbound chan chatClientEvent

	// unregister requests from clients.
	unregister chan uint // the ChatClient id

	geoipClient *geoip.Client

	// a map of user IDs and timers that fire for chat part messages.
	userPartedTimers         map[string]*time.Ticker
	seq                      uint
	maxSocketConnectionLimit uint64

	mu sync.RWMutex
}

// NewChat will return a new instance of the chat server.
func NewChat() *Server {
	maximumConcurrentConnectionLimit := getMaximumConcurrentConnectionLimit()
	setSystemConcurrentConnectionLimit(maximumConcurrentConnectionLimit)

	server := &Server{
		clients:                  map[uint]*Client{},
		outbound:                 make(chan []byte),
		inbound:                  make(chan chatClientEvent),
		unregister:               make(chan uint),
		maxSocketConnectionLimit: maximumConcurrentConnectionLimit,
		geoipClient:              geoip.NewClient(),
		userPartedTimers:         map[string]*time.Ticker{},
	}

	return server
}

// Run will start the chat server.
func (s *Server) Run() {
	for {
		select {
		case clientID := <-s.unregister:
			if client, ok := s.clients[clientID]; ok {
				s.handleClientDisconnected(client)
				s.mu.Lock()
				delete(s.clients, clientID)
				s.mu.Unlock()
			}

		case message := <-s.inbound:
			s.eventReceived(message)
		}
	}
}

// Addclient registers new connection as a User.
func (s *Server) Addclient(conn *websocket.Conn, user *models.User, accessToken string, userAgent string, ipAddress string) *Client {
	client := &Client{
		server:      s,
		conn:        conn,
		User:        user,
		IPAddress:   ipAddress,
		accessToken: accessToken,
		send:        make(chan []byte, 256),
		UserAgent:   userAgent,
		ConnectedAt: time.Now(),
	}

	shouldSendJoinedMessages := data.GetChatJoinPartMessagesEnabled()

	// If there are existing clients connected for this user do not send
	// a user joined message. Do not put this under a mutex, as
	// GetClientsForUser already has a lock.
	if existingConnectedClients, _ := GetClientsForUser(user.ID); len(existingConnectedClients) > 0 {
		shouldSendJoinedMessages = false
	}

	s.mu.Lock()
	{
		// If there is a pending disconnect timer then clear it.
		// Do not send user joined message if enough time hasn't passed where the
		// user chat part message hasn't been sent yet.
		if ticker, ok := s.userPartedTimers[user.ID]; ok {
			ticker.Stop()
			delete(s.userPartedTimers, user.ID)
			shouldSendJoinedMessages = false
		}

		client.Id = s.seq
		s.clients[client.Id] = client
		s.seq++
	}
	s.mu.Unlock()

	log.Traceln("Adding client", client.Id, "total count:", len(s.clients))

	go client.writePump()
	go client.readPump()

	client.sendConnectedClientInfo()

	if getStatus().Online {
		if shouldSendJoinedMessages {
			s.sendUserJoinedMessage(client)
		}
		s.sendWelcomeMessageToClient(client)
	}

	// Asynchronously, optionally, fetch GeoIP data.
	go func(client *Client) {
		client.Geo = s.geoipClient.GetGeoFromIP(ipAddress)
	}(client)

	return client
}

func (s *Server) sendUserJoinedMessage(c *Client) {
	userJoinedEvent := events.UserJoinedEvent{}
	userJoinedEvent.SetDefaults()
	userJoinedEvent.User = c.User
	userJoinedEvent.ClientID = c.Id

	if err := s.Broadcast(userJoinedEvent.GetBroadcastPayload()); err != nil {
		log.Errorln("error adding client to chat server", err)
	}

	// Send chat user joined webhook
	webhooks.SendChatEventUserJoined(userJoinedEvent)
}

func (s *Server) handleClientDisconnected(c *Client) {
	if _, ok := s.clients[c.Id]; ok {
		log.Debugln("Deleting", c.Id)
		delete(s.clients, c.Id)
	}

	additionalClientCheck, _ := GetClientsForUser(c.User.ID)
	if len(additionalClientCheck) > 0 {
		// This user is still connected to chat with another client.
		return
	}

	s.userPartedTimers[c.User.ID] = time.NewTicker(10 * time.Second)

	go func() {
		<-s.userPartedTimers[c.User.ID].C
		s.sendUserPartedMessage(c)
	}()
}

func (s *Server) sendUserPartedMessage(c *Client) {
	s.userPartedTimers[c.User.ID].Stop()
	delete(s.userPartedTimers, c.User.ID)

	userPartEvent := events.UserPartEvent{}
	userPartEvent.SetDefaults()
	userPartEvent.User = c.User
	userPartEvent.ClientID = c.Id

	// If part messages are disabled.
	if data.GetChatJoinPartMessagesEnabled() {
		if err := s.Broadcast(userPartEvent.GetBroadcastPayload()); err != nil {
			log.Errorln("error sending chat part message", err)
		}
	}
	// Send chat user joined webhook
	webhooks.SendChatEventUserParted(userPartEvent)
}

// HandleClientConnection is fired when a single client connects to the websocket.
func (s *Server) HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	if data.GetChatDisabled() {
		_, _ = w.Write([]byte(events.ChatDisabled))
		return
	}

	ipAddress := utils.GetIPAddressFromRequest(r)
	// Check if this client's IP address is banned. If so send a rejection.
	if blocked, err := data.IsIPAddressBanned(ipAddress); blocked {
		log.Debugln("Client ip address has been blocked. Rejecting.")

		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		log.Errorln("error determining if IP address is blocked: ", err)
	}

	// Limit concurrent chat connections
	if uint64(len(s.clients)) >= s.maxSocketConnectionLimit {
		log.Warnln("rejecting incoming client connection as it exceeds the max client count of", s.maxSocketConnectionLimit)
		_, _ = w.Write([]byte(events.ErrorMaxConnectionsExceeded))
		return
	}

	// To allow dev web environments to connect.
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debugln(err)
		return
	}

	accessToken := r.URL.Query().Get("accessToken")
	if accessToken == "" {
		log.Errorln("Access token is required")
		// Return HTTP status code
		_ = conn.Close()
		return
	}

	userRepository := userrepository.Get()

	// A user is required to use the websocket
	user := userRepository.GetUserByToken(accessToken)

	if user == nil {
		// Send error that registration is required
		_ = conn.WriteJSON(events.EventPayload{
			"type": events.ErrorNeedsRegistration,
		})
		_ = conn.Close()
		return
	}

	// User is disabled therefore we should disconnect.
	if user.DisabledAt != nil {
		log.Traceln("Disabled user", user.ID, user.DisplayName, "rejected")
		_ = conn.WriteJSON(events.EventPayload{
			"type": events.ErrorUserDisabled,
		})
		_ = conn.Close()
		return
	}

	userAgent := r.UserAgent()

	s.Addclient(conn, user, accessToken, userAgent, ipAddress)
}

// Broadcast sends message to all connected clients.
func (s *Server) Broadcast(payload events.EventPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		if client == nil {
			continue
		}

		select {
		case client.send <- data:
		default:
			go client.close()
		}
	}

	return nil
}

// Send will send a single payload to a single connected client.
func (s *Server) Send(payload events.EventPayload, client *Client) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Errorln(err)
		return
	}

	client.send <- data
}

// DisconnectClients will forcefully disconnect all clients belonging to a user by ID.
func (s *Server) DisconnectClients(clients []*Client) {
	for _, client := range clients {
		log.Traceln("Disconnecting client", client.User.ID, "owned by", client.User.DisplayName)

		go func(client *Client) {
			event := events.UserDisabledEvent{}
			event.SetDefaults()

			// Send this disabled event specifically to this single connected client
			// to let them know they've been banned.
			_server.Send(event.GetBroadcastPayload(), client)

			// Give the socket time to send out the above message.
			// Unfortunately I don't know of any way to get a real callback to know when
			// the message was successfully sent, so give it a couple seconds.
			time.Sleep(2 * time.Second)

			// Forcefully disconnect if still valid.
			if client != nil {
				client.close()
			}
		}(client)
	}
}

// SendConnectedClientInfoToUser will find all the connected clients assigned to a user
// and re-send each the connected client info.
func SendConnectedClientInfoToUser(userID string) error {
	clients, err := GetClientsForUser(userID)
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
func SendActionToUser(userID string, text string) error {
	clients, err := GetClientsForUser(userID)
	if err != nil {
		return err
	}

	for _, client := range clients {
		_server.sendActionToClient(client, text)
	}

	return nil
}

func (s *Server) eventReceived(event chatClientEvent) {
	c := event.client
	u := c.User

	// If established chat user only mode is enabled and the user is not old
	// enough then reject this event and send them an informative message.
	if u != nil && data.GetChatEstbalishedUsersOnlyMode() && time.Since(event.client.User.CreatedAt) < config.GetDefaults().ChatEstablishedUserModeTimeDuration && !u.IsModerator() {
		s.sendActionToClient(c, "You have not been an established chat participant long enough to take part in chat. Please enjoy the stream and try again later.")
		return
	}

	var typecheck map[string]interface{}
	if err := json.Unmarshal(event.data, &typecheck); err != nil {
		log.Debugln(err)
	}

	eventType := typecheck["type"]

	switch eventType {
	case events.MessageSent:
		s.userMessageSent(event)

	case events.UserNameChanged:
		s.userNameChanged(event)

	case events.UserColorChanged:
		s.userColorChanged(event)
	default:
		log.Debugln(logSanitize(fmt.Sprint(eventType)), "event not found:", logSanitize(fmt.Sprint(typecheck)))
	}
}

func (s *Server) sendWelcomeMessageToClient(c *Client) {
	// Add an artificial delay so people notice this message come in.
	time.Sleep(7 * time.Second)

	welcomeMessage := utils.RenderSimpleMarkdown(data.GetServerWelcomeMessage())

	if welcomeMessage != "" {
		s.sendSystemMessageToClient(c, welcomeMessage)
	}
}

func (s *Server) sendAllWelcomeMessage() {
	welcomeMessage := utils.RenderSimpleMarkdown(data.GetServerWelcomeMessage())

	if welcomeMessage != "" {
		clientMessage := events.SystemMessageEvent{
			Event: events.Event{},
			MessageEvent: events.MessageEvent{
				Body: welcomeMessage,
			},
		}
		clientMessage.SetDefaults()
		_ = s.Broadcast(clientMessage.GetBroadcastPayload())
	}
}

func (s *Server) sendSystemMessageToClient(c *Client, message string) {
	clientMessage := events.SystemMessageEvent{
		Event: events.Event{},
		MessageEvent: events.MessageEvent{
			Body: message,
		},
	}
	clientMessage.SetDefaults()
	clientMessage.RenderBody()
	s.Send(clientMessage.GetBroadcastPayload(), c)
}

func (s *Server) sendActionToClient(c *Client, message string) {
	clientMessage := events.ActionEvent{
		MessageEvent: events.MessageEvent{
			Body: message,
		},
		Event: events.Event{
			Type: events.ChatActionSent,
		},
	}
	clientMessage.SetDefaults()
	clientMessage.RenderBody()
	s.Send(clientMessage.GetBroadcastPayload(), c)
}
