package models

import (
	"net/http"
	"time"

	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/utils"
)

// ConnectedClientsResponse is the response of the currently connected chat clients.
type ConnectedClientsResponse struct {
	Clients []Client `json:"clients"`
}

// Client represents a single chat client.
type Client struct {
	ConnectedAt  time.Time         `json:"connectedAt"`
	LastSeen     time.Time         `json:"-"`
	Username     *string           `json:"username"`
	Geo          *geoip.GeoDetails `json:"geo"`
	UserAgent    string            `json:"userAgent"`
	IPAddress    string            `json:"ipAddress"`
	ClientID     string            `json:"clientID"`
	MessageCount int               `json:"messageCount"`
}

// GenerateClientFromRequest will return a chat client from a http request.
func GenerateClientFromRequest(req *http.Request) Client {
	return Client{
		ConnectedAt:  time.Now(),
		LastSeen:     time.Now(),
		MessageCount: 0,
		UserAgent:    req.UserAgent(),
		IPAddress:    utils.GetIPAddressFromRequest(req),
		Username:     nil,
		ClientID:     utils.GenerateClientIDFromRequest(req),
	}
}
