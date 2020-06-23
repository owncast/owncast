package chat

import (
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"

	"github.com/gabek/owncast/models"
	"github.com/gabek/owncast/utils"
)

const channelBufSize = 100

//Client represents a chat client.
type Client struct {
	ConnectedAt  time.Time
	MessageCount int

	id     string
	ws     *websocket.Conn
	server *Server
	ch     chan models.ChatMessage
	pingch chan models.PingMessage

	doneCh chan bool
}

//NewClient creates a new chat client
func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		log.Panicln("ws cannot be nil")
	}

	if server == nil {
		log.Panicln("server cannot be nil")
	}

	ch := make(chan models.ChatMessage, channelBufSize)
	doneCh := make(chan bool)
	pingch := make(chan models.PingMessage)
	clientID := utils.GenerateClientIDFromRequest(ws.Request())

	return &Client{time.Now(), 0, clientID, ws, server, ch, pingch, doneCh}
}

//GetConnection gets the connection for the client
func (c *Client) GetConnection() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg models.ChatMessage) {
	select {
	case c.ch <- msg:
	default:
		c.server.Remove(c)
		c.server.Err(fmt.Errorf("client %s is disconnected", c.id))
	}
}

//Done marks the client as done
func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	for {
		select {
		// Send a PING keepalive
		case msg := <-c.pingch:
			websocket.JSON.Send(c.ws, msg)
		// send message to the client
		case msg := <-c.ch:
			msg.MessageType = "CHAT"
			// log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			c.server.Remove(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Remove(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg models.ChatMessage

			if err := websocket.JSON.Receive(c.ws, &msg); err == io.EOF {
				c.doneCh <- true
				return
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.MessageCount++
				c.server.SendToAll(msg)
			}
		}
	}
}
