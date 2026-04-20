// This is a centralized place to connect to the database, and hold a reference to it.
// Other packages can share this reference.  This package would also be a place to add any kind of
// persistence-related convenience methods or migrations.

package data

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/persistence/migrations"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

var (
	_db        *sql.DB
	_datastore *Datastore
)

// GetDatabase will return the shared instance of the actual database.
func GetDatabase() *sql.DB {
	return _db
}

// GetStore will return the shared instance of the read/write datastore.
func GetStore() *Datastore {
	return _datastore
}

// SetupPersistence will open the datastore and make it available.
func SetupPersistence(file string) error {
	// Allow support for in-memory databases for tests.

	var db *sql.DB

	if file == ":memory:" {
		inMemoryDb, err := sql.Open("sqlite3", file)
		if err != nil {
			log.Fatal(err.Error())
		}
		db = inMemoryDb
	} else {
		// Create empty DB file if it doesn't exist.
		if !utils.DoesFileExists(file) {
			log.Traceln("Creating new database at", file)

			_, err := os.Create(file) //nolint:gosec
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		onDiskDb, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_cache_size=10000&cache=shared&_journal_mode=WAL", file))
		if err != nil {
			return err
		}
		db = onDiskDb
		db.SetMaxOpenConns(1)
	}
	_db = db

	// Some SQLite optimizations
	_, _ = db.Exec("pragma journal_mode = WAL")
	_, _ = db.Exec("pragma synchronous = normal")
	_, _ = db.Exec("pragma temp_store = memory")
	_, _ = db.Exec("pragma wal_checkpoint(full)")

	// Bring the schema up to date. The migrations package owns all table
	// creation and schema changes; existing pre-goose installs are caught up
	// automatically by its legacy-bridge step.
	if err := migrations.Run(db); err != nil {
		return fmt.Errorf("running database migrations: %w", err)
	}

	_datastore = &Datastore{}
	_datastore.Setup()

	dbBackupTicker := time.NewTicker(1 * time.Hour)
	go func() {
		backupFile := filepath.Join(config.BackupDirectory, "owncastdb.bak")
		for range dbBackupTicker.C {
			utils.Backup(_db, backupFile)
		}
	}()

	return nil
}
