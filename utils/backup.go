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

	"github.com/mattn/go-sqlite3"
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
	// Older Owncast backups can contain trailing bytes from an untruncated overwrite.
	gz.Multistream(false)
	defer gz.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, gz); err != nil { // nolint
		return fmt.Errorf("unable to read backup file %s", err)
	}

	defer gz.Close()

	rawSQL := b.String()

	// Validate the restore without modifying the destination.
	f, err := os.CreateTemp(filepath.Dir(databaseFile), ".owncast-restore-*")
	if err != nil {
		return errors.New("unable to write restored database")
	}
	defer os.Remove(f.Name())
	if err := f.Close(); err != nil {
		return errors.New("unable to write restored database")
	}

	// Create a new database by executing the raw SQL
	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(rawSQL); err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}

	return copyRestoredDatabase(f.Name(), databaseFile)
}

// copyRestoredDatabase uses a SQLite write transaction to preserve the destination
// file's ownership and permissions and roll back an incomplete copy.
func copyRestoredDatabase(source, destination string) error {
	f, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0o600) //nolint:gosec
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	driver := sqlite3.SQLiteDriver{}
	sourceConn, err := driver.Open(source)
	if err != nil {
		return err
	}
	defer sourceConn.Close()
	destinationConn, err := driver.Open(destination)
	if err != nil {
		return err
	}
	defer destinationConn.Close()

	backup, err := destinationConn.(*sqlite3.SQLiteConn).Backup("main", sourceConn.(*sqlite3.SQLiteConn), "main")
	if err != nil {
		return err
	}
	done, stepErr := backup.Step(-1)
	finishErr := backup.Finish()
	if stepErr != nil {
		return stepErr
	}
	if !done {
		return sqlite3.Error{Code: sqlite3.ErrBusy}
	}
	return finishErr
}

// Backup will backup the provided instance of the database to the specified file.
func Backup(db *sql.DB, backupFile string) {
	log.Traceln("Backing up database to", backupFile)

	backupDirectory := filepath.Dir(backupFile)

	if !DoesFileExists(backupDirectory) {
		err := os.MkdirAll(backupDirectory, 0o700)
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
	f, err := os.OpenFile(backupFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600) // nolint
	if err != nil {
		handleError(err)
		return
	}
	defer f.Close()

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
