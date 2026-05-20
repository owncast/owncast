// Package datastore is the database root for the rest of the
// application. SetupPersistence opens the SQLite file, runs migrations,
// and returns a constructed *Datastore so main.go can inject it through
// every *repository.New and into the service Deps structs that need
// direct database access. No package-level handle is retained — the
// composition root owns the lifetime.
package datastore

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/persistence/migrations"
	"github.com/owncast/owncast/utils"
)

// SetupPersistence opens the datastore file (or an in-memory database
// when file == ":memory:"), runs migrations, starts the periodic backup
// goroutine, and returns the constructed *Datastore. backupDirectory is
// the directory the hourly backup goroutine writes owncastdb.bak into.
// The returned value is the sole handle to the database — main.go
// threads it through every consumer.
func SetupPersistence(file, backupDirectory string) (*Datastore, error) {
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
			return nil, err
		}
		db = onDiskDb
		db.SetMaxOpenConns(1)
	}

	// Some SQLite optimizations
	_, _ = db.Exec("pragma journal_mode = WAL")
	_, _ = db.Exec("pragma synchronous = normal")
	_, _ = db.Exec("pragma temp_store = memory")
	_, _ = db.Exec("pragma wal_checkpoint(full)")

	// Bring the schema up to date. The migrations package owns all table
	// creation and schema changes; existing pre-goose installs are caught up
	// automatically by its legacy-bridge step.
	if err := migrations.Run(db, backupDirectory); err != nil {
		return nil, fmt.Errorf("running database migrations: %w", err)
	}

	dataStore := &Datastore{}
	dataStore.Setup(db)

	dbBackupTicker := time.NewTicker(1 * time.Hour)
	go func() {
		backupFile := filepath.Join(backupDirectory, "owncastdb.bak")
		for range dbBackupTicker.C {
			utils.Backup(db, backupFile)
		}
	}()

	return dataStore, nil
}
