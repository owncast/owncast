package plugins

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

// Limits applied to plugin SQL in Go rather than by the SQLite driver, so they
// hold whatever driver a host wires up. pluginhost additionally pins the
// matching SQLite-level limits (SQLITE_LIMIT_LENGTH, SQLITE_LIMIT_VARIABLE_NUMBER,
// …) so an oversized value is refused inside the engine before it is ever
// materialized; these checks are the driver-independent floor under that.
const (
	// MaxSQLStatementBytes caps the JSON request a plugin passes to the host,
	// which bounds the statement text along with it.
	MaxSQLStatementBytes = 64 << 10
	// MaxSQLParams caps bound parameters per statement.
	MaxSQLParams = 64
	// MaxSQLValueBytes caps a single returned column value.
	MaxSQLValueBytes = 1 << 20
	// MaxSQLResultBytes caps the encoded size of a whole query result.
	MaxSQLResultBytes = 1 << 20
	// MaxSQLRows caps rows returned when the plugin does not ask for fewer.
	// Overrunning it is an error, not a silent truncation: a plugin that reads
	// an unbounded table must say so with a LIMIT.
	MaxSQLRows = 10000
	// SQLCallTimeout bounds one exec or query, including time spent waiting on
	// the plugin's serialized connection.
	SQLCallTimeout = 2 * time.Second
)

// sqlStorageUnavailable is what a plugin sees when the host it is running in
// wires no SQL backend at all.
const sqlStorageUnavailable = "SQL storage unavailable"

// errSQLResultTooLarge is the shared "narrow your query" error. Plugins see the
// text, so it names the fix rather than the internal limit that tripped.
func errSQLResultTooLarge(what string, limit int) error {
	return fmt.Errorf("SQL %s exceeds %d bytes; add a LIMIT or select fewer columns", what, limit)
}

// ParseSQLRequest decodes and validates the JSON request a plugin passes to
// owncast_sql_exec / owncast_sql_query.
//
// Numbers are decoded with json.Number so an integral parameter binds as an
// int64: plain interface{} decoding would turn every number into a float64 and
// silently round anything past 2^53 before it reached SQLite, or store it with
// REAL affinity.
func ParseSQLRequest(raw string) (SQLRequest, error) {
	if len(raw) > MaxSQLStatementBytes {
		return SQLRequest{}, fmt.Errorf("SQL request exceeds %d bytes", MaxSQLStatementBytes)
	}
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	var req SQLRequest
	if err := decoder.Decode(&req); err != nil {
		return SQLRequest{}, fmt.Errorf("invalid SQL request: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return SQLRequest{}, errors.New("invalid SQL request: trailing data")
	}
	req.SQL = strings.TrimSpace(req.SQL)
	if req.SQL == "" {
		return SQLRequest{}, errors.New("SQL statement is empty")
	}
	if reason := DeniedSQLReason(req.SQL); reason != "" {
		return SQLRequest{}, errors.New(reason)
	}
	if req.MaxRows < 0 || req.MaxRows > MaxSQLRows {
		return SQLRequest{}, fmt.Errorf("SQL maxRows must be between 0 and %d", MaxSQLRows)
	}
	if len(req.Params) > MaxSQLParams {
		return SQLRequest{}, fmt.Errorf("SQL request has more than %d parameters", MaxSQLParams)
	}
	for i, param := range req.Params {
		value, err := sqlParamValue(param)
		if err != nil {
			return SQLRequest{}, err
		}
		req.Params[i] = value
	}
	return req, nil
}

// sqlParamValue narrows one decoded JSON value to a type SQLite can bind.
func sqlParamValue(param any) (any, error) {
	switch v := param.(type) {
	case nil, bool, string:
		return v, nil
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i, nil
		}
		f, err := v.Float64()
		if err != nil {
			return nil, fmt.Errorf("SQL parameter %q is not a representable number", v.String())
		}
		return f, nil
	default:
		return nil, errors.New("SQL parameters must be null, boolean, number, or string")
	}
}

// SQLRunner executes plugin SQL requests against one plugin's database with the
// host's limits applied. It holds the policy that does not depend on the SQLite
// driver: request validation, parameter typing, the call timeout, transaction
// semantics, and the row, value, and result caps. pluginhost owns the rest, the
// parts that need the driver itself (the authorizer, the SQLite-level limits,
// the per-plugin file, and the connection lifecycle).
//
// The split keeps this half testable without a database file and keeps the
// driver-specific half small enough to audit.
type SQLRunner struct {
	DB *sql.DB
}

// Exec runs the request as one transaction: a multi-statement batch either
// commits whole or leaves the database untouched, and a plugin cannot leave a
// transaction open across calls.
func (r SQLRunner) Exec(ctx context.Context, req SQLRequest) SQLExecResult {
	if r.DB == nil {
		return SQLExecResult{Error: sqlStorageUnavailable}
	}
	ctx, cancel := context.WithTimeout(ctx, SQLCallTimeout)
	defer cancel()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return SQLExecResult{Error: err.Error()}
	}
	result, err := tx.ExecContext(ctx, req.SQL, req.Params...)
	if err != nil {
		_ = tx.Rollback()
		return SQLExecResult{Error: err.Error()}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return SQLExecResult{Error: err.Error()}
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return SQLExecResult{Error: err.Error()}
	}
	if err := tx.Commit(); err != nil {
		return SQLExecResult{Error: err.Error()}
	}
	return SQLExecResult{RowsAffected: rowsAffected, LastInsertID: lastInsertID}
}

// Query runs the request and returns a bounded result set.
//
// MaxRows expresses intent: a plugin that asks for N rows gets at most N and
// learns from Truncated that more matched, which is how queryRow reads one row
// out of a large table without pulling it into the host. A plugin that asks for
// no limit gets an error rather than a silently short answer once the result
// passes MaxSQLRows or MaxSQLResultBytes.
func (r SQLRunner) Query(ctx context.Context, req SQLRequest) SQLQueryResult {
	if r.DB == nil {
		return SQLQueryResult{Error: sqlStorageUnavailable}
	}
	ctx, cancel := context.WithTimeout(ctx, SQLCallTimeout)
	defer cancel()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return SQLQueryResult{Error: err.Error()}
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, req.SQL, req.Params...)
	if err != nil {
		return SQLQueryResult{Error: err.Error()}
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return SQLQueryResult{Error: err.Error()}
	}
	bounded := req.MaxRows > 0
	rowBudget := MaxSQLRows
	if bounded {
		rowBudget = req.MaxRows
	}
	builder, err := newSQLResultBuilder(columns, rowBudget)
	if err != nil {
		return SQLQueryResult{Error: err.Error()}
	}

	for rows.Next() {
		if builder.full() {
			if bounded {
				builder.result.Truncated = true
				break
			}
			return SQLQueryResult{Error: fmt.Sprintf("SQL result exceeds %d rows; add a LIMIT", MaxSQLRows)}
		}
		values, err := scanSQLRow(rows, len(columns))
		if err != nil {
			return SQLQueryResult{Error: err.Error()}
		}
		if err := builder.add(values); err != nil {
			return SQLQueryResult{Error: err.Error()}
		}
	}
	if err := rows.Err(); err != nil {
		return SQLQueryResult{Error: err.Error()}
	}
	if err := rows.Close(); err != nil {
		return SQLQueryResult{Error: err.Error()}
	}
	if err := tx.Commit(); err != nil {
		return SQLQueryResult{Error: err.Error()}
	}
	return builder.result
}

// scanSQLRow reads one row into freshly allocated values. The values become the
// result's own row, so they cannot be reused across iterations.
func scanSQLRow(rows *sql.Rows, columns int) ([]any, error) {
	values := make([]any, columns)
	pointers := make([]any, columns)
	for i := range values {
		pointers[i] = &values[i]
	}
	if err := rows.Scan(pointers...); err != nil {
		return nil, err
	}
	return values, nil
}

// sqlResultBuilder accumulates rows while keeping the JSON the plugin will
// receive inside the host's value, byte, and row budgets.
type sqlResultBuilder struct {
	result    SQLQueryResult
	used      int
	rowBudget int
}

func newSQLResultBuilder(columns []string, rowBudget int) (*sqlResultBuilder, error) {
	encodedColumns, err := json.Marshal(columns)
	if err != nil {
		return nil, err
	}
	return &sqlResultBuilder{
		result: SQLQueryResult{Columns: columns, Rows: make([][]any, 0, min(rowBudget, 64))},
		// The envelope's fixed keys, column names, and optional truncation marker
		// count against the budget so every encoded result stays within the cap.
		used:      len(encodedColumns) + len(`{"columns":,"rows":[]}`) + len(`,"truncated":true`),
		rowBudget: rowBudget,
	}, nil
}

func (b *sqlResultBuilder) full() bool { return len(b.result.Rows) == b.rowBudget }

// add appends one row, or reports why it does not fit.
func (b *sqlResultBuilder) add(values []any) error {
	// Size the values before encoding them: a driver that cannot pin
	// SQLITE_LIMIT_LENGTH could otherwise hand back a single value whose JSON
	// encoding dwarfs the entire result budget.
	estimate := len(`[],`)
	for _, value := range values {
		size := sqlValueSize(value)
		if size > MaxSQLValueBytes {
			return errSQLResultTooLarge("value", MaxSQLValueBytes)
		}
		estimate += size + 1
	}
	if b.used+estimate > MaxSQLResultBytes {
		return errSQLResultTooLarge("result", MaxSQLResultBytes)
	}
	encodedRow, err := json.Marshal(values)
	if err != nil {
		return err
	}
	b.used += len(encodedRow) + 1
	if b.used > MaxSQLResultBytes {
		return errSQLResultTooLarge("result", MaxSQLResultBytes)
	}
	b.result.Rows = append(b.result.Rows, values)
	return nil
}

// sqlValueSize returns the JSON-encoded size of one scanned column value
// without allocating. Strings and blobs dominate; everything else is short.
func sqlValueSize(value any) int {
	switch v := value.(type) {
	case nil:
		return len("null")
	case bool:
		return len("false")
	case string:
		if len(v) > MaxSQLValueBytes {
			return len(v) + 2
		}
		size := 2 // quotes
		for i := 0; i < len(v); {
			r, width := utf8.DecodeRuneInString(v[i:])
			i += width
			if r == utf8.RuneError && width == 1 {
				size += 6
				continue
			}
			switch r {
			case '\\', '"', '\b', '\f', '\n', '\r', '\t':
				size += 2
			case '<', '>', '&', '\u2028', '\u2029':
				size += 6
			default:
				if r < 0x20 {
					size += 6
				} else {
					size += width
				}
			}
		}
		return size
	case []byte:
		// database/sql hands blobs back as []byte, which encoding/json emits as
		// base64: four characters per three bytes, plus quotes.
		return (len(v)+2)/3*4 + 2
	default:
		// Numbers and time values encode to short literals.
		return 32
	}
}
