package pluginhost

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/owncast/owncast/services/plugins"
)

const (
	// pluginSQLiteDriver is a private driver registration: the connection hook
	// below locks every plugin connection down, and Owncast's own datastore must
	// keep using the stock "sqlite3" driver.
	pluginSQLiteDriver = "owncast_plugin_sqlite"
	// pluginSQLFileName is the database inside each plugin's SQL directory.
	pluginSQLFileName = "plugin.db"
	// pluginSQLitePageSize is the page size a freshly created plugin database
	// gets. An existing database keeps whatever it was created with, so the
	// page quota is always derived from the size actually in force.
	pluginSQLitePageSize = 4096
	// pluginSQLiteMaxBytes caps one plugin's database. It is deliberately
	// independent of the storage.fs quota: the two live in separate directory
	// trees, so neither can consume the other's headroom and the limit a plugin
	// hits does not depend on the order it wrote its data.
	pluginSQLiteMaxBytes = 128 << 20 // 128 MiB per plugin
)

var registerPluginSQLiteDriver sync.Once

// deniedSQLiteOps are the authorizer actions a plugin may never perform. The
// authorizer is the security boundary: Exec runs every statement in a batch, so
// a check on statement text alone is bypassable, and a denied PRAGMA is what
// makes max_page_count a quota rather than a suggestion.
//
// Denying PRAGMA also denies PRAGMA reads, so the connection hook reads
// page_size before installing the authorizer.
//
// Temp-schema DDL is denied because a temp object lives outside the main
// database's page quota and persists for the life of the pooled connection,
// which is an unbounded side channel. Sorters, CTEs, subqueries, and json1 are
// unaffected: they use internal ephemeral tables, not temp-schema DDL.
var deniedSQLiteOps = map[int]bool{
	sqlite3.SQLITE_ATTACH:              true,
	sqlite3.SQLITE_DETACH:              true,
	sqlite3.SQLITE_PRAGMA:              true,
	sqlite3.SQLITE_CREATE_TEMP_TABLE:   true,
	sqlite3.SQLITE_CREATE_TEMP_INDEX:   true,
	sqlite3.SQLITE_CREATE_TEMP_TRIGGER: true,
	sqlite3.SQLITE_CREATE_TEMP_VIEW:    true,
	sqlite3.SQLITE_DROP_TEMP_TABLE:     true,
	sqlite3.SQLITE_DROP_TEMP_INDEX:     true,
	sqlite3.SQLITE_DROP_TEMP_TRIGGER:   true,
	sqlite3.SQLITE_DROP_TEMP_VIEW:      true,
}

// temporarySchemaName is the schema a plugin must not reach. Denying the
// SQLITE_*_TEMP_* actions is not enough on its own: SQLite reports
// `CREATE TABLE temp.x` and `CREATE VIEW temp.v` as the ordinary CREATE actions
// and puts "temp" in the authorizer's database-name argument instead, so those
// statements slip past an op-only deny list and put an unquota'd table in the
// temp schema. Denying the schema by name closes every spelling of it,
// including the quoted and mixed-case ones.
const temporarySchemaName = "temp"

// pluginSQLStore owns one database file and one serialized connection per
// plugin. root is the shared plugin-storage root, and each database lives in
// that plugin's db/ subdirectory, which sits outside the storage.fs sandbox so
// no plugin-facing API can reach it.
type pluginSQLStore struct {
	root string

	mu      sync.Mutex
	dbs     map[string]*pluginSQLDatabase
	leases  map[string]int
	closing map[string]int
}

type pluginSQLDatabase struct {
	db *sql.DB
}

func newPluginSQLStore(root string) (*pluginSQLStore, error) {
	if err := os.MkdirAll(root, 0o700); err != nil {
		return nil, fmt.Errorf("create plugin SQL root: %w", err)
	}
	registerPluginSQLiteDriver.Do(func() {
		sql.Register(pluginSQLiteDriver, &sqlite3.SQLiteDriver{
			ConnectHook: configurePluginSQLiteConnection,
		})
	})
	return &pluginSQLStore{
		root:    root,
		dbs:     make(map[string]*pluginSQLDatabase),
		leases:  make(map[string]int),
		closing: make(map[string]int),
	}, nil
}

// configurePluginSQLiteConnection applies the sandbox to one connection. Order
// matters: the PRAGMAs and the page-size read happen first, then the limits and
// the authorizer lock the connection for the rest of its life.
func configurePluginSQLiteConnection(conn *sqlite3.SQLiteConn) error {
	for _, pragma := range []string{
		fmt.Sprintf("PRAGMA page_size=%d", pluginSQLitePageSize),
		// TRUNCATE rather than WAL: a rollback journal is truncated to zero on
		// commit, so a plugin's on-disk footprint stays at its page quota
		// instead of overshooting it by a WAL that is only reclaimed at
		// checkpoint. With one connection there is no concurrency to gain.
		"PRAGMA journal_mode=TRUNCATE",
		// FULL, not NORMAL: a rollback journal only guarantees durability when
		// the journal is synced, and plugin databases are not in Owncast's
		// backup, so a torn write has nothing to restore from.
		"PRAGMA synchronous=FULL",
		"PRAGMA cache_size=-2000", // ~2 MiB, not the datastore's 10000 pages
		"PRAGMA foreign_keys=ON",
		"PRAGMA trusted_schema=OFF",
		// Keep temp data off disk entirely; temp-schema DDL is denied below, so
		// what remains is bounded sorter and CTE scratch space.
		"PRAGMA temp_store=MEMORY",
	} {
		if _, err := conn.Exec(pragma, nil); err != nil {
			return fmt.Errorf("%s: %w", pragma, err)
		}
	}

	// An existing database keeps the page size it was created with, so read the
	// size in force and derive the page cap from it. Assuming 4 KiB would turn
	// the quota into 16x its intended size on a 64 KiB-page database.
	pageSize, err := sqliteScalarInt(conn, "PRAGMA page_size")
	if err != nil {
		return fmt.Errorf("read plugin database page size: %w", err)
	}
	if pageSize <= 0 {
		return fmt.Errorf("plugin database reports page size %d", pageSize)
	}
	maxPages := pluginSQLiteMaxBytes / pageSize
	if _, err := conn.Exec(fmt.Sprintf("PRAGMA max_page_count=%d", maxPages), nil); err != nil {
		return fmt.Errorf("set plugin database page cap: %w", err)
	}

	// SQLITE_LIMIT_LENGTH refuses an oversized string, blob, or row inside the
	// engine, so a plugin cannot make the host allocate hundreds of megabytes
	// with SELECT randomblob(...) only to have the result rejected afterwards.
	conn.SetLimit(sqlite3.SQLITE_LIMIT_LENGTH, plugins.MaxSQLValueBytes)
	conn.SetLimit(sqlite3.SQLITE_LIMIT_SQL_LENGTH, plugins.MaxSQLStatementBytes)
	conn.SetLimit(sqlite3.SQLITE_LIMIT_VARIABLE_NUMBER, plugins.MaxSQLParams)
	// Belt for the ATTACH denial: with zero attached databases permitted, an
	// ATTACH cannot succeed even if the authorizer is ever relaxed.
	conn.SetLimit(sqlite3.SQLITE_LIMIT_ATTACHED, 0)

	// The fourth callback argument is SQLite's database-name argument, which
	// carries the schema an action targets.
	conn.RegisterAuthorizer(func(op int, arg1, arg2, schema string) int {
		if strings.EqualFold(schema, temporarySchemaName) {
			return sqlite3.SQLITE_DENY
		}
		if deniedSQLiteOps[op] {
			return sqlite3.SQLITE_DENY
		}
		if op == sqlite3.SQLITE_FUNCTION &&
			(strings.EqualFold(arg1, "load_extension") || strings.EqualFold(arg2, "load_extension")) {
			return sqlite3.SQLITE_DENY
		}
		return sqlite3.SQLITE_OK
	})
	return nil
}

// sqliteScalarInt runs a single-value query on a raw connection. Used for the
// page-size read that has to happen before the authorizer denies PRAGMA.
func sqliteScalarInt(conn *sqlite3.SQLiteConn, query string) (int, error) {
	rows, err := conn.Query(query, nil)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	values := make([]driver.Value, 1)
	if err := rows.Next(values); err != nil {
		return 0, err
	}
	value, ok := values[0].(int64)
	if !ok {
		return 0, fmt.Errorf("expected an integer, got %T", values[0])
	}
	return int(value), nil
}

// acquire gives one loaded plugin instance a claim on its slug's database.
// Multiple instances can briefly coexist during package preflight or reload,
// so only the last release closes the shared handle.
func (s *pluginSQLStore) acquire(pluginName string) func() {
	s.mu.Lock()
	s.leases[pluginName]++
	s.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() { s.release(pluginName) })
	}
}

func (s *pluginSQLStore) release(pluginName string) {
	s.mu.Lock()
	if s.leases[pluginName] > 1 {
		s.leases[pluginName]--
		s.mu.Unlock()
		return
	}
	delete(s.leases, pluginName)
	s.closing[pluginName]++
	entry := s.dbs[pluginName]
	delete(s.dbs, pluginName)
	s.mu.Unlock()

	closeSQLDatabase(entry)

	s.mu.Lock()
	if s.closing[pluginName] > 1 {
		s.closing[pluginName]--
		s.mu.Unlock()
		return
	}
	delete(s.closing, pluginName)
	// A handle here would mean a host call bypassed the closing block. Sweep it
	// rather than retain a resource with no loaded instance.
	entry = s.dbs[pluginName]
	delete(s.dbs, pluginName)
	s.mu.Unlock()
	closeSQLDatabase(entry)
}

// get returns the plugin's database, opening it on first use.
func (s *pluginSQLStore) get(pluginName string) (*pluginSQLDatabase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closing[pluginName] > 0 {
		return nil, fmt.Errorf("plugin %q is unloading; SQL storage is closed", pluginName)
	}
	if db := s.dbs[pluginName]; db != nil {
		return db, nil
	}

	path, err := resolvePluginDatabasePath(s.root, pluginName, pluginSQLFileName)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create plugin SQL directory: %w", err)
	}
	db, err := sql.Open(pluginSQLiteDriver, "file:"+path+"?_busy_timeout=5000&_txlock=immediate")
	if err != nil {
		return nil, fmt.Errorf("open plugin SQL database: %w", err)
	}
	// One connection: plugin statements serialize against each other rather
	// than contending for locks on the same file.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize plugin SQL database: %w", err)
	}
	entry := &pluginSQLDatabase{db: db}
	s.dbs[pluginName] = entry
	return entry, nil
}

// runner resolves the plugin's database and wraps it in the shared runner that
// applies the host's statement, timeout, and result limits.
func (s *pluginSQLStore) runner(pluginName string) (plugins.SQLRunner, error) {
	entry, err := s.get(pluginName)
	if err != nil {
		return plugins.SQLRunner{}, err
	}
	return plugins.SQLRunner{DB: entry.db}, nil
}

func (s *pluginSQLStore) exec(ctx context.Context, pluginName string, req plugins.SQLRequest) plugins.SQLExecResult {
	runner, err := s.runner(pluginName)
	if err != nil {
		return plugins.SQLExecResult{Error: err.Error()}
	}
	return runner.Exec(ctx, req)
}

func (s *pluginSQLStore) query(ctx context.Context, pluginName string, req plugins.SQLRequest) plugins.SQLQueryResult {
	runner, err := s.runner(pluginName)
	if err != nil {
		return plugins.SQLQueryResult{Error: err.Error()}
	}
	return runner.Query(ctx, req)
}

// closeAll releases every open database. The final sweep on host shutdown.
func (s *pluginSQLStore) closeAll() {
	s.mu.Lock()
	entries := s.dbs
	s.dbs = make(map[string]*pluginSQLDatabase)
	s.leases = make(map[string]int)
	s.mu.Unlock()
	for _, entry := range entries {
		closeSQLDatabase(entry)
	}
}

// closeSQLDatabase closes one plugin database. database/sql prevents new work
// and waits for statements already using the connection.
func closeSQLDatabase(entry *pluginSQLDatabase) {
	if entry == nil {
		return
	}
	_ = entry.db.Close()
}
