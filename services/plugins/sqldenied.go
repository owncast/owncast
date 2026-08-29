package plugins

import "strings"

// Owncast refuses a handful of SQL operations for every plugin, and refuses them
// twice over.
//
// The authoritative refusal is the SQLite authorizer pluginhost installs on each
// connection. It is the security boundary, because it sees the compiled
// statement rather than its text.
//
// This file is the second refusal, and it exists for a different reason: every
// host that runs plugins has to reject the same statements, or a plugin's own
// scenario tests are worthless. The SDK's test runner and dev server use a
// pure-Go SQLite driver that has no authorizer API, so without this check a
// plugin could exercise `CREATE TEMP TABLE` locally, watch its tests pass, and
// then fail on a real server. Refusing in Go, above the driver, makes every host
// agree.
//
// A text check is not a boundary and is not treated as one. It runs in front of
// the authorizer, never instead of it.

// DeniedSQLStatementExamples is the cross-host parity fixture: every statement
// here must be refused by Owncast and by the SDK's hosts alike. Both repositories
// assert against this list, so a form discovered later gets fixed everywhere at
// once instead of drifting.
var DeniedSQLStatementExamples = []string{
	// Reaching another database file.
	"ATTACH DATABASE 'other.db' AS other",
	"attach database 'other.db' as other",
	"DETACH DATABASE other",
	"VACUUM",
	"VACUUM INTO 'copy.db'",
	// Reading or changing connection configuration, which is what keeps the
	// page cap a cap.
	"PRAGMA max_page_count=100000000",
	"PRAGMA page_size",
	"pragma   journal_mode = delete",
	"/* comment */ PRAGMA page_size",
	"-- comment\nPRAGMA page_size",
	"SELECT 1; PRAGMA page_size",
	"EXPLAIN PRAGMA page_size",
	"EXPLAIN QUERY PLAN ATTACH DATABASE 'other.db' AS other",
	// Temp-schema objects, which live outside the page cap. The keyword forms
	// and the schema-qualified forms both have to go: SQLite reports the latter
	// as ordinary DDL with "temp" in the authorizer's database-name argument.
	"CREATE TEMP TABLE scratch (v TEXT)",
	"CREATE TEMPORARY TABLE scratch (v TEXT)",
	"CREATE TEMP VIEW scratch AS SELECT 1",
	"DROP TEMP TABLE scratch",
	"CREATE TABLE temp.scratch (v TEXT)",
	`CREATE TABLE "temp".scratch (v TEXT)`,
	"CREATE TABLE TeMp.scratch (v TEXT)",
	"CREATE TABLE [temp].scratch (v TEXT)",
	"CREATE VIEW temp.scratch AS SELECT 1",
	"CREATE INDEX temp.i ON scratch (v)",
	"SELECT load_extension/* comment */('evil.so')",
	"CREATE TABLE temp/* comment */.scratch (v TEXT)",
	"SELECT load_extension\f('evil.so')",
	"CREATE TABLE temp\f.scratch (v TEXT)",
	"INSERT INTO temp.scratch VALUES ('x')",
	"SELECT * FROM temp.scratch",
	"SELECT * FROM sqlite_temp_master",
	"SELECT * FROM sqlite_temp_schema",
	// Loading native code.
	"SELECT load_extension('evil.so')",
	"SELECT LOAD_EXTENSION('evil.so')",
}

const deniedTransactionControl = "transaction control is not available to plugins"

// deniedStatementHeads are keywords a plugin's top-level statement may not
// begin with. Transaction controls are also refused so the host owns every
// call's commit or rollback; CREATE TRIGGER bodies are tracked separately
// because their terminating END follows an internal semicolon.
var deniedStatementHeads = map[string]string{
	"pragma":    "PRAGMA is not available to plugins",
	"attach":    "ATTACH is not available to plugins",
	"detach":    "DETACH is not available to plugins",
	"vacuum":    "VACUUM is not available to plugins",
	"begin":     deniedTransactionControl,
	"commit":    deniedTransactionControl,
	"end":       deniedTransactionControl,
	"rollback":  deniedTransactionControl,
	"savepoint": deniedTransactionControl,
	"release":   deniedTransactionControl,
}

// deniedIdentifiers name the temp schema without the TEMP keyword.
var deniedIdentifiers = map[string]string{
	"sqlite_temp_master": "the temp schema is not available to plugins",
	"sqlite_temp_schema": "the temp schema is not available to plugins",
}

// deniedFunctions may not be called. Unlike the statement heads these appear
// mid-expression, so they are matched as a name followed by an open parenthesis.
var deniedFunctions = map[string]string{
	"load_extension": "load_extension is not available to plugins",
}

const deniedTempSchema = "temporary tables, indexes, views, and triggers are not available to plugins"

// DeniedSQLReason reports why a plugin's SQL may not run, or an empty string
// when nothing in it is refused. The message is what the plugin sees, so it says
// what is unavailable rather than which check tripped.
func DeniedSQLReason(statement string) string {
	var previous string
	explainPrefix := 0
	atStatementStart := true
	var trigger sqlTriggerState
	for _, token := range scanSQLTokens(statement) {
		switch {
		case token.terminator:
			if !trigger.endStatement() {
				continue
			}
			atStatementStart = true
			previous = ""
			explainPrefix = 0
			continue
		case atStatementStart:
			if consumeSQLExplainPrefix(token.text, &explainPrefix) {
				continue
			}
			if reason, ok := deniedStatementHeads[token.text]; ok {
				return reason
			}
			atStatementStart = false
			explainPrefix = 0
		}
		trigger.observe(previous, token)
		if reason := deniedSQLTokenReason(previous, token); reason != "" {
			return reason
		}
		previous = token.text
	}
	return ""
}

func deniedSQLTokenReason(previous string, token sqlToken) string {
	if reason, ok := deniedIdentifiers[token.text]; ok {
		return reason
	}
	if token.call {
		if reason, ok := deniedFunctions[token.text]; ok {
			return reason
		}
	}
	// A schema qualifier is an identifier followed by a dot.
	if token.qualifier && token.text == "temp" {
		return deniedTempSchema
	}
	if (previous == "create" || previous == "drop") && (token.text == "temp" || token.text == "temporary") {
		return deniedTempSchema
	}
	return ""
}

// sqlTriggerState distinguishes semicolons inside a CREATE TRIGGER body from
// top-level statement boundaries. CASE has its own END token, so its nesting
// must reach zero before END can close the trigger body.
type sqlTriggerState struct {
	inTrigger bool
	body      bool
	ended     bool
	caseDepth int
}

func (s *sqlTriggerState) observe(previous string, token sqlToken) {
	if token.quoted {
		return
	}
	if previous == "create" && token.text == "trigger" {
		s.inTrigger = true
	}
	if !s.inTrigger {
		return
	}
	switch token.text {
	case "begin":
		s.body = true
	case "case":
		if s.body {
			s.caseDepth++
		}
	case "end":
		if s.caseDepth > 0 {
			s.caseDepth--
		} else if s.body {
			s.ended = true
		}
	}
}

func (s *sqlTriggerState) endStatement() bool {
	if s.inTrigger && !s.ended {
		return false
	}
	*s = sqlTriggerState{}
	return true
}

// consumeSQLExplainPrefix keeps the following statement head visible to the
// denial rules while skipping SQLite's optional EXPLAIN QUERY PLAN prefix.
func consumeSQLExplainPrefix(token string, state *int) bool {
	switch {
	case *state == 0 && token == "explain":
		*state = 1
	case *state == 1 && token == "query":
		*state = 2
	case *state == 2 && token == "plan":
		*state = 3
	default:
		return false
	}
	return true
}

// sqlToken is one lexical item of interest: a lower-cased keyword or identifier,
// a statement terminator, or an identifier used as a schema qualifier.
type sqlToken struct {
	text       string
	terminator bool
	qualifier  bool
	call       bool
	quoted     bool
}

// scanSQLTokens reduces SQL to the tokens DeniedSQLReason needs to see. It skips
// string literals and comments so their contents cannot trip a rule, and
// unquotes delimited identifiers so `"temp"` and `[temp]` read the same as
// `temp`. It is not a SQL parser and does not try to be: it recognizes
// statement boundaries, keywords, identifiers, and schema qualifiers, and
// ignores everything else.
func scanSQLTokens(statement string) []sqlToken {
	var tokens []sqlToken
	runes := []rune(statement)
	for i := 0; i < len(runes); {
		switch c := runes[i]; {
		case c == '-' && i+1 < len(runes) && runes[i+1] == '-':
			i = skipSQLLineComment(runes, i)
		case c == '/' && i+1 < len(runes) && runes[i+1] == '*':
			i = skipSQLBlockComment(runes, i)
		case c == '\'':
			i = skipSQLDelimited(runes, i, '\'')
		case c == '"', c == '[', c == '`':
			end := skipSQLDelimited(runes, i, sqlIdentifierCloser(c))
			// Read the identifier without its delimiters, so `"temp"` and
			// `[temp]` are the same token as `temp`.
			text := strings.ToLower(strings.TrimSpace(string(runes[i+1 : max(end-1, i+1)])))
			tokens = append(tokens, newSQLToken(text, runes, end, true))
			i = end
		case c == ';':
			tokens = append(tokens, sqlToken{terminator: true})
			i++
		case isSQLWordRune(c):
			start := i
			for i < len(runes) && isSQLWordRune(runes[i]) {
				i++
			}
			tokens = append(tokens, newSQLToken(strings.ToLower(string(runes[start:i])), runes, i, false))
		default:
			i++
		}
	}
	return tokens
}

func sqlIdentifierCloser(opener rune) rune {
	if opener == '[' {
		return ']'
	}
	return opener
}

func skipSQLLineComment(runes []rune, i int) int {
	for ; i < len(runes) && runes[i] != '\n'; i++ {
	}
	return i
}

func skipSQLBlockComment(runes []rune, i int) int {
	for i += 2; i+1 < len(runes); i++ {
		if runes[i] == '*' && runes[i+1] == '/' {
			return i + 2
		}
	}
	return len(runes)
}

// skipSQLDelimited returns the index just past the delimited run beginning at
// start, treating a doubled closer as an escape the way SQLite does.
func skipSQLDelimited(runes []rune, start int, closer rune) int {
	for i := start + 1; i < len(runes); i++ {
		if runes[i] != closer {
			continue
		}
		if i+1 < len(runes) && runes[i+1] == closer {
			i++
			continue
		}
		return i + 1
	}
	return len(runes)
}

// newSQLToken records what follows a word: a dot makes it a schema or table
// qualifier, an open parenthesis makes it a function call.
func newSQLToken(text string, runes []rune, end int, quoted bool) sqlToken {
	next := nextSQLRune(runes, end)
	return sqlToken{text: text, qualifier: next == '.', call: next == '(', quoted: quoted}
}

// nextSQLRune returns the next rune that is not whitespace or part of a
// comment, or 0 at the end. SQLite treats comments as whitespace, including
// between a function name and its parenthesis or a schema and its dot.
func nextSQLRune(runes []rune, i int) rune {
	for i < len(runes) {
		switch {
		case runes[i] == '-' && i+1 < len(runes) && runes[i+1] == '-':
			i = skipSQLLineComment(runes, i)
		case runes[i] == '/' && i+1 < len(runes) && runes[i+1] == '*':
			i = skipSQLBlockComment(runes, i)
		case strings.ContainsRune(" \t\n\r\f", runes[i]):
			i++
		default:
			return runes[i]
		}
	}
	return 0
}

func isSQLWordRune(c rune) bool {
	return c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
