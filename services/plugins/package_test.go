package plugins

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// findExampleWasm returns the path to a real example wasm we can use to make
// happy-path package tests realistic. Tests that need actual wasm bytes
// t.Skip() if no example has been built yet. Looks in the in-tree examples
// directory first, then falls back to the sibling owncast-plugin-sdk
// checkout (where examples live after the SDK split).
func findExampleWasm(t *testing.T) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	candidates := []string{
		filepath.Join(repoRoot, "examples", "js", "hello-world", "hello-world.wasm"),
		filepath.Join(repoRoot, "..", "owncast-plugin-sdk", "examples", "js", "hello-world", "hello-world.wasm"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	t.Skipf("example wasm not built in any of %v; run npm run build in examples/js/hello-world first", candidates)
	return ""
}

// buildPkg writes a .ocpkg-shaped zip to a temp file and returns the path.
func buildPkg(t *testing.T, entries map[string][]byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.ocpkg")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for name, data := range entries {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("zip create %s: %v", name, err)
		}
		if _, err := w.Write(data); err != nil {
			t.Fatalf("zip write %s: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	return path
}

func validManifestBytes() []byte {
	return []byte(`{
		"api": "1",
		"name": "hello-world",
		"version": "0.1.0",
		"description": "test",
		"permissions": []
	}`)
}

func TestLoadPackage_MissingManifest(t *testing.T) {
	path := buildPkg(t, map[string][]byte{
		pkgWasmFilename: {0x00}, // not parsed; the manifest check fails first
	})
	_, err := LoadPackage(context.Background(), &HostEnv{}, path)
	if err == nil {
		t.Fatal("expected error for missing manifest, got nil")
	}
	if !strings.Contains(err.Error(), pkgManifestFilename) {
		t.Errorf("error mentions manifest filename: got %v", err)
	}
}

func TestLoadPackage_MissingWasm(t *testing.T) {
	path := buildPkg(t, map[string][]byte{
		pkgManifestFilename: validManifestBytes(),
	})
	_, err := LoadPackage(context.Background(), &HostEnv{}, path)
	if err == nil {
		t.Fatal("expected error for missing wasm, got nil")
	}
	if !strings.Contains(err.Error(), pkgWasmFilename) {
		t.Errorf("error mentions wasm filename: got %v", err)
	}
}

func TestLoadPackage_CorruptZip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.ocpkg")
	if err := os.WriteFile(path, []byte("this is not a zip file at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadPackage(context.Background(), &HostEnv{}, path)
	if err == nil {
		t.Fatal("expected error for corrupt zip, got nil")
	}
}

func TestLoadPackage_MalformedManifestJSON(t *testing.T) {
	path := buildPkg(t, map[string][]byte{
		pkgManifestFilename: []byte(`{not valid json`),
		pkgWasmFilename:     {0x00},
	})
	_, err := LoadPackage(context.Background(), &HostEnv{}, path)
	if err == nil {
		t.Fatal("expected error for malformed manifest, got nil")
	}
}

func TestLoadPackage_MissingPathReturnsError(t *testing.T) {
	_, err := LoadPackage(context.Background(), &HostEnv{}, "/nonexistent/path/foo.ocpkg")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadPackage_AssetsExposedAsFS(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		t.Fatalf("read example wasm: %v", err)
	}
	path := buildPkg(t, map[string][]byte{
		pkgManifestFilename:               validManifestBytes(),
		pkgWasmFilename:                   wasmBytes,
		pkgAssetsPrefix + "index.html":    []byte("<h1>hi</h1>"),
		pkgAssetsPrefix + "sub/style.css": []byte("body{}"),
	})

	loaded, err := LoadPackage(context.Background(), &HostEnv{}, path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	defer loaded.Close(context.Background())

	if loaded.AssetsFS == nil {
		t.Fatal("expected AssetsFS to be set when assets/ present in archive")
	}

	got, err := fs.ReadFile(loaded.AssetsFS, "index.html")
	if err != nil {
		t.Fatalf("read assets/index.html: %v", err)
	}
	if !bytes.Equal(got, []byte("<h1>hi</h1>")) {
		t.Errorf("index.html: got %q want %q", got, "<h1>hi</h1>")
	}

	got, err = fs.ReadFile(loaded.AssetsFS, "sub/style.css")
	if err != nil {
		t.Fatalf("read nested asset: %v", err)
	}
	if !bytes.Equal(got, []byte("body{}")) {
		t.Errorf("sub/style.css: got %q", got)
	}
}

func TestLoadPackage_NoAssetsLeavesFSNil(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	path := buildPkg(t, map[string][]byte{
		pkgManifestFilename: validManifestBytes(),
		pkgWasmFilename:     wasmBytes,
	})

	loaded, err := LoadPackage(context.Background(), &HostEnv{}, path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	defer loaded.Close(context.Background())

	if loaded.AssetsFS != nil {
		t.Error("expected AssetsFS to be nil when archive has no assets/")
	}
}

// Verify that arbitrary file content (incl. binary) round-trips through the
// zip reader. Belt-and-suspenders for the in-memory zip approach.
func TestLoadPackage_BinaryAssetsIntact(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	// Random-ish binary payload — a sentinel image-like byte sequence.
	binary := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x01, 0x02, 0x03}
	path := buildPkg(t, map[string][]byte{
		pkgManifestFilename:          validManifestBytes(),
		pkgWasmFilename:              wasmBytes,
		pkgAssetsPrefix + "logo.png": binary,
	})

	loaded, err := LoadPackage(context.Background(), &HostEnv{}, path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	defer loaded.Close(context.Background())

	got, err := fs.ReadFile(loaded.AssetsFS, "logo.png")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !bytes.Equal(got, binary) {
		t.Errorf("binary asset corrupted:\n  want %v\n  got  %v", binary, got)
	}
}

// readZipFile is exercised heavily above, but verify it directly with a tiny
// constructed archive so its error path is independently tested.
func TestReadZipFile_Missing(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("present.txt")
	io.WriteString(w, "hi")
	zw.Close()

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("zip reader: %v", err)
	}
	if _, err := readZipFile(zr, "missing.txt"); err == nil {
		t.Error("expected error for missing entry, got nil")
	}
}

func TestHasPluginInstructions_Package(t *testing.T) {
	with := buildPkg(t, map[string][]byte{
		pkgManifestFilename:        validManifestBytes(),
		pkgWasmFilename:            {0x00},
		PluginInstructionsFilename: []byte("# Hello"),
	})
	if !hasPluginInstructions(with) {
		t.Error("expected hasPluginInstructions=true for package with INSTRUCTIONS.md")
	}

	without := buildPkg(t, map[string][]byte{
		pkgManifestFilename: validManifestBytes(),
		pkgWasmFilename:     {0x00},
	})
	if hasPluginInstructions(without) {
		t.Error("expected hasPluginInstructions=false for package without INSTRUCTIONS.md")
	}
}

func TestHasPluginInstructions_LooseFiles(t *testing.T) {
	dir := t.TempDir()
	wasmPath := filepath.Join(dir, "demo.wasm")
	if err := os.WriteFile(wasmPath, []byte{0x00}, 0o644); err != nil {
		t.Fatal(err)
	}
	if hasPluginInstructions(wasmPath) {
		t.Error("expected false when sibling INSTRUCTIONS.md absent")
	}
	if err := os.WriteFile(filepath.Join(dir, "demo."+PluginInstructionsFilename), []byte("# hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !hasPluginInstructions(wasmPath) {
		t.Error("expected true when sibling <base>.INSTRUCTIONS.md present")
	}
}

func TestInstructionsBytes_Package(t *testing.T) {
	body := []byte("# Setup\n\nDo the thing.")
	path := buildPkg(t, map[string][]byte{
		pkgManifestFilename:        validManifestBytes(),
		pkgWasmFilename:            {0x00},
		PluginInstructionsFilename: body,
	})

	mgr := NewManager(filepath.Dir(path), &HostEnv{})
	if err := mgr.scan(context.Background()); err != nil {
		t.Fatalf("scan: %v", err)
	}

	entries := mgr.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 discovered, got %d", len(entries))
	}
	if !entries[0].HasInstructions {
		t.Error("expected HasInstructions=true in discovered entry")
	}

	got, err := mgr.InstructionsBytes("hello-world")
	if err != nil {
		t.Fatalf("InstructionsBytes: %v", err)
	}
	if !bytes.Equal(got, body) {
		t.Errorf("instructions: got %q want %q", got, body)
	}
}

func TestInstructionsBytes_NoInstructions(t *testing.T) {
	path := buildPkg(t, map[string][]byte{
		pkgManifestFilename: validManifestBytes(),
		pkgWasmFilename:     {0x00},
	})
	mgr := NewManager(filepath.Dir(path), &HostEnv{})
	if err := mgr.scan(context.Background()); err != nil {
		t.Fatalf("scan: %v", err)
	}
	if _, err := mgr.InstructionsBytes("hello-world"); err == nil {
		t.Error("expected error for plugin without instructions")
	}
}

func TestInstructionsBytes_NotDiscovered(t *testing.T) {
	mgr := NewManager(t.TempDir(), &HostEnv{})
	if _, err := mgr.InstructionsBytes("ghost"); err == nil {
		t.Error("expected error for undiscovered plugin")
	}
}
