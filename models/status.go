package models

import (
	"time"
)

//Status represents the status of the system
type Status struct {
	Online                bool `json:"online"`
	ViewerCount           int  `json:"viewerCount"`
	OverallMaxViewerCount int  `json:"overallMaxViewerCount"`
	SessionMaxViewerCount int  `json:"sessionMaxViewerCount"`

	LastConnectTime    time.Time `json:"lastConnectTime"`
	LastDisconnectTime time.Time `json:"lastDisconnectTime"`
}
