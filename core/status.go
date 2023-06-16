package core

import (
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
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

	c := config.GetConfig()

	return models.Status{
		Online:                IsStreamConnected(),
		ViewerCount:           viewerCount,
		OverallMaxViewerCount: _stats.OverallMaxViewerCount,
		SessionMaxViewerCount: _stats.SessionMaxViewerCount,
		LastDisconnectTime:    _stats.LastDisconnectTime,
		LastConnectTime:       _stats.LastConnectTime,
		VersionNumber:         c.VersionNumber,
		StreamTitle:           data.GetStreamTitle(),
	}
}

// GetCurrentBroadcast will return the currently active broadcast.
func GetCurrentBroadcast() *models.CurrentBroadcast {
	return _currentBroadcast
}

// setBroadcaster will store the current inbound broadcasting details.
func setBroadcaster(broadcaster models.Broadcaster) {
	_broadcaster = &broadcaster
}

// GetBroadcaster will return the details of the currently active broadcaster.
func GetBroadcaster() *models.Broadcaster {
	return _broadcaster
}
