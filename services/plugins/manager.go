package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	extism "github.com/extism/go-sdk"
	"github.com/gobwas/glob"
)

// Loaded represents a successfully-loaded plugin. The sidecar manifest is the
// source of truth for identity and permissions; subscriptions are populated
// from the runtime register() call (the SDK derives them from the plugin's
// handlers, so authors don't maintain a duplicate list).
//
// Extism plugin instances are not safe for concurrent calls. The mutex
// serializes calls to a single plugin while still allowing different plugins
// to run in parallel.
//
// AssetsFS is the static-asset root for this plugin, or nil if the plugin
// ships no assets. For loose-files deployments it's an os.DirFS; for .ocpkg
// deployments it's a sub-FS of the in-memory zip reader. plugin.Server reads
// through this interface so both layouts work the same way.
type Loaded struct {
	Manifest    *Manifest
	WasmPath    string
	AssetsFS    fs.FS
	adminGlobs  []glob.Glob // compiled from manifest.admin.pages[].path
	plugin      *extism.Plugin
	mu          sync.Mutex
	failureMu   sync.Mutex
	filterFails int
	disabled    atomic.Bool
	// pkgCloser holds the file-backed zip reader for .ocpkg plugins so the
	// underlying file stays open for AssetsFS reads. nil for loose-files
	// plugins. Closed by Loaded.Close.
	pkgCloser io.Closer
}

// FilterStrikeThreshold is the number of consecutive filter failures a
// plugin can rack up before the dispatcher auto-disables it for the rest
// of the session. The fail-open semantics still apply on the path to the
// strike — events flow normally; the strike just prevents the host from
// drowning in log noise from a permanently-broken plugin.
const FilterStrikeThreshold = 5

// Sandbox caps. A misbehaving plugin should fail its own call; the host
// stays up. These are deliberately generous — realistic plugins won't
// come close. Per-plugin manifest overrides are a future TODO.
const (
	// MaxWasmPages caps a plugin's wasm linear memory. 1 page = 64 KiB,
	// so 1024 = 64 MiB. QuickJS itself takes a few MB; this leaves
	// comfortable room for plugin state.
	MaxWasmPages = 1024

	// MaxExtismHTTPResponseBytes caps the body of any outbound HTTP
	// request extism's built-in http_request makes on the plugin's
	// behalf. Matches the inbound HTTP response cap.
	MaxExtismHTTPResponseBytes = 10 << 20 // 10 MiB

	// MaxExtismVarBytes caps extism's internal per-plugin Var KV (a
	// separate store from our owncast.kv namespace). We don't expose
	// it but defense in depth.
	MaxExtismVarBytes = 1 << 20 // 1 MiB

	// MaxRegisterOutputBytes caps the JSON the SDK emits from register().
	// In practice this is a kilobyte or two (manifest echo) — the cap is
	// just to prevent a buggy or malicious plugin from causing a huge
	// allocation at load time.
	MaxRegisterOutputBytes = 256 << 10 // 256 KiB

	// MaxFilterOutputBytes caps the JSON a plugin's on_filter returns.
	// Filter results carry the (possibly modified) event payload —
	// chat messages, etc. — which are small in any realistic case.
	MaxFilterOutputBytes = 1 << 20 // 1 MiB

	// MaxHTTPHandlerOutputBytes caps the JSON envelope a plugin returns
	// from on_http_request (status + headers + body). Sized to leave
	// headroom over MaxHTTPResponseBodyBytes (server.go); the inner body
	// is then checked again post-unmarshal.
	MaxHTTPHandlerOutputBytes = 12 << 20 // 12 MiB

	// NotifyTimeout caps a single on_event call. Notification handlers
	// can do real work (kv writes, owncast.* host calls), but they
	// shouldn't stall — events fire on the chat hot path.
	NotifyTimeout = 500 * time.Millisecond

	// HTTPHandlerTimeout caps a single on_http_request call. HTTP
	// handlers may legitimately do work (fetch upstream, compute), so
	// this is looser than NotifyTimeout but still bounded.
	HTTPHandlerTimeout = 5 * time.Second
)

// IsDisabled reports whether the plugin has been auto-disabled by the
// strike system. Disabled plugins are skipped by both the filter chain
// and the notification dispatcher.
func (l *Loaded) IsDisabled() bool {
	return l.disabled.Load()
}

func (l *Loaded) recordFilterFailure() bool {
	l.failureMu.Lock()
	defer l.failureMu.Unlock()
	l.filterFails++
	if l.filterFails >= FilterStrikeThreshold && !l.disabled.Load() {
		l.disabled.Store(true)
		return true
	}
	return false
}

func (l *Loaded) recordFilterSuccess() {
	l.failureMu.Lock()
	defer l.failureMu.Unlock()
	l.filterFails = 0
}

// IsAdminPath reports whether the request path (relative to the plugin's
// namespace, e.g. "/admin/foo") matches any of the declared admin page
// globs. Used by Server to require authentication on admin-only routes.
func (l *Loaded) IsAdminPath(path string) bool {
	for _, g := range l.adminGlobs {
		if g.Match(path) {
			return true
		}
	}
	return false
}

// Close releases the underlying wasm instance and any retained file handles
// (the .ocpkg zip reader for packaged plugins). Safe to call multiple times.
func (l *Loaded) Close(ctx context.Context) {
	if l.plugin != nil {
		_ = l.plugin.Close(ctx)
		l.plugin = nil
	}
	if l.pkgCloser != nil {
		_ = l.pkgCloser.Close()
		l.pkgCloser = nil
	}
}

// Manager tracks plugins across two states:
//
//   - Discovered: a file in the plugins directory whose manifest parsed
//     successfully. The host knows it exists and can show it to an admin.
//     No wasm instance, no events delivered.
//
//   - Loaded: discovered + an admin has explicitly enabled it. Wasm
//     instantiated, events flow.
//
// The enabled set persists via an EnabledStore so admin choices survive
// host restarts (a JSON file by default; Owncast backs it with native
// config). Files appearing in the plugins directory are auto-detected (scan
// every ScanInterval) but never auto-loaded — the admin clicks Enable.
type Manager struct {
	pluginsDir   string
	enabledStore EnabledStore
	env          *HostEnv

	mu         sync.RWMutex
	discovered map[string]*DiscoveredEntry // keyed by manifest.name
	loaded     map[string]*Loaded          // subset of discovered that's currently running
	enabledSet map[string]bool             // names the admin has enabled

	scanInterval time.Duration
	cancel       context.CancelFunc // stops the scan loop
	scanCh       chan struct{}      // pings to force a scan (testing / admin trigger)
}

// DiscoveredEntry is the public view of a discovered plugin — what the
// admin UI lists.
type DiscoveredEntry struct {
	Name         string    `json:"name"`
	Version      string    `json:"version,omitempty"`
	Description  string    `json:"description,omitempty"`
	Permissions  []string  `json:"permissions,omitempty"`
	Path         string    `json:"path"`
	Enabled      bool      `json:"enabled"`
	Loaded       bool      `json:"loaded"`
	LastError    string    `json:"lastError,omitempty"`
	DiscoveredAt time.Time `json:"discoveredAt"`
}

// ScanInterval is how often the manager re-scans the plugins directory.
const ScanInterval = 2 * time.Second

// NewManager constructs a Manager that persists its enabled set to a
// .enabled.json file in pluginsDir. Owncast wires NewManagerWithStore
// instead, backing the enabled set with native config storage.
func NewManager(pluginsDir string, env *HostEnv) *Manager {
	store := newFileEnabledStore(filepath.Join(pluginsDir, ".enabled.json"))
	return NewManagerWithStore(pluginsDir, env, store)
}

// NewManagerWithStore is NewManager with an explicit EnabledStore, letting
// the host persist the enabled set wherever it likes (e.g. Owncast's config
// datastore) instead of a JSON file in the plugins directory.
func NewManagerWithStore(pluginsDir string, env *HostEnv, store EnabledStore) *Manager {
	return &Manager{
		pluginsDir:   pluginsDir,
		enabledStore: store,
		env:          env,
		discovered:   make(map[string]*DiscoveredEntry),
		loaded:       make(map[string]*Loaded),
		enabledSet:   make(map[string]bool),
		scanInterval: ScanInterval,
		scanCh:       make(chan struct{}, 1),
	}
}

// Start does the initial scan, loads everything in the enabled set, and
// begins a background scan loop. Stop() cancels the loop.
func (m *Manager) Start(ctx context.Context) error {
	if err := m.loadEnabledSet(); err != nil {
		return fmt.Errorf("load enabled set: %w", err)
	}
	if err := m.scan(ctx); err != nil {
		return fmt.Errorf("initial scan: %w", err)
	}
	// Auto-load anything in the enabled set that isn't already loaded.
	for name, enabled := range m.enabledSet {
		if !enabled {
			continue
		}
		if err := m.loadInternal(ctx, name); err != nil {
			fmt.Fprintf(os.Stderr, "plugin %s: load failed: %v\n", name, err)
		}
	}
	scanCtx, cancel := context.WithCancel(ctx)
	m.cancel = cancel
	go m.scanLoop(scanCtx)
	return nil
}

// Stop cancels the scan loop and closes all loaded plugins.
func (m *Manager) Stop(ctx context.Context) {
	if m.cancel != nil {
		m.cancel()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, l := range m.loaded {
		l.Close(ctx)
	}
	m.loaded = map[string]*Loaded{}
}

// List returns a snapshot of all discovered plugins for admin UI.
func (m *Manager) List() []DiscoveredEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]DiscoveredEntry, 0, len(m.discovered))
	for name, d := range m.discovered {
		entry := *d
		entry.Enabled = m.enabledSet[name]
		_, isLoaded := m.loaded[name]
		entry.Loaded = isLoaded
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Snapshot returns the currently-loaded plugins. Dispatcher and Server call
// this on every operation so changes from Enable/Disable take effect
// without restarting anything.
func (m *Manager) Snapshot() []*Loaded {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Loaded, 0, len(m.loaded))
	for _, l := range m.loaded {
		out = append(out, l)
	}
	return out
}

// Enable marks a discovered plugin as enabled, persists the choice, and
// loads it. No-op if already loaded.
func (m *Manager) Enable(ctx context.Context, name string) error {
	m.mu.Lock()
	if _, ok := m.discovered[name]; !ok {
		m.mu.Unlock()
		return fmt.Errorf("plugin %q not discovered", name)
	}
	if m.enabledSet[name] {
		// Already enabled in the persisted set; just make sure it's loaded.
		if _, ok := m.loaded[name]; ok {
			m.mu.Unlock()
			return nil
		}
	}
	m.enabledSet[name] = true
	m.mu.Unlock()
	if err := m.saveEnabledSet(); err != nil {
		return fmt.Errorf("persist enabled set: %w", err)
	}
	err := m.loadInternal(ctx, name)
	return err
}

// Disable unloads a plugin and persists the choice. No-op if already disabled.
func (m *Manager) Disable(ctx context.Context, name string) error {
	m.mu.Lock()
	if _, ok := m.discovered[name]; !ok {
		m.mu.Unlock()
		return fmt.Errorf("plugin %q not discovered", name)
	}
	delete(m.enabledSet, name)
	loaded := m.loaded[name]
	delete(m.loaded, name)
	m.mu.Unlock()
	if err := m.saveEnabledSet(); err != nil {
		// Don't bail — we've already removed from the in-memory set.
		fmt.Fprintf(os.Stderr, "persist enabled set: %v\n", err)
	}
	if loaded != nil {
		loaded.Close(ctx)
	}
	return nil
}

// Reload unloads and reloads a plugin. Plugin author rebuilt → admin
// triggers a reload to pick up the new wasm without restarting the host.
// Plugin must currently be enabled (otherwise call Enable instead).
func (m *Manager) Reload(ctx context.Context, name string) error {
	m.mu.Lock()
	if _, ok := m.discovered[name]; !ok {
		m.mu.Unlock()
		return fmt.Errorf("plugin %q not discovered", name)
	}
	if !m.enabledSet[name] {
		m.mu.Unlock()
		return fmt.Errorf("plugin %q is not enabled; use Enable to load it", name)
	}
	loaded := m.loaded[name]
	delete(m.loaded, name)
	m.mu.Unlock()
	if loaded != nil {
		loaded.Close(ctx)
	}
	err := m.loadInternal(ctx, name)
	return err
}

// loadInternal performs the actual load and inserts into m.loaded. Assumes
// the caller has already verified the plugin is discovered + enabled.
func (m *Manager) loadInternal(ctx context.Context, name string) error {
	m.mu.RLock()
	d, ok := m.discovered[name]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %q not discovered", name)
	}

	loaded, err := loadByPath(ctx, m.env, d.Path)
	if err != nil {
		m.mu.Lock()
		m.discovered[name].LastError = err.Error()
		m.mu.Unlock()
		return err
	}

	m.mu.Lock()
	m.loaded[name] = loaded
	m.discovered[name].LastError = ""
	m.mu.Unlock()
	return nil
}

// loadByPath dispatches to LoadPlugin or LoadPackage based on file suffix.
// Sets AssetsFS for loose-files plugins.
func loadByPath(ctx context.Context, env *HostEnv, path string) (*Loaded, error) {
	switch {
	case strings.HasSuffix(path, packageSuffix):
		return LoadPackage(ctx, env, path)
	case strings.HasSuffix(path, ".wasm"):
		manifestPath := strings.TrimSuffix(path, ".wasm") + ".manifest.json"
		loaded, err := LoadPlugin(ctx, env, path, manifestPath)
		if err != nil {
			return nil, err
		}
		assetsDir := strings.TrimSuffix(path, ".wasm") + "-assets"
		if info, err := os.Stat(assetsDir); err == nil && info.IsDir() {
			loaded.AssetsFS = os.DirFS(assetsDir)
		}
		return loaded, nil
	}
	return nil, fmt.Errorf("unsupported plugin file: %s", path)
}

// scan re-reads the plugins directory, updates the discovered map, and
// unloads anything whose underlying file has gone away.
func (m *Manager) scan(ctx context.Context) error {
	entries, err := os.ReadDir(m.pluginsDir)
	if err != nil {
		return fmt.Errorf("read plugins dir %q: %w", m.pluginsDir, err)
	}

	seen := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".wasm") && !strings.HasSuffix(name, packageSuffix) {
			continue
		}
		path := filepath.Join(m.pluginsDir, name)
		manifest, err := readManifestForPath(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "discover %s: %v\n", name, err)
			continue
		}
		seen[manifest.Name] = true

		m.mu.Lock()
		if existing, ok := m.discovered[manifest.Name]; ok {
			// Already discovered; refresh manifest metadata in case it changed.
			existing.Version = manifest.Version
			existing.Description = manifest.Description
			existing.Permissions = manifest.Permissions
			existing.Path = path
		} else {
			m.discovered[manifest.Name] = &DiscoveredEntry{
				Name:         manifest.Name,
				Version:      manifest.Version,
				Description:  manifest.Description,
				Permissions:  manifest.Permissions,
				Path:         path,
				DiscoveredAt: time.Now(),
			}
		}
		m.mu.Unlock()
	}

	// Anything we knew about but didn't see this scan: gone from disk.
	var removed []string
	m.mu.RLock()
	for name := range m.discovered {
		if !seen[name] {
			removed = append(removed, name)
		}
	}
	m.mu.RUnlock()

	for _, name := range removed {
		m.mu.Lock()
		delete(m.discovered, name)
		loaded := m.loaded[name]
		delete(m.loaded, name)
		m.mu.Unlock()
		if loaded != nil {
			loaded.Close(ctx)
		}
	}
	return nil
}

func (m *Manager) scanLoop(ctx context.Context) {
	ticker := time.NewTicker(m.scanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		case <-m.scanCh:
		}
		if err := m.scan(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		}
	}
}

// readManifestForPath parses the manifest of a discovered file without
// instantiating the wasm. Used during scan; cheap enough to run every
// few seconds.
func readManifestForPath(path string) (*Manifest, error) {
	switch {
	case strings.HasSuffix(path, packageSuffix):
		return readManifestFromPackage(path)
	case strings.HasSuffix(path, ".wasm"):
		manifestPath := strings.TrimSuffix(path, ".wasm") + ".manifest.json"
		data, err := os.ReadFile(manifestPath) //nolint:gosec // G304: plugin paths are admin-controlled, not user input
		if err != nil {
			return nil, fmt.Errorf("read sidecar manifest: %w", err)
		}
		return ParseManifest(data)
	}
	return nil, fmt.Errorf("unsupported file type: %s", path)
}

// Persistence — a tiny JSON file under the plugins directory listing the
// names the admin has enabled. Survives restarts.

// EnabledStore persists the set of enabled plugin names so admin choices
// survive host restarts. The PoC's default is a JSON file in the plugins
// directory (fileEnabledStore); integrated into Owncast it is backed by the
// native config datastore (see NewManagerWithStore).
type EnabledStore interface {
	// LoadEnabled returns the persisted enabled plugin names. A store with
	// nothing persisted yet returns an empty slice and a nil error.
	LoadEnabled() ([]string, error)
	// SaveEnabled replaces the persisted set with names (already sorted).
	SaveEnabled(names []string) error
}

func (m *Manager) loadEnabledSet() error {
	names, err := m.enabledStore.LoadEnabled()
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, name := range names {
		m.enabledSet[name] = true
	}
	return nil
}

func (m *Manager) saveEnabledSet() error {
	m.mu.RLock()
	names := make([]string, 0, len(m.enabledSet))
	for name, enabled := range m.enabledSet {
		if enabled {
			names = append(names, name)
		}
	}
	m.mu.RUnlock()
	sort.Strings(names)
	return m.enabledStore.SaveEnabled(names)
}

// enabledFileContents is the on-disk JSON shape for fileEnabledStore.
type enabledFileContents struct {
	Enabled []string `json:"enabled"`
}

// fileEnabledStore persists the enabled set to a JSON file. It's the default
// for the standalone runtime and tests; Owncast supplies a config-backed
// EnabledStore instead.
type fileEnabledStore struct {
	path string
}

func newFileEnabledStore(path string) *fileEnabledStore {
	return &fileEnabledStore{path: path}
}

func (s *fileEnabledStore) LoadEnabled() ([]string, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // start with empty set
		}
		return nil, err
	}
	var f enabledFileContents
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return f.Enabled, nil
}

func (s *fileEnabledStore) SaveEnabled(names []string) error {
	data, err := json.MarshalIndent(enabledFileContents{Enabled: names}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// LoadPlugin loads a single plugin given explicit wasm and manifest paths
// (the loose-files layout). Used by the test runner so it shares the exact
// same load + register + validate path that production uses via Start.
//
// AssetsFS on the returned Loaded is left nil — callers that want static
// asset serving should populate it themselves.
func LoadPlugin(ctx context.Context, env *HostEnv, wasmPath, manifestPath string) (*Loaded, error) {
	manifestBytes, err := os.ReadFile(manifestPath) //nolint:gosec // G304: plugin paths are admin-controlled, not user input
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", manifestPath, err)
	}
	wasmBytes, err := os.ReadFile(wasmPath) //nolint:gosec // G304: plugin paths are admin-controlled, not user input
	if err != nil {
		return nil, fmt.Errorf("read wasm %s: %w", wasmPath, err)
	}
	displayName := strings.TrimSuffix(filepath.Base(wasmPath), ".wasm")
	loaded, err := loadFromBytes(ctx, env, manifestBytes, wasmBytes, displayName)
	if err != nil {
		return nil, err
	}
	loaded.WasmPath = wasmPath
	return loaded, nil
}

// loadFromBytes is the shared core of LoadPlugin and LoadPackage.
func loadFromBytes(ctx context.Context, env *HostEnv, manifestBytes, wasmBytes []byte, displayName string) (*Loaded, error) {
	manifest, err := ParseManifest(manifestBytes)
	if err != nil {
		return nil, err
	}

	hostFns := BuildHostFunctions(env, manifest)

	extismManifest := extism.Manifest{
		Wasm:    []extism.Wasm{extism.WasmData{Data: wasmBytes, Name: displayName}},
		Timeout: 10_000, // milliseconds; enables Wazero's WithCloseOnContextDone
		// Sandbox caps. A plugin that exceeds these gets an error from the
		// next Call; the host stays up. Defaults are generous enough for
		// realistic plugins; per-plugin manifest overrides are a future TODO.
		Memory: &extism.ManifestMemory{
			MaxPages:             MaxWasmPages,               // wasm linear memory cap
			MaxHttpResponseBytes: MaxExtismHTTPResponseBytes, // outbound http body cap
			MaxVarBytes:          MaxExtismVarBytes,          // extism's internal Var KV
		},
	}
	for _, p := range manifest.Permissions {
		if p == PermNetworkFetch {
			// Manifest validation already required AllowedHosts to be
			// non-empty when network.fetch is granted, so passing the
			// list through is safe — admins explicitly authorized this
			// scope by approving the manifest at install time.
			extismManifest.AllowedHosts = append([]string(nil), manifest.Network.AllowedHosts...)
			break
		}
	}
	p, err := extism.NewPlugin(ctx, extismManifest, extism.PluginConfig{EnableWasi: true}, hostFns)
	if err != nil {
		return nil, fmt.Errorf("instantiate wasm: %w", err)
	}
	p.SetLogger(func(level extism.LogLevel, message string) {
		fmt.Fprintf(os.Stderr, "[%s] %s\n", displayName, message)
	})

	if !p.FunctionExists("register") {
		_ = p.Close(ctx)
		return nil, fmt.Errorf("plugin does not export register()")
	}
	_, regOut, err := p.Call("register", nil)
	if err != nil {
		_ = p.Close(ctx)
		return nil, fmt.Errorf("call register(): %w", err)
	}
	if len(regOut) > MaxRegisterOutputBytes {
		_ = p.Close(ctx)
		return nil, fmt.Errorf("register() output too large: %d bytes (max %d)", len(regOut), MaxRegisterOutputBytes)
	}
	var runtime Manifest
	if err := json.Unmarshal(regOut, &runtime); err != nil {
		_ = p.Close(ctx)
		return nil, fmt.Errorf("parse register() output: %w", err)
	}
	if err := manifest.AgreesWith(&runtime); err != nil {
		_ = p.Close(ctx)
		return nil, fmt.Errorf("manifest/runtime mismatch: %w", err)
	}
	manifest.Subscriptions = runtime.Subscriptions

	var adminGlobs []glob.Glob
	for _, page := range manifest.Admin.Pages {
		g, err := glob.Compile(page.Path)
		if err != nil {
			_ = p.Close(ctx)
			return nil, fmt.Errorf("manifest.admin.pages: invalid path glob %q: %w", page.Path, err)
		}
		adminGlobs = append(adminGlobs, g)
	}

	return &Loaded{Manifest: manifest, plugin: p, adminGlobs: adminGlobs}, nil
}
