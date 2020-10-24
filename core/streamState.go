package core

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/ffmpeg"
	"github.com/owncast/owncast/utils"

	"github.com/grafov/m3u8"
)

// After the stream goes offline this timer fires a full cleanup after N min.
var _offlineCleanupTimer *time.Timer

// While a stream takes place cleanup old HLS content every N min.
var _onlineCleanupTicker *time.Ticker

//SetStreamAsConnected sets the stream as connected
func SetStreamAsConnected() {
	_stats.StreamConnected = true
	_stats.LastConnectTime = utils.NullTime{time.Now(), true}
	_stats.LastDisconnectTime = utils.NullTime{time.Now(), false}

	StopOfflineCleanupTimer()
	startOnlineCleanupTimer()

	if _yp != nil {
		_yp.Start()
	}

	segmentPath := config.PublicHLSStoragePath
	if config.Config.S3.Enabled {
		segmentPath = config.PrivateHLSStoragePath
	}

	go func() {
		_transcoder = ffmpeg.NewTranscoder()
		_transcoder.TranscoderCompleted = func(error) {

			SetStreamAsDisconnected()
		}
		_transcoder.Start()
	}()

	ffmpeg.StartThumbnailGenerator(segmentPath, config.Config.VideoSettings.HighestQualityStreamIndex)
}

//SetStreamAsDisconnected sets the stream as disconnected.
func SetStreamAsDisconnected() {
	_stats.StreamConnected = false
	_stats.LastDisconnectTime = utils.NullTime{time.Now(), true}
	_broadcaster = nil

	offlineFilename := "offline.ts"
	offlineFilePath := "static/" + offlineFilename

	ffmpeg.StopThumbnailGenerator()
	if _yp != nil {
		_yp.Stop()
	}

	for index := range config.Config.GetVideoStreamQualities() {
		playlistFilePath := fmt.Sprintf(filepath.Join(config.PrivateHLSStoragePath, "%d/stream.m3u8"), index)
		segmentFilePath := fmt.Sprintf(filepath.Join(config.PrivateHLSStoragePath, "%d/%s"), index, offlineFilename)

		utils.Copy(offlineFilePath, segmentFilePath)
		_storage.Save(segmentFilePath, 0)

		if utils.DoesFileExists(playlistFilePath) {
			f, err := os.OpenFile(playlistFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err != nil {
				log.Errorln(err)
			}
			defer f.Close()

			playlist, _, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
			variantPlaylist := playlist.(*m3u8.MediaPlaylist)
			if len(variantPlaylist.Segments) > config.Config.GetMaxNumberOfReferencedSegmentsInPlaylist() {
				variantPlaylist.Segments = variantPlaylist.Segments[:len(variantPlaylist.Segments)]
			}

			err = variantPlaylist.Append(offlineFilename, 8.0, "")
			variantPlaylist.SetDiscontinuity()
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
		_storage.Save(playlistFilePath, 0)
	}

	StartOfflineCleanupTimer()
	stopOnlineCleanupTimer()
}

// StartOfflineCleanupTimer will fire a cleanup after n minutes being disconnected
func StartOfflineCleanupTimer() {
	_offlineCleanupTimer = time.NewTimer(5 * time.Minute)
	go func() {
		for {
			select {
			case <-_offlineCleanupTimer.C:
				// Reset the session count since the session is over
				_stats.SessionMaxViewerCount = 0
				resetDirectories()
				transitionToOfflineVideoStreamContent()
			}
		}
	}()
}

// StopOfflineCleanupTimer will stop the previous cleanup timer
func StopOfflineCleanupTimer() {
	if _offlineCleanupTimer != nil {
		_offlineCleanupTimer.Stop()
	}
}

func startOnlineCleanupTimer() {
	_onlineCleanupTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-_onlineCleanupTicker.C:
				ffmpeg.CleanupOldContent(config.PrivateHLSStoragePath)
			}
		}
	}()
}

func stopOnlineCleanupTimer() {
	if _onlineCleanupTicker != nil {
		_onlineCleanupTicker.Stop()
	}
}
