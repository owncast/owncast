package chat

import (
	"bytes"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/gorilla/websocket"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/geoip"
)

type ChatClient struct {
	id          uint
	accessToken string
	conn        *websocket.Conn
	User        *user.User `json:"user"`
	server      *ChatServer
	IPAddress   string
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
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *ChatClient) readPump() {
	c.rateLimiter = rate.NewLimiter(0.6, 5)

	defer func() {
		c.close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
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

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write(newline); err != nil {
					log.Debugln(err)
				}
				if _, err := w.Write(<-c.send); err != nil {
					log.Debugln(err)
				}
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
	log.Traceln("client closed:", c.User.DisplayName, c.id, c.IPAddress)

	c.conn.Close()
	c.server.unregister <- c
}

func (c *ChatClient) passesRateLimit() bool {
	if !c.rateLimiter.Allow() {
		log.Debugln("Client", c.id, c.User.DisplayName, "has exceeded the messaging rate limiting thresholds.")
		return false
	}

	return true
}
