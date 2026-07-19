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
		Slug: "demo",
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

func TestRegistryBase_UsesDefaultWhenUnset(t *testing.T) {
	// Unset OWNCAST_PLUGIN_REGISTRY falls through to the public catalog
	// so every Owncast instance gets a working Browse tab out of the
	// box without per-deployment configuration.
	t.Setenv("OWNCAST_PLUGIN_REGISTRY", "")
	if got := registryBase(); got != DefaultPluginRegistry {
		t.Errorf("registryBase unset = %q, want %q", got, DefaultPluginRegistry)
	}
}

// --- /api/admin/plugin-registry/list ---

func TestHandleRegistryList_ProxiesUpstream(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/plugins" {
			t.Errorf("upstream path = %q, want /plugins", r.URL.Path)
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
	for _, body := range []string{`{}`, `{"slug":""}`, `{"version":"0.1.0"}`} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install", strings.NewReader(body))
		host.handleRegistryInstall(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("body=%s: status = %d, want 400", body, rec.Code)
		}
	}
}

func TestHandleRegistryInstall_VersionNotFound(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Registry detail returns versions, but not the one we asked for.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"slug":"demo","versions":[{"version":"0.2.0","sha256":"x","downloadURL":"u"}]}`))
	}))
	t.Cleanup(upstream.Close)

	withRegistryEnv(t, upstream.URL)
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
		strings.NewReader(`{"slug":"demo","version":"0.1.0"}`))
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
		strings.NewReader(`{"slug":"hello-world","version":"0.1.0"}`))
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
		strings.NewReader(`{"slug":"hello-world","version":"0.1.0"}`))
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
	if entries[0].Slug != "hello-world" {
		t.Errorf("entry slug = %q, want hello-world", entries[0].Slug)
	}
	if entries[0].Version != "0.1.0" {
		t.Errorf("entry version = %q, want 0.1.0", entries[0].Version)
	}
}

func TestHandleRegistryInstall_RejectsUnloadablePackage(t *testing.T) {
	wasmPath := findExampleWasm(t)
	wasmBytes, _ := os.ReadFile(wasmPath)
	// Package the hello-world wasm under a slug it does NOT bake into its
	// register() output, so the manifest/runtime agreement check fails (the
	// wasm reports slug "hello-world"; the package claims "ghost-plugin").
	pkg := buildPackageBytes(t, []byte(`{
		"api": "1",
		"name": "Ghost Plugin",
		"slug": "ghost-plugin",
		"version": "9.9.9",
		"description": "mismatched manifest/runtime for registry install",
		"permissions": []
	}`), wasmBytes)
	sum := sha256.Sum256(pkg)
	sha := hex.EncodeToString(sum[:])

	upstream := newRegistryStub(t, "ghost-plugin", "9.9.9", pkg, sha)
	withRegistryEnv(t, upstream.URL)
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
		strings.NewReader(`{"slug":"ghost-plugin","version":"9.9.9"}`))
	host.handleRegistryInstall(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body=%s)", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "plugin cannot be installed") {
		t.Errorf("error should say plugin cannot be installed, got %s", body)
	}
	if !strings.Contains(body, "manifest/runtime mismatch") {
		t.Errorf("error should surface the load problem, got %s", body)
	}
	if len(host.manager.List()) != 0 {
		t.Fatalf("unloadable package should not be installed")
	}
}

// --- /api/admin/plugin-registry/<unknown> ---

func TestHandleRegistryRoute_DispatchesActions(t *testing.T) {
	// Point at a stub that 404s the install detail lookup. That keeps
	// the install case as a clean "registry reachable but rejects" path
	// without depending on the (removed) unconfigured-503 behavior.
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/plugins" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`))
			return
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(upstream.Close)
	withRegistryEnv(t, upstream.URL)
	host := newTestHost(t)

	cases := []struct {
		path       string
		method     string
		wantStatus int
	}{
		{"/api/admin/plugin-registry/list", http.MethodGet, http.StatusOK},
		{"/api/admin/plugin-registry/list/", http.MethodGet, http.StatusOK}, // trailing-slash variant
		{"/api/admin/plugin-registry/install", http.MethodPost, http.StatusBadGateway},
		{"/api/admin/plugin-registry/unknown", http.MethodGet, http.StatusNotFound},
	}
	for _, tc := range cases {
		var body io.Reader
		if tc.method == http.MethodPost {
			body = strings.NewReader(`{"slug":"x","version":"y"}`)
		}
		rec := httptest.NewRecorder()
		host.handleRegistryRoute(rec, httptest.NewRequest(tc.method, tc.path, body))
		if rec.Code != tc.wantStatus {
			t.Errorf("%s %s: status = %d, want %d", tc.method, tc.path, rec.Code, tc.wantStatus)
		}
	}
}

// --- test helpers ---

// stubRegistryHomepage is the homepage link the stub registry's detail
// payload advertises. Successful installs persist it host-side.
const stubRegistryHomepage = "https://example.com/docs"

// newRegistryStub returns an httptest server that emulates the
// directory's /plugins/<slug> detail endpoint plus a download URL
// the install handler will GET to retrieve the .ocpkg bytes. The
// stubbed sha256 is whatever the caller passes, so tests can force
// a mismatch by passing a wrong digest. Note the bare /plugins
// path: the host proxy appends `/plugins/<slug>` to its configured
// OWNCAST_PLUGIN_REGISTRY base, treating that base as the API root
// the same way the directory frontend treats its API_HOST.
//
// The detail handler reflects the incoming request's Host header back
// into the downloadURL so we don't have to know httptest's port in
// advance.
func newRegistryStub(t *testing.T, slug, version string, pkg []byte, sha string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/plugins/"+slug, func(w http.ResponseWriter, r *http.Request) {
		downloadURL := fmt.Sprintf("http://%s/ocpkg/%s", r.Host, slug)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"slug":     slug,
			"homepage": stubRegistryHomepage,
			"versions": []map[string]any{
				{
					"version":     version,
					"sha256":      sha,
					"downloadURL": downloadURL,
				},
			},
		})
	})
	mux.HandleFunc("/ocpkg/"+slug, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		_, _ = w.Write(pkg)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// buildJSPackageBytes assembles an in-memory .ocpkg whose code entry is
// author JavaScript run on the embedded shared engine, so install-path
// tests don't depend on the SDK example wasm being built.
func buildJSPackageBytes(t *testing.T, manifest []byte) []byte {
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
	jw, err := zw.Create("plugin.js")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := jw.Write([]byte(`const { definePlugin } = require("@owncast/plugin-sdk"); module.exports = definePlugin({});`)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// TestHandleRegistryInstall_PersistsHomepageSidecar covers the registry
// metadata flow: the registry's detail payload carries a homepage link,
// the install handler persists it via SetPluginHomepage (as a
// <base>.registry.json sidecar next to the installed .ocpkg), and the
// install response's entry JSON reflects it.
func TestHandleRegistryInstall_PersistsHomepageSidecar(t *testing.T) {
	pkg := buildJSPackageBytes(t, []byte(`{
		"api": "1",
		"name": "Hello World",
		"slug": "hello-world",
		"version": "0.1.0",
		"description": "registry homepage test",
		"permissions": []
	}`))
	sum := sha256.Sum256(pkg)
	sha := hex.EncodeToString(sum[:])

	upstream := newRegistryStub(t, "hello-world", "0.1.0", pkg, sha)
	withRegistryEnv(t, upstream.URL)
	host := newTestHost(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
		strings.NewReader(`{"slug":"hello-world","version":"0.1.0"}`))
	host.handleRegistryInstall(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	var entry plugins.DiscoveredEntry
	if err := json.Unmarshal(rec.Body.Bytes(), &entry); err != nil {
		t.Fatalf("decode install response: %v", err)
	}
	if entry.Homepage != stubRegistryHomepage {
		t.Errorf("install response homepage = %q, want %q", entry.Homepage, stubRegistryHomepage)
	}

	// The homepage is persisted as a sidecar next to the installed
	// package so it survives host restarts.
	entries := host.manager.List()
	if len(entries) != 1 {
		t.Fatalf("manager has %d entries, want 1", len(entries))
	}
	sidecar := strings.TrimSuffix(entries[0].Path, filepath.Ext(entries[0].Path)) + ".registry.json"
	raw, err := os.ReadFile(sidecar)
	if err != nil {
		t.Fatalf("read homepage sidecar: %v", err)
	}
	var meta struct {
		Homepage string `json:"homepage"`
	}
	if err := json.Unmarshal(raw, &meta); err != nil {
		t.Fatalf("decode sidecar %s: %v", raw, err)
	}
	if meta.Homepage != stubRegistryHomepage {
		t.Errorf("sidecar homepage = %q, want %q", meta.Homepage, stubRegistryHomepage)
	}
}

// TestHandleRegistryInstall_SamePermissionUpdateNeedsNoReEnable is the
// regression test for the admin-facing re-enable prompt: updating an
// installed, enabled, loaded plugin to a version that declares the SAME
// permission set must come back from the install endpoint with
// enabled=true and no pendingPermissions — the exact JSON fields the
// admin frontend checks to decide whether to pop the enable/approve
// modal (web/pages/admin/plugins.tsx installFromRegistry).
func TestHandleRegistryInstall_SamePermissionUpdateNeedsNoReEnable(t *testing.T) {
	manifest := func(version string) []byte {
		// Permissions deliberately not in sorted order: Enable stores the
		// approved baseline sorted, so an order-sensitive comparison would
		// misreport these as pending.
		return fmt.Appendf(nil, `{
			"api": "1",
			"name": "Perm Bot",
			"slug": "perm-bot",
			"version": %q,
			"description": "same-permission update test",
			"permissions": ["storage.kv", "chat.send"]
		}`, version)
	}
	install := func(host *Host, version string) map[string]any {
		t.Helper()
		pkg := buildJSPackageBytes(t, manifest(version))
		sum := sha256.Sum256(pkg)
		upstream := newRegistryStub(t, "perm-bot", version, pkg, hex.EncodeToString(sum[:]))
		withRegistryEnv(t, upstream.URL)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
			strings.NewReader(fmt.Sprintf(`{"slug":"perm-bot","version":%q}`, version)))
		host.handleRegistryInstall(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("install %s: status = %d, want 200 (body=%s)", version, rec.Code, rec.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode install %s response: %v", version, err)
		}
		return body
	}

	host := newTestHost(t)

	// Fresh install is not enabled: the modal SHOULD show here.
	body := install(host, "0.1.0")
	if enabled, _ := body["enabled"].(bool); enabled {
		t.Fatalf("fresh install must not be enabled, body=%v", body)
	}
	// Admin clicks Enable in the modal (same Manager op the enable
	// action endpoint dispatches).
	if err := host.manager.Enable(context.Background(), "perm-bot"); err != nil {
		t.Fatalf("enable: %v", err)
	}

	// Update to v2 with an identical permission set: no prompt.
	body = install(host, "0.2.0")
	if enabled, _ := body["enabled"].(bool); !enabled {
		t.Errorf("same-permission update must report enabled=true, body=%v", body)
	}
	if pending, ok := body["pendingPermissions"]; ok && pending != nil {
		if list, isList := pending.([]any); !isList || len(list) > 0 {
			t.Errorf("same-permission update must report no pendingPermissions, got %v", pending)
		}
	}
	if loaded, _ := body["loaded"].(bool); !loaded {
		t.Errorf("same-permission update must stay loaded, body=%v", body)
	}
	if v, _ := body["version"].(string); v != "0.2.0" {
		t.Errorf("version after update = %q, want 0.2.0", v)
	}
}

// TestHandleRegistryInstall_SamePermUpdateSurvivesRestart is the same
// scenario across an Owncast restart, with the production
// configEnabledStore persisting the enabled set and approved-permission
// baseline in the datastore. If the baseline is lost in the
// Save/Load round trip, the update reports every permission as pending
// and the admin gets re-prompted to approve an unchanged set.
func TestHandleRegistryInstall_SamePermUpdateSurvivesRestart(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	store := &configEnabledStore{datastore: newTestDatastore(t)}

	manifest := func(version string) []byte {
		return fmt.Appendf(nil, `{
			"api": "1",
			"name": "Perm Bot",
			"slug": "perm-bot",
			"version": %q,
			"description": "restart persistence test",
			"permissions": ["storage.kv", "chat.send"]
		}`, version)
	}
	install := func(host *Host, version string) map[string]any {
		t.Helper()
		pkg := buildJSPackageBytes(t, manifest(version))
		sum := sha256.Sum256(pkg)
		upstream := newRegistryStub(t, "perm-bot", version, pkg, hex.EncodeToString(sum[:]))
		withRegistryEnv(t, upstream.URL)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/admin/plugin-registry/install",
			strings.NewReader(fmt.Sprintf(`{"slug":"perm-bot","version":%q}`, version)))
		host.handleRegistryInstall(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("install %s: status = %d, want 200 (body=%s)", version, rec.Code, rec.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode install %s response: %v", version, err)
		}
		return body
	}

	mgr1 := plugins.NewManagerWithStore(dir, &plugins.HostEnv{}, store)
	if err := mgr1.Start(ctx); err != nil {
		t.Fatalf("manager 1 start: %v", err)
	}
	install(&Host{manager: mgr1}, "0.1.0")
	if err := mgr1.Enable(ctx, "perm-bot"); err != nil {
		t.Fatalf("enable: %v", err)
	}
	mgr1.Stop(ctx)

	// "Restart": a fresh manager over the same plugins dir + datastore.
	mgr2 := plugins.NewManagerWithStore(dir, &plugins.HostEnv{}, store)
	if err := mgr2.Start(ctx); err != nil {
		t.Fatalf("manager 2 start: %v", err)
	}
	t.Cleanup(func() { mgr2.Stop(ctx) })

	body := install(&Host{manager: mgr2}, "0.2.0")
	if enabled, _ := body["enabled"].(bool); !enabled {
		t.Errorf("post-restart same-permission update must report enabled=true, body=%v", body)
	}
	if pending, ok := body["pendingPermissions"]; ok && pending != nil {
		if list, isList := pending.([]any); !isList || len(list) > 0 {
			t.Errorf("post-restart same-permission update must report no pendingPermissions, got %v", pending)
		}
	}
	if loaded, _ := body["loaded"].(bool); !loaded {
		t.Errorf("post-restart same-permission update must stay loaded, body=%v", body)
	}
}
