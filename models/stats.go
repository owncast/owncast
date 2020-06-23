package models

import (
	"time"
)

//Stats holds the stats for the system
type Stats struct {
	SessionMaxViewerCount int       `json:"sessionMaxViewerCount"`
	OverallMaxViewerCount int       `json:"overallMaxViewerCount"`
	LastDisconnectTime    time.Time `json:"lastDisconnectTime"`

	StreamConnected bool                 `json:"-"`
	LastConnectTime time.Time            `json:"-"`
	Clients         map[string]time.Time `json:"-"`
}
