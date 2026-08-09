package models

import "github.com/owncast/owncast/utils"

// Status represents the status of the system.
type Status struct {
	LastConnectTime    *utils.NullTime `json:"lastConnectTime"`
	LastDisconnectTime *utils.NullTime `json:"lastDisconnectTime"`

	VersionNumber         string `json:"versionNumber"`
	StreamTitle           string `json:"streamTitle"`
	ViewerCount           int    `json:"viewerCount"`
	OverallMaxViewerCount int    `json:"overallMaxViewerCount"`
	SessionMaxViewerCount int    `json:"sessionMaxViewerCount"`

	// StreamID identifies the currently-live broadcast. Empty when offline.
	// A viewer needs it to clip what is happening right now.
	StreamID string `json:"streamId,omitempty"`

	Online bool `json:"online"`
}
