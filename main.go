package main

import (
	"flag"
	"fmt"

	"github.com/owncast/owncast/logging"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/router"
)

// the following are injected at build-time
var (
	//GitCommit is the commit which this version of owncast is running
	GitCommit = "unknown"
	//BuildVersion is the version
	BuildVersion = "0.0.0"
	//BuildType is the type of build
	BuildType = "localdev"
)

func main() {
	configureLogging()

	log.Infoln(getReleaseString())

	configFile := flag.String("configFile", "config.yaml", "Config File full path. Defaults to current folder")
	dbFile := flag.String("database", "", "Path to the database file.")
	enableDebugOptions := flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")
	enableVerboseLogging := flag.Bool("enableVerboseLogging", false, "Enable additional logging.")

	flag.Parse()

	if *enableDebugOptions {
		logrus.SetReportCaller(true)
	}

	if *enableVerboseLogging {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err := config.Load(*configFile, getReleaseString(), getVersionNumber()); err != nil {
		panic(err)
	}
	config.Config.EnableDebugFeatures = *enableDebugOptions

	if *dbFile != "" {
		config.Config.DatabaseFilePath = *dbFile
	} else if config.Config.DatabaseFilePath == "" {
		config.Config.DatabaseFilePath = config.Config.GetDataFilePath()
	}

	go metrics.Start()

	data.SetupPersistence()

	// starts the core
	if err := core.Start(); err != nil {
		log.Error("failed to start the core package")
		panic(err)
	}

	if err := router.Start(); err != nil {
		log.Error("failed to start/run the router")
		panic(err)
	}

}

//getReleaseString gets the version string
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
