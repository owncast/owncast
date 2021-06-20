package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/core/transcoder"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"

	"github.com/grafov/m3u8"
)

// After the stream goes offline this timer fires a full cleanup after N min.
var _offlineCleanupTimer *time.Timer

// While a stream takes place cleanup old HLS content every N min.
var _onlineCleanupTicker *time.Ticker

var _currentBroadcast *models.CurrentBroadcast

// setStreamAsConnected sets the stream as connected.
func setStreamAsConnected() {
	_stats.StreamConnected = true
	_stats.LastConnectTime = utils.NullTime{Time: time.Now(), Valid: true}
	_stats.LastDisconnectTime = utils.NullTime{Time: time.Now(), Valid: false}
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

	segmentPath := config.PublicHLSStoragePath
	s3Config := data.GetS3Config()

	if err := setupStorage(); err != nil {
		log.Fatalln("failed to setup the storage", err)
	}

	if s3Config.Enabled {
		segmentPath = config.PrivateHLSStoragePath
	}

	go func() {
		_transcoder = transcoder.NewTranscoder()
		_transcoder.TranscoderCompleted = func(error) {
			SetStreamAsDisconnected()
			_transcoder = nil
			_currentBroadcast = nil
		}
		_transcoder.Start()
	}()

	go webhooks.SendStreamStatusEvent(models.StreamStarted)
	transcoder.StartThumbnailGenerator(segmentPath, data.FindHighestVideoQualityIndex(_currentBroadcast.OutputSettings))
}

// SetStreamAsDisconnected sets the stream as disconnected.
func SetStreamAsDisconnected() {
	_stats.StreamConnected = false
	_stats.LastDisconnectTime = utils.NullTime{Time: time.Now(), Valid: true}
	_broadcaster = nil

	offlineFilename := "offline.ts"
	offlineFilePath := "static/" + offlineFilename

	transcoder.StopThumbnailGenerator()
	rtmp.Disconnect()

	if _yp != nil {
		_yp.Stop()
	}

	for index := range data.GetStreamOutputVariants() {
		playlistFilePath := fmt.Sprintf(filepath.Join(config.PrivateHLSStoragePath, "%d/stream.m3u8"), index)
		segmentFilePath := fmt.Sprintf(filepath.Join(config.PrivateHLSStoragePath, "%d/%s"), index, offlineFilename)

		err := utils.Copy(offlineFilePath, segmentFilePath)
		if err != nil {
			log.Warnln(err)
		}
		_, err = _storage.Save(segmentFilePath, 0)
		if err != nil {
			log.Warnln(err)
		}
		if utils.DoesFileExists(playlistFilePath) {
			f, err := os.OpenFile(playlistFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err != nil {
				log.Errorln(err)
			}
			defer f.Close()

			playlist, _, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
			if err != nil {
				log.Fatalln(err)
			}

			variantPlaylist := playlist.(*m3u8.MediaPlaylist)
			if len(variantPlaylist.Segments) > int(data.GetStreamLatencyLevel().SegmentCount) {
				variantPlaylist.Segments = variantPlaylist.Segments[:len(variantPlaylist.Segments)]
			}

			err = variantPlaylist.Append(offlineFilename, 8.0, "")
			if err != nil {
				log.Fatalln(err)
			}
			err = variantPlaylist.SetDiscontinuity()
			if err != nil {
				log.Fatalln(err)
			}
			_, err = f.WriteAt(variantPlaylist.Encode().Bytes(), 0)
			if err != nil {
				log.Errorln(err)
			}
		} else {
			p, err := m3u8.NewMediaPlaylist(1, 1)
			if err != nil {
				log.Errorln(err)
			}

			// If "offline" content gets changed then change the duration below
			err = p.Append(offlineFilename, 8.0, "")
			if err != nil {
				log.Errorln(err)
			}

			p.Close()
			f, err := os.Create(playlistFilePath)
			if err != nil {
				log.Errorln(err)
			}
			defer f.Close()
			_, err = f.Write(p.Encode().Bytes())
			if err != nil {
				log.Errorln(err)
			}
		}
		_, err = _storage.Save(playlistFilePath, 0)
		if err != nil {
			log.Warnln(err)
		}
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
			transcoder.CleanupOldContent(config.PrivateHLSStoragePath)
		}
	}()
}

func stopOnlineCleanupTimer() {
	if _onlineCleanupTicker != nil {
		_onlineCleanupTicker.Stop()
	}
}
