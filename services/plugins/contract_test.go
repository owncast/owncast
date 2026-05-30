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

// This guards against this vendored copy of the plugin runtime silently
// drifting from the upstream SDK's wire contract. services/plugins is a copy
// of github.com/owncast/plugin-sdk's host-runtime/plugin package; the
// implementation is allowed to fork for Owncast integration, but the *wire
// contract* (host-function names, permission identifiers, serialized type
// shapes) must stay identical, or plugins and the TypeScript SDK break.
//
// plugin-contract.json is the SDK's published contract, copied in verbatim.
// This test re-derives the contract from THIS repo's hostfns.go and compares.
// If it fails, this copy has fallen behind (or ahead of) the SDK — re-sync the
// runtime and copy the SDK's host-runtime/plugin/plugin-contract.json here.
//
// The extractor below is a verbatim copy of the SDK's
// host-runtime/plugin/contract_test.go extractor; keep it identical so the
// produced contract is comparable byte-for-byte.

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

func TestPluginContractMatchesSDK(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")

	src, err := os.ReadFile(filepath.Join(repoRoot, "services/plugins/hostfns.go"))
	if err != nil {
		t.Fatalf("read hostfns.go: %v", err)
	}
	gotJSON := marshalContract(t, buildContractFromSource(t, string(src)))

	want, err := os.ReadFile(filepath.Join(repoRoot, "services/plugins/plugin-contract.json"))
	if err != nil {
		t.Fatalf("read plugin-contract.json: %v", err)
	}

	if string(gotJSON) != string(want) {
		t.Errorf(`Owncast's bundled plugin runtime is out of sync with the plugin SDK.

The plugin host functions, permissions, or data shapes in services/plugins no
longer match what the SDK (github.com/owncast/plugin-sdk) publishes, so plugins
built against the SDK may not work correctly against this build of Owncast.

This usually means a change landed in the SDK but wasn't copied here yet (or a
change was made here without updating the SDK). To fix:

  1. Make the change in the SDK first and regenerate its snapshot there:
       UPDATE_CONTRACT=1 go test ./plugin/ -run TestPluginContract
  2. Copy the SDK's host-runtime/plugin/plugin-contract.json over
     services/plugins/plugin-contract.json here, and bring
     services/plugins/hostfns.go in line with the SDK.
  3. Re-run this test — it passes once the two sides match.`)
	}
}
