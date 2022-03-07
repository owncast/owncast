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
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/core/webhooks"
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

var (
	handler    transcoder.HLSHandler
	fileWriter = transcoder.FileWriterReceiverService{}
)

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

	user.SetupUsers()

	fileWriter.SetupFileWriterReceiverService(&handler)

	if err := createInitialOfflineState(); err != nil {
		log.Error("failed to create the initial offline state")
		return err
	}

	_yp = yp.NewYP(GetStatus)

	if err := chat.Start(GetStatus); err != nil {
		log.Errorln(err)
	}

	// start the rtmp server
	go rtmp.Start(setStreamAsConnected, setBroadcaster)
	logRtmpStart()

	webhooks.InitWorkerPool()

	return nil
}

func logRtmpStart() {
	rtmpPort := data.GetRTMPPortNumber()
	tlsEnabled := data.GetRTMPTLSEnabled()
	tlsString := ""
	if tlsEnabled {
		tlsString = " tls"
	}
	log.Infof("RTMP is accepting inbound"+tlsString+" streams on port %d.", rtmpPort)
}

func createInitialOfflineState() error {
	transitionToOfflineVideoStreamContent()

	return nil
}

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
	go _transcoder.Start()

	// Copy the logo to be the thumbnail
	logo := data.GetLogoPath()
	if err = utils.Copy(filepath.Join("data", logo), "webroot/thumbnail.jpg"); err != nil {
		log.Warnln(err)
	}

	// Delete the preview Gif
	_ = os.Remove(path.Join(config.WebRoot, "preview.gif"))
}

func resetDirectories() {
	log.Trace("Resetting file directories to a clean slate.")

	// Wipe hls data directory
	utils.CleanupDirectory(config.HLSStoragePath)

	// Remove the previous thumbnail
	logo := data.GetLogoPath()
	if utils.DoesFileExists(logo) {
		err := utils.Copy(path.Join("data", logo), filepath.Join(config.WebRoot, "thumbnail.jpg"))
		if err != nil {
			log.Warnln(err)
		}
	}
}
