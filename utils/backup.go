package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/sqlite3dump"
	log "github.com/sirupsen/logrus"
)

// Restore will attempt to restore the database using a specified backup file.
func Restore(backupFile string, databaseFile string) error {
	log.Printf("Restoring database backup %s to %s", backupFile, databaseFile)

	data, err := os.ReadFile(backupFile) // nolint
	if err != nil {
		return fmt.Errorf("unable to read backup file %s", err)
	}

	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("unable to read backup file %s", err)
	}
	defer gz.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, gz); err != nil { // nolint
		return fmt.Errorf("unable to read backup file %s", err)
	}

	defer gz.Close()

	rawSQL := b.String()

	if _, err := os.Create(databaseFile); err != nil {
		return errors.New("unable to write restored database")
	}

	// Create a new database by executing the raw SQL
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return err
	}
	if _, err := db.Exec(rawSQL); err != nil {
		return err
	}

	return nil
}

// Backup will backup the provided instance of the database to the specified file.
func Backup(db *sql.DB, backupFile string) {
	log.Traceln("Backing up database to", backupFile)

	backupDirectory := filepath.Dir(backupFile)

	if !DoesFileExists(backupDirectory) {
		err := os.MkdirAll(backupDirectory, 0700)
		if err != nil {
			log.Errorln("unable to create backup directory. check permissions and ownership.", backupDirectory, err)
			return
		}
	}

	// Dump the entire database as plain text sql
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	if err := sqlite3dump.DumpDB(db, out); err != nil {
		handleError(err)
		return
	}
	_ = out.Flush()

	// Create a new backup file
	f, err := os.OpenFile(backupFile, os.O_WRONLY|os.O_CREATE, 0600) // nolint
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
