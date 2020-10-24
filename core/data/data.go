// This is a centralized place to connect to the database, and hold a reference to it.
// Other packages can share this reference.  This package would also be a place to add any kind of
// persistence-related convenience methods or migrations.

package data

import (
	"database/sql"
	"os"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

var _db *sql.DB

func GetDatabase() *sql.DB {
	return _db
}

func SetupPersistence() {
	file := config.Config.DatabaseFilePath

	// Create empty DB file if it doesn't exist.
	if !utils.DoesFileExists(file) {
		log.Traceln("Creating new database at", file)

		_, err := os.Create(file)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	sqliteDatabase, _ := sql.Open("sqlite3", file)
	_db = sqliteDatabase
}
