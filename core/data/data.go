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
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

const (
	schemaVersion = 3
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

			_, err := os.Create(file)
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

	createWebhooksTable()
	createUsersTable(db)

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS config (
		"key" string NOT NULL PRIMARY KEY,
		"value" TEXT
	);`); err != nil {
		return err
	}

	var version int
	err := db.QueryRow("SELECT value FROM config WHERE key='version'").
		Scan(&version)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}

		// fresh database: initialize it with the current schema version
		_, err := db.Exec("INSERT INTO config(key, value) VALUES(?, ?)", "version", schemaVersion)
		if err != nil {
			return err
		}
		version = schemaVersion
	}

	// is database from a newer Owncast version?
	if version > schemaVersion {
		return fmt.Errorf("incompatible database version %d (versions up to %d are supported)",
			version, schemaVersion)
	}

	// is database schema outdated?
	if version < schemaVersion {
		if err := migrateDatabase(db, version, schemaVersion); err != nil {
			return err
		}
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

func migrateDatabase(db *sql.DB, from, to int) error {
	log.Printf("Migrating database from version %d to %d", from, to)
	dbBackupFile := filepath.Join(config.BackupDirectory, fmt.Sprintf("owncast-v%d.bak", from))
	utils.Backup(db, dbBackupFile)
	for v := from; v < to; v++ {
		log.Tracef("Migration step from %d to %d\n", v, v+1)
		switch v {
		case 0:
			migrateToSchema1(db)
		case 1:
			migrateToSchema2(db)
		case 2:
			migrateToSchema3(db)
		default:
			log.Fatalln("missing database migration step")
		}
	}

	_, err := db.Exec("UPDATE config SET value = ? WHERE key = ?", to, "version")
	if err != nil {
		return err
	}

	return nil
}
