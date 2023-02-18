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

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

const (
	schemaVersion = 7
)

func New(dbPath string) (*Service, error) {
	s := &Service{}

	if err := s.Init(dbPath); err != nil {
		return nil, fmt.Errorf("initializing the data service: %v", err)
	}

	s.Setup()

	return s, nil
}

type Service struct {
	DB    *sql.DB
	Store *Datastore
}

// GetDatabase will return the shared instance of the actual database.
func (s *Service) GetDatabase() *sql.DB {
	return s.DB
}

// GetStore will return the shared instance of the read/write Datastore.
func (s *Service) GetStore() *Datastore {
	return s.Store
}

// Init will open the Datastore and make it available.
func (s *Service) Init(file string) error {
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
	s.DB = db

	s.Store = &Datastore{}
	s.Setup()

	// Some SQLite optimizations
	_, _ = db.Exec("pragma journal_mode = WAL")
	_, _ = db.Exec("pragma synchronous = normal")
	_, _ = db.Exec("pragma temp_store = memory")
	_, _ = db.Exec("pragma wal_checkpoint(full)")

	s.createWebhooksTable()
	s.createUsersTable()
	s.createAccessTokenTable()

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
		if err := s.migrateDatabaseSchema(db, version, schemaVersion); err != nil {
			return err
		}
	}

	dbBackupTicker := time.NewTicker(1 * time.Hour)
	go func() {
		backupFile := filepath.Join(config.BackupDirectory, "owncastdb.bak")
		for range dbBackupTicker.C {
			utils.Backup(s.DB, backupFile)
		}
	}()

	return nil
}
