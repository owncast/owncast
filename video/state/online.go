package state

import (
	"context"
	"io"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/apfederation/outbox"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/services/notifications"
	"github.com/owncast/owncast/services/webhooks"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/data"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/video/transcoder"
	log "github.com/sirupsen/logrus"
)

// setStreamAsConnected sets the stream as connected.
func (vs *VideoState) setStreamAsConnected(rtmpOut *io.PipeReader) {
	now := utils.NullTime{Time: time.Now(), Valid: true}
	vs.status.StreamConnected = true
	vs.status.Status.LastDisconnectTime = nil
	vs.status.Status.LastConnectTime = &now
	vs.status.SessionMaxViewerCount = 0

	configRepository := configrepository.Get()

	vs.status.SetCurrentBroadcast(&models.CurrentBroadcast{
		LatencyLevel:   configRepository.GetStreamLatencyLevel(),
		OutputSettings: configRepository.GetStreamOutputVariants(),
	})

	StopOfflineCleanupTimer()
	startOnlineCleanupTimer()

	if _yp != nil {
		go _yp.Start()
	}

	c := config.Get()
	segmentPath := c.HLSStoragePath

	if err := setupStorage(); err != nil {
		log.Fatalln("failed to setup the storage", err)
	}

	go func() {
		vs.transcoder = transcoder.NewTranscoder()
		vs.transcoder.TranscoderCompleted = func(error) {
			SetStreamAsDisconnected()
			vs.transcoder = nil
			vs.status.SetCurrentBroadcast(nil)
		}
		vs.transcoder.SetStdin(rtmpOut)
		vs.transcoder.Start(true)
	}()

	webhookManager := webhooks.Get()
	go webhookManager.SendStreamStatusEvent(models.StreamStarted)
	transcoder.StartThumbnailGenerator(segmentPath, configRepository.FindHighestVideoQualityIndex(vs.status.GetCurrentBroadcast().OutputSettings))

	_ = vs.chatService.SendSystemAction("Stay tuned, the stream is **starting**!", true)
	vs.chatService.SendAllWelcomeMessage()

	// Send delayed notification messages.
	_onlineTimerCancelFunc = startLiveStreamNotificationsTimer()
}

func startOnlineCleanupTimer() {
	_onlineCleanupTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for range _onlineCleanupTicker.C {
			if err := _storage.Cleanup(); err != nil {
				log.Errorln(err)
			}
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

			configRepository := configrepository.Get()

			// Send Fediverse message.
			if configRepository.GetFederationEnabled() {
				ob := outbox.Get()
				log.Traceln("Sending Federated Go Live message.")
				if err := ob.SendLive(); err != nil {
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
