package core

import (
	"time"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core/ffmpeg"
	"github.com/gabek/owncast/utils"
)

//SetStreamAsConnected sets the stream as connected
func SetStreamAsConnected() {
	_stats.StreamConnected = true
	_stats.LastConnectTime = utils.NullTime{time.Now(), true}
	_stats.LastDisconnectTime = utils.NullTime{time.Now(), false}

	timeSinceDisconnect := time.Since(_stats.LastDisconnectTime.Time).Minutes()
	if timeSinceDisconnect > 15 {
		_stats.SessionMaxViewerCount = 0
	}

	segmentPath := config.Config.GetPublicHLSSavePath()
	if config.Config.S3.Enabled {
		segmentPath = config.Config.GetPrivateHLSSavePath()
	}

	go func() {
		_transcoder = ffmpeg.NewTranscoder()
		_transcoder.Start()
	}()

	ffmpeg.StartThumbnailGenerator(segmentPath, config.Config.VideoSettings.HighestQualityStreamIndex)
}

//SetStreamAsDisconnected sets the stream as disconnected
func SetStreamAsDisconnected() {
	_stats.StreamConnected = false
	_stats.LastDisconnectTime = utils.NullTime{time.Now(), true}

	go func() {
		_transcoder := ffmpeg.NewTranscoder()
		_transcoder.SetSegmentLength(10)
		_transcoder.SetAppendToStream(false)
		_transcoder.SetInput(config.Config.GetOfflineContentPath())
		_transcoder.Start()
	}()
}
