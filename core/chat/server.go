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
	"github.com/owncast/owncast/utils"
)

var _server *ChatServer

type ChatServer struct {
	mu             sync.RWMutex
	seq            uint
	clients        map[uint]*ChatClient
	maxClientCount uint

	// send outbound message payload to all clients
	outbound chan []byte

	// receive inbound message payload from all clients
	inbound chan chatClientEvent

	// Unregister requests from clients.
	unregister chan *ChatClient
}

func NewChat() *ChatServer {
	server := &ChatServer{
		clients:        map[uint]*ChatClient{},
		outbound:       make(chan []byte),
		inbound:        make(chan chatClientEvent),
		unregister:     make(chan *ChatClient),
		maxClientCount: getMaxConnectionCount(),
	}

	return server
}

func (s *ChatServer) Run() {
	for {
		select {
		case client := <-s.unregister:
			if _, ok := s.clients[client.id]; ok {
				s.mu.Lock()
				delete(s.clients, client.id)
				close(client.send)
				s.mu.Unlock()
			}

		case message := <-s.inbound:
			s.eventReceived(message)
		}
	}
}

// Addclient registers new connection as a User.
func (s *ChatServer) Addclient(conn *websocket.Conn, user *user.User, accessToken string, userAgent string) *ChatClient {
	client := &ChatClient{
		server:      s,
		conn:        conn,
		User:        user,
		IPAddress:   conn.RemoteAddr().String(),
		accessToken: accessToken,
		send:        make(chan []byte, 256),
		UserAgent:   userAgent,
		ConnectedAt: time.Now(),
	}

	s.mu.Lock()
	{
		client.id = s.seq
		s.clients[client.id] = client
		s.seq++
	}
	s.mu.Unlock()

	fmt.Println("Adding client", client.id)

	go client.writePump()
	go client.readPump()

	userJoinedEvent := events.UserJoinedEvent{}
	userJoinedEvent.SetDefaults()
	userJoinedEvent.User = user

	s.Broadcast(userJoinedEvent.GetBroadcastPayload())
	s.sendWelcomeMessageToClient(client)

	return client
}

func (s *ChatServer) ClientClosed(c *ChatClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.close()

	if _, ok := s.clients[c.id]; ok {
		fmt.Println("Deleting", c.id)
		delete(s.clients, c.id)
	}
}

func (s *ChatServer) HandleClientConnection(w http.ResponseWriter, r *http.Request) {
	if data.GetChatDisabled() {
		w.Write([]byte(events.Event_Chat_Disabled))
		return
	}

	// Limit concurrent chat connections
	if uint(len(s.clients)) >= s.maxClientCount {
		log.Warnln("rejecting incoming client connection as it exceeds the max client count of", s.maxClientCount)
		w.Write([]byte(events.Event_Error_Max_Connections_Exceeded))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	accessToken := r.URL.Query().Get("accessToken")
	if accessToken == "" {
		log.Errorln("Access token is required")
		// Return HTTP status code
		conn.Close()
		return
	}

	// A user is required to use the websocket
	user := user.GetUserByToken(accessToken)
	if user == nil {
		log.Errorln(accessToken, "has no user")
		conn.WriteJSON(events.EventPayload{
			"type": events.Event_Error_Needs_Registration,
		})
		// Send error that registration is required
		conn.Close()
		return
	}

	// User is disabled therefore we should disconnect.
	if user.DisabledAt != nil {
		log.Warnln("Disabled user", user.Id, user.DisplayName, "rejected")
		conn.WriteJSON(events.EventPayload{
			"type": events.Event_Error_User_Disabled,
		})
		conn.Close()
		return
	}

	userAgent := r.UserAgent()

	s.Addclient(conn, user, accessToken, userAgent)
}

// Broadcast sends message to all connected clients.
func (s *ChatServer) Broadcast(payload events.EventPayload) error {
	// fmt.Println("sending", payload)
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, client := range s.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(s.clients, client.id)
		}
	}

	return nil
}

func (s *ChatServer) Send(payload events.EventPayload, client *ChatClient) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Errorln(err)
		return
	}

	client.send <- data
}

// DisconnectUser will forcefully disconnect all clients belonging to a user by ID.
func (s *ChatServer) DisconnectUser(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, client := range s.clients {
		if client.User.Id == userID {
			go func() {
				event := events.UserDisabledEvent{}
				event.SetDefaults()
				_server.Send(event.GetBroadcastPayload(), client)
				time.Sleep(1 * time.Second) // Allow the socket to send out the above message
				client.close()
			}()
		}
	}
}

func (s *ChatServer) eventReceived(event chatClientEvent) {
	var typecheck map[string]interface{}
	if err := json.Unmarshal(event.data, &typecheck); err != nil {
		panic(err)
	}

	eventType := typecheck["type"]

	// fmt.Println("handleEvent", eventType)
	switch eventType {
	case events.Event_MessageSent:
		s.userMessageSent(event)

	case events.Event_UserNameChanged:
		s.userNameChanged(event)

	default:
		fmt.Println(eventType, "event not found:", typecheck)
	}
}

func (s *ChatServer) sendWelcomeMessageToClient(c *ChatClient) {
	// Add an artificial delay so people notice this message come in.
	time.Sleep(7 * time.Second)

	welcomeMessage := utils.RenderSimpleMarkdown(data.GetServerWelcomeMessage())

	if welcomeMessage != "" {
		//s.sendSystemMessageToClient(c, welcomeMessage)
		s.sendActionToClient(c, welcomeMessage)
	}
}

func (s *ChatServer) sendSystemMessageToClient(c *ChatClient, message string) {
	clientMessage := events.SystemMessageEvent{
		Event: events.Event{},
		MessageEvent: events.MessageEvent{
			Body: message,
		},
	}
	clientMessage.SetDefaults()
	s.Send(clientMessage.GetBroadcastPayload(), c)
}

func (s *ChatServer) sendActionToClient(c *ChatClient, message string) {
	clientMessage := events.ActionEvent{
		MessageEvent: events.MessageEvent{
			Body: message,
		},
	}
	clientMessage.SetDefaults()
	clientMessage.RenderBody()
	s.Send(clientMessage.GetBroadcastPayload(), c)
}
