package core

import (
	"context"
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/core/transcoder"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications"
	"github.com/owncast/owncast/utils"
)

// After the stream goes offline this timer fires a full cleanup after N min.
var _offlineCleanupTimer *time.Timer

// While a stream takes place cleanup old HLS content every N min.
var _onlineCleanupTicker *time.Ticker

var _currentBroadcast *models.CurrentBroadcast

var _onlineTimerCancelFunc context.CancelFunc

var _lastNotified *time.Time

// setStreamAsConnected sets the stream as connected.
func setStreamAsConnected(rtmpOut *io.PipeReader) {
	now := utils.NullTime{Time: time.Now(), Valid: true}
	_stats.StreamConnected = true
	_stats.LastDisconnectTime = nil
	_stats.LastConnectTime = &now
	_stats.SessionMaxViewerCount = 0

	_currentBroadcast = &models.CurrentBroadcast{
		LatencyLevel:   data.GetStreamLatencyLevel(),
		OutputSettings: data.GetStreamOutputVariants(),
	}

	StopOfflineCleanupTimer()
	startOnlineCleanupTimer()

	if _yp != nil {
		go _yp.Start()
	}

	segmentPath := config.HLSStoragePath

	if err := setupStorage(); err != nil {
		log.Fatalln("failed to setup the storage", err)
	}

	go func() {
		_transcoder = transcoder.NewTranscoder()
		_transcoder.TranscoderCompleted = func(error) {
			SetStreamAsDisconnected()
			_transcoder = nil
			_currentBroadcast = nil
		}
		_transcoder.SetStdin(rtmpOut)
		_transcoder.Start()
	}()

	go webhooks.SendStreamStatusEvent(models.StreamStarted)
	transcoder.StartThumbnailGenerator(segmentPath, data.FindHighestVideoQualityIndex(_currentBroadcast.OutputSettings))

	_ = chat.SendSystemAction("Stay tuned, the stream is **starting**!", true)
	chat.SendAllWelcomeMessage()

	// Send delayed notification messages.
	_onlineTimerCancelFunc = startLiveStreamNotificationsTimer()
}

// SetStreamAsDisconnected sets the stream as disconnected.
func SetStreamAsDisconnected() {
	_ = chat.SendSystemAction("The stream is ending.", true)

	now := utils.NullTime{Time: time.Now(), Valid: true}
	if _onlineTimerCancelFunc != nil {
		_onlineTimerCancelFunc()
	}

	_stats.StreamConnected = false
	_stats.LastDisconnectTime = &now
	_stats.LastConnectTime = nil
	_broadcaster = nil

	offlineFilename := "offline.ts"

	offlineFilePath, err := saveOfflineClipToDisk(offlineFilename)
	if err != nil {
		log.Errorln(err)
		return
	}

	transcoder.StopThumbnailGenerator()
	rtmp.Disconnect()

	if _yp != nil {
		_yp.Stop()
	}

	// If there is no current broadcast available the previous stream
	// likely failed for some reason. Don't try to append to it.
	// Just transition to offline.
	if _currentBroadcast == nil {
		stopOnlineCleanupTimer()
		transitionToOfflineVideoStreamContent()
		log.Errorln("unexpected nil _currentBroadcast")
		return
	}

	for index := range _currentBroadcast.OutputSettings {
		makeVariantIndexOffline(index, offlineFilePath, offlineFilename)
	}

	StartOfflineCleanupTimer()
	stopOnlineCleanupTimer()
	saveStats()

	go webhooks.SendStreamStatusEvent(models.StreamStopped)
}

// StartOfflineCleanupTimer will fire a cleanup after n minutes being disconnected.
func StartOfflineCleanupTimer() {
	_offlineCleanupTimer = time.NewTimer(5 * time.Minute)
	go func() {
		for range _offlineCleanupTimer.C {
			// Set video to offline state
			resetDirectories()
			transitionToOfflineVideoStreamContent()
		}
	}()
}

// StopOfflineCleanupTimer will stop the previous cleanup timer.
func StopOfflineCleanupTimer() {
	if _offlineCleanupTimer != nil {
		_offlineCleanupTimer.Stop()
	}
}

func startOnlineCleanupTimer() {
	_onlineCleanupTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for range _onlineCleanupTicker.C {
			transcoder.CleanupOldContent(config.HLSStoragePath)
		}
	}()
}

func stopOnlineCleanupTimer() {
	if _onlineCleanupTicker != nil {
		_onlineCleanupTicker.Stop()
	}
}

func startLiveStreamNotificationsTimer() context.CancelFunc {
	// Send delayed notification messages.
	c, cancelFunc := context.WithCancel(context.Background())
	_onlineTimerCancelFunc = cancelFunc
	go func(c context.Context) {
		select {
		case <-time.After(time.Minute * 2.0):
			if _lastNotified != nil && time.Since(*_lastNotified) < 10*time.Minute {
				return
			}

			// Send Fediverse message.
			if data.GetFederationEnabled() {
				log.Traceln("Sending Federated Go Live message.")
				if err := activitypub.SendLive(); err != nil {
					log.Errorln(err)
				}
			}

			// Send notification to those who have registered for them.
			if notifier, err := notifications.New(data.GetDatastore()); err != nil {
				log.Errorln(err)
			} else {
				notifier.Notify()
			}

			now := time.Now()
			_lastNotified = &now
		case <-c.Done():
		}
	}(c)

	return cancelFunc
}
