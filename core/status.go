package core

import (
	"time"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core/ffmpeg"
	"github.com/gabek/owncast/models"
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
	}
}

//SetStreamAsConnected sets the stream as connected
func SetStreamAsConnected() {
	_stats.StreamConnected = true
	_stats.LastConnectTime = time.Now()

	timeSinceDisconnect := time.Since(_stats.LastDisconnectTime).Minutes()
	if timeSinceDisconnect > 15 {
		_stats.SessionMaxViewerCount = 0
	}

	chunkPath := config.Config.PublicHLSPath
	if usingExternalStorage {
		chunkPath = config.Config.PrivateHLSPath
	}

	ffmpeg.StartThumbnailGenerator(chunkPath)
}

//SetStreamAsDisconnected sets the stream as disconnected
func SetStreamAsDisconnected() {
	_stats.StreamConnected = false
	_stats.LastDisconnectTime = time.Now()

	if config.Config.EnableOfflineImage {
		ffmpeg.ShowStreamOfflineState()
	}
}
