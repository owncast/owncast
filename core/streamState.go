package core

import (
	"context"
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/transcoder"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications"
	"github.com/owncast/owncast/utils"
)

// setStreamAsConnected sets the stream as connected.
func (s *Service) setStreamAsConnected(rtmpOut *io.PipeReader) {
	s.stats.StreamConnected = true
	s.stats.LastDisconnectTime = nil
	s.stats.LastConnectTime = &utils.NullTime{Time: time.Now(), Valid: true}
	s.stats.SessionMaxViewerCount = 0

	s.streamState._currentBroadcast = &models.CurrentBroadcast{
		LatencyLevel:   s.Data.GetStreamLatencyLevel(),
		OutputSettings: s.Data.GetStreamOutputVariants(),
	}

	s.StopOfflineCleanupTimer()
	s.startOnlineCleanupTimer()

	if s.YP != nil {
		go s.YP.Start()
	}

	segmentPath := config.HLSStoragePath

	if err := s.setupStorage(); err != nil {
		log.Fatalln("failed to setup the Storage", err)
	}

	go func() {
		s.Transcoder = transcoder.NewTranscoder(s.Data)
		s.Transcoder.TranscoderCompleted = func(error) {
			s.SetStreamAsDisconnected()
			s.Transcoder = nil
			s.streamState._currentBroadcast = nil
		}
		s.Transcoder.SetStdin(rtmpOut)
		s.Transcoder.Start()
	}()

	go s.Webhooks.SendStreamStatusEvent(models.StreamStarted)
	s.Transcoder.StartThumbnailGenerator(segmentPath, s.Data.FindHighestVideoQualityIndex(s.streamState._currentBroadcast.OutputSettings))

	_ = s.Chat.SendSystemAction("Stay tuned, the stream is **starting**!", true)
	s.Chat.SendAllWelcomeMessage()

	// Send delayed notification messages.
	s.streamState._onlineTimerCancelFunc = s.startLiveStreamNotificationsTimer()
}

// SetStreamAsDisconnected sets the stream as disconnected.
func (s *Service) SetStreamAsDisconnected() {
	_ = s.Chat.SendSystemAction("The stream is ending.", true)

	now := utils.NullTime{Time: time.Now(), Valid: true}
	if s.streamState._onlineTimerCancelFunc != nil {
		s.streamState._onlineTimerCancelFunc()
	}

	s.stats.StreamConnected = false
	s.stats.LastDisconnectTime = &now
	s.stats.LastConnectTime = nil
	s.broadcaster = nil

	offlineFilename := "offline.ts"

	offlineFilePath, err := s.saveOfflineClipToDisk(offlineFilename)
	if err != nil {
		log.Errorln(err)
		return
	}

	s.Transcoder.StopThumbnailGenerator()
	s.Rtmp.Disconnect()

	if s.YP != nil {
		s.YP.Stop()
	}

	// If there is no current broadcast available the previous stream
	// likely failed for some reason. Don't try to append to it.
	// Just transition to offline.
	if s.streamState._currentBroadcast == nil {
		s.stopOnlineCleanupTimer()
		s.transitionToOfflineVideoStreamContent()
		log.Errorln("unexpected nil _currentBroadcast")
		return
	}

	for index := range s.streamState._currentBroadcast.OutputSettings {
		s.makeVariantIndexOffline(index, offlineFilePath, offlineFilename)
	}

	s.StartOfflineCleanupTimer()
	s.stopOnlineCleanupTimer()
	s.saveStats()

	go s.Webhooks.SendStreamStatusEvent(models.StreamStopped)
}

// StartOfflineCleanupTimer will fire a cleanup after n minutes being disconnected.
func (s *Service) StartOfflineCleanupTimer() {
	s.streamState._offlineCleanupTimer = time.NewTimer(5 * time.Minute)
	go func() {
		for range s.streamState._offlineCleanupTimer.C {
			// Set video to offline state
			s.resetDirectories()
			s.transitionToOfflineVideoStreamContent()
		}
	}()
}

// StopOfflineCleanupTimer will stop the previous cleanup timer.
func (s *Service) StopOfflineCleanupTimer() {
	if s.streamState._offlineCleanupTimer != nil {
		s.streamState._offlineCleanupTimer.Stop()
	}
}

func (s *Service) startOnlineCleanupTimer() {
	s.streamState._onlineCleanupTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for range s.streamState._onlineCleanupTicker.C {
			s.Transcoder.CleanupOldContent(config.HLSStoragePath)
		}
	}()
}

func (s *Service) stopOnlineCleanupTimer() {
	if s.streamState._onlineCleanupTicker != nil {
		s.streamState._onlineCleanupTicker.Stop()
	}
}

func (s *Service) startLiveStreamNotificationsTimer() context.CancelFunc {
	// Send delayed notification messages.
	c, cancelFunc := context.WithCancel(context.Background())
	s.streamState._onlineTimerCancelFunc = cancelFunc
	go func(c context.Context) {
		select {
		case <-time.After(time.Minute * 2.0):
			if s.streamState._lastNotified != nil && time.Since(*s.streamState._lastNotified) < 10*time.Minute {
				return
			}

			// Send Fediverse message.
			if s.Data.GetFederationEnabled() {
				log.Traceln("Sending Federated Go Live message.")
				if err := s.ActivityPub.SendLive(); err != nil {
					log.Errorln(err)
				}
			}

			// Send notification to those who have registered for them.
			if notifier, err := notifications.New(s.Data); err != nil {
				log.Errorln(err)
			} else {
				notifier.Notify()
			}

			now := time.Now()
			s.streamState._lastNotified = &now
		case <-c.Done():
		}
	}(c)

	return cancelFunc
}
