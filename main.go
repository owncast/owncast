package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/markbates/pkger"
	"github.com/owncast/owncast/logging"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/router"
	"github.com/owncast/owncast/utils"
)

func main() {

	// Enable bundling of admin assets
	_ = pkger.Include("/admin")

	dbFile := flag.String("database", "", "Path to the database file.")
	logDirectory := flag.String("logdir", "", "Directory where logs will be written to")
	backupDirectory := flag.String("backupdir", "", "Directory where backups will be written to")
	enableDebugOptions := flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")
	enableVerboseLogging := flag.Bool("enableVerboseLogging", false, "Enable additional logging.")
	restoreDatabaseFile := flag.String("restoreDatabase", "", "Restore an Owncast database backup")
	newStreamKey := flag.String("streamkey", "", "Set your stream key/admin password")
	webServerPortOverride := flag.String("webserverport", "", "Force the web server to listen on a specific port")
	webServerIPOverride := flag.String("webserverip", "", "Force web server to listen on this IP address")
	rtmpPortOverride := flag.Int("rtmpport", 0, "Set listen port for the RTMP server")

	flag.Parse()

	if *logDirectory != "" {
		config.LogDirectory = *logDirectory
	}

	configureLogging(*enableDebugOptions, *enableVerboseLogging)
	log.Infoln(config.GetReleaseString())

	if *backupDirectory != "" {
		config.BackupDirectory = *backupDirectory
	}

	// Create the data directory if needed
	if !utils.DoesFileExists("data") {
		os.Mkdir("./data", 0700)
	}

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

		log.Println("Saving new web server port number to", portNumber)
		data.SetHTTPPortNumber(float64(portNumber))
	}
	config.WebServerPort = data.GetHTTPPortNumber()

	// Set the web server ip
	if *webServerIPOverride != "" {
		log.Println("Saving new web server listen IP address to", *webServerIPOverride)
		data.SetHTTPListenAddress(string(*webServerIPOverride))
	}
	config.WebServerIP = data.GetHTTPListenAddress()

	// Set the rtmp server port
	if *rtmpPortOverride > 0 {
		log.Println("Saving new RTMP server port number to", *rtmpPortOverride)
		data.SetRTMPPortNumber(float64(*rtmpPortOverride))
	}

	// starts the core
	if err := core.Start(); err != nil {
		log.Fatalln("failed to start the core package", err)
	}

	if err := router.Start(); err != nil {
		log.Fatalln("failed to start/run the router", err)
	}
}

func configureLogging(enableDebugFeatures bool, enableVerboseLogging bool) {
	logging.Setup(enableDebugFeatures, enableVerboseLogging)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
