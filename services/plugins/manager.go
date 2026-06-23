package plugins

import (
	"archive/zip"
	"bytes"
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
	log "github.com/sirupsen/logrus"
	"github.com/tetratelabs/wazero"
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
// PublicFS is the web-served file root for this plugin at
// /plugins/<slug>/<file>, or nil when the plugin ships nothing under
// public/. Populated from a `public/` directory in the .ocpkg (or a
// <name>-public/ sibling for loose-files plugins). plugin.Server reads
// through this interface.
//
// AssetsFS is the internal-only file root the host reads for manifest
// fields that inline file contents (styles, scripts, extraPageContent),
// or nil when the plugin ships nothing under assets/. Files under
// assets/ are never reachable through the plugin's URL space; authors
// can therefore put source CSS/JS/HTML the host inlines into
// /api/config or /customjavascript without the same files also being
// downloadable directly. Populated from an `assets/` directory in the
// .ocpkg (or a <name>-assets/ sibling for loose-files plugins).
type Loaded struct {
	Manifest    *Manifest
	WasmPath    string
	PublicFS    fs.FS
	AssetsFS    fs.FS
	adminGlobs  []glob.Glob // compiled from manifest.admin.pages[].path
	adminPaths  []string    // original path strings, used for "page gates descendants" prefix-matching
	plugin      *extism.Plugin
	mu          sync.Mutex
	failureMu   sync.Mutex
	filterFails int
	disabled    atomic.Bool
	// pkgCloser holds the file-backed zip reader for .ocpkg plugins so the
	// underlying file stays open for PublicFS / AssetsFS reads. nil for
	// loose-files plugins. Closed by Loaded.Close.
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

// wasmCompilationCache is shared across every plugin instance so wazero
// compiles a given wasm module to native code at most once for the lifetime
// of the process. The cache is keyed on module content, so it dedupes the
// repeated work of reloading the same plugin (dev churn, reload-on-change,
// or running multiple instances of one plugin). It is safe for concurrent
// use across runtimes.
//
// NOTE: this does NOT dedupe the interpreter engine across *different*
// plugins — each plugin's wasm bundles the QuickJS/CPython engine together
// with the author's code, so every plugin's bytes differ and hash to a
// distinct cache key. Sharing the engine across distinct plugins requires
// splitting it out of the per-plugin module; tracked separately.
var (
	wasmCompilationCache     wazero.CompilationCache
	wasmCompilationCacheOnce sync.Once
)

func sharedCompilationCache() wazero.CompilationCache {
	wasmCompilationCacheOnce.Do(func() {
		wasmCompilationCache = wazero.NewCompilationCache()
	})
	return wasmCompilationCache
}

// IsDisabled reports whether the plugin has been auto-disabled by the
// strike system. Disabled plugins are omitted from Manager.Snapshot, so they
// are skipped by every dispatch and serve path (filter chain, notifications,
// HTTP/SSE routing, timers, and CSS/JS/tab/action contributions).
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
// namespace, e.g. "/admin/foo") is gated by one of the declared admin
// pages. A declared page path gates itself, anything under it as a path
// prefix (so "/admin" gates "/admin/api/x"), and anything its glob
// matches (so "/admin/*" still works as it did when authors wrote an
// explicit wildcard). Used by Server to require authentication on
// admin-only routes.
func (l *Loaded) IsAdminPath(path string) bool {
	for i, g := range l.adminGlobs {
		if g.Match(path) {
			return true
		}
		// Page paths gate their descendants too, so a declaration of
		// "/admin" automatically covers "/admin/api/settings" without
		// the author having to spell out a glob.
		if lit := l.adminPaths[i]; lit != "" && strings.HasPrefix(path, lit+"/") {
			return true
		}
	}
	return false
}

// CallTabContent invokes the plugin's on_tab_content export with the
// given slug and optional user, returning the rendered HTML. Returns
// empty string when the plugin does not export on_tab_content.
func (p *Loaded) CallTabContent(ctx context.Context, slug string, user *HostUser) (string, error) {
	return p.callContentExport(ctx, "on_tab_content", slug, user)
}

// CallPageContent invokes the plugin's on_page_content export with
// the given slug and optional user, returning the rendered HTML.
// Returns empty string when the plugin does not export on_page_content.
func (p *Loaded) CallPageContent(ctx context.Context, slug string, user *HostUser) (string, error) {
	return p.callContentExport(ctx, "on_page_content", slug, user)
}

func (p *Loaded) callContentExport(ctx context.Context, export, slug string, user *HostUser) (string, error) {
	p.mu.Lock()
	pl := p.plugin
	p.mu.Unlock()
	if pl == nil || !pl.FunctionExists(export) {
		return "", nil
	}
	req := map[string]any{"slug": slug}
	if user != nil {
		req["user"] = user
	}
	input, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("%s: marshal input: %w", export, err)
	}
	callCtx, cancel := context.WithTimeout(ctx, HTTPHandlerTimeout)
	defer cancel()
	p.mu.Lock()
	_, out, err := pl.CallWithContext(callCtx, export, input)
	p.mu.Unlock()
	if err != nil {
		return "", fmt.Errorf("%s: %w", export, err)
	}
	return string(out), nil
}

// CallPageStyles invokes the plugin's on_page_styles export, returning the
// CSS the plugin wants inlined into the viewer page's customStyles. The host
// calls this for every enabled plugin that holds ui.modify; a plugin opts in
// simply by exporting the function. Returns empty string when the plugin does
// not export on_page_styles. Unlike on_page_content, the call takes no slug or
// user — dynamic styles are global so the /api/config response stays cacheable.
func (p *Loaded) CallPageStyles(ctx context.Context) (string, error) {
	return p.callVoidExport(ctx, "on_page_styles")
}

// CallPageScripts mirrors CallPageStyles for the plugin's on_page_scripts
// export, returning JavaScript to append to the viewer's customJavascript.
func (p *Loaded) CallPageScripts(ctx context.Context) (string, error) {
	return p.callVoidExport(ctx, "on_page_scripts")
}

// callVoidExport invokes an export that takes no input and returns a string.
// Used by the global, user-agnostic dynamic-content hooks (on_page_styles,
// on_page_scripts). Returns empty string when the plugin doesn't export it.
func (p *Loaded) callVoidExport(ctx context.Context, export string) (string, error) {
	p.mu.Lock()
	pl := p.plugin
	p.mu.Unlock()
	if pl == nil || !pl.FunctionExists(export) {
		return "", nil
	}
	callCtx, cancel := context.WithTimeout(ctx, HTTPHandlerTimeout)
	defer cancel()
	p.mu.Lock()
	_, out, err := pl.CallWithContext(callCtx, export, nil)
	p.mu.Unlock()
	if err != nil {
		return "", fmt.Errorf("%s: %w", export, err)
	}
	return string(out), nil
}

// Close releases the underlying wasm instance and any retained file handles
// (the .ocpkg zip reader for packaged plugins). Safe to call multiple times.
func (l *Loaded) Close(ctx context.Context) {
	// Clear the instance pointer under the lock so concurrent dispatch
	// goroutines (which read p.plugin under the same lock) never observe a torn
	// instance. Close the wasm instance OUTSIDE the lock: l.plugin.Close can
	// block on an in-flight call that itself holds l.mu, so closing under the
	// lock would deadlock.
	l.mu.Lock()
	pl := l.plugin
	l.plugin = nil
	l.mu.Unlock()
	if pl != nil {
		// Closes this instance only. For shared-engine plugins the underlying
		// CompiledPlugin (the engine) is process-global and is never closed
		// here; for legacy wasm plugins this closes their own runtime.
		_ = pl.Close(ctx)
	}
	// Drop the call-time identity so any in-flight or stale host call from this
	// plugin resolves to "not found" and fails cleanly rather than acting on a
	// torn-down plugin.
	if l.Manifest != nil {
		globalPluginRegistry.remove(l.Manifest.Slug)
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

	mu          sync.RWMutex
	discovered  map[string]*DiscoveredEntry // keyed by manifest.name
	loaded      map[string]*Loaded          // subset of discovered that's currently running
	enabledSet  map[string]bool             // names the admin has enabled
	approvedSet map[string][]string         // plugin name -> sorted approved permission set

	scanInterval time.Duration
	cancel       context.CancelFunc // stops the scan loop
	scanCh       chan struct{}      // pings to force a scan (testing / admin trigger)

	// onUnload, if set, is called with a plugin's slug just before its
	// instance is closed (disable, reload, disk-removal, or host stop), so
	// host subsystems can release per-plugin resources. Timers use it to
	// cancel a plugin's pending callbacks.
	onUnload func(slug string)
}

// SetOnUnload registers a callback invoked with a plugin's slug right before
// its instance is closed. Used to cancel per-plugin host resources (timers).
func (m *Manager) SetOnUnload(fn func(slug string)) { m.onUnload = fn }

// notifyUnload fires the onUnload hook for a plugin about to be closed.
func (m *Manager) notifyUnload(l *Loaded) {
	if m.onUnload != nil && l != nil && l.Manifest != nil {
		m.onUnload(l.Manifest.Slug)
	}
}

// DiscoveredEntry is the public view of a discovered plugin: what the
// admin UI lists, and what the registry's install endpoint returns
// to the host. Two name-like fields:
//
//   - Slug is the canonical identifier (URL segment, KV namespace,
//     file path, registry primary key). Stable, lowercase, hyphenated.
//   - DisplayName is the human-readable name the admin sees in lists.
//     Set by the plugin author; can contain any characters.
//
// BotDisplayName, when non-empty, overrides DisplayName as the chat
// identity for plugins that post to chat. Empty means "use
// DisplayName" (resolved at chat-send time, not here).
type DiscoveredEntry struct {
	Slug           string   `json:"slug"`
	DisplayName    string   `json:"name"`
	BotDisplayName string   `json:"botDisplayName,omitempty"`
	Version        string   `json:"version,omitempty"`
	Description    string   `json:"description,omitempty"`
	Permissions    []string `json:"permissions,omitempty"`
	Path           string   `json:"path"`
	Enabled        bool     `json:"enabled"`
	Loaded         bool     `json:"loaded"`
	// AutoDisabled is set when the dispatcher's strike system stopped
	// invoking the plugin (too many consecutive filter failures).
	// Enabled stays true so the admin can see what they originally chose,
	// but the plugin isn't doing any work; reload or fix-and-rebuild to
	// reset the strike counter.
	AutoDisabled bool `json:"autoDisabled,omitempty"`
	// HasIcon reports whether the plugin ships an icon.png alongside its
	// manifest (top-level in the .ocpkg, or <base>.icon.png next to a
	// loose .wasm). The admin UI fetches the bytes from
	// /api/plugins/<name>/icon when this is true; no http.serve
	// permission is required to ship one.
	HasIcon bool `json:"hasIcon,omitempty"`
	// HasInstructions reports whether the plugin ships an INSTRUCTIONS.md
	// alongside its manifest (top-level in the .ocpkg, or
	// <base>.INSTRUCTIONS.md next to a loose .wasm). The admin UI fetches
	// the markdown from /api/admin/plugins/<name>/instructions when this is
	// true and renders it in a details tab; no permission is required to
	// ship one.
	HasInstructions bool `json:"hasInstructions,omitempty"`
	// PendingPermissions lists permissions the current manifest declares
	// that the admin has not yet approved. Non-empty means the plugin was
	// updated on disk to request more access than was originally granted;
	// the plugin will not load until the admin re-enables it (which
	// captures a fresh approval snapshot covering the new set).
	PendingPermissions []string `json:"pendingPermissions,omitempty"`
	// AllowedHosts mirrors manifest.network.allowedHosts so the admin
	// UI can show the host scope alongside the network.fetch
	// permission entry. Empty for plugins that don't request
	// network.fetch; the host's load-time validator rejects
	// network.fetch without an entry, so a present permission
	// without a non-empty list never occurs.
	AllowedHosts []string    `json:"allowedHosts,omitempty"`
	LastError    string      `json:"lastError,omitempty"`
	DiscoveredAt time.Time   `json:"discoveredAt"`
	AdminPages   []AdminPage `json:"adminPages,omitempty"`
	// Commands lists the plugin's chat commands for the admin details view.
	// Derived by the SDK and reported via register(), so it's only known once
	// the plugin is loaded — empty for discovered-but-not-loaded plugins.
	Commands []CommandInfo `json:"commands,omitempty"`
	// Config is the plugin's manifest-declared config schema (key → field with
	// type/default/description). The admin UI auto-renders an editable settings
	// form from it; the saved values are read back by owncast.config.get().
	Config map[string]ConfigField `json:"config,omitempty"`
}

// ConfigSchema returns a discovered plugin's manifest config schema (nil if the
// plugin is unknown or declares no config).
func (m *Manager) ConfigSchema(slug string) map[string]ConfigField {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if d, ok := m.discovered[slug]; ok {
		return d.Config
	}
	return nil
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
		approvedSet:  make(map[string][]string),
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
	// Capture an approval baseline for any plugin enabled before the
	// approved-permissions snapshot existed (older persisted state, or a
	// fresh first run that pre-seeded enabled names). Without this,
	// existing installs would see every permission as "pending" and the
	// plugin would auto-disable on the next start.
	if m.captureMissingApprovals() {
		if err := m.saveEnabledSet(); err != nil {
			fmt.Fprintf(os.Stderr, "persist enabled set: %v\n", err)
		}
	}
	// Auto-load anything in the enabled set that isn't already loaded
	// AND whose approved-permission set covers the current manifest.
	// Every reason a load doesn't complete is surfaced both in the
	// server log (so an operator inspecting `journalctl` sees it) and
	// in the DiscoveredEntry's LastError (so the admin UI's plugin
	// list shows a non-empty status hint). Without that, the entry
	// sits at "enabled, not loaded" with no indication of why.
	for name, enabled := range m.enabledSet {
		if !enabled {
			continue
		}
		if pending := m.pendingForLocked(name); len(pending) > 0 {
			log.Warnf("plugin %s: not loaded at startup; declares unapproved permissions %s; admin needs to re-approve in the plugin's detail view",
				name, strings.Join(pending, ", "))
			m.mu.Lock()
			if d, ok := m.discovered[name]; ok {
				d.PendingPermissions = pending
				d.LastError = "permissions need admin approval; re-enable from the plugin's detail view"
			}
			m.mu.Unlock()
			continue
		}
		if err := m.loadInternal(ctx, name); err != nil {
			log.Warnf("plugin %s: load failed at startup: %v", name, err)
			// loadInternal already set LastError on the entry; the
			// log line above adds the operator-visible trail.
			continue
		}
		if _, loaded := m.loaded[name]; !loaded {
			// Defensive: loadInternal returned nil but the plugin
			// didn't land in m.loaded. Shouldn't happen, but if it
			// ever does, the admin UI's "enabled, not loaded" entry
			// gets a non-empty LastError so the bug surfaces.
			log.Warnf("plugin %s: load returned no error but plugin is not in the loaded set", name)
			m.mu.Lock()
			if d, ok := m.discovered[name]; ok {
				d.LastError = "internal: load returned no error but plugin failed to register; check the server log"
			}
			m.mu.Unlock()
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
		m.notifyUnload(l)
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
		l, isLoaded := m.loaded[name]
		entry.Loaded = isLoaded
		if isLoaded {
			entry.AutoDisabled = l.IsDisabled()
			// Command metadata comes from register(), so it's only available
			// once the plugin is loaded.
			entry.Commands = l.Manifest.Commands
		}
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Slug < out[j].Slug })
	return out
}

// ActiveAuthGate returns the slug of the enabled auth.gate plugin (if any) and
// whether it is currently loaded. The viewer-auth gate middleware calls this on
// every request, so it scans under a read lock without allocating the full
// List(). Only one auth.gate plugin can be enabled at a time (see Enable), so
// the first match is authoritative.
func (m *Manager) ActiveAuthGate() (slug string, loaded bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for name, on := range m.enabledSet {
		if !on {
			continue
		}
		d := m.discovered[name]
		if d == nil || !hasPerm(d.Permissions, PermAuthGate) {
			continue
		}
		_, isLoaded := m.loaded[name]
		return d.Slug, isLoaded
	}
	return "", false
}

// hasPerm reports whether perms contains want.
func hasPerm(perms []string, want string) bool {
	for _, p := range perms {
		if p == want {
			return true
		}
	}
	return false
}

// Snapshot returns the currently-loaded, active plugins. Dispatcher and Server
// call this on every operation so changes from Enable/Disable take effect
// without restarting anything.
//
// Strike-disabled plugins are omitted: a plugin auto-disabled by the strike
// system must stop doing ALL work — events, HTTP/SSE, timers, and CSS/JS/tab
// contributions all dispatch off this snapshot — not just drop out of the
// filter chain. (The admin UI still sees them via List(), which reports
// AutoDisabled separately.)
func (m *Manager) Snapshot() []*Loaded {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Loaded, 0, len(m.loaded))
	for _, l := range m.loaded {
		if l.IsDisabled() {
			continue
		}
		out = append(out, l)
	}
	return out
}

// Enable marks a discovered plugin as enabled, captures the current
// manifest's permission set as the approved baseline (so any later
// expansion triggers a re-approval flow), persists the choice, and
// loads the plugin. No-op if already loaded.
func (m *Manager) Enable(ctx context.Context, name string) error {
	m.mu.Lock()
	d, ok := m.discovered[name]
	if !ok {
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
	// Only one auth.gate plugin may be enabled at a time — the viewer-auth gate
	// is singular (two gates would be ambiguous). Refuse to enable a second.
	if hasPerm(d.Permissions, PermAuthGate) {
		for other, on := range m.enabledSet {
			if !on || other == name {
				continue
			}
			if od := m.discovered[other]; od != nil && hasPerm(od.Permissions, PermAuthGate) {
				m.mu.Unlock()
				return fmt.Errorf("cannot enable %q: another authentication-gate plugin (%q) is already enabled; disable it first", d.Slug, od.Slug)
			}
		}
	}
	m.enabledSet[name] = true
	snapshot := append([]string(nil), d.Permissions...)
	sort.Strings(snapshot)
	m.approvedSet[name] = snapshot
	d.PendingPermissions = nil
	m.mu.Unlock()
	if err := m.saveEnabledSet(); err != nil {
		return fmt.Errorf("persist enabled set: %w", err)
	}
	err := m.loadInternal(ctx, name)
	return err
}

// validateUploadedPackage checks an uploaded .ocpkg's bytes, verifies the
// wasm can complete the same register()/manifest-agreement path used by a
// real load, and returns its manifest. Pulled out of Install so the
// validation steps and the (lock + write + scan) plumbing don't combine into
// one long function.
func validateUploadedPackage(ctx context.Context, env *HostEnv, packageBytes []byte) (*Manifest, error) {
	if len(packageBytes) == 0 {
		return nil, fmt.Errorf("upload is empty")
	}
	if len(packageBytes) > MaxUploadBytes {
		return nil, fmt.Errorf("upload exceeds %d-byte cap", MaxUploadBytes)
	}
	zr, err := zip.NewReader(bytes.NewReader(packageBytes), int64(len(packageBytes)))
	if err != nil {
		return nil, fmt.Errorf("not a valid .ocpkg: %w", err)
	}
	manifestBytes, err := readZipFile(zr, pkgManifestFilename, maxManifestEntryBytes)
	if err != nil {
		return nil, fmt.Errorf("missing manifest: %w", err)
	}
	manifest, err := ParseManifest(manifestBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid manifest: %w", err)
	}
	codeName, runtimeType, ok := detectPackageCode(zr)
	if !ok {
		return nil, fmt.Errorf("missing plugin code (expected one of plugin.js, plugin.py, plugin.wasm)")
	}
	codeBytes, err := readZipFile(zr, codeName, maxCodeEntryBytes)
	if err != nil {
		return nil, fmt.Errorf("missing plugin code (%s): %w", codeName, err)
	}
	// Preflight the package through the real load path before we write it into
	// the plugins directory. Without this, an .ocpkg whose manifest parses but
	// whose wasm fails register() or disagrees with the manifest appears to
	// "install" successfully from the catalog, then only surfaces as a
	// discovered-but-broken plugin later. The admin clicked Install, so the
	// operation should fail up front if the package cannot be loaded.
	// Extract assetsFS from the zip before calling loadFromBytes so register()
	// sees the same owncast_asset_read host function behavior as a real load.
	var assetsFS fs.FS
	if hasZipDir(zr, pkgAssetsPrefix) {
		if sub, err := fs.Sub(zr, strings.TrimSuffix(pkgAssetsPrefix, "/")); err == nil {
			assetsFS = sub
		}
	}
	// The plugin registry is process-global and keyed by slug: loadFromBytes
	// put()s the slug and the trailing Close() remove()s it. If a plugin of this
	// slug is already loaded and enabled, this throwaway preflight would delete
	// the live instance's identity, after which every gated host call from the
	// still-running plugin resolves as "not found" and silently no-ops until a
	// reload. Snapshot the existing entry and restore it after the preflight.
	//
	// ponytail: this still points the slug at the preflight identity *during*
	// the load itself, which is harmless in practice (same slug and KV
	// namespace; the new manifest's permission set differs only for the tens of
	// ms of a manual install). If that window ever matters, load the preflight
	// under an isolated registry key instead of the real slug.
	prevID, hadPrev := globalPluginRegistry.get(manifest.Slug)
	restoreLiveIdentity := func() {
		if hadPrev {
			globalPluginRegistry.put(prevID)
		}
	}
	loaded, err := loadFromBytes(ctx, env, manifestBytes, codeBytes, runtimeType, manifest.Slug, assetsFS)
	if err != nil {
		restoreLiveIdentity()
		return nil, err
	}
	loaded.Close(ctx)
	restoreLiveIdentity()
	return manifest, nil
}

// destPathForInstall picks the on-disk path for an installed plugin.
// Prefers an existing discovered entry's path so an update replaces the
// same file even if the previous file was named non-canonically;
// otherwise falls back to <pluginsDir>/<name>.ocpkg.
func (m *Manager) destPathForInstall(name string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if existing, ok := m.discovered[name]; ok && strings.HasSuffix(existing.Path, packageSuffix) {
		return existing.Path
	}
	return filepath.Join(m.pluginsDir, name+packageSuffix)
}

// atomicWritePackage writes packageBytes to destPath via a sibling temp
// file + rename, so a partial write never leaves a half-baked .ocpkg in
// the plugins directory for scan to trip over.
//
// destPath is computed from the uploaded plugin's manifest.slug
// (validated by slugPattern at parse time, so no slashes or dots), but
// the safety isn't visible at the rename site. Re-validating that the
// resolved destination sits under pluginsDir keeps the property local
// to this function and defends against a future caller that constructs
// destPath some other way.
func atomicWritePackage(pluginsDir, destPath string, packageBytes []byte) error {
	cleanDir, err := filepath.Abs(filepath.Clean(pluginsDir))
	if err != nil {
		return fmt.Errorf("resolve plugins dir: %w", err)
	}
	cleanDest, err := filepath.Abs(filepath.Clean(destPath))
	if err != nil {
		return fmt.Errorf("resolve install path: %w", err)
	}
	if !strings.HasPrefix(cleanDest, cleanDir+string(os.PathSeparator)) {
		return fmt.Errorf("refusing to install outside plugins dir: %s", destPath)
	}

	tmpFile, err := os.CreateTemp(pluginsDir, ".upload-*"+packageSuffix)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	cleanup := tmpPath
	defer func() {
		if cleanup != "" {
			_ = os.Remove(cleanup)
		}
	}()
	if _, err := tmpFile.Write(packageBytes); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("write upload: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close upload: %w", err)
	}
	if err := os.Rename(tmpPath, cleanDest); err != nil {
		return fmt.Errorf("install at %s: %w", cleanDest, err)
	}
	cleanup = ""
	return nil
}

// MaxUploadBytes is the cap on package bytes Install will accept. Larger
// uploads are rejected before any file is written so a malformed or hostile
// upload can't fill the plugins directory.
const MaxUploadBytes = 50 * 1024 * 1024

// Install validates a .ocpkg's contents, writes it into the plugins
// directory, and forces a scan so the new (or updated) plugin shows up
// in the next List() without waiting for the scan tick. The destination
// filename is derived from the manifest's name, not from the uploaded
// name, so an admin can't drop a file outside the plugins directory by
// abusing the upload filename, and a plugin update from a differently
// named .ocpkg ends up replacing the right file. Returns the discovered
// entry for the installed plugin.
func (m *Manager) Install(ctx context.Context, packageBytes []byte) (*DiscoveredEntry, error) {
	manifest, err := validateUploadedPackage(ctx, m.env, packageBytes)
	if err != nil {
		return nil, err
	}
	destPath := m.destPathForInstall(manifest.Slug)
	if err := atomicWritePackage(m.pluginsDir, destPath, packageBytes); err != nil {
		return nil, err
	}
	// Force an immediate scan so the upload appears in List() before
	// this call returns.
	if err := m.scan(ctx); err != nil {
		return nil, fmt.Errorf("scan after install: %w", err)
	}

	m.mu.RLock()
	entry, ok := m.discovered[manifest.Slug]
	var snapshot DiscoveredEntry
	if ok {
		snapshot = *entry
	}
	m.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("scan did not pick up the installed plugin %q", manifest.Slug)
	}
	return &snapshot, nil
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
		m.notifyUnload(loaded)
		loaded.Close(ctx)
	}
	return nil
}

// Uninstall unloads a plugin, deletes its file from the plugins
// directory, and clears the admin's persisted state (enabled flag and
// approved-permissions snapshot) for that plugin. The plugin's
// per-plugin config store is intentionally preserved so an accidental
// uninstall followed by a reinstall doesn't lose the streamer's
// settings.
func (m *Manager) Uninstall(ctx context.Context, name string) error {
	m.mu.Lock()
	entry, ok := m.discovered[name]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("plugin %q not discovered", name)
	}
	path := entry.Path
	loaded := m.loaded[name]
	delete(m.loaded, name)
	delete(m.discovered, name)
	delete(m.enabledSet, name)
	delete(m.approvedSet, name)
	m.mu.Unlock()

	if loaded != nil {
		m.notifyUnload(loaded)
		loaded.Close(ctx)
	}

	// Persist the cleared enabled + approved state before touching the
	// file, so a crash between the two leaves at worst a still-present
	// file with no in-memory references (the next scan will rediscover
	// it as a fresh install, which is the safer of the two outcomes).
	if err := m.saveEnabledSet(); err != nil {
		fmt.Fprintf(os.Stderr, "persist enabled set after uninstall: %v\n", err)
	}

	// Delete the .ocpkg (or the loose .wasm + sidecar manifest pair).
	// A missing file is fine; the scan that picks up our state change
	// would have dropped the entry anyway.
	if err := removePluginFiles(path); err != nil {
		return fmt.Errorf("delete plugin file: %w", err)
	}
	return nil
}

// removePluginFiles deletes the on-disk artifacts that correspond to a
// discovered plugin's path. For a .ocpkg this is just the one file. For
// the loose-files layout it's the .wasm, the sibling .manifest.json, the
// optional icon/instructions siblings, and the optional <base>-assets/
// directory.
func removePluginFiles(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	if strings.HasSuffix(path, ".wasm") {
		base := strings.TrimSuffix(path, ".wasm")
		_ = os.Remove(base + ".manifest.json")
		_ = os.Remove(base + ".icon.png")
		_ = os.Remove(base + "." + PluginInstructionsFilename)
		_ = os.RemoveAll(base + "-assets")
	}
	return nil
}

// Reload unloads and reloads a plugin. Plugin author rebuilt → admin
// triggers a reload to pick up the new wasm AND the new manifest without
// restarting the host or waiting for the next discovery scan tick.
// Plugin must currently be enabled (otherwise call Enable instead).
func (m *Manager) Reload(ctx context.Context, name string) error {
	m.mu.Lock()
	entry, ok := m.discovered[name]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("plugin %q not discovered", name)
	}
	if !m.enabledSet[name] {
		m.mu.Unlock()
		return fmt.Errorf("plugin %q is not enabled; use Enable to load it", name)
	}
	path := entry.Path
	loaded := m.loaded[name]
	delete(m.loaded, name)
	m.mu.Unlock()
	if loaded != nil {
		m.notifyUnload(loaded)
		loaded.Close(ctx)
	}

	// Re-read the manifest from disk so any author-side edits (new admin
	// pages, version bump, permission changes) surface immediately. A
	// failure here is logged but not fatal; the wasm reload below can
	// still succeed and the stale metadata will refresh on the next
	// scanLoop tick.
	if manifest, err := readManifestForPath(path); err == nil {
		m.mu.Lock()
		if existing, ok := m.discovered[name]; ok {
			existing.Version = manifest.Version
			existing.Description = manifest.Description
			existing.Permissions = manifest.Permissions
			existing.AdminPages = manifest.Admin.Pages
			existing.PendingPermissions = pendingPermissions(manifest.Permissions, m.approvedSet[name])
		}
		m.mu.Unlock()
	} else {
		fmt.Fprintf(os.Stderr, "plugin %q reload: re-reading manifest failed: %v\n", name, err)
	}

	return m.loadInternal(ctx, name)
}

// loadInternal performs the actual load and inserts into m.loaded. Assumes
// the caller has already verified the plugin is discovered + enabled.
// Refuses to load a plugin whose current manifest declares permissions
// the admin has not yet approved.
func (m *Manager) loadInternal(ctx context.Context, name string) error {
	m.mu.RLock()
	d, ok := m.discovered[name]
	approved := m.approvedSet[name]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %q not discovered", name)
	}
	if pending := pendingPermissions(d.Permissions, approved); len(pending) > 0 {
		err := fmt.Errorf("plugin %q declares new permissions that need admin approval: %s",
			name, strings.Join(pending, ", "))
		m.mu.Lock()
		if existing, ok := m.discovered[name]; ok {
			existing.PendingPermissions = pending
			existing.LastError = err.Error()
		}
		m.mu.Unlock()
		return err
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
// For loose-files plugins it reads PublicFS and AssetsFS from `public/`
// and `assets/` directories in the same parent directory as the wasm,
// matching the layout authors keep in their plugin source tree.
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
		parent := filepath.Dir(path)
		if pub := filepath.Join(parent, "public"); dirExists(pub) {
			loaded.PublicFS = os.DirFS(pub)
		}
		// AssetsFS is already set by LoadPlugin from the assets/ sibling dir.
		return loaded, nil
	}
	return nil, fmt.Errorf("unsupported plugin file: %s", path)
}

// dirExists reports whether path is an existing directory. Pulled out
// so loadByPath reads cleanly when it checks two sibling directories.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
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
		seen[manifest.Slug] = true

		hasIcon := hasPluginIcon(path)
		hasInstructions := hasPluginInstructions(path)

		m.mu.Lock()
		if existing, ok := m.discovered[manifest.Slug]; ok {
			// Already discovered; refresh manifest metadata in case it changed.
			existing.DisplayName = manifest.DisplayName
			existing.BotDisplayName = manifest.Bot.DisplayName
			existing.Version = manifest.Version
			existing.Description = manifest.Description
			existing.Permissions = manifest.Permissions
			existing.Path = path
			existing.AdminPages = manifest.Admin.Pages
			existing.HasIcon = hasIcon
			existing.HasInstructions = hasInstructions
			existing.PendingPermissions = pendingPermissions(manifest.Permissions, m.approvedSet[manifest.Slug])
			existing.AllowedHosts = manifest.Network.AllowedHosts
			existing.Config = manifest.Config
		} else {
			m.discovered[manifest.Slug] = &DiscoveredEntry{
				Slug:               manifest.Slug,
				DisplayName:        manifest.DisplayName,
				BotDisplayName:     manifest.Bot.DisplayName,
				Version:            manifest.Version,
				Description:        manifest.Description,
				Permissions:        manifest.Permissions,
				AllowedHosts:       manifest.Network.AllowedHosts,
				Path:               path,
				DiscoveredAt:       time.Now(),
				AdminPages:         manifest.Admin.Pages,
				HasIcon:            hasIcon,
				HasInstructions:    hasInstructions,
				PendingPermissions: pendingPermissions(manifest.Permissions, m.approvedSet[manifest.Slug]),
				Config:             manifest.Config,
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
			m.notifyUnload(loaded)
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

// Persistence: a tiny JSON file under the plugins directory listing the
// names the admin has enabled along with the permission set the admin
// approved at enable time. Survives restarts.

// EnabledStore persists per-plugin admin choices: the enabled set, plus
// the permission set the admin approved the last time each plugin was
// enabled. The default is a JSON file in the plugins directory
// (fileEnabledStore); integrated into Owncast it is backed by the native
// config datastore (see NewManagerWithStore).
type EnabledStore interface {
	// Load returns the persisted state. A fresh store returns an empty
	// StoreData and a nil error.
	Load() (StoreData, error)
	// Save replaces the persisted state atomically.
	Save(StoreData) error
}

// StoreData is the persisted admin state for the plugin manager.
type StoreData struct {
	// Enabled is the sorted set of plugin names the admin has enabled.
	Enabled []string
	// ApprovedPermissions maps plugin name to the sorted permission set
	// the admin approved the last time that plugin was enabled. If a
	// plugin appears in Enabled but not here, the next manager startup
	// captures its current manifest perms as the baseline.
	ApprovedPermissions map[string][]string
}

func (m *Manager) loadEnabledSet() error {
	data, err := m.enabledStore.Load()
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, name := range data.Enabled {
		m.enabledSet[name] = true
	}
	for name, perms := range data.ApprovedPermissions {
		clone := append([]string(nil), perms...)
		sort.Strings(clone)
		m.approvedSet[name] = clone
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
	approvals := make(map[string][]string, len(m.approvedSet))
	for name, perms := range m.approvedSet {
		clone := append([]string(nil), perms...)
		approvals[name] = clone
	}
	m.mu.RUnlock()
	sort.Strings(names)
	return m.enabledStore.Save(StoreData{Enabled: names, ApprovedPermissions: approvals})
}

// enabledFileContents is the on-disk JSON shape for fileEnabledStore.
type enabledFileContents struct {
	Enabled             []string            `json:"enabled"`
	ApprovedPermissions map[string][]string `json:"approvedPermissions,omitempty"`
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

func (s *fileEnabledStore) Load() (StoreData, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return StoreData{}, nil
		}
		return StoreData{}, err
	}
	var f enabledFileContents
	if err := json.Unmarshal(data, &f); err != nil {
		return StoreData{}, err
	}
	return StoreData(f), nil
}

func (s *fileEnabledStore) Save(d StoreData) error {
	contents := enabledFileContents(d)
	data, err := json.MarshalIndent(contents, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// pendingForLocked returns the permissions the current manifest declares
// that the admin has not yet approved. Caller need not hold the
// manager's lock; this helper acquires its own RLock.
func (m *Manager) pendingForLocked(name string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.discovered[name]
	if !ok {
		return nil
	}
	return pendingPermissions(d.Permissions, m.approvedSet[name])
}

// captureMissingApprovals fills in approval baselines for any plugin in
// the enabled set that doesn't already have one. Used on startup so an
// existing install (where the approved-permissions field didn't exist
// in the persisted state) doesn't suddenly see every permission as
// pending. After capturing, clears PendingPermissions on the affected
// entries so List() reflects the silent baseline. Returns true when any
// new baseline was captured (caller is expected to persist).
func (m *Manager) captureMissingApprovals() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	changed := false
	for name, enabled := range m.enabledSet {
		if !enabled {
			continue
		}
		if _, hasApproval := m.approvedSet[name]; hasApproval {
			continue
		}
		d, ok := m.discovered[name]
		if !ok {
			continue
		}
		snapshot := append([]string(nil), d.Permissions...)
		sort.Strings(snapshot)
		m.approvedSet[name] = snapshot
		d.PendingPermissions = nil
		changed = true
	}
	return changed
}

// pendingPermissions returns the permissions in manifestPerms that aren't
// in approved. Both inputs are treated as case-sensitive strings; the
// result is sorted and may be nil when there's no gap.
func pendingPermissions(manifestPerms, approved []string) []string {
	if len(manifestPerms) == 0 {
		return nil
	}
	approvedIdx := make(map[string]bool, len(approved))
	for _, p := range approved {
		approvedIdx[p] = true
	}
	var pending []string
	for _, p := range manifestPerms {
		if !approvedIdx[p] {
			pending = append(pending, p)
		}
	}
	sort.Strings(pending)
	return pending
}

// LoadPlugin loads a single plugin given explicit wasm and manifest paths
// (the loose-files layout). Used by the test runner so it shares the exact
// same load + register + validate path that production uses via Start.
// Assets are discovered from an `assets/` sibling directory when present and
// wired into the returned Loaded so the same asset-backed host APIs work for
// both loose-file and packaged plugins.
func LoadPlugin(ctx context.Context, env *HostEnv, artifactPath, manifestPath string) (*Loaded, error) {
	manifestBytes, err := os.ReadFile(manifestPath) //nolint:gosec // G304: plugin paths are admin-controlled, not user input
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", manifestPath, err)
	}
	artifactBytes, err := os.ReadFile(artifactPath) //nolint:gosec // G304: plugin paths are admin-controlled, not user input
	if err != nil {
		return nil, fmt.Errorf("read plugin code %s: %w", artifactPath, err)
	}
	// Infer the runtime from the artifact's extension (.js / .py / .wasm); the
	// author doesn't declare it in the manifest.
	runtimeType := runtimeForExt(artifactPath)
	base := filepath.Base(artifactPath)
	displayName := strings.TrimSuffix(base, filepath.Ext(base))
	var assetsFS fs.FS
	if as := filepath.Join(filepath.Dir(artifactPath), "assets"); dirExists(as) {
		assetsFS = os.DirFS(as)
	}
	loaded, err := loadFromBytes(ctx, env, manifestBytes, artifactBytes, runtimeType, displayName, assetsFS)
	if err != nil {
		return nil, err
	}
	loaded.WasmPath = artifactPath
	return loaded, nil
}

// runtimeForExt maps a code artifact's file extension to its runtime:
// .js → javascript, .py → python, anything else (.wasm) → wasm.
func runtimeForExt(path string) string {
	switch filepath.Ext(path) {
	case ".js":
		return RuntimeJavaScript
	case ".py":
		return RuntimePython
	default:
		return RuntimeWasm
	}
}

// loadFromBytes is the shared core of LoadPlugin and LoadPackage. artifactBytes
// is the plugin's wasm (for "wasm"/legacy) or its source (for the shared-engine
// runtimes "javascript"/"python"); runtime — inferred by the caller from the
// code artifact (its filename/extension), not authored in the manifest —
// selects which. It is authoritative over any manifest-declared type.
func loadFromBytes(ctx context.Context, env *HostEnv, manifestBytes, artifactBytes []byte, runtimeType, displayName string, assetsFS fs.FS) (*Loaded, error) {
	manifest, err := ParseManifest(manifestBytes)
	if err != nil {
		return nil, err
	}
	manifest.Type = runtimeType

	p, err := instantiate(ctx, env, manifest, manifestBytes, artifactBytes, displayName)
	if err != nil {
		return nil, err
	}

	// Register the plugin's identity so shared host functions can resolve it at
	// call time (by the slug stashed in the instance's Extism config). Done
	// before register() because the engine bootstrap evals author top-level
	// code during register(), which may call host functions. On any failure
	// below, drop the identity again.
	globalPluginRegistry.put(buildIdentity(env, manifest, assetsFS))
	loaded := false
	defer func() {
		if !loaded {
			globalPluginRegistry.remove(manifest.Slug)
		}
	}()

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
	if err := requireChatFilterPermission(manifest, runtime.Subscriptions); err != nil {
		_ = p.Close(ctx)
		return nil, err
	}
	manifest.Subscriptions = runtime.Subscriptions
	// Command metadata is derived by the SDK and reported via register() (like
	// subscriptions); the host aggregates it for !help. Informational only, so
	// it isn't part of AgreesWith's security checks.
	manifest.Commands = runtime.Commands

	var adminGlobs []glob.Glob
	var adminPaths []string
	for _, page := range manifest.Admin.Pages {
		g, err := glob.Compile(page.Path)
		if err != nil {
			_ = p.Close(ctx)
			return nil, fmt.Errorf("manifest.admin.pages: invalid path glob %q: %w", page.Path, err)
		}
		adminGlobs = append(adminGlobs, g)
		adminPaths = append(adminPaths, page.Path)
	}

	loaded = true
	return &Loaded{Manifest: manifest, plugin: p, adminGlobs: adminGlobs, adminPaths: adminPaths, AssetsFS: assetsFS}, nil
}

// instantiate creates the extism plugin instance for a manifest: a per-plugin
// instance of the shared engine for "js"/"python", or a self-contained module
// for "wasm"/legacy. Either way the result is a *extism.Plugin with its
// identity slug stashed in Config so shared host functions can resolve the
// caller, and its per-plugin network scope applied.
func instantiate(ctx context.Context, env *HostEnv, manifest *Manifest, manifestBytes, artifactBytes []byte, displayName string) (*extism.Plugin, error) {
	// Give the guest the real host wall and monotonic clocks (wazero's default
	// is a frozen 2022 clock). Nanosleep is deliberately NOT wired so a plugin
	// can't block inside a call and burn its timeout budget.
	moduleConfig := wazero.NewModuleConfig().WithSysWalltime().WithSysNanotime()

	if manifest.usesSharedEngine() {
		engine, err := compiledEngines.get(ctx, env, manifest.Type)
		if err != nil {
			return nil, err
		}
		inst, err := engine.Instance(ctx, extism.PluginInstanceConfig{ModuleConfig: moduleConfig})
		if err != nil {
			return nil, fmt.Errorf("instantiate %s engine: %w", manifest.Type, err)
		}
		// Inject the per-plugin script + manifest the engine bootstrap reads at
		// runtime, plus the slug shared host functions resolve identity by.
		inst.Config = map[string]string{
			configKeySlug:     manifest.Slug,
			configKeyScript:   string(artifactBytes),
			configKeyManifest: string(manifestBytes),
		}
		applyNetworkScope(inst, manifest)
		return inst, nil
	}

	// Legacy self-contained wasm: each plugin compiles into its own runtime.
	// Share the compilation cache so reloads / repeated loads of the same
	// module don't recompile its embedded engine.
	extismManifest := extism.Manifest{
		Wasm:    []extism.Wasm{extism.WasmData{Data: artifactBytes, Name: displayName}},
		Timeout: 10_000, // milliseconds; enables Wazero's WithCloseOnContextDone
		Memory: &extism.ManifestMemory{
			MaxPages:             MaxWasmPages,
			MaxHttpResponseBytes: MaxExtismHTTPResponseBytes,
			MaxVarBytes:          MaxExtismVarBytes,
		},
	}
	for _, perm := range manifest.Permissions {
		if perm == PermNetworkFetch {
			extismManifest.AllowedHosts = append([]string(nil), manifest.Network.AllowedHosts...)
			break
		}
	}
	pc := extism.PluginConfig{
		EnableWasi:    true,
		ModuleConfig:  moduleConfig,
		RuntimeConfig: wazero.NewRuntimeConfig().WithCompilationCache(sharedCompilationCache()),
	}
	p, err := extism.NewPlugin(ctx, extismManifest, pc, BuildHostFunctions(env))
	if err != nil {
		return nil, fmt.Errorf("instantiate wasm: %w", err)
	}
	// Even self-contained plugins now use the shared, identity-resolving host
	// functions, so stash the slug for the registry lookup.
	p.Config = map[string]string{configKeySlug: manifest.Slug}
	return p, nil
}

// applyNetworkScope sets a shared-engine instance's per-plugin AllowedHosts.
// The compiled engine carries no host allowlist (it's one engine for many
// plugins); each instance's network scope comes from its own manifest. Manifest
// validation already requires AllowedHosts to be non-empty when network.fetch
// is granted.
func applyNetworkScope(inst *extism.Plugin, manifest *Manifest) {
	for _, perm := range manifest.Permissions {
		if perm == PermNetworkFetch {
			inst.AllowedHosts = append([]string(nil), manifest.Network.AllowedHosts...)
			return
		}
	}
}

// buildIdentity assembles the call-time identity shared host functions resolve
// for a loaded plugin (see registry.go).
func buildIdentity(env *HostEnv, manifest *Manifest, assetsFS fs.FS) *pluginIdentity {
	id := &pluginIdentity{
		slug:        manifest.Slug,
		granted:     stringSet(manifest.Permissions),
		chatDisplay: manifest.ChatDisplayName(),
		assetsFS:    assetsFS,
		manifest:    manifest,
	}
	if env != nil && env.KV != nil {
		id.kvNamespace = env.KV.Namespace(manifest.Slug)
	}
	return id
}

// requireChatFilterPermission rejects a runtime registration that
// subscribes to filterChatMessage without declaring the chat.filter
// permission. Modifying or dropping every chat message is a meaningful
// side-effect, so an admin must opt in at install time by granting it.
func requireChatFilterPermission(manifest *Manifest, subs Subscriptions) error {
	if !manifest.hasPermission(PermChatFilter) {
		for _, s := range subs.Filter {
			if s.Event == EventChatMessageReceived {
				return fmt.Errorf(
					"plugin subscribes to filterChatMessage but does not declare "+
						"the %q permission. Add %q to the manifest's permissions "+
						"so an admin can see at install time that this plugin "+
						"reads or modifies every chat message",
					PermChatFilter, PermChatFilter,
				)
			}
		}
	}
	return nil
}

// PluginIconFilename is the conventional filename a plugin uses to ship
// an icon. For .ocpkg plugins it lives at the root of the zip; for loose
// .wasm files it sits next to the wasm as "<base>.icon.png" so multiple
// plugins in one directory don't collide.
const PluginIconFilename = "icon.png"

// hasPluginIcon reports whether the plugin at path bundles an icon.
// Called per-scan, so it has to be cheap — a single zip central-
// directory read for packaged plugins, or a stat for loose files.
func hasPluginIcon(path string) bool {
	switch {
	case strings.HasSuffix(path, packageSuffix):
		zr, err := zip.OpenReader(path)
		if err != nil {
			return false
		}
		defer zr.Close()
		for _, f := range zr.File {
			if f.Name == PluginIconFilename {
				return true
			}
		}
		return false
	case strings.HasSuffix(path, ".wasm"):
		iconPath := strings.TrimSuffix(path, ".wasm") + ".icon.png"
		_, err := os.Stat(iconPath)
		return err == nil
	}
	return false
}

// IconBytes returns the raw bytes of the plugin's icon.png, or an error
// if the plugin isn't discovered or doesn't ship one. Callers should
// serve the bytes with content-type image/png.
//
// Read on-demand from disk on each call: icons are tiny, requests are
// rare (admin UI only), and reading fresh means an admin who swaps the
// icon between two scans gets the new one without a host restart.
func (m *Manager) IconBytes(name string) ([]byte, error) {
	m.mu.RLock()
	entry, ok := m.discovered[name]
	if !ok {
		m.mu.RUnlock()
		return nil, fmt.Errorf("plugin %q not discovered", name)
	}
	path := entry.Path
	hasIcon := entry.HasIcon
	m.mu.RUnlock()
	if !hasIcon {
		return nil, fmt.Errorf("plugin %q has no icon", name)
	}

	switch {
	case strings.HasSuffix(path, packageSuffix):
		zr, err := zip.OpenReader(path)
		if err != nil {
			return nil, fmt.Errorf("open package %s: %w", path, err)
		}
		defer zr.Close()
		for _, f := range zr.File {
			if f.Name != PluginIconFilename {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("open icon entry: %w", err)
			}
			defer rc.Close()
			// Bound the decompressed read so a zip-bomb icon entry can't OOM the host.
			return io.ReadAll(io.LimitReader(rc, maxIconEntryBytes))
		}
		return nil, fmt.Errorf("plugin %q icon entry vanished between scan and read", name)
	case strings.HasSuffix(path, ".wasm"):
		iconPath := strings.TrimSuffix(path, ".wasm") + ".icon.png"
		return os.ReadFile(iconPath) //nolint:gosec // G304: plugin paths are admin-controlled, not user input
	}
	return nil, fmt.Errorf("plugin %q: unsupported file type for icon lookup", name)
}

// PluginInstructionsFilename is the conventional filename a plugin uses to
// ship author-written usage instructions. For .ocpkg plugins it lives at
// the root of the zip; for loose .wasm files it sits next to the wasm as
// "<base>.INSTRUCTIONS.md" so multiple plugins in one directory don't
// collide. The content is markdown, rendered by the admin UI.
const PluginInstructionsFilename = "INSTRUCTIONS.md"

// hasPluginInstructions reports whether the plugin at path bundles an
// INSTRUCTIONS.md. Called per-scan, so it has to be cheap — a single zip
// central-directory read for packaged plugins, or a stat for loose files.
func hasPluginInstructions(path string) bool {
	switch {
	case strings.HasSuffix(path, packageSuffix):
		zr, err := zip.OpenReader(path)
		if err != nil {
			return false
		}
		defer zr.Close()
		for _, f := range zr.File {
			if f.Name == PluginInstructionsFilename {
				return true
			}
		}
		return false
	case strings.HasSuffix(path, ".wasm"):
		instructionsPath := strings.TrimSuffix(path, ".wasm") + "." + PluginInstructionsFilename
		_, err := os.Stat(instructionsPath)
		return err == nil
	}
	return false
}

// InstructionsBytes returns the raw markdown of the plugin's
// INSTRUCTIONS.md, or an error if the plugin isn't discovered or doesn't
// ship one. Callers should serve the bytes with content-type text/markdown.
//
// Read on-demand from disk on each call: instructions are small, requests
// are rare (admin UI only), and reading fresh means an admin who swaps the
// file between two scans gets the new content without a host restart.
func (m *Manager) InstructionsBytes(name string) ([]byte, error) {
	m.mu.RLock()
	entry, ok := m.discovered[name]
	if !ok {
		m.mu.RUnlock()
		return nil, fmt.Errorf("plugin %q not discovered", name)
	}
	path := entry.Path
	has := entry.HasInstructions
	m.mu.RUnlock()
	if !has {
		return nil, fmt.Errorf("plugin %q has no instructions", name)
	}

	switch {
	case strings.HasSuffix(path, packageSuffix):
		zr, err := zip.OpenReader(path)
		if err != nil {
			return nil, fmt.Errorf("open package %s: %w", path, err)
		}
		defer zr.Close()
		for _, f := range zr.File {
			if f.Name != PluginInstructionsFilename {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("open instructions entry: %w", err)
			}
			defer rc.Close()
			// Bound the decompressed read so a zip-bomb entry can't OOM the host.
			return io.ReadAll(io.LimitReader(rc, maxInstructionsEntryBytes))
		}
		return nil, fmt.Errorf("plugin %q instructions entry vanished between scan and read", name)
	case strings.HasSuffix(path, ".wasm"):
		instructionsPath := strings.TrimSuffix(path, ".wasm") + "." + PluginInstructionsFilename
		return os.ReadFile(instructionsPath) //nolint:gosec // G304: plugin paths are admin-controlled, not user input
	}
	return nil, fmt.Errorf("plugin %q: unsupported file type for instructions lookup", name)
}
