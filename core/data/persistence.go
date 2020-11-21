package data

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Datastore struct {
	db    *sql.DB
	cache map[string][]byte
}

// Get will query the database for the key and return the entry
func (ds *Datastore) Get(key string) (ConfigEntry, error) {
	cachedValue, err := ds.GetCachedValue(key)
	if err == nil {
		return ConfigEntry{
			Key:   key,
			Value: cachedValue,
		}, nil
	}

	var resultKey string
	var resultValue []byte

	row := ds.db.QueryRow("SELECT key, value FROM datastore WHERE key = ? LIMIT 1", key)
	if err := row.Scan(&resultKey, &resultValue); err != nil {
		return ConfigEntry{}, err
	}

	result := ConfigEntry{
		Key:   resultKey,
		Value: resultValue,
	}

	return result, nil
}

// Save will save the ConfigEntry to the database.
func (ds *Datastore) Save(e ConfigEntry) error {
	var dataGob bytes.Buffer
	enc := gob.NewEncoder(&dataGob)
	if err := enc.Encode(e.Value); err != nil {
		return err
	}

	tx, err := ds.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO datastore(key, value) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Key, dataGob.Bytes())
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Fatalln(err)
	}

	ds.SetCachedValue(e.Key, dataGob.Bytes())

	return nil
}

func (ds *Datastore) Setup() {
	ds.cache = make(map[string][]byte)
	ds.db = GetDatabase()

	createTableSQL := `CREATE TABLE IF NOT EXISTS datastore (
		"key" string NOT NULL PRIMARY KEY,
		"value" BLOB,
		"timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
	);`

	stmt, err := ds.db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalln(err)
	}
}
