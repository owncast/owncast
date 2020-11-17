package core

import (
	"os"
	"path"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/ffmpeg"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/yp"
)

var (
	_stats       *models.Stats
	_storage     models.StorageProvider
	_transcoder  *ffmpeg.Transcoder
	_yp          *yp.YP
	_broadcaster *models.Broadcaster
)

var handler ffmpeg.HLSHandler
var fileWriter = ffmpeg.FileWriterReceiverService{}

// Start starts up the core processing.
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

	// The HLS handler takes the written HLS playlists and segments
	// and makes storage decisions.  It's rather simple right now
	// but will play more useful when recordings come into play.
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

	chat.Setup(ChatListenerImpl{})

	// start the rtmp server
	go rtmp.Start(setStreamAsConnected, setBroadcaster)

	port := config.Config.GetPublicWebServerPort()
	log.Infof("Web server is listening on port %d, RTMP is accepting inbound streams on port 1935.", port)
	log.Infoln("The web admin interface is available at /admin.")

	return nil
}

func createInitialOfflineState() error {
	// Provide default files
	if !utils.DoesFileExists(filepath.Join(config.WebRoot, "thumbnail.jpg")) {
		if err := utils.Copy("static/logo.png", filepath.Join(config.WebRoot, "thumbnail.jpg")); err != nil {
			return err
		}
	}

	transitionToOfflineVideoStreamContent()

	return nil
}

// transitionToOfflineVideoStreamContent will overwrite the current stream with the
// offline video stream state only.  No live stream HLS segments will continue to be
// referenced.
func transitionToOfflineVideoStreamContent() {
	log.Traceln("Firing transcoder with offline stream state")

	offlineFilename := "offline.ts"
	offlineFilePath := "static/" + offlineFilename
	_transcoder := ffmpeg.NewTranscoder()
	_transcoder.SetSegmentLength(10)
	_transcoder.SetInput(offlineFilePath)
	_transcoder.Start()

	// Copy the logo to be the thumbnail
	err := utils.Copy(filepath.Join("webroot", config.Config.InstanceDetails.Logo), "webroot/thumbnail.jpg")
	if err != nil {
		log.Warnln(err)
	}

	// Delete the preview Gif
	os.Remove(path.Join(config.WebRoot, "preview.gif"))
}

func resetDirectories() {
	log.Trace("Resetting file directories to a clean slate.")

	// Wipe the public, web-accessible hls data directory
	os.RemoveAll(config.PublicHLSStoragePath)
	os.RemoveAll(config.PrivateHLSStoragePath)
	err := os.MkdirAll(config.PublicHLSStoragePath, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll(config.PrivateHLSStoragePath, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	// Remove the previous thumbnail
	os.Remove(filepath.Join(config.WebRoot, "thumbnail.jpg"))

	// Create private hls data dirs
	if len(config.Config.VideoSettings.StreamQualities) != 0 {
		for index := range config.Config.VideoSettings.StreamQualities {
			err = os.MkdirAll(path.Join(config.PrivateHLSStoragePath, strconv.Itoa(index)), 0777)
			if err != nil {
				log.Fatalln(err)
			}

			err = os.MkdirAll(path.Join(config.PublicHLSStoragePath, strconv.Itoa(index)), 0777)
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		err = os.MkdirAll(path.Join(config.PrivateHLSStoragePath, strconv.Itoa(0)), 0777)
		if err != nil {
			log.Fatalln(err)
		}

		err = os.MkdirAll(path.Join(config.PublicHLSStoragePath, strconv.Itoa(0)), 0777)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Remove the previous thumbnail
	err = utils.Copy(path.Join(config.WebRoot, config.Config.InstanceDetails.Logo), "webroot/thumbnail.jpg")
	if err != nil {
		log.Warnln(err)
	}
}
