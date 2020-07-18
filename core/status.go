package core

import (
	"time"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core/ffmpeg"
	"github.com/gabek/owncast/models"
	"github.com/gabek/owncast/utils"
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

//SetStreamAsConnected sets the stream as connected
func SetStreamAsConnected() {
	_stats.StreamConnected = true
	_stats.LastConnectTime = utils.NullTime{time.Now(), true}

	timeSinceDisconnect := time.Since(_stats.LastDisconnectTime.Time).Minutes()
	if timeSinceDisconnect > 15 {
		_stats.SessionMaxViewerCount = 0
	}

	chunkPath := config.Config.GetPublicHLSSavePath()
	if usingExternalStorage {
		chunkPath = config.Config.GetPrivateHLSSavePath()
	}

	ffmpeg.StartThumbnailGenerator(chunkPath, config.Config.VideoSettings.HighestQualityStreamIndex)
}

//SetStreamAsDisconnected sets the stream as disconnected
func SetStreamAsDisconnected() {
	_stats.StreamConnected = false
	_stats.LastDisconnectTime = utils.NullTime{time.Now(), true}

	ffmpeg.ShowStreamOfflineState()
}
