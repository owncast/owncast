package chat

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/gorilla/websocket"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat/events"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/geoip"
)

// Client represents a single chat client.
type Client struct {
	id          uint
	accessToken string
	conn        *websocket.Conn
	User        *user.User `json:"user"`
	server      *Server
	ipAddress   string `json:"-"`
	// Buffered channel of outbound messages.
	send         chan []byte
	rateLimiter  *rate.Limiter
	timeoutTimer *time.Timer
	inTimeout    bool
	Geo          *geoip.GeoDetails `json:"geo"`
	MessageCount int               `json:"messageCount"`
	UserAgent    string            `json:"userAgent"`
	ConnectedAt  time.Time         `json:"connectedAt"`
	mu           sync.Mutex
}

type chatClientEvent struct {
	data   []byte
	client *Client
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

func (c *Client) sendConnectedClientInfo() {
	payload := events.EventPayload{
		"type": events.ConnectedUserInfo,
		"user": c.User,
	}

	c.sendPayload(payload)
}

func (c *Client) readPump() {
	// Allow 3 messages every two seconds.
	limit := rate.Every(2 * time.Second / 3)
	c.rateLimiter = rate.NewLimiter(limit, 1)

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

		// Check if this client is temporarily blocked from sending messages.
		if c.inTimeout {
			continue
		}

		// Guard against floods.
		if !c.passesRateLimit() {
			log.Warnln("Client", c.id, c.User.DisplayName, "has exceeded the messaging rate limiting thresholds and messages are being rejected temporarily.")
			c.startChatRejectionTimeout()

			continue
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.handleEvent(message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
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

			// Optimization: Send multiple events in a single websocket message.
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.send)
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

func (c *Client) handleEvent(data []byte) {
	c.server.inbound <- chatClientEvent{data: data, client: c}
}

func (c *Client) close() {
	defer func() {
		if a := recover(); a != nil {
			log.Println("RECOVER", a)
		}
	}()

	log.Traceln("client closed:", c.User.DisplayName, c.id, c.ipAddress)

	_ = c.conn.Close()
	c.server.unregister <- c.id
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.send != nil {
		close(c.send)
		c.send = nil
	}
}

func (c *Client) passesRateLimit() bool {
	return c.rateLimiter.Allow() && !c.inTimeout
}

func (c *Client) startChatRejectionTimeout() {
	if c.timeoutTimer != nil {
		return
	}

	c.inTimeout = true
	c.timeoutTimer = time.NewTimer(10 * time.Second)
	go func(c *Client) {
		for range c.timeoutTimer.C {
			c.inTimeout = false
			c.timeoutTimer = nil
		}
	}(c)

	c.sendAction("You are temporarily blocked from sending chat messages due to perceived flooding.")
}

func (c *Client) sendPayload(payload events.EventPayload) {
	var data []byte
	data, err := json.Marshal(payload)
	if err != nil {
		log.Errorln(err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.send != nil {
		c.send <- data
	}
}

func (c *Client) sendAction(message string) {
	clientMessage := events.ActionEvent{
		MessageEvent: events.MessageEvent{
			Body: message,
		},
	}
	clientMessage.SetDefaults()
	clientMessage.RenderBody()
	c.sendPayload(clientMessage.GetBroadcastPayload())
}
