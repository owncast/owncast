package utils

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRestoreAcceptsTrailingDataFromLegacyBackup(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE restored (value TEXT); INSERT INTO restored VALUES ('ok');"); err != nil {
		t.Fatal(err)
	}

	backupFile := filepath.Join(t.TempDir(), "owncastdb.bak")
	Backup(db, backupFile)

	compressed, err := os.ReadFile(backupFile)
	if err != nil {
		t.Fatal(err)
	}
	compressed = append(compressed, []byte("trailing garbage from an older backup")...)
	if err := os.WriteFile(backupFile, compressed, 0o600); err != nil {
		t.Fatal(err)
	}

	databaseFile := filepath.Join(t.TempDir(), "owncast.db")
	if err := Restore(backupFile, databaseFile); err != nil {
		t.Fatalf("Restore() error = %v", err)
	}

	restoredDB, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		t.Fatal(err)
	}
	defer restoredDB.Close()

	var value string
	if err := restoredDB.QueryRow("SELECT value FROM restored").Scan(&value); err != nil {
		t.Fatal(err)
	}
	if value != "ok" {
		t.Fatalf("restored value = %q, want %q", value, "ok")
	}
}

func TestBackupTruncatesExistingFile(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE restored (value TEXT)"); err != nil {
		t.Fatal(err)
	}
	for i := range 512 {
		if _, err := db.Exec("INSERT INTO restored VALUES (?)", fmt.Sprintf("row-%08d-%08d", i, i*7919)); err != nil {
			t.Fatal(err)
		}
	}

	backupFile := filepath.Join(t.TempDir(), "owncastdb.bak")
	Backup(db, backupFile)
	if _, err := db.Exec("DELETE FROM restored"); err != nil {
		t.Fatal(err)
	}
	Backup(db, backupFile)
	secondBackup, err := os.ReadFile(backupFile)
	if err != nil {
		t.Fatal(err)
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(secondBackup))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(io.Discard, gzipReader); err != nil {
		t.Fatalf("new backup contains trailing data: %v", err)
	}
	if err := gzipReader.Close(); err != nil {
		t.Fatal(err)
	}

	databaseFile := filepath.Join(t.TempDir(), "restored.db")
	if err := Restore(backupFile, databaseFile); err != nil {
		t.Fatalf("round-trip Restore() error = %v", err)
	}

	restoredDB, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		t.Fatal(err)
	}
	defer restoredDB.Close()

	var count int
	if err := restoredDB.QueryRow("SELECT COUNT(*) FROM restored").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("round-trip row count = %d, want 0", count)
	}
}
