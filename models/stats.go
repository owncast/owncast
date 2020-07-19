package models

import (
	"time"

	"github.com/gabek/owncast/utils"
)

//Stats holds the stats for the system
type Stats struct {
	SessionMaxViewerCount int            `json:"sessionMaxViewerCount"`
	OverallMaxViewerCount int            `json:"overallMaxViewerCount"`
	LastDisconnectTime    utils.NullTime `json:"lastDisconnectTime"`

	StreamConnected bool                 `json:"-"`
	LastConnectTime utils.NullTime       `json:"-"`
	Clients         map[string]time.Time `json:"-"`
}
