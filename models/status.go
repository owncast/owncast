package models

import "github.com/owncast/owncast/utils"

// Status represents the status of the system.
type Status struct {
	Online                bool `json:"online"`
	ViewerCount           int  `json:"viewerCount"`
	OverallMaxViewerCount int  `json:"overallMaxViewerCount,omitempty"`
	SessionMaxViewerCount int  `json:"sessionMaxViewerCount,omitempty"`

	LastConnectTime    utils.NullTime `json:"lastConnectTime"`
	LastDisconnectTime utils.NullTime `json:"lastDisconnectTime"`

	VersionNumber string `json:"versionNumber"`
	StreamTitle   string `json:"streamTitle"`
}
