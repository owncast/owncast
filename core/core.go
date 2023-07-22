package core

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/services/notifications"
	"github.com/owncast/owncast/services/status"
	"github.com/owncast/owncast/services/webhooks"
	"github.com/owncast/owncast/services/yp"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/data"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/video/rtmp"
	"github.com/owncast/owncast/video/transcoder"
)

var (
	_stats       *models.Stats
	_storage     models.StorageProvider
	_transcoder  *transcoder.Transcoder
	_yp          *yp.YP
	_broadcaster *models.Broadcaster
	handler      transcoder.HLSHandler
	fileWriter   = transcoder.FileWriterReceiverService{}
)

// Start starts up the core processing.
func Start() error {
	resetDirectories()
	configRepository := configrepository.Get()

	configRepository.PopulateDefaults()

	if err := configRepository.VerifySettings(); err != nil {
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

	// user.SetupUsers()
	// auth.Setup(data.GetDatastore())

	fileWriter.SetupFileWriterReceiverService(&handler)

	if err := createInitialOfflineState(); err != nil {
		log.Error("failed to create the initial offline state")
		return err
	}

	s := status.Get()
	gsf := func() *models.Status {
		s := status.Get()
		return &s.Status
	}

	_yp = yp.NewYP(gsf)

	if err := chat.Start(gsf); err != nil {
		log.Errorln(err)
	}

	// start the rtmp server
	go rtmp.Start(setStreamAsConnected, s.SetBroadcaster)

	rtmpPort := configRepository.GetRTMPPortNumber()
	if rtmpPort != 1935 {
		log.Infof("RTMP is accepting inbound streams on port %d.", rtmpPort)
	}

	webhooks.InitTemporarySingleton(gsf)

	notifications.Setup(data.GetDatastore())

	return nil
}
