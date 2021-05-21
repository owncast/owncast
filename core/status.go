package core

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

// GetStatus gets the status of the system.
func GetStatus() models.Status {
	if _stats == nil {
		return models.Status{}
	}

	viewerCount := 0
	if IsStreamConnected() {
		viewerCount = len(_stats.Viewers)
	}

	return models.Status{
		Online:                IsStreamConnected(),
		ViewerCount:           viewerCount,
		OverallMaxViewerCount: _stats.OverallMaxViewerCount,
		SessionMaxViewerCount: _stats.SessionMaxViewerCount,
		LastDisconnectTime:    _stats.LastDisconnectTime,
		LastConnectTime:       _stats.LastConnectTime,
		VersionNumber:         config.VersionNumber,
		StreamTitle:           data.GetStreamTitle(),
	}
}

func GetCurrentBroadcast() *models.CurrentBroadcast {
	return _currentBroadcast
}

// setBroadcaster will store the current inbound broadcasting details.
func setBroadcaster(broadcaster models.Broadcaster) {
	_broadcaster = &broadcaster
}

func GetBroadcaster() *models.Broadcaster {
	return _broadcaster
}
