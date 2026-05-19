package main

import (
	"context"
	"flag"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/logging"
	"github.com/owncast/owncast/persistence/configrepository"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/cache"
	"github.com/owncast/owncast/services/rtmp"
	"github.com/owncast/owncast/services/stream"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers"
	"github.com/owncast/owncast/webserver/handlers/admin"
	"github.com/owncast/owncast/webserver/handlers/auth/fediverse"
	"github.com/owncast/owncast/webserver/router"
)

var (
	dbFile                         = flag.String("database", "", "Path to the database file.")
	logDirectory                   = flag.String("logdir", "", "Directory where logs will be written to")
	backupDirectory                = flag.String("backupdir", "", "Directory where backups will be written to")
	enableDebugOptions             = flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")
	enableVerboseLogging           = flag.Bool("enableVerboseLogging", false, "Enable additional logging.")
	restoreDatabaseFile            = flag.String("restoreDatabase", "", "Restore an Owncast database backup")
	newAdminPassword               = flag.String("adminpassword", "", "Set your admin password")
	newStreamKey                   = flag.String("streamkey", "", "Set a temporary stream key for this session")
	webServerPortOverride          = flag.String("webserverport", "", "Force the web server to listen on a specific port")
	webServerIPOverride            = flag.String("webserverip", "", "Force web server to listen on this IP address")
	rtmpPortOverride               = flag.Int("rtmpport", 0, "Set listen port for the RTMP server")
	followerValidationIntervalSecs = flag.Int("followervalidationinterval", 0, "Set follower validation interval in seconds")
)

// nolint:cyclop
func main() {
	flag.Parse()

	if *logDirectory != "" {
		config.LogDirectory = *logDirectory
	}

	if *backupDirectory != "" {
		config.BackupDirectory = *backupDirectory
	}

	// Create the data directory if needed
	if !utils.DoesFileExists("data") {
		if err := os.Mkdir("./data", 0o700); err != nil {
			log.Fatalln("Cannot create data directory", err)
		}
	}

	// Migrate old (pre 0.1.0) emoji to new location if they exist.
	utils.MigrateCustomEmojiLocations()

	// Otherwise save the default emoji to the data directory.
	if err := data.SetupEmojiDirectory(); err != nil {
		log.Fatalln("Cannot set up emoji directory", err)
	}

	// Recreate the temp dir
	if utils.DoesFileExists(config.TempDir) {
		err := os.RemoveAll(config.TempDir)
		if err != nil {
			log.Fatalln("Unable to remove temp dir! Check permissions.", config.TempDir, err)
		}
	}
	if err := os.Mkdir(config.TempDir, 0o700); err != nil {
		log.Fatalln("Unable to create temp dir!", err)
	}

	configureLogging(*enableDebugOptions, *enableVerboseLogging)
	log.Infoln(config.GetReleaseString())

	// Allows a user to restore a specific database backup
	if *restoreDatabaseFile != "" {
		databaseFile := config.DatabaseFilePath
		if *dbFile != "" {
			databaseFile = *dbFile
		}

		if err := utils.Restore(*restoreDatabaseFile, databaseFile); err != nil {
			log.Fatalln(err)
		}

		log.Println("Database has been restored.  Restart Owncast.")
		log.Exit(0)
	}

	config.EnableDebugFeatures = *enableDebugOptions

	if *dbFile != "" {
		config.DatabaseFilePath = *dbFile
	}

	if err := data.SetupPersistence(config.DatabaseFilePath); err != nil {
		log.Fatalln("failed to open database", err)
	}

	handleCommandLineFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Composition root: construct services here and inject them into the
	// components that consume them. As more packages migrate off package-
	// level singletons into services/<domain>/, their constructors join
	// this block.
	cacheContainer := cache.New()
	defer cacheContainer.Stop()

	rtmpSvc := rtmp.New(rtmp.Deps{})

	apSvc := activitypub.New(activitypub.Deps{Datastore: data.GetDatastore()})
	apSvc.Start()

	streamSvc := stream.New(stream.Deps{
		Rtmp:        rtmpSvc,
		Activitypub: apSvc,
	})
	if err := streamSvc.Start(ctx); err != nil {
		log.Fatalln("failed to start the stream service", err)
	}
	defer streamSvc.Stop(ctx)

	go metrics.Start(streamSvc)

	adminHandlers := admin.New(admin.Deps{
		Stream:      streamSvc,
		Rtmp:        rtmpSvc,
		Activitypub: apSvc,
	})

	fediverseHandler := fediverse.New(fediverse.Deps{
		Activitypub: apSvc,
	})

	h := handlers.NewHandlers(handlers.Deps{
		Cache:       cacheContainer,
		Stream:      streamSvc,
		Admin:       adminHandlers,
		Activitypub: apSvc,
		Fediverse:   fediverseHandler,
	})

	if err := router.Start(*enableVerboseLogging, h, apSvc.Controllers()); err != nil {
		log.Fatalln("failed to start/run the router", err)
	}
}

func handleCommandLineFlags() {
	configRepository := configrepository.Get()

	if *newAdminPassword != "" {
		if err := configRepository.SetAdminPassword(*newAdminPassword); err != nil {
			log.Errorln("Error setting your admin password.", err)
			log.Exit(1)
		} else {
			log.Infoln("Admin password changed")
		}
	}

	if *newStreamKey != "" {
		log.Println("Temporary stream key is set for this session.")
		config.TemporaryStreamKey = *newStreamKey
	}

	// Set the web server port
	if *webServerPortOverride != "" {
		portNumber, err := strconv.Atoi(*webServerPortOverride)
		if err != nil {
			log.Warnln(err)
			return
		}

		log.Println("Saving new web server port number to", portNumber)
		if err := configRepository.SetHTTPPortNumber(float64(portNumber)); err != nil {
			log.Errorln(err)
		}
	}
	config.WebServerPort = configRepository.GetHTTPPortNumber()

	// Set the web server ip
	if *webServerIPOverride != "" {
		log.Println("Saving new web server listen IP address to", *webServerIPOverride)
		if err := configRepository.SetHTTPListenAddress(*webServerIPOverride); err != nil {
			log.Errorln(err)
		}
	}
	config.WebServerIP = configRepository.GetHTTPListenAddress()

	// Set the rtmp server port
	if *rtmpPortOverride > 0 {
		log.Println("Saving new RTMP server port number to", *rtmpPortOverride)
		if err := configRepository.SetRTMPPortNumber(float64(*rtmpPortOverride)); err != nil {
			log.Errorln(err)
		}
	}

	// Set the follower validation interval
	if *followerValidationIntervalSecs > 0 {
		config.FollowerValidationInterval = time.Duration(*followerValidationIntervalSecs) * time.Second
		log.Printf("Follower validation interval set to %v", config.FollowerValidationInterval)
	}
}

func configureLogging(enableDebugFeatures bool, enableVerboseLogging bool) {
	logging.Setup(enableDebugFeatures, enableVerboseLogging)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
