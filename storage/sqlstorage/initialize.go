package sqlstorage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// InitializeDatabase will open the datastore and make it available.
func InitializeDatabase(file string, schemaVersion int) (*sql.DB, error) {
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

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS config (
		"key" string NOT NULL PRIMARY KEY,
		"value" TEXT
	);`); err != nil {
		return nil, err
	}

	CreateAllTables(db)

	var version int
	err := db.QueryRow("SELECT value FROM config WHERE key='version'").
		Scan(&version)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}

		// fresh database: initialize it with the current schema version
		_, err := db.Exec("INSERT INTO config(key, value) VALUES(?, ?)", "version", schemaVersion)
		if err != nil {
			return nil, err
		}
		version = schemaVersion
	}

	// is database from a newer Owncast version?
	if version > schemaVersion {
		return nil, fmt.Errorf("incompatible database version %d (versions up to %d are supported)",
			version, schemaVersion)
	}

	// is database schema outdated?
	if version < schemaVersion {
		migrations := NewSqlMigrations(db)

		if err := migrations.MigrateDatabaseSchema(db, version, schemaVersion); err != nil {
			return nil, err
		}
	}

	dbBackupTicker := time.NewTicker(1 * time.Hour)
	go func() {
		c := config.Get()
		backupFile := filepath.Join(c.BackupDirectory, "owncastdb.bak")
		for range dbBackupTicker.C {
			utils.Backup(db, backupFile)
		}
	}()

	return db, nil
}
