package core

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/core/transcoder"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/yp"
)

var (
	_stats       *models.Stats
	_storage     models.StorageProvider
	_transcoder  *transcoder.Transcoder
	_yp          *yp.YP
	_broadcaster *models.Broadcaster
)

var handler transcoder.HLSHandler
var fileWriter = transcoder.FileWriterReceiverService{}

// Start starts up the core processing.
func Start() error {
	resetDirectories()

	data.PopulateDefaults()

	if err := data.VerifySettings(); err != nil {
		log.Error(err)
		return err
	}

	if err := setupStats(); err != nil {
		log.Error("failed to setup the stats")
		return err
	}

	// The HLS handler takes the written HLS playlists and segments
	// and makes storage decisions.  It's rather simple right now
	// but will play more useful when recordings come into play.
	handler = transcoder.HLSHandler{}

	if err := setupStorage(); err != nil {
		log.Errorln("storage error", err)
	}

	fileWriter.SetupFileWriterReceiverService(&handler)

	if err := createInitialOfflineState(); err != nil {
		log.Error("failed to create the initial offline state")
		return err
	}

	_yp = yp.NewYP(GetStatus)

	chat.Setup(ChatListenerImpl{})

	// start the rtmp server
	go rtmp.Start(setStreamAsConnected, setBroadcaster)

	rtmpPort := data.GetRTMPPortNumber()
	log.Infof("RTMP is accepting inbound streams on port %d.", rtmpPort)

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
	_transcoder := transcoder.NewTranscoder()
	_transcoder.SetInput(offlineFilePath)
	_transcoder.SetIdentifier("offline")
	_transcoder.Start()

	// Copy the logo to be the thumbnail
	logo := data.GetLogoPath()
	err := utils.Copy(filepath.Join("data", logo), "webroot/thumbnail.jpg")
	if err != nil {
		log.Warnln(err)
	}

	// Delete the preview Gif
	os.Remove(path.Join(config.WebRoot, "preview.gif"))
}

func resetDirectories() {
	log.Trace("Resetting file directories to a clean slate.")

	// Wipe the public, web-accessible hls data directory
	err := os.RemoveAll(config.PublicHLSStoragePath)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.RemoveAll(config.PrivateHLSStoragePath)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.MkdirAll(config.PublicHLSStoragePath, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll(config.PrivateHLSStoragePath, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	// Remove the previous thumbnail
	logo := data.GetLogoPath()
	if utils.DoesFileExists(logo) {
		err = utils.Copy(path.Join("data", logo), filepath.Join(config.WebRoot, "thumbnail.jpg"))
		if err != nil {
			log.Warnln(err)
		}
	}
}
