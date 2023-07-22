package data

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"sync"

	// sqlite requires a blank import.
	_ "github.com/mattn/go-sqlite3"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/storage/sqlstorage"
	log "github.com/sirupsen/logrus"
)

const (
	schemaVersion = 7
)

// Store is the global key/value store for configuration values.
type Store struct {
	DB     *sql.DB
	cache  map[string][]byte
	DbLock *sync.Mutex
}

// NewStore creates a new datastore.
func NewStore(file string) (*Store, error) {
	s := &Store{
		cache:  make(map[string][]byte),
		DbLock: &sync.Mutex{},
	}

	db, err := sqlstorage.InitializeDatabase(file, schemaVersion)
	if err != nil {
		return nil, err
	}
	s.DB = db
	s.warmCache()
	temporaryGlobalDatastoreInstance = s
	return s, nil
}

var temporaryGlobalDatastoreInstance *Store

// GetDatastore returns the shared instance of the owncast datastore.
func GetDatastore() *Store {
	if temporaryGlobalDatastoreInstance == nil {
		c := config.Get()
		i, err := NewStore(c.DatabaseFilePath)
		if err != nil {
			log.Fatal(err)
		}
		temporaryGlobalDatastoreInstance = i
	}
	return temporaryGlobalDatastoreInstance
}

// GetQueries will return the shared instance of the SQL query generator.
func (ds *Store) GetQueries() *sqlstorage.Queries {
	return sqlstorage.New(ds.DB)
}

// Get will query the database for the key and return the entry.
func (ds *Store) Get(key string) (models.ConfigEntry, error) {
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
func (ds *Store) Save(e models.ConfigEntry) error {
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
