package chat

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"

	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/models"
)

//Server represents the server which handles the chat
type Server struct {
	Messages []models.ChatMessage
	Clients  map[string]*Client

	pattern   string
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan models.ChatMessage
	pingCh    chan models.PingMessage
	doneCh    chan bool
	errCh     chan error
}

//NewServer creates a new chat server
func NewServer() *Server {
	messages := []models.ChatMessage{}
	clients := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan models.ChatMessage)
	pingCh := make(chan models.PingMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)

	// Demo messages only.  Remove me eventually!!!
	messages = append(messages, models.ChatMessage{"Tom Nook", "I'll be there with Bells on! Ho ho!", "https://gamepedia.cursecdn.com/animalcrossingpocketcamp_gamepedia_en/thumb/4/4f/Timmy_Icon.png/120px-Timmy_Icon.png?version=87b38d7d6130411d113486c2db151385", "demo-message-1", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"Redd", "Fool me once, shame on you. Fool me twice, stop foolin' me.", "https://vignette.wikia.nocookie.net/animalcrossing/images/3/3d/Redd2.gif/revision/latest?cb=20100710004252", "demo-message-2", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"Kevin", "You just caught me before I was about to go work out weeweewee!", "https://vignette.wikia.nocookie.net/animalcrossing/images/2/20/NH-Kevin_poster.png/revision/latest/scale-to-width-down/100?cb=20200410185817", "demo-message-3", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"Isabelle", " Isabelle is the mayor's highly capable secretary. She can be forgetful sometimes, but you can always count on her for information about the town. She wears her hair up in a bun that makes her look like a shih tzu. Mostly because she is one! She also has a twin brother named Digby.", "https://dodo.ac/np/images/thumb/7/7b/IsabelleTrophyWiiU.png/200px-IsabelleTrophyWiiU.png", "demo-message-4", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"Judy", "myohmy, I'm dancing my dreams away.", "https://vignette.wikia.nocookie.net/animalcrossing/images/5/50/NH-Judy_poster.png/revision/latest/scale-to-width-down/100?cb=20200522063219", "demo-message-5", "ChatMessage"})
	messages = append(messages, models.ChatMessage{"Blathers", "Blathers is an owl with brown feathers. His face is white and he has a yellow beak. His arms are wing shaped and he has yellow talons. His eyes are very big with small black irises. He also has big pink cheek circles on his cheeks. His belly appears to be checkered in diamonds with light brown and white squares, similar to an argyle vest, which is traditionally associated with academia. His green bowtie further alludes to his academic nature.", "https://vignette.wikia.nocookie.net/animalcrossing/images/b/b3/NH-character-Blathers.png/revision/latest?cb=20200229053519", "demo-message-6", "ChatMessage"})

	server := &Server{
		messages,
		clients,
		"/entry", //hardcoded due to the UI requiring this and it is not configurable
		addCh,
		delCh,
		sendAllCh,
		pingCh,
		doneCh,
		errCh,
	}

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				server.ping()
			}
		}
	}()

	return server
}

//Add adds a client to the server
func (s *Server) Add(c *Client) {
	s.addCh <- c
}

//Remove removes a client from the server
func (s *Server) Remove(c *Client) {
	s.delCh <- c
}

//SendToAll sends a message to all of the connected clients
func (s *Server) SendToAll(msg models.ChatMessage) {
	s.sendAllCh <- msg
}

//Done marks the server as done
func (s *Server) Done() {
	s.doneCh <- true
}

//Err handles an error
func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.Messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg models.ChatMessage) {
	for _, c := range s.Clients {
		c.Write(msg)
	}
}

func (s *Server) ping() {
	// fmt.Println("Start pinging....", len(s.clients))

	ping := models.PingMessage{MessageType: "PING"}
	for _, c := range s.Clients {
		c.pingch <- ping
	}
}

func (s *Server) onConnection(ws *websocket.Conn) {
	client := NewClient(ws, s)

	defer func() {
		log.Printf("The client was connected for %s and sent %d messages (%s)", time.Since(client.ConnectedAt), client.MessageCount, client.id)

		if err := ws.Close(); err != nil {
			s.errCh <- err
		}
	}()

	s.Add(client)
	client.Listen()
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	http.Handle(s.pattern, websocket.Handler(s.onConnection))

	log.Printf("Starting the websocket listener on: %s", s.pattern)

	for {
		select {
		// add new a client
		case c := <-s.addCh:
			s.Clients[c.id] = c

			core.SetClientActive(c.id)
			s.sendPastMessages(c)

		// remove a client
		case c := <-s.delCh:
			delete(s.Clients, c.id)
			core.RemoveClient(c.id)

		// broadcast a message to all clients
		case msg := <-s.sendAllCh:
			log.Println("Sending a message to all:", msg)
			s.Messages = append(s.Messages, msg)
			s.sendAll(msg)

		case ping := <-s.pingCh:
			fmt.Println("PING?", ping)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
