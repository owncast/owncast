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
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/persistence/migrations"
	"github.com/owncast/owncast/utils"
)

// SetupPersistence opens the datastore file (or an in-memory database
// when file == ":memory:"), runs migrations, starts the periodic backup
// goroutine, and returns the constructed *Datastore. backupDirectory is
// the directory where the current and retained hourly backups are written.
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

			createdFile, err := os.Create(file) //nolint:gosec
			if err != nil {
				log.Fatal(err.Error())
			}
			if err := createdFile.Close(); err != nil {
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

	dbBackupTicker := time.NewTicker(time.Hour)
	go func() {
		for backupTime := range dbBackupTicker.C {
			backupDatabaseAt(db, backupDirectory, backupTime)
		}
	}()

	return dataStore, nil
}

const (
	databaseBackupFileName      = "owncastdb.bak"
	databaseBackupHistoryPrefix = "owncastdb-"
	databaseBackupHistorySuffix = ".bak"
	maxDatabaseBackups          = 24
)

func backupDatabaseAt(db *sql.DB, backupDirectory string, backupTime time.Time) {
	if err := os.MkdirAll(backupDirectory, 0o700); err != nil {
		log.Errorln("unable to create backup directory", backupDirectory, err)
		return
	}

	backupFile := filepath.Join(backupDirectory, databaseBackupFileName)
	if utils.DoesFileExists(backupFile) {
		historyFile := filepath.Join(
			backupDirectory,
			fmt.Sprintf("%s%s%s", databaseBackupHistoryPrefix, backupTime.UTC().Format("20060102T150405Z"), databaseBackupHistorySuffix),
		)
		if err := utils.Copy(backupFile, historyFile); err != nil {
			log.Errorln("unable to preserve previous database backup", err)
			return
		}
	}

	utils.Backup(db, backupFile)

	entries, err := os.ReadDir(backupDirectory)
	if err != nil {
		log.Errorln("unable to read database backup directory", err)
		return
	}

	var historyFiles []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), databaseBackupHistoryPrefix) || !strings.HasSuffix(entry.Name(), databaseBackupHistorySuffix) {
			continue
		}
		historyFiles = append(historyFiles, entry.Name())
	}
	sort.Strings(historyFiles)

	for len(historyFiles) > maxDatabaseBackups-1 {
		oldest := filepath.Join(backupDirectory, historyFiles[0])
		if err := os.Remove(oldest); err != nil {
			log.Errorln("unable to remove old database backup", oldest, err)
			return
		}
		historyFiles = historyFiles[1:]
	}
}
