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

// TestRestorePreservesDestinationOnFailure covers failed restores and successful replacement.
func TestRestorePreservesDestinationOnFailure(t *testing.T) {
	for _, tc := range []struct {
		name      string
		sql       string
		existing  bool
		wantError bool
		locked    bool
	}{
		{"invalid SQL", "invalid SQL", true, true, false},
		{"partially executed SQL", "CREATE TABLE partial (value TEXT); invalid SQL", true, true, false},
		{"new destination", "invalid SQL", false, true, false},
		{"successful replacement", "CREATE TABLE restored (value TEXT); INSERT INTO restored VALUES ('restored');", true, false, false},
		{"locked destination", "CREATE TABLE restored (value TEXT);", true, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			databaseFile := filepath.Join(dir, "owncast.db")
			if tc.existing {
				db, err := sql.Open("sqlite3", databaseFile)
				if err != nil {
					t.Fatal(err)
				}
				_, execErr := db.Exec("CREATE TABLE original (value TEXT); INSERT INTO original VALUES ('keep me');")
				closeErr := db.Close()
				if execErr != nil || closeErr != nil {
					t.Fatalf("create original database: exec = %v, close = %v", execErr, closeErr)
				}
			}
			original, err := os.ReadFile(databaseFile)
			if err != nil && tc.existing {
				t.Fatal(err)
			}
			var originalInfo os.FileInfo
			if tc.existing {
				if err := os.Chmod(databaseFile, 0o640); err != nil {
					t.Fatal(err)
				}
				originalInfo, err = os.Stat(databaseFile)
				if err != nil {
					t.Fatal(err)
				}
			}
			if tc.locked {
				db, err := sql.Open("sqlite3", databaseFile)
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					_, _ = db.Exec("ROLLBACK")
					_ = db.Close()
				})
				if _, err := db.Exec("BEGIN EXCLUSIVE"); err != nil {
					t.Fatal(err)
				}
			}

			var compressed bytes.Buffer
			writer := gzip.NewWriter(&compressed)
			if _, err := writer.Write([]byte(tc.sql)); err != nil {
				t.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			backupFile := filepath.Join(dir, "backup.gz")
			if err := os.WriteFile(backupFile, compressed.Bytes(), 0o600); err != nil {
				t.Fatal(err)
			}

			if err := Restore(backupFile, databaseFile); (err != nil) != tc.wantError {
				t.Fatalf("Restore() error = %v, want error = %t", err, tc.wantError)
			}
			if tc.wantError {
				if tc.existing {
					current, err := os.ReadFile(databaseFile)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(current, original) {
						t.Error("failed restore changed the existing database")
					}
				} else if _, err := os.Stat(databaseFile); !os.IsNotExist(err) {
					t.Errorf("failed restore left a destination database: %v", err)
				}
			} else {
				restoredInfo, err := os.Stat(databaseFile)
				if err != nil {
					t.Fatal(err)
				}
				if !os.SameFile(originalInfo, restoredInfo) {
					t.Error("restore replaced the original file and its ownership")
				}
				if restoredInfo.Mode() != originalInfo.Mode() {
					t.Errorf("restored mode = %v, want %v", restoredInfo.Mode(), originalInfo.Mode())
				}
				db, err := sql.Open("sqlite3", databaseFile)
				if err != nil {
					t.Fatal(err)
				}
				defer db.Close()
				var value string
				if err := db.QueryRow("SELECT value FROM restored").Scan(&value); err != nil || value != "restored" {
					t.Fatalf("restored value = %q, error = %v", value, err)
				}
				var originalTables int
				if err := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE name = 'original'").Scan(&originalTables); err != nil || originalTables != 0 {
					t.Fatalf("original tables = %d, error = %v", originalTables, err)
				}
			}

			files, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, file := range files {
				if file.Name() != "backup.gz" && file.Name() != "owncast.db" {
					t.Errorf("restore left an unexpected file: %s", file.Name())
				}
			}
		})
	}
}
