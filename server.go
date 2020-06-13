package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	clients   map[string]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*Message{}
	clients := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	// Demo messages only.  Remove me eventually!!!
	messages = append(messages, &Message{"Tom Nook", "I'll be there with Bells on! Ho ho!", "https://gamepedia.cursecdn.com/animalcrossingpocketcamp_gamepedia_en/thumb/4/4f/Timmy_Icon.png/120px-Timmy_Icon.png?version=87b38d7d6130411d113486c2db151385", "demo-message-1"})
	messages = append(messages, &Message{"Redd", "Fool me once, shame on you. Fool me twice, stop foolin' me.", "https://vignette.wikia.nocookie.net/animalcrossing/images/3/3d/Redd2.gif/revision/latest?cb=20100710004252", "demo-message-2"})
	messages = append(messages, &Message{"Kevin", "You just caught me before I was about to go work out weeweewee!", "https://vignette.wikia.nocookie.net/animalcrossing/images/2/20/NH-Kevin_poster.png/revision/latest/scale-to-width-down/100?cb=20200410185817", "demo-message-3"})
	messages = append(messages, &Message{"Isabelle", " Isabelle is the mayor's highly capable secretary. She can be forgetful sometimes, but you can always count on her for information about the town. She wears her hair up in a bun that makes her look like a shih tzu. Mostly because she is one! She also has a twin brother named Digby.", "https://dodo.ac/np/images/thumb/7/7b/IsabelleTrophyWiiU.png/200px-IsabelleTrophyWiiU.png", "demo-message-4"})
	messages = append(messages, &Message{"Judy", "myohmy, I'm dancing my dreams away.", "https://vignette.wikia.nocookie.net/animalcrossing/images/5/50/NH-Judy_poster.png/revision/latest/scale-to-width-down/100?cb=20200522063219", "demo-message-5"})
	messages = append(messages, &Message{"Blathers", "Blathers is an owl with brown feathers. His face is white and he has a yellow beak. His arms are wing shaped and he has yellow talons. His eyes are very big with small black irises. He also has big pink cheek circles on his cheeks. His belly appears to be checkered in diamonds with light brown and white squares, similar to an argyle vest, which is traditionally associated with academia. His green bowtie further alludes to his academic nature.", "https://vignette.wikia.nocookie.net/animalcrossing/images/b/b3/NH-character-Blathers.png/revision/latest?cb=20200229053519", "demo-message-6"})

	return &Server{
		pattern,
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) ClientCount() int {
	return len(s.clients)
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			s.clients[c.id] = c
			viewerAdded(c.id)
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			delete(s.clients, c.id)
			viewerRemoved(c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
