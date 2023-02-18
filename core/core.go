package core

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/auth"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/core/transcoder"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/geoip"
	"github.com/owncast/owncast/internal/activitypub"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/yp"
)

func New(ap *activitypub.Service, wh *webhooks.Service, yp *yp.YP, n *notifications.Notifier, c *chat.Service, r *rtmp.Service) (*Service, error) {
	s := &Service{
		ActivityPub: ap,
		Data:        ap.Persistence.Data,
		Webhooks:    wh,
		YP:          yp,
		Notifier:    n,
		Chat:        c,
		Rtmp:        r,
		Transcoder:  transcoder.NewTranscoder(ap.Persistence.Data),
	}

	return s, nil
}

type Service struct {
	mux        sync.RWMutex
	Storage    models.StorageProvider
	Transcoder *transcoder.Transcoder
	handler    transcoder.HLSHandler
	fileWriter transcoder.FileWriterReceiverService

	YP          *yp.YP
	broadcaster *models.Broadcaster
	stats       *models.Stats
	geoIpClient *geoip.Client
	Notifier    *notifications.Notifier
	Data        *data.Service
	Chat        *chat.Service
	Rtmp        *rtmp.Service

	ActivityPub *activitypub.Service
	Webhooks    *webhooks.Service

	streamState struct {
		_offlineCleanupTimer   *time.Timer  // After the stream goes offline this timer fires a full cleanup after N min.
		_onlineCleanupTicker   *time.Ticker // While a stream takes place cleanup old HLS content every N min.
		_currentBroadcast      *models.CurrentBroadcast
		_onlineTimerCancelFunc context.CancelFunc
		_lastNotified          *time.Time
	}
}

// Start starts up the core processing.
func (s *Service) Start() (err error) {
	s.resetDirectories()

	s.Data.PopulateDefaults()

	if err := s.Data.VerifySettings(); err != nil {
		log.Error(err)
		return err
	}

	if err := s.setupStats(); err != nil {
		log.Error("failed to setup the stats")
		return err
	}

	// The HLS handler takes the written HLS playlists and segments
	// and makes Storage decisions.  It's rather simple right now
	// but will play more useful when recordings come into play.
	s.handler = transcoder.HLSHandler{}

	if err := s.setupStorage(); err != nil {
		log.Errorln("Storage error", err)
	}

	auth.Setup(s.Data.Store)

	s.fileWriter.SetupFileWriterReceiverService(&s.handler)

	if err := s.createInitialOfflineState(); err != nil {
		log.Error("failed to create the initial offline state")
		return err
	}

	s.YP = yp.NewYP(s.GetStatus, s.Data)

	if s.Rtmp, err = rtmp.New(s.Data); err != nil {
		return fmt.Errorf("starting Rtmp service: %v", err)
	}

	// start the Rtmp server
	go s.Rtmp.Start(s.setStreamAsConnected, s.SetBroadcaster)

	rtmpPort := s.Data.GetRTMPPortNumber()
	log.Infof("RTMP is accepting inbound streams on port %d.", rtmpPort)

	s.Webhooks.InitWorkerPool()

	s.Notifier.Setup(s.Data.GetStore())

	if err = s.Chat.Start(s.GetStatus); err != nil {
		log.Errorln(err)
	}

	return nil
}

func (s *Service) createInitialOfflineState() error {
	s.transitionToOfflineVideoStreamContent()

	return nil
}

// transitionToOfflineVideoStreamContent will overwrite the current stream with the
// offline video stream state only.  No live stream HLS segments will continue to be
// referenced.
func (s *Service) transitionToOfflineVideoStreamContent() {
	log.Traceln("Firing Transcoder with offline stream state")

	_transcoder := transcoder.NewTranscoder(s.Data)
	_transcoder.SetIdentifier("offline")
	_transcoder.SetLatencyLevel(models.GetLatencyLevel(4))
	_transcoder.SetIsEvent(true)

	offlineFilePath, err := s.saveOfflineClipToDisk("offline.ts")
	if err != nil {
		log.Fatalln("unable to save offline clip:", err)
	}

	_transcoder.SetInput(offlineFilePath)
	go _transcoder.Start()

	// Copy the logo to be the thumbnail
	logo := s.Data.GetLogoPath()
	dst := filepath.Join(config.TempDir, "thumbnail.jpg")
	if err = utils.Copy(filepath.Join("data", logo), dst); err != nil {
		log.Warnln(err)
	}

	// Delete the preview Gif
	_ = os.Remove(path.Join(config.DataDirectory, "preview.gif"))
}

func (s *Service) resetDirectories() {
	log.Trace("Resetting file directories to a clean slate.")

	// Wipe hls Data directory
	utils.CleanupDirectory(config.HLSStoragePath)

	// Remove the previous thumbnail
	logo := s.Data.GetLogoPath()
	if utils.DoesFileExists(logo) {
		err := utils.Copy(path.Join("Data", logo), filepath.Join(config.DataDirectory, "thumbnail.jpg"))
		if err != nil {
			log.Warnln(err)
		}
	}
}
