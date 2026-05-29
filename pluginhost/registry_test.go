package pluginhost

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/owncast/owncast/services/plugins"
)

// Tests covering /api/admin/plugin-registry/* — the in-Owncast browse
// + install endpoints. Each test that needs the install path uses a
// real Manager pointing at a temp plugins directory and a fake
// registry HTTP server (httptest.NewServer) so the SHA256 + download
// + Manager.Install chain runs end to end against trusted bytes.

// findExampleWasm returns the bundled hello-world example wasm path,
// skipping if it isn't built. The example lives in the sibling
// owncast-plugin-sdk repo since the SDK split.
func findExampleWasm(t *testing.T) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	candidates := []string{
		filepath.Join(repoRoot, "owncast-plugin-sdk", "examples", "js", "hello-world", "hello-world.wasm"),
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

// buildPackageBytes assembles an in-memory .ocpkg with the given
// manifest + wasm + optional assets. Mirrors the host's accepted
// shape (plugin.manifest.json + plugin.wasm at the zip root).
func buildPackageBytes(t *testing.T, manifest []byte, wasm []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	mw, err := zw.Create("plugin.manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mw.Write(manifest); err != nil {
		t.Fatal(err)
	}
	ww, err := zw.Create("plugin.wasm")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ww.Write(wasm); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func helloWorldManifest() []byte {
	return []byte(`{
		"api": "1",
		"name": "hello-world",
		"version": "0.1.0",
		"description": "hello-world plugin for registry tests",
		"permissions": []
	}`)
}

// newTestHost spins up a Manager rooted at a fresh temp dir and wraps
// it in a Host with only the fields the registry handlers touch. The
// rest of the Host's deps (kv, sse, etc.) stay nil because the
// browse/install paths don't reach them.
func newTestHost(t *testing.T) *Host {
	t.Helper()
	dir := t.TempDir()
	mgr := plugins.NewManager(dir, &plugins.HostEnv{})
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("manager start: %v", err)
	}
	t.Cleanup(func() { mgr.Stop(context.Background()) })
	return &Host{manager: mgr}
}

// withRegistryEnv sets OWNCAST_PLUGIN_REGISTRY for the duration of a
// test and restores it afterward. Hides the t.Setenv ceremony so
// individual tests stay short.
func withRegistryEnv(t *testing.T, url string) {
	t.Helper()
	t.Setenv("OWNCAST_PLUGIN_REGISTRY", url)
}

// --- pure unit tests ---

func TestFindVersion(t *testing.T) {
	detail := &registryDetail{
		Name: "demo",
		Versions: []registryVersion{
			{Version: "0.2.0", SHA256: "b", DownloadURL: "u-b"},
			{Version: "0.1.0", SHA256: "a", DownloadURL: "u-a"},
		},
	}
	if v := findVersion(detail, "0.1.0"); v == nil || v.SHA256 != "a" {
		t.Errorf("findVersion(0.1.0) = %v, want SHA=a", v)
	}
	if v := findVersion(detail, "0.2.0"); v == nil || v.SHA256 != "b" {
		t.Errorf("findVersion(0.2.0) = %v, want SHA=b", v)
	}
	if v := findVersion(detail, "9.9.9"); v != nil {
		t.Errorf("findVersion of unknown version should be nil, got %v", v)
	}
}

func TestRegistryBase_TrimsTrailingSlash(t *testing.T) {
	t.Setenv("OWNCAST_PLUGIN_REGISTRY", "https://owncast.directory/")
	if got := registryBase(); got != "https://owncast.directory" {
		t.Errorf("registryBase trimmed = %q, want %q", got, "https://owncast.directory")
	}
}

func TestRegistryBase_EmptyWhenUnset(t *testing.T) {
	t.Setenv("OWNCAST_PLUGIN_REGISTRY", "")
	if got := registryBase(); got != "" {
		t.Errorf("registryBase unset = %q, want empty", got)
	}
}

// --- /api/admin/plugin-registry/list ---

func TestHandleRegistryList_EmptyWhenUnconfigured(t *testing.T) {
	withRegistryEnv(t, "")
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	host.handleRegistryList(rec, httptest.NewRequest(http.MethodGet, "/api/admin/plugin-registry/list", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	var body []any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v (body=%q)", err, rec.Body.String())
	}
	if len(body) != 0 {
		t.Errorf("body = %v, want empty array", body)
	}
}

func TestHandleRegistryList_ProxiesUpstream(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/plugins" {
			t.Errorf("upstream path = %q, want /api/plugins", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"demo","summary":"hi"}]`))
	}))
	t.Cleanup(upstream.Close)

	withRegistryEnv(t, upstream.URL)
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	host.handleRegistryList(rec, httptest.NewRequest(http.MethodGet, "/api/admin/plugin-registry/list", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != `[{"name":"demo","summary":"hi"}]` {
		t.Errorf("body = %q", got)
	}
}

func TestHandleRegistryList_BadGatewayOnUpstreamDown(t *testing.T) {
	// Pointing at a port that's certainly closed (port 1 is reserved).
	withRegistryEnv(t, "http://127.0.0.1:1")
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	host.handleRegistryList(rec, httptest.NewRequest(http.MethodGet, "/api/admin/plugin-registry/list", nil))

	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want 502", rec.Code)
	}
}

// --- /api/admin/plugin-registry/install ---

func TestHandleRegistryInstall_RejectsNonPost(t *testing.T) {
	host := newTestHost(t)
	rec := httptest.NewRecorder()
	host.handleRegistryInstall(rec, httptest.NewRequest(http.MethodGet, "/api/admin/plugin-registry/install", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}

func TestHandleRegistryInstall_RejectsMissingFields(t *testing.T) {
	host := newTestHost(t)
	for _, body := range []string{`{}`, `{"name":""}`, `{"version":"0.1.0"}`} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install", strings.NewReader(body))
		host.handleRegistryInstall(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("body=%s: status = %d, want 400", body, rec.Code)
		}
	}
}

func TestHandleRegistryInstall_ServiceUnavailableWhenUnconfigured(t *testing.T) {
	withRegistryEnv(t, "")
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
		strings.NewReader(`{"name":"demo","version":"0.1.0"}`))
	host.handleRegistryInstall(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", rec.Code)
	}
}

func TestHandleRegistryInstall_VersionNotFound(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Registry detail returns versions, but not the one we asked for.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"demo","versions":[{"version":"0.2.0","sha256":"x","downloadURL":"u"}]}`))
	}))
	t.Cleanup(upstream.Close)

	withRegistryEnv(t, upstream.URL)
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
		strings.NewReader(`{"name":"demo","version":"0.1.0"}`))
	host.handleRegistryInstall(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404 (body=%s)", rec.Code, rec.Body.String())
	}
}

func TestHandleRegistryInstall_SHA256Mismatch(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	pkg := buildPackageBytes(t, helloWorldManifest(), wasmBytes)

	upstream := newRegistryStub(t, "hello-world", "0.1.0", pkg, "deadbeef-wrong-hash")
	withRegistryEnv(t, upstream.URL)
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
		strings.NewReader(`{"name":"hello-world","version":"0.1.0"}`))
	host.handleRegistryInstall(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want 502 (body=%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "hash") {
		t.Errorf("error should mention hash mismatch, got %s", rec.Body.String())
	}
}

func TestHandleRegistryInstall_Success(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	pkg := buildPackageBytes(t, helloWorldManifest(), wasmBytes)
	sum := sha256.Sum256(pkg)
	sha := hex.EncodeToString(sum[:])

	upstream := newRegistryStub(t, "hello-world", "0.1.0", pkg, sha)
	withRegistryEnv(t, upstream.URL)
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
		strings.NewReader(`{"name":"hello-world","version":"0.1.0"}`))
	host.handleRegistryInstall(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	// The Manager.Install path drops the .ocpkg into the manager's
	// plugins directory. Confirm by walking List() rather than poking
	// at the filesystem so we don't couple to the path layout.
	entries := host.manager.List()
	if len(entries) != 1 {
		t.Fatalf("manager has %d entries, want 1", len(entries))
	}
	if entries[0].Name != "hello-world" {
		t.Errorf("entry name = %q, want hello-world", entries[0].Name)
	}
	if entries[0].Version != "0.1.0" {
		t.Errorf("entry version = %q, want 0.1.0", entries[0].Version)
	}
}

// --- /api/admin/plugin-registry/<unknown> ---

func TestHandleRegistryRoute_DispatchesActions(t *testing.T) {
	withRegistryEnv(t, "")
	host := newTestHost(t)

	cases := []struct {
		path       string
		method     string
		wantStatus int
	}{
		{"/api/admin/plugin-registry/list", http.MethodGet, http.StatusOK},
		{"/api/admin/plugin-registry/list/", http.MethodGet, http.StatusOK}, // trailing-slash variant
		{"/api/admin/plugin-registry/install", http.MethodPost, http.StatusServiceUnavailable},
		{"/api/admin/plugin-registry/unknown", http.MethodGet, http.StatusNotFound},
	}
	for _, tc := range cases {
		var body io.Reader
		if tc.method == http.MethodPost {
			body = strings.NewReader(`{"name":"x","version":"y"}`)
		}
		rec := httptest.NewRecorder()
		host.handleRegistryRoute(rec, httptest.NewRequest(tc.method, tc.path, body))
		if rec.Code != tc.wantStatus {
			t.Errorf("%s %s: status = %d, want %d", tc.method, tc.path, rec.Code, tc.wantStatus)
		}
	}
}

// --- test helpers ---

// newRegistryStub returns an httptest server that emulates the
// directory's /api/plugins/<name> detail endpoint plus a download URL
// the install handler will GET to retrieve the .ocpkg bytes. The
// stubbed sha256 is whatever the caller passes, so tests can force a
// mismatch by passing a wrong digest.
//
// The detail handler reflects the incoming request's Host header back
// into the downloadURL so we don't have to know httptest's port in
// advance.
func newRegistryStub(t *testing.T, name, version string, pkg []byte, sha string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/plugins/"+name, func(w http.ResponseWriter, r *http.Request) {
		downloadURL := fmt.Sprintf("http://%s/ocpkg/%s", r.Host, name)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"name": name,
			"versions": []map[string]any{
				{
					"version":     version,
					"sha256":      sha,
					"downloadURL": downloadURL,
				},
			},
		})
	})
	mux.HandleFunc("/ocpkg/"+name, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write(pkg)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}
