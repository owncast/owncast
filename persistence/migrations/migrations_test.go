package migrations

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func tableExists(t *testing.T, db *sql.DB, name string) bool {
	t.Helper()
	var n string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, name).Scan(&n)
	return err == nil
}

func gooseVersion(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	var v int64
	err := db.QueryRow(`SELECT COALESCE(MAX(version_id), 0) FROM goose_db_version WHERE is_applied=1`).Scan(&v)
	if err != nil {
		return 0
	}
	return v
}

func mustExec(t *testing.T, db *sql.DB, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("exec %q: %v", query, err)
	}
}

// TestRun_FreshDatabase verifies that Run on an empty database creates all
// expected tables and records the goose baseline migration.
func TestRun_FreshDatabase(t *testing.T) {
	db := openTestDB(t)

	if err := Run(db, t.TempDir()); err != nil {
		t.Fatalf("Run on fresh DB: %v", err)
	}

	expectedTables := []string{
		"datastore", "webhooks", "users", "user_access_tokens",
		"ap_followers", "ap_outbox", "ap_accepted_activities",
		"notifications", "messages", "auth", "ip_bans",
		"federated_servers",
		"goose_db_version",
	}
	for _, name := range expectedTables {
		if !tableExists(t, db, name) {
			t.Errorf("expected table %q to exist after fresh migration", name)
		}
	}

	// Fresh installs should not have the legacy config table.
	if tableExists(t, db, "config") {
		t.Error("fresh install should not have legacy config table")
	}

	if v := gooseVersion(t, db); v != 4 {
		t.Errorf("goose version = %d, want 4", v)
	}

	// Calling Run a second time should be a no-op (idempotent).
	if err := Run(db, t.TempDir()); err != nil {
		t.Fatalf("second Run: %v", err)
	}
}

// TestRun_LegacyDatabaseAtV9 verifies that an existing v9 install transitions
// to goose without invoking legacy migrations and without altering the schema.
func TestRun_LegacyDatabaseAtV9(t *testing.T) {
	db := openTestDB(t)
	createV9Schema(t, db)

	// Snapshot a row count to verify the schema isn't altered.
	var tableCount int
	mustScan(t, db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table'`), &tableCount)

	if err := Run(db, t.TempDir()); err != nil {
		t.Fatalf("Run on v9 legacy DB: %v", err)
	}

	// Goose should record the latest migration.
	if v := gooseVersion(t, db); v != 4 {
		t.Errorf("goose version = %d, want 4", v)
	}

	// Config version should still be 9, the legacy bridge was not invoked.
	var version int
	mustScan(t, db.QueryRow(`SELECT value FROM config WHERE key='version'`), &version)
	if version != 9 {
		t.Errorf("config.version = %d, want 9", version)
	}

	// goose_db_version plus the federated_servers table were added by
	// the featured-streams migration; legacy schemas already have an
	// ap_followers.owncast_server column applied via the same migration.
	var newTableCount int
	mustScan(t, db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table'`), &newTableCount)
	if newTableCount != tableCount+2 { // +1 goose_db_version, +1 federated_servers
		t.Errorf("table count changed from %d to %d (expected +2 for goose_db_version + federated_servers)", tableCount, newTableCount)
	}
}

// TestRun_LegacyDatabasePreV9 verifies that a pre-v9 install runs the legacy
// bridge to reach v9, then goose records the baseline.
func TestRun_LegacyDatabasePreV9(t *testing.T) {
	// The legacy migration code writes a backup file under the provided
	// backupDirectory; the test passes a temp dir so the side effect is
	// contained.

	db := openTestDB(t)

	// Create all v9 tables (the legacy code tolerates columns that already
	// exist — ALTERs become no-ops with warnings) but set config.version=7
	// so the bridge runs migrateToSchema8 and migrateToSchema9.
	createV9Schema(t, db)
	mustExec(t, db, `UPDATE config SET value = 7 WHERE key = 'version'`)

	if err := Run(db, t.TempDir()); err != nil {
		t.Fatalf("Run on v7 legacy DB: %v", err)
	}

	// Legacy bridge should have updated config.version to 9.
	var version int
	mustScan(t, db.QueryRow(`SELECT value FROM config WHERE key='version'`), &version)
	if version != 9 {
		t.Errorf("config.version = %d after legacy bridge, want 9", version)
	}

	// Goose should have recorded the latest migration.
	if v := gooseVersion(t, db); v != 4 {
		t.Errorf("goose version = %d, want 4", v)
	}
}

// TestRun_BackfillsLinkedIdentities verifies that the 00004 migration upgrades
// existing auth rows in place: token becomes auth_key, provider is set, and the
// IndieAuth/Fediverse freebies (profile_url / handle) are backfilled — so
// returning users keep their linked identity across the upgrade.
func TestRun_BackfillsLinkedIdentities(t *testing.T) {
	db := openTestDB(t)
	createV9Schema(t, db) // auth table here still has the old "token" column

	// Seed one identity per built-in provider, as a pre-goose install would.
	mustExec(t, db, `INSERT INTO auth(user_id, token, type) VALUES('u-indie', 'https://me.example.com', 'indieauth')`)
	mustExec(t, db, `INSERT INTO auth(user_id, token, type) VALUES('u-fedi', '@me@host.example', 'fediverse')`)

	if err := Run(db, t.TempDir()); err != nil {
		t.Fatalf("Run: %v", err)
	}

	// IndieAuth: auth_key carried over, provider=type, profile_url=auth_key,
	// no handle, consent off.
	var authKey, provider string
	var profileURL, handle sql.NullString
	var isPublic bool
	mustScan(t, db.QueryRow(`SELECT auth_key, provider, profile_url, handle, is_public FROM auth WHERE user_id='u-indie'`),
		&authKey, &provider, &profileURL, &handle, &isPublic)
	if authKey != "https://me.example.com" || provider != "indieauth" {
		t.Errorf("indieauth row: auth_key=%q provider=%q", authKey, provider)
	}
	if !profileURL.Valid || profileURL.String != "https://me.example.com" {
		t.Errorf("indieauth profile_url = %v, want the me URL", profileURL)
	}
	if handle.Valid {
		t.Errorf("indieauth handle should be NULL, got %q", handle.String)
	}
	if isPublic {
		t.Errorf("indieauth is_public = %v, want false", isPublic)
	}

	// Fediverse: handle backfilled from auth_key, no profile_url (resolved later).
	mustScan(t, db.QueryRow(`SELECT auth_key, provider, profile_url, handle, is_public FROM auth WHERE user_id='u-fedi'`),
		&authKey, &provider, &profileURL, &handle, &isPublic)
	if authKey != "@me@host.example" || provider != "fediverse" {
		t.Errorf("fediverse row: auth_key=%q provider=%q", authKey, provider)
	}
	if profileURL.Valid {
		t.Errorf("fediverse profile_url should be NULL, got %q", profileURL.String)
	}
	if !handle.Valid || handle.String != "@me@host.example" {
		t.Errorf("fediverse handle = %v, want the account", handle)
	}
}

// createV9Schema creates the full v9 schema plus the legacy config table with
// version=9, simulating what an existing install looked like before goose.
func createV9Schema(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExec(t, db, `CREATE TABLE config ("key" TEXT NOT NULL PRIMARY KEY, "value" TEXT)`)
	mustExec(t, db, `INSERT INTO config(key, value) VALUES('version', 9)`)

	mustExec(t, db, `CREATE TABLE datastore ("key" TEXT NOT NULL PRIMARY KEY, "value" BLOB, "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL)`)
	mustExec(t, db, `CREATE TABLE webhooks ("id" INTEGER PRIMARY KEY AUTOINCREMENT, "url" TEXT NOT NULL, "events" TEXT NOT NULL, "timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP, "last_used" DATETIME)`)
	mustExec(t, db, `CREATE TABLE users ("id" TEXT NOT NULL, "display_name" TEXT NOT NULL, "display_color" INTEGER NOT NULL, "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP, "disabled_at" TIMESTAMP, "previous_names" TEXT DEFAULT '', "namechanged_at" TIMESTAMP, "authenticated_at" TIMESTAMP, "scopes" TEXT, "type" TEXT DEFAULT 'STANDARD', "last_used" DATETIME DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id))`)
	mustExec(t, db, `CREATE TABLE user_access_tokens ("token" TEXT NOT NULL PRIMARY KEY, "user_id" TEXT NOT NULL, "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL)`)
	mustExec(t, db, `CREATE TABLE ap_followers ("iri" TEXT NOT NULL, "inbox" TEXT NOT NULL, "shared_inbox" TEXT, "name" TEXT, "username" TEXT NOT NULL, "image" TEXT, "request" TEXT NOT NULL, "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP, "approved_at" TIMESTAMP, "disabled_at" TIMESTAMP, "request_object" BLOB, "last_validated_at" TIMESTAMP, "first_validation_failure_at" TIMESTAMP, PRIMARY KEY (iri))`)
	mustExec(t, db, `CREATE TABLE ap_outbox ("iri" TEXT NOT NULL, "value" BLOB, "type" TEXT NOT NULL, "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP, "live_notification" BOOLEAN DEFAULT FALSE, PRIMARY KEY (iri))`)
	mustExec(t, db, `CREATE TABLE ap_accepted_activities ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "iri" TEXT NOT NULL, "actor" TEXT NOT NULL, "type" TEXT NOT NULL, "timestamp" TIMESTAMP NOT NULL)`)
	mustExec(t, db, `CREATE TABLE notifications ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "channel" TEXT NOT NULL, "destination" TEXT NOT NULL, "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
	mustExec(t, db, `CREATE TABLE messages ("id" TEXT NOT NULL, "user_id" TEXT, "body" TEXT, "eventType" TEXT, "hidden_at" DATETIME, "timestamp" DATETIME, "title" TEXT, "subtitle" TEXT, "image" TEXT, "link" TEXT, PRIMARY KEY (id))`)
	mustExec(t, db, `CREATE TABLE auth ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "user_id" TEXT NOT NULL, "token" TEXT NOT NULL, "type" TEXT NOT NULL, "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL)`)
	mustExec(t, db, `CREATE TABLE ip_bans ("ip_address" TEXT NOT NULL PRIMARY KEY, "notes" TEXT, "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
}

func mustScan(t *testing.T, row *sql.Row, dest ...any) {
	t.Helper()
	if err := row.Scan(dest...); err != nil {
		t.Fatal(err)
	}
}
