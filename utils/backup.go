package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
  "path/filepath"
  
	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/sqlite3dump"
	log "github.com/sirupsen/logrus"
)

func Restore(backupFile string, databaseFile string) error {
	log.Printf("Restoring database backup %s to %s", backupFile, databaseFile)

	data, err := ioutil.ReadFile(backupFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to read backup file %s", err))
	}

	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to read backup file %s", err))
	}
	defer gz.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, gz); err != nil {
		return errors.New(fmt.Sprintf("Unable to read backup file %s", err))
	}

	defer gz.Close()

	rawSql := b.String()

	if _, err := os.Create(databaseFile); err != nil {
		return errors.New("Unable to write restored database")
	}

	// Create a new database by executing the raw SQL
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return err
	}
	_, err = db.Exec(rawSql)
	if err != nil {
		return err
	}

	return nil
}

func Backup(db *sql.DB, backupFile string) {
	log.Traceln("Backing up database to", backupFile)

  BackupDirectory := filepath.Dir(backupFile)
  _, err := os.Stat(BackupDirectory)
  if os.IsNotExist(err) {
    err = os.MkdirAll(BackupDirectory, 0777)
    if err != nil {
      log.Fatalln(err)
    }
  }
  
	// Dump the entire database as plain text sql
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	if err := sqlite3dump.DumpDB(db, out); err != nil {
		handleError(err)
		return
	}
	out.Flush()

	// Create a new backup file
	f, err := os.OpenFile(backupFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		handleError(err)
		return
	}

	// Create a gzip compression writer
	w, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		handleError(err)
		return
	}

	// Write compressed data
	if _, err := w.Write(b.Bytes()); err != nil {
		handleError(err)
		return
	}

	if err := w.Close(); err != nil {
		handleError(err)
		return
	}
}

func handleError(err error) {
	log.Errorln("unable to backup owncast database to file", err)
}
