package models

import (
	"net/http"
	"time"

	"github.com/owncast/owncast/services/geoip"
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

// PlaybackSource describes where a client's playback measurements came
// from.
type PlaybackSource string

const (
	// PlaybackSourceClient means the player reported the measurements
	// itself.
	PlaybackSourceClient PlaybackSource = "client"
	// PlaybackSourceServer means the server observed the measurements
	// while serving video to a player that reports nothing.
	PlaybackSourceServer PlaybackSource = "server"
)

// PlaybackClient is a single client currently playing back video and the
// most recent playback health measured for it.
type PlaybackClient struct {
	FirstSeen time.Time             `json:"firstSeen"`
	Geo       *geoip.GeoDetails     `json:"geo"`
	Playback  *PlaybackClientHealth `json:"playback"`
	UserAgent string                `json:"userAgent"`
	IPAddress string                `json:"ipAddress"`
	// ClientID identifies the playback client: the player-supplied CMCD
	// session ID when it provides one, which distinguishes players
	// sharing an address, otherwise the viewer identity.
	ClientID string `json:"clientID"`
	// ViewerID identifies the viewer this client belongs to.
	ViewerID string `json:"viewerID"`
}

// PlaybackClientHealth is the latest playback health of one client. A nil
// measurement is unknown for that client rather than measured as zero:
// not every player reports every value, and a client the server can only
// observe reveals nothing about its latency, quality, or errors.
type PlaybackClientHealth struct {
	LastUpdate            time.Time      `json:"lastUpdate"`
	BandwidthKbps         *float64       `json:"bandwidthKbps"`
	LatencySeconds        *float64       `json:"latencySeconds"`
	DownloadSeconds       *float64       `json:"downloadSeconds"`
	BitrateKbps           *float64       `json:"bitrateKbps"`
	ErrorCount            *float64       `json:"errorCount"`
	QualityVariantChanges *float64       `json:"qualityVariantChanges"`
	PlayerState           string         `json:"playerState,omitempty"`
	MeasurementStatus     string         `json:"measurementStatus,omitempty"`
	Source                PlaybackSource `json:"source"`
}
