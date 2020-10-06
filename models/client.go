package models

import (
	"net/http"
	"time"

	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/utils"
)

type ConnectedClientsResponse struct {
	Clients []Client `json:"clients"`
}

type Client struct {
	ConnectedAt  time.Time         `json:"connectedAt"`
	LastSeen     time.Time         `json:"-"`
	MessageCount int               `json:"messageCount"`
	UserAgent    string            `json:"userAgent"`
	IPAddress    string            `json:"ipAddress"`
	Username     *string           `json:"username"`
	ClientID     string            `json:"clientID"`
	Geo          *geoip.GeoDetails `json:"geo"`
}

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
