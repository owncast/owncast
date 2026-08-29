package datastore

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/owncast/owncast/utils"
)

func TestBackupDatabaseAtKeepsCurrentAndHistory(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE values_table (value TEXT); INSERT INTO values_table VALUES ('old');"); err != nil {
		t.Fatal(err)
	}

	backupDirectory := t.TempDir()
	firstBackupTime := time.Date(2026, 8, 17, 18, 0, 0, 0, time.UTC)
	backupDatabaseAt(db, backupDirectory, firstBackupTime)

	if _, err := db.Exec("UPDATE values_table SET value = 'new'"); err != nil {
		t.Fatal(err)
	}
	secondBackupTime := firstBackupTime.Add(time.Hour)
	backupDatabaseAt(db, backupDirectory, secondBackupTime)

	historyFile := filepath.Join(backupDirectory, "owncastdb-20260817T190000Z.bak")
	if _, err := os.Stat(historyFile); err != nil {
		t.Fatalf("history backup missing: %v", err)
	}

	historyDatabase := filepath.Join(t.TempDir(), "history.db")
	if err := utils.Restore(historyFile, historyDatabase); err != nil {
		t.Fatal(err)
	}
	assertBackupValue(t, historyDatabase, "old")

	currentDatabase := filepath.Join(t.TempDir(), "current.db")
	if err := utils.Restore(filepath.Join(backupDirectory, databaseBackupFileName), currentDatabase); err != nil {
		t.Fatal(err)
	}
	assertBackupValue(t, currentDatabase, "new")
}

func TestBackupDatabaseAtPrunesOldHistory(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE values_table (value TEXT)"); err != nil {
		t.Fatal(err)
	}

	backupDirectory := t.TempDir()
	firstBackupTime := time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC)
	for i := range 25 {
		if _, err := db.Exec("DELETE FROM values_table"); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec("INSERT INTO values_table VALUES (?)", fmt.Sprintf("backup-%d", i)); err != nil {
			t.Fatal(err)
		}
		backupDatabaseAt(db, backupDirectory, firstBackupTime.Add(time.Duration(i)*time.Hour))
	}

	entries, err := os.ReadDir(backupDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != maxDatabaseBackups {
		t.Fatalf("backup file count = %d, want %d", len(entries), maxDatabaseBackups)
	}
	if _, err := os.Stat(filepath.Join(backupDirectory, "owncastdb-20260817T010000Z.bak")); !os.IsNotExist(err) {
		t.Fatalf("oldest history backup still exists, stat error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(backupDirectory, "owncastdb-20260817T020000Z.bak")); err != nil {
		t.Fatalf("newest retained history backup missing: %v", err)
	}
}

func assertBackupValue(t *testing.T, databaseFile, want string) {
	t.Helper()
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var got string
	if err := db.QueryRow("SELECT value FROM values_table").Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("backup value = %q, want %q", got, want)
	}
}
