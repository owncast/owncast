package chat

import (
	"bytes"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/gorilla/websocket"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/geoip"
)

type ChatClient struct {
	id          uint
	accessToken string
	conn        *websocket.Conn
	User        *user.User `json:"user"`
	server      *ChatServer
	ipAddress   string `json:"-"`
	// Buffered channel of outbound messages.
	send         chan []byte
	rateLimiter  *rate.Limiter
	Geo          *geoip.GeoDetails `json:"geo"`
	MessageCount int               `json:"messageCount"`
	UserAgent    string            `json:"userAgent"`
	ConnectedAt  time.Time         `json:"connectedAt"`
}

type chatClientEvent struct {
	data   []byte
	client *ChatClient
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	// Larger messages get thrown away.
	// Messages > *2 the socket gets closed.
	maxMessageSize = config.MaxSocketPayloadSize
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *ChatClient) sendConnectedClientInfo() {
	payload := events.EventPayload{
		"type": events.ConnectedUserInfo,
		"user": c.User,
	}

	c.sendPayload(payload)
}

func (c *ChatClient) readPump() {
	c.rateLimiter = rate.NewLimiter(0.6, 5)

	defer func() {
		c.close()
	}()

	// If somebody is sending 2x the max message size they're likely a bad actor
	// and should be disconnected. Below we throw away messages > max size.
	c.conn.SetReadLimit(maxMessageSize * 2)

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.close()
			}
			break
		}

		// Throw away messages greater than max message size.
		if len(message) > maxMessageSize {
			c.sendAction("Sorry, that message exceeded the maximum size and can't be delivered.")
			continue
		}

		// Guard against floods.
		if !c.passesRateLimit() {
			continue
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.handleEvent(message)
	}
}

func (c *ChatClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The server closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				log.Debugln(err)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *ChatClient) handleEvent(data []byte) {
	c.server.inbound <- chatClientEvent{data: data, client: c}
}

func (c *ChatClient) close() {
	log.Traceln("client closed:", c.User.DisplayName, c.id, c.ipAddress)

	c.conn.Close()
	c.server.unregister <- c
	if c.send != nil {
		close(c.send)
		c.send = nil
	}
}

func (c *ChatClient) passesRateLimit() bool {
	if !c.rateLimiter.Allow() {
		log.Debugln("Client", c.id, c.User.DisplayName, "has exceeded the messaging rate limiting thresholds.")
		return false
	}

	return true
}

func (c *ChatClient) sendPayload(payload events.EventPayload) {
	var data []byte
	data, err := json.Marshal(payload)
	if err != nil {
		log.Errorln(err)
		return
	}

	if c.send != nil {
		c.send <- data
	}
}

func (c *ChatClient) sendAction(message string) {
	clientMessage := events.ActionEvent{
		MessageEvent: events.MessageEvent{
			Body: message,
		},
	}
	clientMessage.SetDefaults()
	clientMessage.RenderBody()
	c.sendPayload(clientMessage.GetBroadcastPayload())
}
