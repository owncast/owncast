package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/auth"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/core/transcoder"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
	"github.com/owncast/owncast/persistence/tables"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/yp"
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
	// configRepository.PopulateDefaults()

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

	tables.SetupUsers(data.GetDatastore().DB)
	auth.Setup(data.GetDatastore())

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

	rtmpPort := configRepository.GetRTMPPortNumber()
	if rtmpPort != 1935 {
		log.Infof("RTMP is accepting inbound streams on port %d.", rtmpPort)
	}

	webhooks.SetupWebhooks(GetStatus)

	notificationsrepository.Setup(data.GetStore())

	return nil
}

func createInitialOfflineState() error {
	transitionToOfflineVideoStreamContent()

	return nil
}

// transitionToOfflineVideoStreamContent will overwrite the current stream with the
// offline video stream state only.  No live stream HLS segments will continue to be
// referenced.
func transitionToOfflineVideoStreamContent() {
	log.Traceln("Placing offline fMP4 content into HLS directories")

	configRepository := configrepository.Get()

	offlineInitPath, offlineSegmentPath, err := saveOfflineFMP4ToDisk()
	if err != nil {
		log.Fatalln("unable to save offline fMP4 files:", err)
	}

	variants := configRepository.GetStreamOutputVariants()
	if len(variants) == 0 {
		// Ensure at least one variant directory exists.
		variants = make([]models.StreamOutputVariant, 1)
	}

	for index := range variants {
		variantDir := filepath.Join(config.HLSStoragePath, fmt.Sprintf("%d", index))
		if err := os.MkdirAll(variantDir, 0o750); err != nil {
			log.Errorln("unable to create variant directory:", err)
			continue
		}
		makeVariantIndexOffline(index, offlineInitPath, offlineSegmentPath)
	}

	// Write the master playlist that points to each variant's playlist.
	masterPlaylistPath := filepath.Join(config.HLSStoragePath, "stream.m3u8")
	masterFile, err := os.Create(masterPlaylistPath) //nolint:gosec
	if err != nil {
		log.Errorln("unable to create master playlist:", err)
	} else {
		defer masterFile.Close()
		_, _ = masterFile.WriteString("#EXTM3U\n")
		_, _ = masterFile.WriteString("#EXT-X-VERSION:7\n")
		_, _ = masterFile.WriteString("#EXT-X-INDEPENDENT-SEGMENTS\n")
		for index := range variants {
			_, _ = fmt.Fprintf(masterFile, "#EXT-X-STREAM-INF:BANDWIDTH=0\n")
			_, _ = fmt.Fprintf(masterFile, "%d/stream.m3u8\n", index)
		}
	}

	// Copy the logo to be the thumbnail
	logo := configRepository.GetLogoPath()
	dst := filepath.Join(config.TempDir, "thumbnail.jpg")
	if err = utils.Copy(filepath.Join("data", logo), dst); err != nil {
		log.Warnln(err)
	}

	// Delete the preview Gif
	_ = os.Remove(path.Join(config.DataDirectory, "preview.gif"))
}

func resetDirectories() {
	log.Trace("Resetting file directories to a clean slate.")

	// Wipe hls data directory
	utils.CleanupDirectory(config.HLSStoragePath)

	// Remove the previous thumbnail
	configRepository := configrepository.Get()
	logo := configRepository.GetLogoPath()
	if utils.DoesFileExists(logo) {
		err := utils.Copy(path.Join("data", logo), filepath.Join(config.DataDirectory, "thumbnail.jpg"))
		if err != nil {
			log.Warnln(err)
		}
	}
}
