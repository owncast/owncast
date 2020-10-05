package models

import "github.com/owncast/owncast/utils"

//Status represents the status of the system
type Status struct {
	Online                bool `json:"online"`
	ViewerCount           int  `json:"viewerCount"`
	OverallMaxViewerCount int  `json:"overallMaxViewerCount"`
	SessionMaxViewerCount int  `json:"sessionMaxViewerCount"`

	LastConnectTime    utils.NullTime `json:"lastConnectTime"`
	LastDisconnectTime utils.NullTime `json:"lastDisconnectTime"`
}
