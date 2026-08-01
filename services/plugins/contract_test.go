package plugins

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// This guards the plugin wire contract against accidental drift. hostfns.go is
// the Go source of truth for host-function names, permission identifiers, and
// serialized type shapes. The TypeScript SDK and wire-protocol documentation
// mirror it.
//
// plugin-contract.json is the generated fingerprint. Regenerate it only for an
// intentional contract change by running this test with UPDATE_CONTRACT=1.

// hostFnNameRe matches an actual host-function name exactly — not the error
// strings ("owncast_x from %s: ...") that also begin with the prefix.
var hostFnNameRe = regexp.MustCompile(`^owncast_[a-z_]+$`)

type contract struct {
	Permissions   map[string]string            `json:"permissions"`
	HostFunctions []string                     `json:"hostFunctions"`
	WireTypes     map[string]map[string]string `json:"wireTypes"`
}

func buildContractFromSource(t *testing.T, src string) contract {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hostfns.go", src, 0)
	if err != nil {
		t.Fatalf("parse hostfns.go: %v", err)
	}

	c := contract{
		Permissions: map[string]string{},
		WireTypes:   map[string]map[string]string{},
	}

	fnSet := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING {
			if s, err := strconv.Unquote(lit.Value); err == nil && hostFnNameRe.MatchString(s) {
				fnSet[s] = true
			}
		}
		return true
	})
	for name := range fnSet {
		c.HostFunctions = append(c.HostFunctions, name)
	}
	sort.Strings(c.HostFunctions)

	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch gd.Tok {
		case token.CONST:
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for i, name := range vs.Names {
					if !strings.HasPrefix(name.Name, "Perm") || i >= len(vs.Values) {
						continue
					}
					if lit, ok := vs.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						if v, err := strconv.Unquote(lit.Value); err == nil {
							c.Permissions[name.Name] = v
						}
					}
				}
			}
		case token.TYPE:
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				fields := map[string]string{}
				for _, field := range st.Fields.List {
					tag := jsonTagName(field.Tag)
					if tag == "" || tag == "-" {
						continue
					}
					fields[tag] = exprTypeString(field.Type)
				}
				if len(fields) > 0 {
					c.WireTypes[ts.Name.Name] = fields
				}
			}
		}
	}
	return c
}

func jsonTagName(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	raw, err := strconv.Unquote(tag.Value)
	if err != nil {
		return ""
	}
	for _, part := range strings.Fields(raw) {
		if strings.HasPrefix(part, "json:") {
			v, err := strconv.Unquote(strings.TrimPrefix(part, "json:"))
			if err != nil {
				return ""
			}
			return strings.Split(v, ",")[0]
		}
	}
	return ""
}

func exprTypeString(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprTypeString(t.X)
	case *ast.ArrayType:
		return "[]" + exprTypeString(t.Elt)
	case *ast.SelectorExpr:
		return exprTypeString(t.X) + "." + t.Sel.Name
	case *ast.MapType:
		return "map[" + exprTypeString(t.Key) + "]" + exprTypeString(t.Value)
	default:
		return "?"
	}
}

func marshalContract(t *testing.T, c contract) []byte {
	t.Helper()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		t.Fatalf("marshal contract: %v", err)
	}
	return append(data, '\n')
}

func normalizeNewlines(b []byte) []byte {
	return []byte(strings.ReplaceAll(string(b), "\r\n", "\n"))
}

func TestOperationResultsUseErrorAsTheOnlyFailureSignal(t *testing.T) {
	tests := map[string]struct {
		result any
		want   string
	}{
		"filesystem write success":   {FSResult{}, `{}`},
		"filesystem write failure":   {FSResult{Error: "failed"}, `{"error":"failed"}`},
		"filesystem delete success":  {FSResult{}, `{}`},
		"filesystem delete failure":  {FSResult{Error: "failed"}, `{"error":"failed"}`},
		"SQL exec success":           {SQLExecResult{RowsAffected: 1, LastInsertID: 2}, `{"rowsAffected":1,"lastInsertId":2}`},
		"SQL exec failure":           {SQLExecResult{Error: "failed"}, `{"error":"failed","rowsAffected":0,"lastInsertId":0}`},
		"SQL query success":          {SQLQueryResult{Columns: []string{}, Rows: [][]any{}}, `{"columns":[],"rows":[]}`},
		"SQL query failure":          {SQLQueryResult{Error: "failed"}, `{"error":"failed","columns":null,"rows":null}`},
		"video config write success": {VideoConfigWriteResult{}, `{}`},
		"video config write failure": {VideoConfigWriteResult{Error: "failed"}, `{"error":"failed"}`},
		"auth session success":       {GrantSessionResult{}, `{}`},
		"auth session failure":       {GrantSessionResult{Error: "failed"}, `{"error":"failed"}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(test.result)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != test.want {
				t.Fatalf("got %s, want %s", got, test.want)
			}
		})
	}
}

func TestPluginContractMatchesSDK(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")

	src, err := os.ReadFile(filepath.Join(repoRoot, "services/plugins/hostfns.go"))
	if err != nil {
		t.Fatalf("read hostfns.go: %v", err)
	}
	gotJSON := marshalContract(t, buildContractFromSource(t, string(src)))
	contractPath := filepath.Join(repoRoot, "services/plugins/plugin-contract.json")
	if os.Getenv("UPDATE_CONTRACT") == "1" {
		if err := os.WriteFile(contractPath, gotJSON, 0o644); err != nil {
			t.Fatalf("update plugin-contract.json: %v", err)
		}
		return
	}

	want, err := os.ReadFile(contractPath)
	if err != nil {
		t.Fatalf("read plugin-contract.json: %v", err)
	}

	if string(normalizeNewlines(gotJSON)) != string(normalizeNewlines(want)) {
		t.Errorf(`The plugin wire contract snapshot is stale.

The host functions, permissions, or data shapes derived from
services/plugins/hostfns.go no longer match plugin-contract.json.

If the change is intentional, regenerate the snapshot:

  UPDATE_CONTRACT=1 go test ./services/plugins/ -run TestPluginContractMatchesSDK

Then rerun this test without UPDATE_CONTRACT.`)
	}
}
