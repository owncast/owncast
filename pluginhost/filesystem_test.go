package pluginhost

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/owncast/owncast/services/plugins"
)

// resolvePluginSandboxPath is the security boundary for the storage.fs
// host functions: every plugin-supplied path must resolve to somewhere
// inside that plugin's own sandbox, no matter how the plugin tries to
// climb out. These tests pin that contract.

func TestResolvePluginSandboxPath_StaysInsideSandbox(t *testing.T) {
	root := t.TempDir()
	sandbox := filepath.Join(root, "my-plugin")

	cases := []struct {
		name string
		rel  string
		want string // path relative to the sandbox, after resolution
	}{
		{"simple file", "notes.txt", "notes.txt"},
		{"nested file", "cache/today.json", "cache/today.json"},
		{"empty path is the sandbox root", "", "."},
		{"dot is the sandbox root", ".", "."},
		{"leading slash is rooted at the sandbox", "/etc/passwd", "etc/passwd"},
		{"single parent collapses inside", "../secret", "secret"},
		{"many parents collapse inside", "../../../../etc/passwd", "etc/passwd"},
		{"interior parent is normalized", "a/b/../c", "a/c"},
		{"sneaky mixed traversal", "foo/../../bar", "bar"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolvePluginSandboxPath(root, "my-plugin", tc.rel)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			want := filepath.Join(sandbox, tc.want)
			if got != want {
				t.Fatalf("resolvePluginSandboxPath(%q) = %q, want %q", tc.rel, got, want)
			}
			if got != sandbox && !strings.HasPrefix(got, sandbox+string(os.PathSeparator)) {
				t.Fatalf("resolved path %q escaped sandbox %q", got, sandbox)
			}
		})
	}
}

func TestResolvePluginSandboxPath_IsolatesPlugins(t *testing.T) {
	root := t.TempDir()
	// One plugin must never resolve a path into another plugin's sandbox,
	// even by spelling the sibling's name through a traversal sequence.
	got, err := resolvePluginSandboxPath(root, "plugin-a", "../plugin-b/data.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sibling := filepath.Join(root, "plugin-b")
	if got == sibling || strings.HasPrefix(got, sibling+string(os.PathSeparator)) {
		t.Fatalf("plugin-a resolved into plugin-b's sandbox: %q", got)
	}
	if want := filepath.Join(root, "plugin-a", "plugin-b", "data.json"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// TestFilesystemHostFns_RoundTrip exercises the wired closures end to end
// against a sandbox rooted at a temp dir, confirming write/read/list/
// exists/delete agree and that a traversal attempt is refused.
func TestFilesystemHostFns_RoundTrip(t *testing.T) {
	root := t.TempDir()
	env := &plugins.HostEnv{}
	wireFilesystemHostFnsWithRoot(env, root)

	const slug = "demo"
	if err := env.FSWrite(slug, "sub/greeting.txt", []byte("hello")); err != nil {
		t.Fatalf("FSWrite: %v", err)
	}

	// The bytes landed inside the plugin's sandbox, nowhere else.
	onDisk := filepath.Join(root, slug, "sub", "greeting.txt")
	if data, err := os.ReadFile(onDisk); err != nil || string(data) != "hello" {
		t.Fatalf("expected file at %q with %q, got data=%q err=%v", onDisk, "hello", data, err)
	}

	if data, err := env.FSRead(slug, "sub/greeting.txt"); err != nil || string(data) != "hello" {
		t.Fatalf("FSRead = %q, %v; want %q", data, err, "hello")
	}

	names, err := env.FSList(slug, "sub")
	if err != nil {
		t.Fatalf("FSList: %v", err)
	}
	if len(names) != 1 || names[0] != "greeting.txt" {
		t.Fatalf("FSList = %v, want [greeting.txt]", names)
	}

	if ok, err := env.FSExists(slug, "sub/greeting.txt"); err != nil || !ok {
		t.Fatalf("FSExists = %v, %v; want true", ok, err)
	}
	if ok, err := env.FSExists(slug, "sub/missing.txt"); err != nil || ok {
		t.Fatalf("FSExists(missing) = %v, %v; want false", ok, err)
	}

	if err := env.FSDelete(slug, "sub/greeting.txt"); err != nil {
		t.Fatalf("FSDelete: %v", err)
	}
	if ok, _ := env.FSExists(slug, "sub/greeting.txt"); ok {
		t.Fatalf("file still exists after delete")
	}

	// Listing a directory that doesn't exist is empty, not an error.
	if names, err := env.FSList(slug, "never-created"); err != nil || len(names) != 0 {
		t.Fatalf("FSList(missing dir) = %v, %v; want empty, nil", names, err)
	}

	// A traversal attempt is neutralized, not leaked: "../escape.txt"
	// collapses to a file inside the plugin's own sandbox, and nothing is
	// written to the parent (sibling-plugin / data-dir) level.
	if err := env.FSWrite(slug, "../escape.txt", []byte("x")); err != nil {
		t.Fatalf("FSWrite (traversal) should be contained, not error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "escape.txt")); !os.IsNotExist(err) {
		t.Fatalf("traversal write leaked a file outside the sandbox")
	}
	if _, err := os.Stat(filepath.Join(root, slug, "escape.txt")); err != nil {
		t.Fatalf("traversal write should have landed inside the sandbox: %v", err)
	}
}

func TestFilesystemHostFns_RejectsOversizedWrite(t *testing.T) {
	root := t.TempDir()
	env := &plugins.HostEnv{}
	wireFilesystemHostFnsWithRoot(env, root)

	tooBig := make([]byte, maxPluginFileBytes+1)
	if err := env.FSWrite("demo", "big.bin", tooBig); err == nil {
		t.Fatalf("expected oversized write to be rejected")
	}
}
