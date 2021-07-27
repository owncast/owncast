package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"

	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/utils"
)

var _server *Server

// a map of user IDs and when they last were active.
var _lastSeenCache = map[string]time.Time{}

// Server represents an instance of the chat server.
type Server struct {
	mu                       sync.RWMutex
	seq                      uint
	clients                  map[uint]*Client
	maxSocketConnectionLimit int64

	// send outbound message payload to all clients
	outbound chan []byte

	// receive inbound message payload from all clients
	inbound chan chatClientEvent

	// unregister requests from clients.
	unregister chan uint // the ChatClient id

	geoipClient *geoip.Client
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
	}

	return server
}

// Run will start the chat server.
func (s *Server) Run() {
	for {
		select {
		case clientID := <-s.unregister:
			if _, ok := s.clients[clientID]; ok {
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
func (s *Server) Addclient(conn *websocket.Conn, user *user.User, accessToken string, userAgent string, ipAddress string) *Client {
	client := &Client{
		server:      s,
		conn:        conn,
		User:        user,
		ipAddress:   ipAddress,
		accessToken: accessToken,
		send:        make(chan []byte, 256),
		UserAgent:   userAgent,
		ConnectedAt: time.Now(),
	}

	// Do not send user re-joined broadcast message if they've been active within 5 minutes.
	shouldSendJoinedMessages := true
	if previouslyLastSeen, ok := _lastSeenCache[user.ID]; ok && time.Since(previouslyLastSeen) < time.Minute*5 {
		shouldSendJoinedMessages = false
	}

	s.mu.Lock()
	{
		client.id = s.seq
		s.clients[client.id] = client
		s.seq++
		_lastSeenCache[user.ID] = time.Now()
	}
	s.mu.Unlock()

	log.Traceln("Adding client", client.id, "total count:", len(s.clients))

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
	userJoinedEvent.ClientID = c.id

	if err := s.Broadcast(userJoinedEvent.GetBroadcastPayload()); err != nil {
		log.Errorln("error adding client to chat server", err)
	}

	// Send chat user joined webhook
	webhooks.SendChatEventUserJoined(userJoinedEvent)
}

// ClientClosed is fired when a client disconnects or connection is dropped.
func (s *Server) ClientClosed(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.close()

	if _, ok := s.clients[c.id]; ok {
		log.Debugln("Deleting", c.id)
		delete(s.clients, c.id)
	}
}

// HandleClientConnection is fired when a single client connects to the websocket.
func (s *Server) HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	if data.GetChatDisabled() {
		_, _ = w.Write([]byte(events.ChatDisabled))
		return
	}

	// Limit concurrent chat connections
	if int64(len(s.clients)) >= s.maxSocketConnectionLimit {
		log.Warnln("rejecting incoming client connection as it exceeds the max client count of", s.maxSocketConnectionLimit)
		_, _ = w.Write([]byte(events.ErrorMaxConnectionsExceeded))
		return
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

	// A user is required to use the websocket
	user := user.GetUserByToken(accessToken)
	if user == nil {
		_ = conn.WriteJSON(events.EventPayload{
			"type": events.ErrorNeedsRegistration,
		})
		// Send error that registration is required
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
	ipAddress := utils.GetIPAddressFromRequest(r)

	s.Addclient(conn, user, accessToken, userAgent, ipAddress)
}

// Broadcast sends message to all connected clients.
func (s *Server) Broadcast(payload events.EventPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		defer func() {
			if a := recover(); a != nil {
				fmt.Println("RECOVER", a)
			}
		}()

		for _, client := range s.clients {
			if client == nil {
				continue
			}

			if client.send != nil {
				select {
				case client.send <- data:
				default:
					client.close()
					delete(s.clients, client.id)
				}
			}
		}
	}()

	return nil
}

// Send will send a single payload to a single connected client.
func (s *Server) Send(payload events.EventPayload, client *Client) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Errorln(err)
		return
	}

	if client.send != nil {
		client.send <- data
	}
}

// DisconnectUser will forcefully disconnect all clients belonging to a user by ID.
func (s *Server) DisconnectUser(userID string) {
	s.mu.Lock()
	clients, err := GetClientsForUser(userID)
	s.mu.Unlock()

	if err != nil || clients == nil || len(clients) == 0 {
		log.Debugln("Requested to disconnect user", userID, err)
		return
	}

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

func (s *Server) eventReceived(event chatClientEvent) {
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

	default:
		log.Debugln(eventType, "event not found:", typecheck)
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
	}
	clientMessage.SetDefaults()
	clientMessage.RenderBody()
	s.Send(clientMessage.GetBroadcastPayload(), c)
}
