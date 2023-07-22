package cmd

import (
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

func (app *Application) handleRestoreBackup(restoreDatabaseFile *string) {
	// Allows a user to restore a specific database backup
	databaseFile := app.configservice.DatabaseFilePath
	if *dbFile != "" {
		databaseFile = *dbFile
	}

	if err := utils.Restore(*restoreDatabaseFile, databaseFile); err != nil {
		log.Fatalln(err)
	}

	log.Println("Database has been restored.  Restart Owncast.")
}
