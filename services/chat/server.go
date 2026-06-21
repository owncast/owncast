package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
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
	"github.com/owncast/owncast/services/geoip"
	plugins "github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/services/webhooks"
	"github.com/owncast/owncast/utils"
)

// Service is an instance of the chat server. Construct via New(Deps).
type Service struct {
	clients map[uint]*Client

	// send outbound message payload to all clients
	outbound chan []byte

	// receive inbound message payload from all clients
	inbound chan chatClientEvent

	// unregister requests from clients.
	unregister chan uint // the ChatClient id

	geoipClient *geoip.Client

	// webhooks dispatches chat events (messages, joins, parts,
	// renames, visibility toggles) to configured webhook destinations.
	webhooks *webhooks.Service

	// datastore is the database handle the chat pruner uses to
	// trim aged messages on a recurring timer.
	datastore *datastore.Datastore

	// getStatus returns the current stream status; used by the
	// inbound-message path to make stream-state-aware decisions
	// (welcome messages, recent-disconnect heuristics).
	getStatus func() models.Status

	// events is the shared dispatcher. Inbound user messages are run through
	// its filter chain before broadcast, letting consumers (the plugin host)
	// rewrite or drop them. nil means no filtering.
	events *dispatcher.Dispatcher

	// configRepository provides server-side chat settings consulted on each
	// inbound message (slur/spam filters, established-users-only mode,
	// auth requirement, welcome message text, etc.).
	configRepository configrepository.ConfigRepository

	// authRepository is consulted on each inbound websocket connection to
	// reject banned IP addresses before upgrading to a chat client.
	authRepository authrepository.AuthRepository

	// chatMessageRepository persists user/system/action/federated chat
	// events and answers history + moderation visibility queries.
	chatMessageRepository chatmessagerepository.ChatMessageRepository

	// userRepository provides access to user records (lookups by ID/token,
	// display-name availability checks, color/name updates) consulted
	// throughout the chat message lifecycle.
	userRepository userrepository.UserRepository

	// chatMessagesSentCounter is the Prometheus counter incremented on
	// each accepted inbound message.
	chatMessagesSentCounter prometheus.Gauge

	// a map of user IDs and timers that fire for chat part messages.
	userPartedTimers         map[string]*time.Ticker
	seq                      uint
	maxSocketConnectionLimit uint64

	mu sync.RWMutex
}

// NewChat will return a new instance of the chat server.
func newServer(webhooksSvc *webhooks.Service) *Service {
	maximumConcurrentConnectionLimit := getMaximumConcurrentConnectionLimit()
	setSystemConcurrentConnectionLimit(maximumConcurrentConnectionLimit)

	server := &Service{
		clients:                  map[uint]*Client{},
		outbound:                 make(chan []byte),
		inbound:                  make(chan chatClientEvent),
		unregister:               make(chan uint),
		maxSocketConnectionLimit: maximumConcurrentConnectionLimit,
		geoipClient:              geoip.NewClient(),
		webhooks:                 webhooksSvc,
		userPartedTimers:         map[string]*time.Ticker{},
	}

	return server
}

// Run will start the chat server.
func (s *Service) Run() {
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
func (s *Service) Addclient(conn *websocket.Conn, user *models.User, accessToken string, userAgent string, ipAddress string) *Client {
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

	// isNewJoin tracks whether this connection is a genuinely new chat join, as
	// opposed to an additional client for an already-connected user or a quick
	// reconnect. The admin "show join/part messages" setting is intentionally
	// NOT consulted here: it governs only the visible broadcast in
	// sendUserJoinedMessage, not whether the join occurred (so webhooks still
	// fire when join/part messages are hidden). See #4950.
	isNewJoin := true

	// If there are existing clients connected for this user this is not a new
	// join. Do not put this under a mutex, as GetClientsForUser already has a
	// lock.
	if existingConnectedClients, _ := s.GetClientsForUser(user.ID); len(existingConnectedClients) > 0 {
		isNewJoin = false
	}

	s.mu.Lock()
	{
		// If there is a pending disconnect timer then clear it. A quick
		// reconnect before the part was sent is not a new join, so neither the
		// join broadcast nor the webhook should fire for it.
		if ticker, ok := s.userPartedTimers[user.ID]; ok {
			ticker.Stop()
			delete(s.userPartedTimers, user.ID)
			isNewJoin = false
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

	if s.getStatus().Online {
		if isNewJoin {
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

func (s *Service) sendUserJoinedMessage(c *Client) {
	userJoinedEvent := events.UserJoinedEvent{}
	userJoinedEvent.SetDefaults()
	userJoinedEvent.User = c.User
	userJoinedEvent.ClientID = c.Id

	// Only broadcast the visible join message when join/part messages are
	// enabled; the webhook below fires regardless, mirroring sendUserPartedMessage (#4950).
	if s.configRepository.GetChatJoinPartMessagesEnabled() {
		if err := s.Broadcast(userJoinedEvent.GetBroadcastPayload()); err != nil {
			log.Errorln("error adding client to chat server", err)
		}
	}

	// Send chat user joined webhook
	s.webhooks.SendChatEventUserJoined(userJoinedEvent)
}

func (s *Service) handleClientDisconnected(c *Client) {
	if _, ok := s.clients[c.Id]; ok {
		log.Debugln("Deleting", c.Id)
		delete(s.clients, c.Id)
	}

	additionalClientCheck, _ := s.GetClientsForUser(c.User.ID)
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

func (s *Service) sendUserPartedMessage(c *Client) {
	s.userPartedTimers[c.User.ID].Stop()
	delete(s.userPartedTimers, c.User.ID)

	userPartEvent := events.UserPartEvent{}
	userPartEvent.SetDefaults()
	userPartEvent.User = c.User
	userPartEvent.ClientID = c.Id

	// If part messages are disabled.
	if s.configRepository.GetChatJoinPartMessagesEnabled() {
		if err := s.Broadcast(userPartEvent.GetBroadcastPayload()); err != nil {
			log.Errorln("error sending chat part message", err)
		}
	}
	// Send chat user joined webhook
	s.webhooks.SendChatEventUserParted(userPartEvent)
}

// HandleClientConnection is fired when a single client connects to the websocket.
func (s *Service) HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	if s.configRepository.GetChatDisabled() {
		_, _ = w.Write([]byte(events.ChatDisabled))
		return
	}

	ipAddress := utils.GetIPAddressFromRequest(r)
	// Check if this client's IP address is banned. If so send a rejection.
	if blocked, err := s.authRepository.IsIPAddressBanned(ipAddress); blocked {
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

	// Read the access token before upgrading so we can set the chat identity
	// cookie on the handshake response — once the connection is upgraded the
	// underlying ResponseWriter is hijacked and can no longer set headers.
	// This covers returning users whose token predates the cookie, since
	// every chat session opens a websocket with the token in the query.
	accessToken := r.URL.Query().Get("accessToken")
	var responseHeader http.Header
	if accessToken != "" {
		responseHeader = http.Header{}
		utils.AddChatAccessTokenCookieHeader(responseHeader, r, accessToken)
	} else if t := plugins.SessionTokenFromContext(r.Context()); t != "" {
		// Viewer authenticated through the plugin auth gate: the signed gate
		// session cookie (already verified by the gate middleware) carries the
		// access token, surfaced here via the request context. No query param
		// and no extra cookie needed — the gate cookie is the credential.
		accessToken = t
	}

	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		log.Debugln(err)
		return
	}

	if accessToken == "" {
		log.Errorln("Access token is required")
		// Return HTTP status code
		_ = conn.Close()
		return
	}

	// A user is required to use the websocket
	user := s.userRepository.GetUserByToken(accessToken)

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
func (s *Service) Broadcast(payload events.EventPayload) error {
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
func (s *Service) Send(payload events.EventPayload, client *Client) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Errorln(err)
		return
	}

	client.send <- data
}

// DisconnectClients will forcefully disconnect all clients belonging to a user by ID.
func (s *Service) DisconnectClients(clients []*Client) {
	for _, client := range clients {
		log.Traceln("Disconnecting client", client.User.ID, "owned by", client.User.DisplayName)

		go func(client *Client) {
			event := events.UserDisabledEvent{}
			event.SetDefaults()

			// Send this disabled event specifically to this single connected client
			// to let them know they've been banned.
			s.Send(event.GetBroadcastPayload(), client)

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
func (s *Service) SendConnectedClientInfoToUser(userID string) error {
	clients, err := s.GetClientsForUser(userID)
	if err != nil {
		return err
	}

	// Get an updated reference to the user.
	user := s.userRepository.GetUserByID(userID)
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
func (s *Service) SendActionToUser(userID string, text string) error {
	clients, err := s.GetClientsForUser(userID)
	if err != nil {
		return err
	}

	for _, client := range clients {
		s.sendActionToClient(client, text)
	}

	return nil
}

func (s *Service) eventReceived(event chatClientEvent) {
	c := event.client
	u := c.User

	// If established chat user only mode is enabled and the user is not old
	// enough then reject this event and send them an informative message.
	if u != nil && s.configRepository.GetChatEstbalishedUsersOnlyMode() && time.Since(event.client.User.CreatedAt) < config.GetDefaults().ChatEstablishedUserModeTimeDuration && !u.IsModerator() {
		s.sendActionToClient(c, "You have not been an established chat participant long enough to take part in chat. Please enjoy the stream and try again later.")
		return
	}

	// Check if authentication is required for chat
	if s.configRepository.GetChatRequireAuthentication() {
		if u == nil || (!u.Authenticated && !u.IsModerator()) {
			s.sendActionToClient(c, "Authentication is required to participate in chat.")
			return
		}
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

func (s *Service) sendWelcomeMessageToClient(c *Client) {
	// Add an artificial delay so people notice this message come in.
	time.Sleep(7 * time.Second)

	welcomeMessage := utils.RenderSimpleMarkdown(s.configRepository.GetServerWelcomeMessage())

	if welcomeMessage != "" {
		s.sendSystemMessageToClient(c, welcomeMessage)
	}
}

func (s *Service) sendAllWelcomeMessage() {
	welcomeMessage := utils.RenderSimpleMarkdown(s.configRepository.GetServerWelcomeMessage())

	if welcomeMessage != "" {
		clientMessage := events.SystemMessageEvent{
			Event: events.Event{},
			MessageEvent: events.MessageEvent{
				Body: welcomeMessage,
			},
			ServerName: s.configRepository.GetServerName(),
		}
		clientMessage.SetDefaults()
		_ = s.Broadcast(clientMessage.GetBroadcastPayload())
	}
}

func (s *Service) sendSystemMessageToClient(c *Client, message string) {
	clientMessage := events.SystemMessageEvent{
		Event: events.Event{},
		MessageEvent: events.MessageEvent{
			Body: message,
		},
		ServerName: s.configRepository.GetServerName(),
	}
	clientMessage.SetDefaults()
	clientMessage.RenderBody()
	s.Send(clientMessage.GetBroadcastPayload(), c)
}

func (s *Service) sendActionToClient(c *Client, message string) {
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
