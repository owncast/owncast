package chat

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

var (
	_server *server
)

//Server represents the server which handles the chat
type server struct {
	Clients map[string]*Client

	pattern  string
	listener models.ChatListener

	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan models.ChatMessage
	pingCh    chan models.PingMessage
	doneCh    chan bool
	errCh     chan error
}

//Add adds a client to the server
func (s *server) add(c *Client) {
	s.addCh <- c
}

//Remove removes a client from the server
func (s *server) remove(c *Client) {
	s.delCh <- c
}

//SendToAll sends a message to all of the connected clients
func (s *server) SendToAll(msg models.ChatMessage) {
	s.sendAllCh <- msg
}

//Done marks the server as done
func (s *server) done() {
	s.doneCh <- true
}

//Err handles an error
func (s *server) err(err error) {
	s.errCh <- err
}

func (s *server) sendAll(msg models.ChatMessage) {
	for _, c := range s.Clients {
		c.Write(msg)
	}
}

func (s *server) ping() {
	ping := models.PingMessage{MessageType: PING}
	for _, c := range s.Clients {
		c.pingch <- ping
	}
}

func (s *server) usernameChanged(msg models.NameChangeEvent) {
	for _, c := range s.Clients {
		c.usernameChangeChannel <- msg
	}
}

func (s *server) onConnection(ws *websocket.Conn) {
	client := NewClient(ws)

	defer func() {
		log.Tracef("The client was connected for %s and sent %d messages (%s)", time.Since(client.ConnectedAt), client.MessageCount, client.ClientID)

		if err := ws.Close(); err != nil {
			s.errCh <- err
		}
	}()

	s.add(client)
	client.Listen()
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *server) Listen() {
	http.Handle(s.pattern, websocket.Handler(s.onConnection))

	log.Tracef("Starting the websocket listener on: %s", s.pattern)

	for {
		select {
		// add new a client
		case c := <-s.addCh:
			s.Clients[c.socketID] = c
			s.listener.ClientAdded(c.GetViewerClientFromChatClient())
			s.sendWelcomeMessageToClient(c)

		// remove a client
		case c := <-s.delCh:
			delete(s.Clients, c.socketID)
			s.listener.ClientRemoved(c.socketID)

			// message was recieved from a client and should be sanitized, validated
			// and distributed to other clients.
		case msg := <-s.sendAllCh:
			// Will turn markdown into html, sanitize user-supplied raw html
			// and standardize this message into something safe we can send everyone else.
			msg.RenderAndSanitizeMessageBody()

			s.listener.MessageSent(msg)
			s.sendAll(msg)

			// Store in the message history
			addMessage(msg)
		case ping := <-s.pingCh:
			fmt.Println("PING?", ping)

		case err := <-s.errCh:
			log.Error("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s *server) sendWelcomeMessageToClient(c *Client) {
	go func() {
		// Add an artificial delay so people notice this message come in.
		time.Sleep(7 * time.Second)

		initialChatMessageText := fmt.Sprintf("Welcome to %s! %s", config.Config.InstanceDetails.Title, config.Config.InstanceDetails.Summary)
		initialMessage := models.ChatMessage{"owncast-server", config.Config.InstanceDetails.Name, initialChatMessageText, "initial-message-1", "SYSTEM", true, time.Now()}
		c.Write(initialMessage)
	}()

}

func (s *server) getClientForClientID(clientID string) *Client {
	for _, client := range s.Clients {
		if client.ClientID == clientID {
			return client
		}
	}

	return nil
}
