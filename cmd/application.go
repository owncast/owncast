package cmd

import (
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/services/metrics"
	"github.com/owncast/owncast/storage/configrepository"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	configservice    *config.Config
	metricsservice   *metrics.Metrics
	configRepository *configrepository.SqlConfigRepository

	maximumConcurrentConnectionLimit int64
}

/*
The order of this setup matters.
- Parse flags
- Set the session runtime values
- Use the session values to configure data persistence
*/
func (app *Application) Start() {
	app.configservice = config.Get()

	app.parseFlags()
	app.configureLogging(*enableDebugOptions, *enableVerboseLogging, app.configservice.LogDirectory)
	app.showStartupMessage()

	app.setSessionConfig()
	app.createDirectories()

	app.maximumConcurrentConnectionLimit = getMaximumConcurrentConnectionLimit()
	setSystemConcurrentConnectionLimit(app.maximumConcurrentConnectionLimit)

	// If we're restoring a backup, do that and exit.
	if *restoreDatabaseFile != "" {
		app.handleRestoreBackup(restoreDatabaseFile)
		log.Exit(0)
	}

	if *backupDirectory != "" {
		app.configservice.BackupDirectory = *backupDirectory
	}

	app.startServices()
}
