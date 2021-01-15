package main

import (
	"flag"
	"fmt"

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
	GitCommit = "unknown"
	// BuildVersion is the version.
	BuildVersion = config.CurrentBuildString
	// BuildType is the type of build.
	BuildType = "dev"
)

func main() {
	configureLogging()

	log.Infoln(getReleaseString())
	// Enable bundling of admin assets
	_ = pkger.Include("/admin")

	// configFile := flag.String("configFile", "config.yaml", "Config File full path. Defaults to current folder")
	dbFile := flag.String("database", "", "Path to the database file.")
	enableDebugOptions := flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")
	enableVerboseLogging := flag.Bool("enableVerboseLogging", false, "Enable additional logging.")
	restoreDatabaseFile := flag.String("restoreDatabase", "", "Restore an Owncast database backup")
	flag.Parse()

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

	// starts the core
	if err := core.Start(); err != nil {
		log.Fatalln("failed to start the core package", err)
	}

	if err := router.Start(); err != nil {
		log.Fatalln("failed to start/run the router", err)
	}
}

// getReleaseString gets the version string.
func getReleaseString() string {
	return fmt.Sprintf("Owncast v%s-%s (%s)", BuildVersion, BuildType, GitCommit)
}

func getVersionNumber() string {
	return BuildVersion
}

func configureLogging() {
	logging.Setup()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
