package utils

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"os"

	"github.com/schollz/sqlite3dump"
	log "github.com/sirupsen/logrus"
)

func Backup(db *sql.DB, backupFile string) {
	log.Traceln("Backing up database")

	// Dump the entire database as plain text sql
	var b bytes.Buffer
	out := bufio.NewWriter(&b)
	if err := sqlite3dump.DumpDB(db, out); err != nil {
		handleError(err)
		return
	}

	// Create a new backup file
	f, err := os.OpenFile(backupFile, os.O_RDWR|os.O_CREATE, 0600)
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
