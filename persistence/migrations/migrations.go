// Package migrations is the single entry point for evolving the Owncast
// database schema.
//
// On startup, Run is called exactly once by the persistence layer. It brings
// any install (fresh or existing) up to the current schema and then returns.
// All further schema changes must be added as additional numbered SQL files
// in this directory; never amend an existing migration that has shipped.
//
// # Legacy-install support
//
// Before the goose system existed, schema state was tracked in a hand-rolled
// `config.version` row and advanced by a switch statement in
// persistence/legacymigrations. For installs that predate this package, Run
// detects that row and invokes the legacy catch-up path to bring the schema
// to the point where goose can take over. The legacy package is frozen —
// no new cases may be added to it; all new schema work happens here.
package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/persistence/legacymigrations"
)

//go:embed *.sql
var embedMigrations embed.FS

// legacyBaselineVersion is the schema version the legacy hand-rolled
// migration system reached immediately before goose took over. Any install
// whose legacy `config.version` row is below this value is caught up by the
// frozen legacymigrations package before goose runs.
const legacyBaselineVersion = 9

// Run brings the database schema up to date.
//
// For fresh installs, goose runs the baseline migration plus anything newer.
// For installs that predate goose, the legacy catch-up runs first and then
// the baseline is recorded as applied (its statements are all IF NOT EXISTS
// and become no-ops on an already-populated database).
// backupDirectory is where the legacy bridge writes pre-migration backups.
func Run(db *sql.DB, backupDirectory string) error {
	legacyVersion, err := readLegacyVersion(db)
	if err != nil {
		return fmt.Errorf("reading legacy schema version: %w", err)
	}

	if legacyVersion > 0 && legacyVersion < legacyBaselineVersion {
		log.Infof("Legacy schema at version %d; upgrading to %d before handing off to goose",
			legacyVersion, legacyBaselineVersion)
		if err := legacymigrations.MigrateDatabaseSchema(db, backupDirectory, legacyVersion, legacyBaselineVersion); err != nil {
			return fmt.Errorf("legacy schema catch-up: %w", err)
		}
	}

	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(gooseLogger{})
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("setting goose dialect: %w", err)
	}
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}
	return nil
}

// readLegacyVersion returns the schema version stored in the legacy `config`
// table, or 0 if that table (or the version row) does not exist — which is
// how we recognise a fresh install.
func readLegacyVersion(db *sql.DB) (int, error) {
	var name string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='config'`).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	var version int
	err = db.QueryRow(`SELECT value FROM config WHERE key='version'`).Scan(&version)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return version, nil
}

// gooseLogger adapts goose's Logger interface onto logrus so migration
// output flows through the project's normal logging pipeline. It silences
// the "no migrations to run" message that goose emits on every startup
// when the schema is already current.
type gooseLogger struct{}

func (gooseLogger) Fatalf(format string, v ...interface{}) { log.Fatalf(format, v...) }

func (gooseLogger) Printf(format string, v ...interface{}) {
	if strings.Contains(fmt.Sprintf(format, v...), "no migrations to run") {
		return
	}
	log.Infof(format, v...)
}
