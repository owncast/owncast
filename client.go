package main

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

const channelBufSize = 100

// Chat client.
type Client struct {
	id     string
	ws     *websocket.Conn
	server *Server
	ch     chan *ChatMessage
	pingch chan *PingMessage

	doneCh chan bool
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		log.Panicln("ws cannot be nil")
	}

	if server == nil {
		log.Panicln("server cannot be nil")
	}

	ch := make(chan *ChatMessage, channelBufSize)
	doneCh := make(chan bool)
	pingch := make(chan *PingMessage)
	clientID := getClientIDFromRequest(ws.Request())
	return &Client{clientID, ws, server, ch, pingch, doneCh}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *ChatMessage) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

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
			c.server.Del(c)
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
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg ChatMessage
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.SendAll(&msg)
			}
		}
	}
}
