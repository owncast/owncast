package cmd

import (
	"flag"
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

func (app *Application) parseFlags() {
	flag.Parse()
}
