package plugins

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newRunner(t *testing.T) SQLRunner {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	return SQLRunner{DB: db}
}

// An integral JSON parameter must reach SQLite as an int64. Decoding into
// interface{} yields float64, which rounds past 2^53 and can bind with REAL
// affinity, silently corrupting what the plugin stored.
func TestParseSQLRequestKeepsIntegerParametersExact(t *testing.T) {
	req, err := ParseSQLRequest(`{"sql":"SELECT ?, ?, ?, ?, ?","params":[1152921504606846977,-5,1.5,"x",null]}`)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := req.Params[0], int64(1152921504606846977); got != want {
		t.Fatalf("large integer decoded as %#v (%T), want %d", got, got, want)
	}
	if got, want := req.Params[1], int64(-5); got != want {
		t.Fatalf("negative integer decoded as %#v, want %d", got, want)
	}
	if got, want := req.Params[2], 1.5; got != want {
		t.Fatalf("fractional number decoded as %#v, want %v", got, want)
	}
	if req.Params[3] != "x" || req.Params[4] != nil {
		t.Fatalf("string/null parameters decoded as %#v, %#v", req.Params[3], req.Params[4])
	}
}

func TestSQLRunnerStoresLargeIntegersAsIntegers(t *testing.T) {
	runner := newRunner(t)
	ctx := context.Background()
	if result := runner.Exec(ctx, SQLRequest{SQL: "CREATE TABLE t (v)"}); result.Error != "" {
		t.Fatal(result.Error)
	}
	req, err := ParseSQLRequest(`{"sql":"INSERT INTO t VALUES (?)","params":[1152921504606846977]}`)
	if err != nil {
		t.Fatal(err)
	}
	if result := runner.Exec(ctx, req); result.Error != "" {
		t.Fatal(result.Error)
	}
	// A column with no affinity keeps whatever type was bound, so this fails if
	// the parameter arrived as a float.
	result := runner.Query(ctx, SQLRequest{SQL: "SELECT typeof(v), v FROM t"})
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if got := result.Rows[0][0]; got != "integer" {
		t.Fatalf("stored value has type %v, want integer", got)
	}
	if got, want := result.Rows[0][1], int64(1152921504606846977); got != want {
		t.Fatalf("stored value is %v, want %d", got, want)
	}
}

func TestParseSQLRequestRejectsBadRequests(t *testing.T) {
	for name, raw := range map[string]string{
		"empty statement":      `{"sql":"   "}`,
		"whitespace only":      `{"sql":"\n\t"}`,
		"not json":             `{"sql":`,
		"trailing object":      `{"sql":"SELECT 1"}{"sql":"SELECT 2"}`,
		"trailing junk":        `{"sql":"SELECT 1"}junk`,
		"array parameter":      `{"sql":"SELECT ?","params":[[1,2]]}`,
		"object parameter":     `{"sql":"SELECT ?","params":[{"a":1}]}`,
		"too many parameters":  fmt.Sprintf(`{"sql":"SELECT 1","params":[%s]}`, strings.TrimSuffix(strings.Repeat("1,", MaxSQLParams+1), ",")),
		"oversized statement":  fmt.Sprintf(`{"sql":"SELECT '%s'"}`, strings.Repeat("x", MaxSQLStatementBytes)),
		"unrepresentable size": `{"sql":"SELECT ?","params":[1e999]}`,
		"negative max rows":    `{"sql":"SELECT 1","maxRows":-1}`,
		"too many max rows":    fmt.Sprintf(`{"sql":"SELECT 1","maxRows":%d}`, MaxSQLRows+1),
	} {
		if _, err := ParseSQLRequest(raw); err == nil {
			t.Errorf("%s: expected a validation error", name)
		}
	}
	if _, err := ParseSQLRequest(`{"sql":"SELECT 1"}`); err != nil {
		t.Fatalf("a valid request was rejected: %v", err)
	}
}

// An unbounded query that overruns the row cap is an error, not a short answer:
// a plugin must not silently act on a truncated view of its own table.
func TestSQLRunnerRowCapIsAnErrorWhenUnbounded(t *testing.T) {
	runner := newRunner(t)
	ctx := context.Background()
	seed(t, runner, MaxSQLRows+1)

	result := runner.Query(ctx, SQLRequest{SQL: "SELECT n FROM t"})
	if result.Error == "" || !strings.Contains(result.Error, "add a LIMIT") {
		t.Fatalf("expected a row-cap error, got %d rows and error %q", len(result.Rows), result.Error)
	}
}

// A caller-supplied MaxRows is intent, so the host stops there and says it
// truncated. This is what lets queryRow read one row out of a large table.
func TestSQLRunnerMaxRowsTruncatesInsteadOfFailing(t *testing.T) {
	runner := newRunner(t)
	ctx := context.Background()
	seed(t, runner, MaxSQLRows+50)

	result := runner.Query(ctx, SQLRequest{SQL: "SELECT n FROM t ORDER BY n", MaxRows: 1})
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	if len(result.Rows) != 1 {
		t.Fatalf("asked for 1 row, got %d", len(result.Rows))
	}
	if !result.Truncated {
		t.Fatal("expected truncated to report that more rows matched")
	}
	if got, want := result.Rows[0][0], int64(0); got != want {
		t.Fatalf("first row is %v, want %v", got, want)
	}

	exact := runner.Query(ctx, SQLRequest{SQL: "SELECT n FROM t WHERE n < 3 ORDER BY n", MaxRows: 3})
	if exact.Error != "" || len(exact.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %+v", exact)
	}
	if exact.Truncated {
		t.Fatal("a result that fit must not report truncation")
	}
}

// The runner is the portable floor under SQLITE_LIMIT_LENGTH: a driver that
// cannot pin that limit must still not be able to hand a huge value through.
func TestSQLRunnerBoundsValueAndResultSize(t *testing.T) {
	runner := newRunner(t)
	ctx := context.Background()

	oversized := runner.Query(ctx, SQLRequest{SQL: fmt.Sprintf("SELECT hex(zeroblob(%d))", MaxSQLValueBytes)})
	if oversized.Error == "" || !strings.Contains(oversized.Error, "value exceeds") {
		t.Fatalf("expected a per-value limit error, got %q", oversized.Error)
	}

	escaped := runner.Query(ctx, SQLRequest{
		SQL:    "SELECT ?",
		Params: []any{strings.Repeat("<", MaxSQLValueBytes/6+1)},
	})
	if escaped.Error == "" || !strings.Contains(escaped.Error, "value exceeds") {
		t.Fatalf("expected JSON escaping to count against the per-value limit, got %q", escaped.Error)
	}

	if result := runner.Exec(ctx, SQLRequest{SQL: "CREATE TABLE t (n INTEGER, pad TEXT)"}); result.Error != "" {
		t.Fatal(result.Error)
	}
	pad := strings.Repeat("x", 4096)
	for i := 0; i < 400; i++ {
		if result := runner.Exec(ctx, SQLRequest{SQL: "INSERT INTO t VALUES (?, ?)", Params: []any{int64(i), pad}}); result.Error != "" {
			t.Fatal(result.Error)
		}
	}
	result := runner.Query(ctx, SQLRequest{SQL: "SELECT n, pad FROM t"})
	if result.Error == "" || !strings.Contains(result.Error, "add a LIMIT") {
		t.Fatalf("expected a result-size error well under the row cap, got %d rows and error %q", len(result.Rows), result.Error)
	}
	// The same query with a caller limit small enough to fit still works.
	bounded := runner.Query(ctx, SQLRequest{SQL: "SELECT n, pad FROM t", MaxRows: 10})
	if bounded.Error != "" || len(bounded.Rows) != 10 {
		t.Fatalf("a bounded read of the same table should succeed, got %+v", bounded.Error)
	}
}

func TestSQLValueSizeAccountsForJSONEscaping(t *testing.T) {
	for _, value := range []string{
		"plain",
		"\"\\\b\f\n\r\t",
		"\x00\x1f",
		"<>&",
		"\u2028\u2029",
		string([]byte{0xff}),
	} {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		if got := sqlValueSize(value); got != len(encoded) {
			t.Errorf("size of %q is %d, want %d", value, got, len(encoded))
		}
	}
}

// Whatever the runner returns has to survive the trip to the plugin as JSON,
// and it must stay inside the advertised budget.
func TestSQLRunnerResultEncodesWithinBudget(t *testing.T) {
	runner := newRunner(t)
	ctx := context.Background()
	seed(t, runner, 500)

	result := runner.Query(ctx, SQLRequest{SQL: "SELECT n FROM t ORDER BY n", MaxRows: 500})
	if result.Error != "" {
		t.Fatal(result.Error)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) > MaxSQLResultBytes {
		t.Fatalf("encoded result is %d bytes, over the %d-byte budget", len(encoded), MaxSQLResultBytes)
	}
	// The optional truncation marker is part of the same byte budget. This
	// value would fit if the builder forgot to count `,"truncated":true`, but
	// the final JSON would be one byte over the advertised cap.
	probe, err := json.Marshal(SQLQueryResult{
		Columns:   []string{"v"},
		Rows:      [][]any{{""}},
		Truncated: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	value := strings.Repeat("x", MaxSQLResultBytes-len(probe)+1)
	overflow := runner.Query(ctx, SQLRequest{
		SQL:     "SELECT ? AS v UNION ALL SELECT 'overflow'",
		Params:  []any{value},
		MaxRows: 1,
	})
	if overflow.Error == "" {
		encoded, err := json.Marshal(overflow)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("truncated result encoded to %d bytes without a limit error", len(encoded))
	}
}

func TestSQLRunnerWithoutDatabaseReportsUnavailable(t *testing.T) {
	var runner SQLRunner
	ctx := context.Background()
	if result := runner.Exec(ctx, SQLRequest{SQL: "SELECT 1"}); result.Error == "" {
		t.Fatalf("expected exec to report the store unavailable, got %+v", result)
	}
	if result := runner.Query(ctx, SQLRequest{SQL: "SELECT 1"}); result.Error == "" {
		t.Fatalf("expected query to report the store unavailable, got %+v", result)
	}
}

func seed(t *testing.T, runner SQLRunner, rows int) {
	t.Helper()
	ctx := context.Background()
	if result := runner.Exec(ctx, SQLRequest{SQL: "CREATE TABLE t (n INTEGER)"}); result.Error != "" {
		t.Fatal(result.Error)
	}
	statement := fmt.Sprintf("WITH RECURSIVE c(i) AS (SELECT 0 UNION ALL SELECT i + 1 FROM c WHERE i < %d) INSERT INTO t SELECT i FROM c", rows-1)
	if result := runner.Exec(ctx, SQLRequest{SQL: statement}); result.Error != "" {
		t.Fatal(result.Error)
	}
}
