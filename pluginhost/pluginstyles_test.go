package pluginhost

import (
	"bytes"
	"strings"
	"testing"

	plugins "github.com/owncast/owncast/services/plugins"
)

// writeWrappedScript wraps each plugin's JS in a try/catch so one plugin's
// runtime error can't break the shared /customjavascript bundle.
func TestWriteWrappedScript(t *testing.T) {
	var buf bytes.Buffer
	writeWrappedScript(&buf, "demo", "client.js", []byte("doThing();"))
	out := buf.String()

	for _, want := range []string{
		"// plugin: demo — client.js",
		"try {",
		"doThing();",
		`} catch (e) { console.error("owncast plugin demo script error:", e); }`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("wrapped script missing %q\n got:\n%s", want, out)
		}
	}

	// A contribution without a trailing newline still gets one before the
	// catch, so the `}` doesn't land on the plugin's last line.
	if !strings.Contains(out, "doThing();\n} catch") {
		t.Errorf("expected newline inserted before catch; got:\n%s", out)
	}
}

func TestDeclaredThemeVars(t *testing.T) {
	css := []byte(`:root {
		--theme-color-action: #33d17a;
		--theme-color-background-main:#0b1020;
		--theme-rounded-corners: 8px;
		--not-a-theme-var: red;
	}
	.foo { color: var(--theme-color-action); } /* a use, not a decl */
	:root { --theme-color-action: #fff; } /* duplicate decl */`)

	got := declaredThemeVars(css)
	want := []string{
		"theme-color-action",
		"theme-color-background-main",
		"theme-rounded-corners",
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		// declaredThemeVars sorts its output, so positions are stable.
		if got[i] != want[i] {
			t.Errorf("var %d: got %q, want %q", i, got[i], want[i])
		}
	}

	// CSS with no theme declarations (only a var() use) yields nil so the
	// admin UI shows no swatch badges.
	if vars := declaredThemeVars([]byte(`.x { color: var(--theme-color-action); }`)); vars != nil {
		t.Errorf("expected nil for a var() use with no declaration, got %v", vars)
	}
}

func TestManifestHasPermission(t *testing.T) {
	m := &plugins.Manifest{Permissions: []string{"storage.kv", "ui.modify"}}
	if !manifestHasPermission(m, plugins.PermUIModify) {
		t.Error("expected ui.modify to be reported as granted")
	}
	if manifestHasPermission(m, "http.serve") {
		t.Error("did not expect http.serve to be reported as granted")
	}
	if manifestHasPermission(&plugins.Manifest{}, plugins.PermUIModify) {
		t.Error("empty permission list should grant nothing")
	}
}
