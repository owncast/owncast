package pluginhost

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/owncast/owncast/services/plugins"
)

func newTestSQLStore(t *testing.T) (*pluginSQLStore, string) {
	t.Helper()
	root := t.TempDir()
	store, err := newPluginSQLStore(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(store.closeAll)
	return store, root
}

func TestPluginSQLStoreIsolationAndRoundTrip(t *testing.T) {
	store, _ := newTestSQLStore(t)
	ctx := context.Background()

	exec := func(plugin, statement string, params ...any) plugins.SQLExecResult {
		return store.exec(ctx, plugin, plugins.SQLRequest{SQL: statement, Params: params})
	}
	if result := exec("first", "CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT)"); result.Error != "" {
		t.Fatal(result.Error)
	}
	if result := exec("first", "INSERT INTO items (name) VALUES (?)", "one"); result.Error != "" || result.LastInsertID != 1 {
		t.Fatalf("insert result: %+v", result)
	}
	query := store.query(ctx, "first", plugins.SQLRequest{SQL: "SELECT id, name FROM items"})
	if query.Error != "" || len(query.Rows) != 1 || query.Rows[0][1] != "one" {
		t.Fatalf("query result: %+v", query)
	}
	empty := store.query(ctx, "first", plugins.SQLRequest{SQL: "SELECT id FROM items WHERE id = ?", Params: []any{int64(99)}})
	if empty.Error != "" || empty.Columns == nil || empty.Rows == nil || len(empty.Rows) != 0 {
		t.Fatalf("empty query result: %+v", empty)
	}
	other := store.query(ctx, "second", plugins.SQLRequest{SQL: "SELECT name FROM items"})
	if other.Error == "" {
		t.Fatal("plugin databases are not isolated")
	}
}

// A plugin must not be able to reach outside its own database, and must not be
// able to lift the limits that bound it.
//
// The store calls SQLite directly, without the Go-level check ParseSQLRequest
// applies at the host-function boundary, so this exercises the SQLite authorizer
// on its own. That matters twice over: the authorizer is the actual security
// boundary, and it has to refuse everything in the shared parity fixture, since
// the SDK's hosts refuse that same list with the Go check. A statement a plugin
// author sees refused in a scenario test is then refused on a real server too.
func TestPluginSQLStoreDeniesEscapesAndLimitChanges(t *testing.T) {
	store, root := newTestSQLStore(t)
	ctx := context.Background()

	for _, statement := range plugins.DeniedSQLStatementExamples {
		if result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: statement}); result.Error == "" {
			t.Errorf("the SQLite authorizer allowed %q", statement)
		}
	}

	// Forms that need a real path on this machine, so they cannot live in the
	// shared fixture.
	for name, statement := range map[string]string{
		"attach the host database": fmt.Sprintf("ATTACH DATABASE '%s' AS host", filepath.Join(root, "owncast.db")),
		"vacuum into a file":       fmt.Sprintf("VACUUM INTO '%s'", filepath.Join(root, "escape.db")),
	} {
		if result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: statement}); result.Error == "" {
			t.Errorf("%s: the SQLite authorizer allowed %q", name, statement)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "escape.db")); err == nil {
		t.Fatal("VACUUM INTO wrote a file outside the plugin's database")
	}
	// The refused statements left the connection usable.
	if result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: "CREATE TABLE main.after (v TEXT)"}); result.Error != "" {
		t.Fatalf("the connection should still be usable: %v", result.Error)
	}
	// The denials must not cost ordinary SQL anything: sorters, recursive CTEs,
	// subqueries, views, triggers, and json1 all still work.
	for _, statement := range []string{
		"CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT, n INTEGER)",
		"CREATE INDEX items_name ON items (name)",
		"CREATE VIEW item_names AS SELECT name FROM items",
		"CREATE TRIGGER touch AFTER INSERT ON items BEGIN UPDATE items SET n = n; END",
		"INSERT INTO items (name, n) VALUES ('b', 2), ('a', 1)",
	} {
		if result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: statement}); result.Error != "" {
			t.Fatalf("%q: %v", statement, result.Error)
		}
	}
	for _, statement := range []string{
		"SELECT name FROM items ORDER BY name",
		"WITH RECURSIVE counter(i) AS (SELECT 1 UNION ALL SELECT i + 1 FROM counter WHERE i < 500) SELECT count(*) FROM counter",
		"SELECT (SELECT count(*) FROM items) AS total",
		"SELECT json_object('name', name) FROM items",
		"SELECT name FROM item_names UNION SELECT name FROM items ORDER BY 1",
	} {
		if result := store.query(ctx, "plugin", plugins.SQLRequest{SQL: statement}); result.Error != "" {
			t.Errorf("%q should be allowed: %v", statement, result.Error)
		}
	}
}

// The page quota is derived from the page size actually in force. A database
// created with 64 KiB pages must still cap at pluginSQLiteMaxBytes, not at
// sixteen times that. The SQL root is unreachable from storage.fs, so this is
// defence in depth rather than the only barrier.
func TestPluginSQLStoreDerivesPageQuotaFromPageSize(t *testing.T) {
	store, root := newTestSQLStore(t)
	path := filepath.Join(root, "plugin", pluginDatabaseDirName, pluginSQLFileName)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	seed, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := seed.Exec("PRAGMA page_size=65536; VACUUM; CREATE TABLE big (v BLOB)"); err != nil {
		t.Fatal(err)
	}
	if err := seed.Close(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	payload := strings.Repeat("x", 900_000)
	var lastErr string
	inserted := 0
	for i := 0; i < 400; i++ {
		result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: "INSERT INTO big VALUES (?)", Params: []any{payload}})
		if result.Error != "" {
			lastErr = result.Error
			break
		}
		inserted++
	}
	if lastErr == "" {
		t.Fatal("expected the page quota to stop the inserts")
	}
	// Assert the quota is what stopped it: any other failure (a rejected
	// connection hook, a missing table) would also end the loop early and would
	// otherwise pass this test while proving nothing.
	if !strings.Contains(lastErr, "database or disk is full") {
		t.Fatalf("expected the page quota to stop the inserts, got %q after %d rows", lastErr, inserted)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() < pluginSQLiteMaxBytes/2 {
		t.Fatalf("database only reached %d bytes after %d rows; the quota tripped too early to prove anything", info.Size(), inserted)
	}
	if info.Size() > 2*pluginSQLiteMaxBytes {
		t.Fatalf("database grew to %d bytes with 64 KiB pages; the %d-byte quota was derived from the wrong page size", info.Size(), pluginSQLiteMaxBytes)
	}
	t.Logf("64 KiB-page database capped at %d bytes after %d rows", info.Size(), inserted)
}

// An oversized value must be refused inside SQLite rather than materialized and
// then rejected, so a plugin cannot make the host allocate its way to an OOM.
func TestPluginSQLStoreRefusesOversizedValuesCheaply(t *testing.T) {
	store, _ := newTestSQLStore(t)
	ctx := context.Background()

	start := time.Now()
	result := store.query(ctx, "plugin", plugins.SQLRequest{SQL: "SELECT randomblob(100000000)"})
	elapsed := time.Since(start)
	if result.Error == "" {
		t.Fatal("expected a 100 MB blob to be refused")
	}
	if elapsed > plugins.SQLCallTimeout {
		t.Fatalf("refusing the blob took %s; it should fail before materializing", elapsed)
	}
	if result := store.query(ctx, "plugin", plugins.SQLRequest{SQL: "SELECT hex(randomblob(10000000))"}); result.Error == "" {
		t.Fatal("expected a 20 MB hex string to be refused")
	}
	// A value inside the limit is unaffected.
	if result := store.query(ctx, "plugin", plugins.SQLRequest{SQL: "SELECT length(randomblob(1000))"}); result.Error != "" {
		t.Fatalf("a small blob should be allowed: %v", result.Error)
	}
}

// The timeout starts before database/sql waits for the plugin's sole
// connection. Two concurrent runaway calls must therefore expire together
// rather than each getting a fresh timeout after the previous call finishes.
func TestPluginSQLStoreTimeoutIncludesConnectionWait(t *testing.T) {
	store, _ := newTestSQLStore(t)
	const runaway = "WITH RECURSIVE c(i) AS (SELECT 1 UNION ALL SELECT i + 1 FROM c) SELECT sum(i) FROM c"

	start := make(chan struct{})
	results := make(chan plugins.SQLQueryResult, 2)
	for range 2 {
		go func() {
			<-start
			results <- store.query(context.Background(), "plugin", plugins.SQLRequest{SQL: runaway})
		}()
	}
	started := time.Now()
	close(start)
	for range 2 {
		if result := <-results; result.Error == "" {
			t.Fatal("runaway query completed without a timeout")
		}
	}
	if elapsed := time.Since(started); elapsed > plugins.SQLCallTimeout+time.Second {
		t.Fatalf("serialized calls took %s; timeout did not include the connection wait", elapsed)
	}
}

func TestPluginSQLStoreExecIsAtomic(t *testing.T) {
	store, _ := newTestSQLStore(t)
	ctx := context.Background()

	result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: "CREATE TABLE items (name TEXT); INSERT INTO items VALUES ('ok'); INSERT INTO missing VALUES (1)"})
	if result.Error == "" {
		t.Fatal("expected the failing statement to roll back the batch")
	}
	query := store.query(ctx, "plugin", plugins.SQLRequest{SQL: "SELECT count(*) FROM items"})
	if query.Error == "" || !strings.Contains(query.Error, "no such table") {
		t.Fatalf("expected atomic rollback to remove the table, got %+v", query)
	}
	// A plugin cannot leave a transaction open across calls.
	if result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: "BEGIN"}); result.Error == "" {
		t.Fatal("expected BEGIN inside the host transaction to fail")
	}
	if result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: "CREATE TABLE after (v TEXT)"}); result.Error != "" {
		t.Fatalf("the connection should still be usable: %v", result.Error)
	}
}

// A package preflight can load a second instance with the same slug. Releasing
// that throwaway instance must preserve the live instance's open handle; only
// the final instance release closes it.
func TestPluginSQLStoreLeasesPreserveLiveHandle(t *testing.T) {
	store, _ := newTestSQLStore(t)
	ctx := context.Background()

	releaseLive := store.acquire("plugin")
	if result := store.exec(ctx, "plugin", plugins.SQLRequest{SQL: "CREATE TABLE items (value TEXT); INSERT INTO items VALUES ('live')"}); result.Error != "" {
		t.Fatal(result.Error)
	}
	store.mu.Lock()
	liveDB := store.dbs["plugin"]
	store.mu.Unlock()

	releasePreflight := store.acquire("plugin")
	releasePreflight()
	store.mu.Lock()
	afterPreflight := store.dbs["plugin"]
	leases := store.leases["plugin"]
	store.mu.Unlock()
	if afterPreflight != liveDB || leases != 1 {
		t.Fatalf("preflight release changed live SQL state: db=%p want %p, leases=%d want 1", afterPreflight, liveDB, leases)
	}

	releaseLive()
	store.mu.Lock()
	open := len(store.dbs)
	leases = len(store.leases)
	store.mu.Unlock()
	if open != 0 || leases != 0 {
		t.Fatalf("final release left %d handles and %d leases", open, leases)
	}

	releaseReload := store.acquire("plugin")
	defer releaseReload()
	query := store.query(ctx, "plugin", plugins.SQLRequest{SQL: "SELECT value FROM items"})
	if query.Error != "" || len(query.Rows) != 1 || query.Rows[0][0] != "live" {
		t.Fatalf("reload did not reopen the persisted database: %+v", query)
	}
}

// A plugin's files and its database share one storage directory, which is what
// an operator sees, but they are in separate subtrees. The storage.fs sandbox is
// rooted at files/, so db/ is not a path the filesystem API refuses, it is a
// path the filesystem API cannot express. That distinction is the whole point:
// reserving filenames inside one shared directory would have to keep pace with
// SQLite's sidecar naming, and a single missed operation would let a plugin
// pre-seed a database larger than its quota, read its raw bytes, or delete it
// while the host holds it open.
func TestPluginDatabaseIsUnreachableFromFilesystemSandbox(t *testing.T) {
	dataDir := t.TempDir()
	root := filepath.Join(dataDir, pluginStorageRootDirName)

	store, err := newPluginSQLStore(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(store.closeAll)
	if result := store.exec(context.Background(), "plugin", plugins.SQLRequest{SQL: "CREATE TABLE items (v TEXT)"}); result.Error != "" {
		t.Fatal(result.Error)
	}

	dbPath := filepath.Join(root, "plugin", pluginDatabaseDirName, pluginSQLFileName)
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected the database at %s: %v", dbPath, err)
	}
	// Both live under one per-plugin directory, so an operator has a single
	// place to look, back up, or delete.
	storageDir := pluginStorageDir(root, "plugin")
	if filepath.Dir(filepath.Dir(dbPath)) != storageDir {
		t.Fatalf("database at %s is not inside the plugin's storage directory %s", dbPath, storageDir)
	}

	// Every storage.fs path, including traversal attempts and the database's own
	// absolute path, resolves back inside files/ and can never name db/.
	filesDir := filepath.Join(storageDir, pluginFilesDirName)
	for _, rel := range []string{
		pluginSQLFileName,
		pluginDatabaseDirName + "/" + pluginSQLFileName,
		"../" + pluginDatabaseDirName + "/" + pluginSQLFileName,
		"../../plugin/" + pluginDatabaseDirName + "/" + pluginSQLFileName,
		"../../../" + pluginStorageRootDirName + "/plugin/" + pluginDatabaseDirName + "/" + pluginSQLFileName,
		dbPath,
		dbPath + "-journal",
	} {
		resolved, err := resolvePluginSandboxPath(root, "plugin", rel)
		if err != nil {
			continue
		}
		if resolved == dbPath || strings.HasPrefix(resolved, filepath.Join(storageDir, pluginDatabaseDirName)+string(os.PathSeparator)) {
			t.Errorf("storage.fs path %q reaches the database subtree (%s)", rel, resolved)
		}
		if !strings.HasPrefix(resolved, filesDir+string(os.PathSeparator)) && resolved != filesDir {
			t.Errorf("storage.fs path %q resolved outside files/ (%s)", rel, resolved)
		}
	}

	// The quota walk covers files/ only, so a database never consumes the
	// filesystem quota and the two limits stay independent whichever API a
	// plugin writes with first.
	used, err := dirSize(filesDir)
	if err != nil {
		t.Fatal(err)
	}
	if used != 0 {
		t.Fatalf("the plugin database counted %d bytes against the storage.fs quota", used)
	}
	if err := fsQuotaCheck(root, "plugin", filepath.Join(filesDir, "note.txt"), maxPluginDataBytes); err != nil {
		t.Fatalf("a write filling the whole filesystem quota should still be allowed: %v", err)
	}
}
