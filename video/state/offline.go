package state

import (
	"os"
	"path"
	"path/filepath"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/video/transcoder"
	log "github.com/sirupsen/logrus"
)

// transitionToOfflineVideoStreamContent will overwrite the current stream with the
// offline video stream state only.  No live stream HLS segments will continue to be
// referenced.
func transitionToOfflineVideoStreamContent() {
	log.Traceln("Firing transcoder with offline stream state")

	_transcoder := transcoder.NewTranscoder()
	_transcoder.SetIdentifier("offline")
	_transcoder.SetLatencyLevel(models.GetLatencyLevel(4))
	_transcoder.SetIsEvent(true)

	offlineFilePath, err := saveOfflineClipToDisk("offline.ts")
	if err != nil {
		log.Fatalln("unable to save offline clip:", err)
	}

	_transcoder.SetInput(offlineFilePath)
	go _transcoder.Start(false)

	configRepository := configrepository.Get()

	// Copy the logo to be the thumbnail
	logo := configRepository.GetLogoPath()
	c := config.Get()
	dst := filepath.Join(c.TempDir, "thumbnail.jpg")
	if err = utils.Copy(filepath.Join("data", logo), dst); err != nil {
		log.Warnln(err)
	}

	// Delete the preview Gif
	_ = os.Remove(path.Join(config.DataDirectory, "preview.gif"))
}

func createInitialOfflineState() error {
	transitionToOfflineVideoStreamContent()

	return nil
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

	webhookManager := webhooks.Get()
	go webhookManager.SendStreamStatusEvent(models.StreamStopped)
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
