package models

import (
	"github.com/owncast/owncast/utils"
)

// Stats holds the stats for the system.
type Stats struct {
	LastDisconnectTime *utils.NullTime `json:"lastDisconnectTime"`

	LastConnectTime       *utils.NullTime    `json:"-"`
	ChatClients           map[string]Client  `json:"-"`
	Viewers               map[string]*Viewer `json:"-"`
	SessionMaxViewerCount int                `json:"sessionMaxViewerCount"`
	OverallMaxViewerCount int                `json:"overallMaxViewerCount"`

	StreamConnected bool `json:"-"`
}
