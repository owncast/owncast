package main

import (
	"flag"
	"fmt"

	logger "github.com/gabek/owncast/log"
	"github.com/gabek/owncast/termui"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/router"
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

var _logger = logger.Logger{}

func main() {
	log.Infoln(getVersion())

	configFile := flag.String("configFile", "config.yaml", "Config File full path. Defaults to current folder")
	chatDbFile := flag.String("chatDatabase", "", "Path to the chat database file.")
	enableDebugOptions := flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")
	enableVerboseLogging := flag.Bool("enableVerboseLogging", false, "Enable additional logging.")
	enableTerminalUI := flag.Bool("enableTerminalUI", false, "Enable terminal GUI.")

	flag.Parse()

	if err := config.Load(*configFile, getVersion()); err != nil {
		panic(err)
	}

	config.Config.EnableDebugFeatures = *enableDebugOptions
	config.Config.EnableTerminalUI = *enableTerminalUI
	_logger.Setup()

	if *enableDebugOptions {
		logrus.SetReportCaller(true)
	}

	if *enableVerboseLogging {
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if *chatDbFile != "" {
		config.Config.ChatDatabaseFilePath = *chatDbFile
	} else if config.Config.ChatDatabaseFilePath == "" {
		config.Config.ChatDatabaseFilePath = "chat.db"
	}

	if config.Config.EnableTerminalUI {
		go func() {
			err := termui.Setup(getVersion())
			if err != nil {
				log.Infoln("Disabling terminal UI.")
				config.Config.EnableTerminalUI = false
			}
		}()
	}

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

//getVersion gets the version string
func getVersion() string {
	return fmt.Sprintf("Owncast v%s-%s (%s)", BuildVersion, BuildType, GitCommit)
}
