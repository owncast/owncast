package main

import (
	"flag"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/app"
	"github.com/owncast/owncast/logging"
	"github.com/owncast/owncast/utils"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
)

var (
	dbFile                = flag.String("database", "", "Path to the database file.")
	logDirectory          = flag.String("logdir", "", "Directory where logs will be written to")
	backupDirectory       = flag.String("backupdir", "", "Directory where backups will be written to")
	enableDebugOptions    = flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")
	enableVerboseLogging  = flag.Bool("enableVerboseLogging", false, "Enable additional logging.")
	restoreDatabaseFile   = flag.String("restoreDatabase", "", "Restore an Owncast database backup")
	newAdminPassword      = flag.String("adminpassword", "", "Set your admin password")
	newStreamKey          = flag.String("streamkey", "", "Set a temporary stream key for this session")
	webServerPortOverride = flag.String("webserverport", "", "Force the web server to listen on a specific port")
	webServerIPOverride   = flag.String("webserverip", "", "Force web server to listen on this IP address")
	rtmpPortOverride      = flag.Int("rtmpport", 0, "Set listen port for the RTMP server")
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

	instance, err := app.New("data")
	if err != nil {
		log.Fatalf("initializing app instance: %v", err)
	}

	handleCommandLineFlags(instance.Data)

	if err = instance.Serve(); err != nil {
		log.Fatalf("serving app: %v", err)
	}
}

func handleCommandLineFlags(d *data.Service) {
	if *newAdminPassword != "" {
		if err := d.SetAdminPassword(*newAdminPassword); err != nil {
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
		if err := d.SetHTTPPortNumber(float64(portNumber)); err != nil {
			log.Errorln(err)
		}
	}
	config.WebServerPort = d.GetHTTPPortNumber()

	// Set the web server ip
	if *webServerIPOverride != "" {
		log.Println("Saving new web server listen IP address to", *webServerIPOverride)
		if err := d.SetHTTPListenAddress(*webServerIPOverride); err != nil {
			log.Errorln(err)
		}
	}
	config.WebServerIP = d.GetHTTPListenAddress()

	// Set the rtmp server port
	if *rtmpPortOverride > 0 {
		log.Println("Saving new RTMP server port number to", *rtmpPortOverride)
		if err := d.SetRTMPPortNumber(float64(*rtmpPortOverride)); err != nil {
			log.Errorln(err)
		}
	}
}

func configureLogging(enableDebugFeatures bool, enableVerboseLogging bool) {
	logging.Setup(enableDebugFeatures, enableVerboseLogging)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
