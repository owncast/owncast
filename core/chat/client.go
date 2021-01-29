package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"

	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"

	"github.com/teris-io/shortid"
	"golang.org/x/time/rate"
)

const channelBufSize = 100

//Client represents a chat client.
type Client struct {
	ConnectedAt  time.Time
	MessageCount int
	UserAgent    string
	IPAddress    string
	Username     *string
	ClientID     string            // How we identify unique viewers when counting viewer counts.
	Geo          *geoip.GeoDetails `json:"geo"`
	Ignore       bool              // If set to true this will not be treated as a viewer

	socketID              string // How we identify a single websocket client.
	ws                    *websocket.Conn
	ch                    chan models.ChatEvent
	pingch                chan models.PingMessage
	usernameChangeChannel chan models.NameChangeEvent
	userJoinedChannel     chan models.UserJoinedEvent

	doneCh chan bool

	rateLimiter *rate.Limiter
}

// NewClient creates a new chat client.
func NewClient(ws *websocket.Conn) *Client {
	if ws == nil {
		log.Panicln("ws cannot be nil")
	}

	var ignoreClient = false
	for _, extraData := range ws.Config().Protocol {
		if extraData == "IGNORE_CLIENT" {
			ignoreClient = true
		}
	}

	ch := make(chan models.ChatEvent, channelBufSize)
	doneCh := make(chan bool)
	pingch := make(chan models.PingMessage)
	usernameChangeChannel := make(chan models.NameChangeEvent)
	userJoinedChannel := make(chan models.UserJoinedEvent)

	ipAddress := utils.GetIPAddressFromRequest(ws.Request())
	userAgent := ws.Request().UserAgent()
	socketID, _ := shortid.Generate()
	clientID := socketID

	rateLimiter := rate.NewLimiter(0.6, 5)

	return &Client{time.Now(), 0, userAgent, ipAddress, nil, clientID, nil, ignoreClient, socketID, ws, ch, pingch, usernameChangeChannel, userJoinedChannel, doneCh, rateLimiter}
}

func (c *Client) write(msg models.ChatEvent) {
	select {
	case c.ch <- msg:
	default:
		_server.removeClient(c)
		_server.err(fmt.Errorf("client %s is disconnected", c.ClientID))
	}
}

// Listen Write and Read request via channel.
func (c *Client) listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via channel.
func (c *Client) listenWrite() {
	for {
		select {
		// Send a PING keepalive
		case msg := <-c.pingch:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				c.handleClientSocketError(err)
			}
		// send message to the client
		case msg := <-c.ch:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				c.handleClientSocketError(err)
			}
		case msg := <-c.usernameChangeChannel:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				c.handleClientSocketError(err)
			}
		case msg := <-c.userJoinedChannel:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				c.handleClientSocketError(err)
			}

		// receive done request
		case <-c.doneCh:
			_server.removeClient(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

func (c *Client) handleClientSocketError(err error) {
	log.Warnln("Websocket client error: ", err.Error())
	_server.removeClient(c)
}

func (c *Client) passesRateLimit() bool {
	if !c.rateLimiter.Allow() {
		log.Debugln("Client", c.ClientID, "has exceeded the messaging rate limiting thresholds.")
		return false
	}

	return true
}

// Listen read request via channel.
func (c *Client) listenRead() {
	for {
		select {
		// receive done request
		case <-c.doneCh:
			_server.remove(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var data []byte
			err := websocket.Message.Receive(c.ws, &data)
			if err != nil {
				if err == io.EOF {
					c.doneCh <- true
				} else {
					c.handleClientSocketError(err)
				}
				return
			}

			var messageTypeCheck map[string]interface{}
			err = json.Unmarshal(data, &messageTypeCheck)
			if err != nil {
				log.Errorln(err)
			}

			messageType := messageTypeCheck["type"].(string)

			if !c.passesRateLimit() {
				continue
			}

			if messageType == models.MessageSent {
				c.chatMessageReceived(data)
			} else if messageType == models.UserNameChanged {
				c.userChangedName(data)
			} else if messageType == models.UserJoined {
				c.userJoined(data)
			}
		}
	}
}

func (c *Client) userJoined(data []byte) {
	var msg models.UserJoinedEvent
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Errorln(err)
		return
	}

	msg.ID = shortid.MustGenerate()
	msg.Type = models.UserJoined
	msg.Timestamp = time.Now()

	c.Username = &msg.Username

	_server.userJoined(msg)
}

func (c *Client) userChangedName(data []byte) {
	var msg models.NameChangeEvent
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Errorln(err)
	}
	msg.Type = models.UserNameChanged
	msg.ID = shortid.MustGenerate()
	_server.usernameChanged(msg)
	c.Username = &msg.NewName
}

func (c *Client) chatMessageReceived(data []byte) {
	var msg models.ChatEvent
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Errorln(err)
	}

	msg.SetDefaults()

	c.MessageCount++
	c.Username = &msg.Author

	msg.ClientID = c.ClientID
	_server.SendToAll(msg)
}

// GetViewerClientFromChatClient returns a general models.Client from a chat websocket client.
func (c *Client) GetViewerClientFromChatClient() models.Client {
	return models.Client{
		ConnectedAt:  c.ConnectedAt,
		MessageCount: c.MessageCount,
		UserAgent:    c.UserAgent,
		IPAddress:    c.IPAddress,
		Username:     c.Username,
		ClientID:     c.ClientID,
		Geo:          geoip.GetGeoFromIP(c.IPAddress),
	}
}
