package cmd

import (
	"strconv"

	"github.com/owncast/owncast/storage/configrepository"
	log "github.com/sirupsen/logrus"
)

func (app *Application) setSessionConfig() {
	// Stream key
	if *newStreamKey != "" {
		log.Println("Temporary stream key is set for this session.")
		app.configservice.TemporaryStreamKey = *newStreamKey
	}

	app.configservice.EnableDebugFeatures = *enableDebugOptions

	if *dbFile != "" {
		app.configservice.DatabaseFilePath = *dbFile
	}

	if *logDirectory != "" {
		app.configservice.LogDirectory = *logDirectory
	}
}

func (app *Application) saveUpdatedConfig() {
	configRepository := configrepository.Get()

	if *newAdminPassword != "" {
		if err := configRepository.SetAdminPassword(*newAdminPassword); err != nil {
			log.Errorln("Error setting your admin password.", err)
			log.Exit(1)
		} else {
			log.Infoln("Admin password changed")
		}
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
	app.configservice.WebServerPort = configRepository.GetHTTPPortNumber()

	// Set the web server ip
	if *webServerIPOverride != "" {
		log.Println("Saving new web server listen IP address to", *webServerIPOverride)
		if err := configRepository.SetHTTPListenAddress(*webServerIPOverride); err != nil {
			log.Errorln(err)
		}
	}
	app.configservice.WebServerIP = configRepository.GetHTTPListenAddress()

	// Set the rtmp server port
	if *rtmpPortOverride > 0 {
		log.Println("Saving new RTMP server port number to", *rtmpPortOverride)
		if err := configRepository.SetRTMPPortNumber(float64(*rtmpPortOverride)); err != nil {
			log.Errorln(err)
		}
	}
}
