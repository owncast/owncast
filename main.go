package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/markbates/pkger"
	"github.com/owncast/owncast/logging"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/router"
	"github.com/owncast/owncast/utils"
)

// the following are injected at build-time.
var (
	// GitCommit is the commit which this version of owncast is running.
	GitCommit = ""
	// BuildVersion is the version.
	BuildVersion = config.StaticVersionNumber
	// BuildPlatform is the type of build.
	BuildPlatform = "dev"
)

func main() {
	configureLogging()

	// Enable bundling of admin assets
	_ = pkger.Include("/admin")

	configFile := flag.String("configFile", "config.yaml", "Config file path to migrate to the new database")
	dbFile := flag.String("database", "", "Path to the database file.")
	enableDebugOptions := flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")
	enableVerboseLogging := flag.Bool("enableVerboseLogging", false, "Enable additional logging.")
	restoreDatabaseFile := flag.String("restoreDatabase", "", "Restore an Owncast database backup")
	newStreamKey := flag.String("streamkey", "", "Set your stream key/admin password")
	webServerPortOverride := flag.String("webserverport", "", "Force the web server to listen on a specific port")

	flag.Parse()

	config.ConfigFilePath = *configFile
	config.VersionNumber = BuildVersion
	if GitCommit != "" {
		config.GitCommit = GitCommit
	} else {
		config.GitCommit = time.Now().Format("20060102")
	}
	config.BuildPlatform = BuildPlatform

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

	if *enableDebugOptions {
		logrus.SetReportCaller(true)
	}

	if *enableVerboseLogging {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	config.EnableDebugFeatures = *enableDebugOptions

	if *dbFile != "" {
		config.DatabaseFilePath = *dbFile
	}

	go metrics.Start()

	err := data.SetupPersistence(config.DatabaseFilePath)
	if err != nil {
		log.Fatalln("failed to open database", err)
	}

	if *newStreamKey != "" {
		if err := data.SetStreamKey(*newStreamKey); err != nil {
			log.Errorln("Error setting your stream key.", err)
		} else {
			log.Infoln("Stream key changed to", *newStreamKey)
		}

		log.Exit(0)
	}

	// Set the web server port
	if *webServerPortOverride != "" {
		portNumber, err := strconv.Atoi(*webServerPortOverride)
		if err != nil {
			log.Warnln(err)
			return
		}

		config.WebServerPort = portNumber
	} else {
		config.WebServerPort = data.GetHTTPPortNumber()
	}

	// starts the core
	if err := core.Start(); err != nil {
		log.Fatalln("failed to start the core package", err)
	}

	if err := router.Start(); err != nil {
		log.Fatalln("failed to start/run the router", err)
	}

}

func configureLogging() {
	logging.Setup()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
