package models

import (
	"github.com/owncast/owncast/utils"
)

// Stats holds the stats for the system.
type Stats struct {
	SessionMaxViewerCount int             `json:"sessionMaxViewerCount"`
	OverallMaxViewerCount int             `json:"overallMaxViewerCount"`
	LastDisconnectTime    *utils.NullTime `json:"lastDisconnectTime"`

	StreamConnected bool               `json:"-"`
	LastConnectTime *utils.NullTime    `json:"-"`
	ChatClients     map[string]Client  `json:"-"`
	Viewers         map[string]*Viewer `json:"-"`
}
