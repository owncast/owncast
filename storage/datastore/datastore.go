package datastore

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	// sqlite requires a blank import.
	_ "github.com/mattn/go-sqlite3"
	"github.com/owncast/owncast/db"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

const (
	schemaVersion = 7
)

// Datastore is the global key/value store for configuration values.
type Datastore struct {
	DB     *sql.DB
	cache  map[string][]byte
	DbLock *sync.Mutex
}

var temporaryGlobalDatastoreInstance *Datastore

// NewDatastore creates a new datastore.
func NewDatastore(file string) (*Datastore, error) {
	r := &Datastore{
		cache:  make(map[string][]byte),
		DbLock: &sync.Mutex{},
	}

	if err := r.InitializeDatabase(file); err != nil {
		return nil, err
	}
	return r, nil
}

// GetDatastore returns the shared instance of the owncast datastore.
func GetDatastore() *Datastore {
	if temporaryGlobalDatastoreInstance == nil {
		c := config.GetConfig()
		i, err := NewDatastore(c.DatabaseFilePath)
		if err != nil {
			log.Fatal(err)
		}
		temporaryGlobalDatastoreInstance = i
	}
	return temporaryGlobalDatastoreInstance
}

// GetQueries will return the shared instance of the SQL query generator.
func (ds *Datastore) GetQueries() *db.Queries {
	return db.New(ds.DB)
}

// Get will query the database for the key and return the entry.
func (ds *Datastore) Get(key string) (models.ConfigEntry, error) {
	cachedValue, err := ds.GetCachedValue(key)
	if err == nil {
		return models.ConfigEntry{
			Key:   key,
			Value: cachedValue,
		}, nil
	}

	var resultKey string
	var resultValue []byte

	row := ds.DB.QueryRow("SELECT key, value FROM datastore WHERE key = ? LIMIT 1", key)
	if err := row.Scan(&resultKey, &resultValue); err != nil {
		return models.ConfigEntry{}, err
	}

	result := models.ConfigEntry{
		Key:   resultKey,
		Value: resultValue,
	}
	ds.SetCachedValue(resultKey, resultValue)

	return result, nil
}

// Save will save the models.ConfigEntry to the database.
func (ds *Datastore) Save(e models.ConfigEntry) error {
	ds.DbLock.Lock()
	defer ds.DbLock.Unlock()

	var dataGob bytes.Buffer
	enc := gob.NewEncoder(&dataGob)
	if err := enc.Encode(e.Value); err != nil {
		return err
	}

	tx, err := ds.DB.Begin()
	if err != nil {
		return err
	}
	var stmt *sql.Stmt
	stmt, err = tx.Prepare("INSERT INTO datastore (key, value) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(e.Key, dataGob.Bytes())

	if err != nil {
		return err
	}
	defer stmt.Close()

	if err = tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	ds.SetCachedValue(e.Key, dataGob.Bytes())

	return nil
}

// MustExec will execute a SQL statement on a provided database instance.
func (ds *Datastore) MustExec(s string) {
	stmt, err := ds.DB.Prepare(s)
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Warnln(err)
	}
}

// InitializeDatabase will open the datastore and make it available.
func (ds *Datastore) InitializeDatabase(file string) error {
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

	ds.DB = db

	// Some SQLite optimizations
	_, _ = db.Exec("pragma journal_mode = WAL")
	_, _ = db.Exec("pragma synchronous = normal")
	_, _ = db.Exec("pragma temp_store = memory")
	_, _ = db.Exec("pragma wal_checkpoint(full)")

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS config (
		"key" string NOT NULL PRIMARY KEY,
		"value" TEXT
	);`); err != nil {
		return err
	}

	ds.createTables()

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
		if err := migrateDatabaseSchema(db, version, schemaVersion); err != nil {
			return err
		}
	}

	dbBackupTicker := time.NewTicker(1 * time.Hour)
	go func() {
		c := config.GetConfig()
		backupFile := filepath.Join(c.BackupDirectory, "owncastdb.bak")
		for range dbBackupTicker.C {
			utils.Backup(db, backupFile)
		}
	}()

	return nil
}
