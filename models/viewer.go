package models

import (
	"net/http"
	"time"

	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/utils"
)

// Viewer represents a single video viewer.
type Viewer struct {
	FirstSeen time.Time         `json:"firstSeen"`
	LastSeen  time.Time         `json:"-"`
	Geo       *geoip.GeoDetails `json:"geo"`
	UserAgent string            `json:"userAgent"`
	IPAddress string            `json:"ipAddress"`
	ClientID  string            `json:"clientID"`
}

// GenerateViewerFromRequest will return a chat client from a http request.
func GenerateViewerFromRequest(req *http.Request) Viewer {
	return Viewer{
		FirstSeen: time.Now(),
		LastSeen:  time.Now(),
		UserAgent: req.UserAgent(),
		IPAddress: utils.GetIPAddressFromRequest(req),
		ClientID:  utils.GenerateClientIDFromRequest(req),
	}
}
