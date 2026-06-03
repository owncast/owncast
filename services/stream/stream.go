package stream

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/notifications"
	"github.com/owncast/owncast/services/transcoder"
	"github.com/owncast/owncast/utils"
)

// Start brings up the storage backend, transcoder, RTMP listener, chat,
// webhooks, and the directory/notification subsystems. Must be called
// exactly once after New().
func (s *Service) Start(_ context.Context) error {
	s.resetDirectories()

	if err := s.configRepository.VerifySettings(s.cfg.TemporaryStreamKey); err != nil {
		log.Error(err)
		return err
	}

	if err := s.setupStats(); err != nil {
		log.Error("failed to setup the stats")
		return err
	}

	// The HLS handler takes the written HLS playlists and segments and
	// makes storage decisions. Rather simple right now but will play more
	// useful when recordings come into play.
	s.handler = transcoder.HLSHandler{}

	if err := s.setupStorage(); err != nil {
		log.Errorln("storage error", err)
	}

	s.fileWriter.SetupFileWriterReceiverService(&s.handler)

	if err := s.createInitialOfflineState(); err != nil {
		log.Error("failed to create the initial offline state")
		return err
	}

	if err := s.chat.Start(); err != nil {
		log.Errorln(err)
	}

	// Begin accepting an inbound broadcast through the stream engine. The
	// local engine binds the in-process RTMP listener; lifecycle is reported
	// back to this service via the StreamEvents methods.
	if err := s.engine.Start(); err != nil {
		return err
	}

	rtmpPort := s.configRepository.GetRTMPPortNumber()
	if rtmpPort != 1935 {
		log.Infof("RTMP is accepting inbound streams on port %d.", rtmpPort)
	}

	s.webhooks.Start()

	return nil
}

// Stop releases anything the service is holding. Currently a no-op
// because individual goroutines and ffmpeg children are tied to
// stream-connect/disconnect rather than overall service lifetime, but the
// hook is here for graceful-shutdown plumbing later: process exit
// reclaims the goroutines and ffmpeg child today; future work cancels
// onlineTimerCancelFunc, stops tickers, and closes the transcoder
// cleanly.
func (s *Service) Stop(_ context.Context) {
}

func (s *Service) createInitialOfflineState() error {
	s.transitionToOfflineVideoStreamContent()
	return nil
}

// transitionToOfflineVideoStreamContent overwrites the current stream
// with the offline video stream state only. No live stream HLS segments
// will continue to be referenced.
func (s *Service) transitionToOfflineVideoStreamContent() {
	log.Traceln("Placing offline fMP4 content into HLS directories")

	offlineInitPath, offlineSegmentPath, err := saveOfflineFMP4ToDisk()
	if err != nil {
		log.Fatalln("unable to save offline fMP4 files:", err)
	}

	variants := s.configRepository.GetStreamOutputVariants()
	if len(variants) == 0 {
		variants = make([]models.StreamOutputVariant, 1)
	}
	for index := range variants {
		variantDir := filepath.Join(config.HLSStoragePath, fmt.Sprintf("%d", index))
		if err := os.MkdirAll(variantDir, 0o750); err != nil {
			log.Errorln("unable to create variant directory:", err)
			continue
		}
		s.makeVariantIndexOffline(index, offlineInitPath, offlineSegmentPath)
	}

	masterPlaylistPath := filepath.Join(config.HLSStoragePath, "stream.m3u8")
	masterTmp, err := os.CreateTemp(config.HLSStoragePath, "tmp-stream-*.m3u8")
	if err != nil {
		log.Errorln("unable to create master playlist temp file:", err)
	} else {
		_, _ = masterTmp.WriteString("#EXTM3U\n#EXT-X-VERSION:7\n#EXT-X-INDEPENDENT-SEGMENTS\n")
		for index := range variants {
			_, _ = fmt.Fprintf(masterTmp, "#EXT-X-STREAM-INF:BANDWIDTH=0\n%d/stream.m3u8\n", index)
		}
		if err := masterTmp.Close(); err != nil {
			log.Errorln("unable to close master playlist:", err)
		} else if err := utils.Move(masterTmp.Name(), masterPlaylistPath); err != nil {
			log.Errorln("unable to atomically replace master playlist:", err)
		} else if _, err := s.storage.Save(masterPlaylistPath, 0); err != nil {
			log.Errorln("unable to save master playlist:", err)
		}
	}

	logo := s.configRepository.GetLogoPath()
	dst := filepath.Join(s.cfg.TempDir, "thumbnail.jpg")
	if err = utils.Copy(filepath.Join("data", logo), dst); err != nil {
		log.Warnln(err)
	}
	_ = os.Remove(path.Join(config.DataDirectory, "preview.gif"))
}

func (s *Service) resetDirectories() {
	log.Trace("Resetting file directories to a clean slate.")

	// Wipe hls data directory
	utils.CleanupDirectory(config.HLSStoragePath)

	// Remove the previous thumbnail
	logo := s.configRepository.GetLogoPath()
	if utils.DoesFileExists(logo) {
		err := utils.Copy(path.Join("data", logo), filepath.Join(config.DataDirectory, "thumbnail.jpg"))
		if err != nil {
			log.Warnln(err)
		}
	}
}

// setStreamAsConnected is the local engine's RTMP on-connect callback. It
// reports the broadcast online (StreamConnected) and then spins up the
// in-process transcoder + thumbnail generator fed by the inbound RTMP pipe.
func (s *Service) setStreamAsConnected(rtmpOut *io.PipeReader) {
	broadcast := &models.CurrentBroadcast{
		LatencyLevel:   s.configRepository.GetStreamLatencyLevel(),
		OutputSettings: s.configRepository.GetStreamOutputVariants(),
	}
	s.StreamConnected(broadcast)

	segmentPath := config.HLSStoragePath

	go func() {
		s.transcoder = transcoder.NewTranscoder(s.cfg, s.configRepository)
		s.transcoder.TranscoderCompleted = func(err error) {
			s.StreamDisconnected(err)
			s.transcoder = nil
		}
		s.transcoder.SetStdin(rtmpOut)
		s.transcoder.Start(true)
	}()

	selectedThumbnailVideoQualityIndex, isVideoPassthrough := s.configRepository.FindHighestVideoQualityIndex(broadcast.OutputSettings)
	s.thumbnailGen = transcoder.NewThumbnailGenerator(s.cfg, s.configRepository)
	s.thumbnailGen.Start(segmentPath, selectedThumbnailVideoQualityIndex, isVideoPassthrough)
}

// StreamConnected records a broadcast going live and runs the always-on
// go-live side effects. Implements StreamEvents; the local engine calls it
// in-process, a remote engine over its signaling channel.
func (s *Service) StreamConnected(broadcast *models.CurrentBroadcast) {
	s.currentBroadcast = broadcast
	s.applyStreamOnline()
}

// applyStreamOnline runs the always-on side effects of a stream going live:
// stats, cleanup/notification timers, storage setup, and the chat / webhook /
// federation go-live announcements. It deliberately excludes the transcoder
// and thumbnail pipeline (engine work), so a remote engine can drive just
// these reactions through StreamConnected.
func (s *Service) applyStreamOnline() {
	now := utils.NullTime{Time: time.Now(), Valid: true}
	s.statsMu.Lock()
	s.stats.StreamConnected = true
	s.stats.LastDisconnectTime = nil
	s.stats.LastConnectTime = &now
	s.stats.SessionMaxViewerCount = 0
	s.statsMu.Unlock()

	s.StopOfflineCleanupTimer()
	s.startOnlineCleanupTimer()

	if s.yp != nil {
		go s.yp.Start()
	}

	if err := s.setupStorage(); err != nil {
		log.Fatalln("failed to setup the storage", err)
	}

	go s.webhooks.SendStreamStatusEvent(models.StreamStarted)

	_ = s.chat.SendSystemAction("Stay tuned, the stream is **starting**!", true)
	s.chat.SendAllWelcomeMessage()

	// Send delayed notification messages.
	s.onlineTimerCancelFunc = s.startLiveStreamNotificationsTimer()

	// Let peer Owncast servers refresh their featured-streams listing.
	// Send one Offer immediately so a go-live propagates to followers
	// without waiting for the first periodic tick, then start the ticker
	// to keep the listing fresh for the duration of the stream.
	if s.configRepository.GetFederationEnabled() {
		if err := s.activitypub.SendStreamPing(); err != nil {
			log.Errorf("unable to send immediate go-live stream ping: %v", err)
		}
		s.activitypub.StartStreamPingTicker()
	}
}

// StreamDisconnected records a broadcast ending and clears the in-flight
// broadcast settings. Implements StreamEvents; the local engine calls it when
// the transcoder exits, a remote engine when its stream goes offline.
func (s *Service) StreamDisconnected(_ error) {
	s.applyStreamOffline()
	s.currentBroadcast = nil
}

// applyStreamOffline handles the always-on cleanup when a live stream ends.
func (s *Service) applyStreamOffline() {
	_ = s.chat.SendSystemAction("The stream is ending.", true)

	now := utils.NullTime{Time: time.Now(), Valid: true}
	s.statsMu.Lock()
	s.stats.StreamConnected = false
	s.stats.LastDisconnectTime = &now
	s.stats.LastConnectTime = nil
	s.statsMu.Unlock()
	s.broadcaster = nil

	// Stop the federated stream-ping ticker so we don't keep advertising
	// that we are live after the stream has ended, then tell followers we've
	// gone offline so they drop us from the live section of their
	// featured-streams directory immediately instead of waiting for the
	// staleness sweep.
	s.activitypub.StopStreamPingTicker()
	if s.configRepository.GetFederationEnabled() {
		if err := s.activitypub.SendStreamGoingOffline(); err != nil {
			log.Errorf("unable to send stream-offline Leave activity: %v", err)
		}
	}

	offlineInitPath, offlineSegmentPath, err := saveOfflineFMP4ToDisk()
	if err != nil {
		log.Errorln(err)
		return
	}
	// Clean up temp files after all variants have been updated.
	defer func() {
		_ = os.Remove(offlineInitPath)
		_ = os.Remove(offlineSegmentPath)
	}()

	if s.thumbnailGen != nil {
		s.thumbnailGen.Stop()
	}
	s.rtmp.Disconnect()

	if s.yp != nil {
		s.yp.Stop()
	}

	// If there is no current broadcast available the previous stream
	// likely failed for some reason. Don't try to append to it. Just
	// transition to offline.
	if s.currentBroadcast == nil {
		s.stopOnlineCleanupTimer()
		s.transitionToOfflineVideoStreamContent()
		log.Errorln("unexpected nil currentBroadcast")
		return
	}

	for index := range s.currentBroadcast.OutputSettings {
		s.makeVariantIndexOffline(index, offlineInitPath, offlineSegmentPath)
	}

	s.StartOfflineCleanupTimer()
	s.stopOnlineCleanupTimer()
	s.saveStats()

	go s.webhooks.SendStreamStatusEvent(models.StreamStopped)
}

// StartOfflineCleanupTimer fires a cleanup after n minutes being
// disconnected.
func (s *Service) StartOfflineCleanupTimer() {
	s.offlineCleanupTimer = time.NewTimer(5 * time.Minute)
	go func() {
		for range s.offlineCleanupTimer.C {
			// Set video to offline state
			s.resetDirectories()
			s.transitionToOfflineVideoStreamContent()
		}
	}()
}

// StopOfflineCleanupTimer stops the previous offline cleanup timer.
func (s *Service) StopOfflineCleanupTimer() {
	if s.offlineCleanupTimer != nil {
		s.offlineCleanupTimer.Stop()
	}
}

func (s *Service) startOnlineCleanupTimer() {
	s.onlineCleanupTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for range s.onlineCleanupTicker.C {
			if err := s.storage.Cleanup(); err != nil {
				log.Errorln(err)
			}
		}
	}()
}

func (s *Service) stopOnlineCleanupTimer() {
	if s.onlineCleanupTicker != nil {
		s.onlineCleanupTicker.Stop()
	}
}

func (s *Service) startLiveStreamNotificationsTimer() context.CancelFunc {
	// Send delayed notification messages.
	c, cancelFunc := context.WithCancel(context.Background())
	s.onlineTimerCancelFunc = cancelFunc
	go func(c context.Context) {
		select {
		case <-time.After(time.Minute * 2.0):
			if s.lastNotified != nil && time.Since(*s.lastNotified) < 10*time.Minute {
				return
			}

			// Send Fediverse message.
			if s.configRepository.GetFederationEnabled() {
				log.Traceln("Sending Federated Go Live message.")
				if err := s.activitypub.SendLive(); err != nil {
					log.Errorln(err)
				}
			}

			// Send notification to those who have registered for them.
			if notificationService, err := notifications.New(s.datastore, s.configRepository); err != nil {
				log.Errorln(err)
			} else {
				notificationService.Notify()
			}

			now := time.Now()
			s.lastNotified = &now
		case <-c.Done():
		}
	}(c)

	return cancelFunc
}
