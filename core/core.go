package core

import (
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gabek/owncast/core/storageproviders"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/ffmpeg"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/yp"
)

var (
	_stats        *models.Stats
	_storage      models.StorageProvider
	_transcoder   *ffmpeg.Transcoder
	_cleanupTimer *time.Timer
	_yp           *yp.YP
	_broadcaster  *models.Broadcaster
)

var handler ffmpeg.HLSHandler
var fileWriter = ffmpeg.FileWriterReceiverService{}

//Start starts up the core processing
func Start() error {
	resetDirectories()

	if err := setupStats(); err != nil {
		log.Error("failed to setup the stats")
		return err
	}

	if err := setupStorage(); err != nil {
		log.Error("failed to setup the storage")
		return err
	}

	handler = ffmpeg.HLSHandler{}
	handler.Storage = _storage
	fileWriter.SetupFileWriterReceiverService(&handler)

	if err := createInitialOfflineState(); err != nil {
		log.Error("failed to create the initial offline state")
		return err
	}

	if config.Config.YP.Enabled {
		_yp = yp.NewYP(GetStatus)
	} else {
		yp.DisplayInstructions()
	}

	if config.Config.S3.Enabled {
		_storage = &storageproviders.S3Storage{}
	} else {
		_storage = &storageproviders.LocalStorage{}
	}

	chat.Setup(ChatListenerImpl{})

	return nil
}

func createInitialOfflineState() error {
	// Provide default files
	if !utils.DoesFileExists(filepath.Join(config.WebRoot, "thumbnail.jpg")) {
		if err := utils.Copy("static/logo.png", filepath.Join(config.WebRoot, "thumbnail.jpg")); err != nil {
			return err
		}
	}

	SetStreamAsDisconnected()

	return nil
}

func startCleanupTimer() {
	_cleanupTimer = time.NewTimer(5 * time.Minute)
	go func() {
		for {
			select {
			case <-_cleanupTimer.C:
				// Reset the session count since the session is over
				_stats.SessionMaxViewerCount = 0
				resetDirectories()
				ffmpeg.ShowStreamOfflineState()
			}
		}
	}()
}

// StopCleanupTimer will stop the previous cleanup timer
func stopCleanupTimer() {
	if _cleanupTimer != nil {
		_cleanupTimer.Stop()
	}
}

func resetDirectories() {
	log.Trace("Resetting file directories to a clean slate.")

	// Wipe the public, web-accessible hls data directory
	os.RemoveAll(config.PublicHLSStoragePath)
	os.RemoveAll(config.PrivateHLSStoragePath)
	os.MkdirAll(config.PublicHLSStoragePath, 0777)
	os.MkdirAll(config.PrivateHLSStoragePath, 0777)

	// Remove the previous thumbnail
	os.Remove(filepath.Join(config.WebRoot, "thumbnail.jpg"))

	// Create private hls data dirs
	if len(config.Config.VideoSettings.StreamQualities) != 0 {
		for index := range config.Config.VideoSettings.StreamQualities {
			os.MkdirAll(path.Join(config.PrivateHLSStoragePath, strconv.Itoa(index)), 0777)
			os.MkdirAll(path.Join(config.PublicHLSStoragePath, strconv.Itoa(index)), 0777)
		}
	} else {
		os.MkdirAll(path.Join(config.PrivateHLSStoragePath, strconv.Itoa(0)), 0777)
		os.MkdirAll(path.Join(config.PublicHLSStoragePath, strconv.Itoa(0)), 0777)
	}
}
