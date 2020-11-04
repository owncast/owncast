package core

import (
	"github.com/owncast/owncast/config"
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
		VersionNumber:         config.Config.VersionNumber,
	}
}

// setBroadcaster will store the current inbound broadcasting details
func setBroadcaster(broadcaster models.Broadcaster) {
	_broadcaster = &broadcaster
}

func GetBroadcaster() *models.Broadcaster {
	return _broadcaster
}
