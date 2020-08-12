package core

import (
	"github.com/owncast/owncast/models"
)

//GetStatus gets the status of the system
func GetStatus() models.Status {
	if _stats == nil {
		return models.Status{}
	}

	return models.Status{
		Online:                IsStreamConnected(),
		ViewerCount:           len(_stats.Clients),
		OverallMaxViewerCount: _stats.OverallMaxViewerCount,
		SessionMaxViewerCount: _stats.SessionMaxViewerCount,
		LastDisconnectTime:    _stats.LastDisconnectTime,
		LastConnectTime:       _stats.LastConnectTime,
	}
}
