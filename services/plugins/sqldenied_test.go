package plugins

import (
	"strings"
	"testing"
)

// Every statement in the shared parity fixture must be refused before it reaches
// a driver, so a plugin gets the same answer from Owncast and from the SDK's
// hosts. pluginhost has the matching test against the real authorizer.
func TestDeniedSQLReasonRefusesEveryParityExample(t *testing.T) {
	for _, statement := range DeniedSQLStatementExamples {
		if reason := DeniedSQLReason(statement); reason == "" {
			t.Errorf("statement is not refused: %q", statement)
		}
	}
}

// ParseSQLRequest is the single entry point both hosts use, so the refusal has to
// happen there and not only in the helper.
func TestParseSQLRequestRefusesDeniedStatements(t *testing.T) {
	for _, statement := range []string{
		"PRAGMA page_size",
		"CREATE TEMP TABLE t (v)",
		"CREATE TABLE temp.t (v)",
		"BEGIN IMMEDIATE",
		"COMMIT",
		"END",
		"ROLLBACK",
		"SAVEPOINT plugin",
		"RELEASE plugin",
		"CREATE TRIGGER touch AFTER INSERT ON items BEGIN SELECT CASE WHEN 1 THEN 1 END; END; COMMIT",
		`CREATE TRIGGER quoted AFTER INSERT ON items BEGIN UPDATE items SET name = "case"; END; PRAGMA page_size`,
	} {
		raw := `{"sql":` + quoteJSON(statement) + `}`
		if _, err := ParseSQLRequest(raw); err == nil {
			t.Errorf("ParseSQLRequest accepted %q", statement)
		}
	}
}

// The check must not cost ordinary SQL anything. These are the shapes a real
// plugin writes, including ones that mention the denied words inside string
// literals, comments, column names, and identifiers that merely start with
// "temp".
func TestDeniedSQLReasonAllowsOrdinarySQL(t *testing.T) {
	for _, statement := range []string{
		"CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT)",
		"CREATE TABLE main.items (id INTEGER PRIMARY KEY)",
		"CREATE INDEX items_name ON items (name)",
		"CREATE VIEW item_names AS SELECT name FROM items",
		"CREATE TRIGGER touch AFTER INSERT ON items BEGIN UPDATE items SET name = name; END",
		"CREATE TRIGGER choose AFTER INSERT ON items BEGIN UPDATE items SET name = CASE WHEN name = '' THEN 'empty' ELSE name END; END",
		`CREATE TRIGGER quoted AFTER INSERT ON items BEGIN UPDATE items SET "end" = "case"; END`,
		"INSERT INTO items (name) VALUES (?)",
		"SELECT name FROM items ORDER BY name",
		"EXPLAIN SELECT name FROM items",
		"EXPLAIN QUERY PLAN SELECT name FROM items",
		"WITH RECURSIVE c(i) AS (SELECT 1 UNION ALL SELECT i + 1 FROM c WHERE i < 10) SELECT count(*) FROM c",
		"SELECT json_object('name', name) FROM items",
		"UPDATE items SET name = 'PRAGMA page_size' WHERE id = 1",
		"INSERT INTO items (name) VALUES ('CREATE TEMP TABLE x')",
		"INSERT INTO items (name) VALUES ('attach database ''x'' as y')",
		"SELECT 'temp.scratch' AS label FROM items",
		"-- PRAGMA page_size\nSELECT 1",
		"/* ATTACH DATABASE */ SELECT 1",
		"CREATE TABLE temperature (celsius REAL)",
		"SELECT temperature FROM temperature",
		"CREATE TABLE templates (body TEXT)",
		"SELECT items.name FROM items",
		"SELECT main.items.name FROM main.items",
		"CREATE TABLE t (temp INTEGER)",
		"SELECT temp FROM t",
	} {
		if reason := DeniedSQLReason(statement); reason != "" {
			t.Errorf("ordinary SQL was refused: %q (%s)", statement, reason)
		}
	}
}

// The message reaches the plugin author, so it should name the capability rather
// than the internals.
func TestDeniedSQLReasonExplainsWhatIsUnavailable(t *testing.T) {
	for statement, want := range map[string]string{
		"PRAGMA page_size":                 "PRAGMA",
		"ATTACH DATABASE 'x' AS y":         "ATTACH",
		"CREATE TEMP TABLE t (v)":          "temporary tables",
		"CREATE TABLE temp.t (v)":          "temporary tables",
		"SELECT * FROM sqlite_temp_master": "temp schema",
	} {
		reason := DeniedSQLReason(statement)
		if !strings.Contains(reason, want) {
			t.Errorf("%q produced %q, want it to mention %q", statement, reason, want)
		}
	}
}

func quoteJSON(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
