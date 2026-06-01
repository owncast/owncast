package pluginhost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications/browser"
	"github.com/owncast/owncast/notifications/discord"
	"github.com/owncast/owncast/persistence/authrepository"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/dispatcher"
	"github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/services/plugins/kv"
	"github.com/owncast/owncast/services/stream"
)

// pluginsEnabledConfigKey is the datastore key under which the set of
// admin-enabled plugin names is persisted.
const pluginsEnabledConfigKey = "plugins.enabled"

// jsonErrorKey is the JSON map key used to surface a single human-readable
// error string back to the admin UI in error responses.
const jsonErrorKey = "error"

// pluginsApprovedPermsConfigKey is the datastore key under which the
// per-plugin admin-approved permission snapshots are persisted (as a
// JSON-encoded map keyed by plugin name).
const pluginsApprovedPermsConfigKey = "plugins.approvedPermissions"

// Deps bundles the Owncast services the plugin runtime adapts into
// HostEnv host functions. Everything here is already constructed in the main
// composition root and passed in by reference.
type Deps struct {
	Datastore               *datastore.Datastore
	Chat                    *chat.Service
	Stream                  *stream.Service
	Activitypub             *activitypub.Service
	Events                  *dispatcher.Dispatcher
	ConfigRepository        configrepository.ConfigRepository
	UserRepository          userrepository.UserRepository
	AuthRepository          *authrepository.SqlAuthRepository
	NotificationsRepository notificationsrepository.NotificationsRepository
	ChatMessageRepository   chatmessagerepository.ChatMessageRepository

	// RequireAdminAuth wraps a handler with the host's admin Basic Auth
	// middleware. The plugin host uses it on management endpoints and on
	// manifest-declared admin paths inside the plugin static server, so
	// auth (realm, CORS, credential check) comes from one implementation
	// instead of being duplicated here.
	RequireAdminAuth func(http.HandlerFunc) http.HandlerFunc
	// IsAdminRequest is the predicate behind RequireAdminAuth, exposed as a
	// boolean for paths that don't reject unauthenticated requests but need
	// to know whether the caller is an authenticated admin (e.g. the
	// req.authenticated field passed to a plugin's HTTP handler).
	IsAdminRequest func(*http.Request) bool
}

// Host owns the running plugin runtime: the manager (discovery +
// enable/disable lifecycle), the HTTP handler that serves /plugins/<name>/*,
// and the SSE hub backing host-owned event streams.
type Host struct {
	manager          *plugins.Manager
	server           *plugins.Server
	sse              *plugins.SSEHub
	configRepository configrepository.ConfigRepository
	// requireAdminAuth is the host's admin Basic Auth middleware, plumbed
	// in from main.go so the management API and the plugin static server
	// share one credential check (realm, logging, etc.) with the rest of
	// the admin API.
	requireAdminAuth func(http.HandlerFunc) http.HandlerFunc
	// kv is the per-plugin config store the runtime hands plugins. The
	// host reads each plugin's reserved actionsOverrideConfigKey from
	// this store to let a plugin override its manifest-declared action
	// buttons at runtime (e.g. an admin page that lets the streamer
	// rename buttons without rebuilding the plugin).
	kv kv.Store
	// tickCancel stops the once-a-second tick goroutine when the host stops.
	tickCancel context.CancelFunc
}

// runtimeActionsConfigKey is the reserved key inside a plugin's own
// config namespace that holds action buttons the plugin has added at
// runtime via owncast.actions.add. The host's Actions() reader returns
// manifest.actions ++ this list on every /api/config request.
const runtimeActionsConfigKey = "owncast.actions"

// Handler is the http.Handler for /plugins/<name>/* (static assets, dynamic
// on_http_request, and the reserved _sse endpoint).
func (p *Host) Handler() http.Handler { return p.server }

// Stop closes all loaded plugins.
func (p *Host) Stop(ctx context.Context) {
	if p.tickCancel != nil {
		p.tickCancel()
	}
	p.manager.Stop(ctx)
}

// Actions returns every action-button currently contributed by loaded
// plugins, projected into the shape Owncast's web config uses for
// admin-defined externalActions. Empty when no plugins declare actions
// (or no plugins are loaded). Called from handlers.GetWebConfig so the
// viewer-facing config endpoint surfaces plugin and admin actions in a
// single list, without each consumer needing to know about plugins.
//
// For each plugin, the effective list is manifest.Actions ++ any
// buttons the plugin has added at runtime via owncast.actions.add
// (read from runtimeActionsConfigKey in the plugin's own config).
func (p *Host) Actions() []models.ExternalAction {
	if p == nil || p.manager == nil {
		return nil
	}
	loaded := p.manager.Snapshot()
	var out []models.ExternalAction
	for _, l := range loaded {
		actions := p.actionsForPlugin(l)
		for _, a := range actions {
			out = append(out, models.ExternalAction{
				URL:            a.Url,
				HTML:           a.Html,
				Title:          a.Title,
				Description:    a.Description,
				Icon:           a.Icon,
				Color:          a.Color,
				OpenExternally: a.OpenExternally,
			})
		}
	}
	return out
}

// StylesContent returns the concatenated CSS bytes contributed by
// every currently-loaded plugin. Each plugin's manifest.styles files
// are read from the plugin's AssetsFS (the static-asset root inside
// the ocpkg) and joined with delimiter comments naming the source
// plugin, so devtools "view source" still attributes a rule to the
// plugin that shipped it. The bytes are folded into the existing
// customStyles config field at request time, so the viewer page
// renders a single inline <style> block instead of one <link> tag
// per plugin asset. Disabled or not-loaded plugins contribute
// nothing.
func (p *Host) StylesContent() []byte {
	return p.readManifestAssets(
		func(m *plugins.Manifest) []string { return m.Styles },
		"/* plugin: %s — %s */\n",
	)
}

// ScriptsContent mirrors StylesContent for manifest.scripts: returns
// the concatenated JavaScript bytes contributed by every loaded
// plugin, with `// plugin: <slug>` delimiters between contributions.
// Appended to the admin's customJavascript by the /customjavascript
// handler so the viewer page loads one script tag for both sources.
func (p *Host) ScriptsContent() []byte {
	return p.readManifestAssets(
		func(m *plugins.Manifest) []string { return m.Scripts },
		"// plugin: %s — %s\n",
	)
}

// PageContent returns the concatenated HTML bytes contributed by
// every loaded plugin's manifest.extraPageContent file, each preceded
// by an `<!-- plugin: <slug> ... -->` delimiter for in-page
// attribution. The /api/config handler prepends these bytes to the
// admin's rendered extraPageContent so plugin HTML lands at the top
// of the extra-content block. Markdown rendering is skipped for
// plugin HTML so the markdown processor can't mangle it.
func (p *Host) PageContent() []byte {
	return p.readManifestAssets(
		func(m *plugins.Manifest) []string {
			if m.ExtraPageContent == "" {
				return nil
			}
			return []string{m.ExtraPageContent}
		},
		"<!-- plugin: %s — %s -->\n",
	)
}

// Tabs returns every viewer-page tab contributed by loaded plugins
// via manifest.tabs. Each tab's content file is read from the
// plugin's assets/ directory once per call; tabs whose file can't be
// read are skipped (logged) so a broken file doesn't take down the
// rest of the list. /api/config emits the result as `pluginTabs`,
// which the viewer page renders alongside the built-in tabs.
func (p *Host) Tabs() []models.PluginTab {
	if p == nil || p.manager == nil {
		return nil
	}
	var out []models.PluginTab
	for _, l := range p.manager.Snapshot() {
		if l == nil || l.Manifest == nil || l.AssetsFS == nil {
			continue
		}
		slug := l.Manifest.Slug
		pluginPrefix := "/plugins/" + slug + "/"
		for _, tab := range l.Manifest.Tabs {
			relPath := strings.TrimPrefix(tab.Content, pluginPrefix)
			data, err := fs.ReadFile(l.AssetsFS, relPath)
			if err != nil {
				log.Warnf("plugin %s: skipping tab %q (%s): %v", slug, tab.Title, relPath, err)
				continue
			}
			out = append(out, models.PluginTab{
				Slug:  slug,
				Title: tab.Title,
				HTML:  string(data),
			})
		}
	}
	return out
}

// readManifestAssets walks every loaded plugin, calls pick() to get
// the manifest entries to include (Styles or Scripts), and reads each
// file from the plugin's AssetsFS. delimiter is a Printf format
// taking (slug, relPath) used as a per-contribution header. Files
// that can't be read are skipped (logged) so a broken asset doesn't
// take down the rest of the bundle.
//
// The manifest entries arrive as rewritten plugin-namespace URLs
// (e.g. "/plugins/styles-demo/theme.css"); the file inside the
// AssetsFS lives at the path beneath the plugin prefix
// ("theme.css"), so the lookup is a prefix strip.
func (p *Host) readManifestAssets(pick func(*plugins.Manifest) []string, delimiter string) []byte {
	if p == nil || p.manager == nil {
		return nil
	}
	var buf bytes.Buffer
	for _, l := range p.manager.Snapshot() {
		if l == nil || l.Manifest == nil || l.AssetsFS == nil {
			continue
		}
		entries := pick(l.Manifest)
		if len(entries) == 0 {
			continue
		}
		pluginPrefix := "/plugins/" + l.Manifest.Slug + "/"
		for _, entry := range entries {
			relPath := strings.TrimPrefix(entry, pluginPrefix)
			data, err := fs.ReadFile(l.AssetsFS, relPath)
			if err != nil {
				log.Warnf("plugin %s: skipping asset %s: %v", l.Manifest.Slug, relPath, err)
				continue
			}
			fmt.Fprintf(&buf, delimiter, l.Manifest.Slug, relPath)
			buf.Write(data)
			if len(data) > 0 && data[len(data)-1] != '\n' {
				buf.WriteByte('\n')
			}
		}
	}
	return buf.Bytes()
}

// actionsForPlugin returns the plugin's manifest action list with any
// runtime additions appended. Kept separate so Actions() reads cleanly
// and the append logic can be tested independently.
func (p *Host) actionsForPlugin(l *plugins.Loaded) []plugins.ActionButton {
	combined := l.Manifest.Actions
	if p.kv == nil {
		return combined
	}
	raw, err := p.kv.Namespace(l.Manifest.Slug).Get(runtimeActionsConfigKey)
	if err != nil || len(raw) == 0 {
		return combined
	}
	var extras []plugins.ActionButton
	if err := json.Unmarshal(raw, &extras); err != nil {
		return combined
	}
	return append(combined, extras...)
}

// New builds the HostEnv from Owncast services, constructs and
// starts the plugin manager, and returns the assembled host. Plugins are
// optional infrastructure: callers should log and continue on error rather
// than aborting startup.
func New(ctx context.Context, deps Deps) (*Host, error) {
	pluginsDir := filepath.Join(config.DataDirectory, "plugins")
	if err := os.MkdirAll(pluginsDir, 0o700); err != nil {
		return nil, fmt.Errorf("create plugins directory: %w", err)
	}

	env := &plugins.HostEnv{KV: newDatastoreKVStore(deps.Datastore)}
	wirePluginHostEnv(env, deps)

	sseHub := plugins.NewSSEHub()
	env.SSE = sseHub

	enabledStore := &configEnabledStore{datastore: deps.Datastore}
	manager := plugins.NewManagerWithStore(pluginsDir, env, enabledStore)
	if err := manager.Start(ctx); err != nil {
		return nil, fmt.Errorf("start plugin manager: %w", err)
	}

	// Boot summary — show every discovered plugin with its status and the
	// permissions it asks for. Per-plugin "loaded" lines come from the
	// manager itself; this gives the full picture (including disabled ones)
	// so the admin doesn't need to hit the API to see what's installed.
	logPluginSummary(manager.List())

	// Emit delivers plugin-published custom events to other plugins'
	// subscribers. Wired post-Start because it reads the live plugin set.
	pluginDispatcher := plugins.NewLiveDispatcher(manager.Snapshot)
	env.Emit = pluginDispatcher.Dispatch

	// Subscribe to the shared event dispatcher: deliver Owncast's events
	// (chat, stream lifecycle, moderation, …) to plugins' notify handlers,
	// and run plugin filterChatMessage handlers on inbound chat messages
	// before they're broadcast.
	deps.Events.AddListener(newPluginEventListener(pluginDispatcher))
	deps.Events.AddFilter(newPluginChatFilter(pluginDispatcher))

	// Host-driven timers: plugins can't setTimeout in the sandbox, so
	// owncast.timer.* asks the host to schedule callbacks. The hub resolves a
	// plugin slug to its live instance to call back; cancelling a plugin's
	// timers on unload is wired through the manager's onUnload hook.
	timerHub := plugins.NewTimerHub(func(slug string) *plugins.Loaded {
		for _, l := range manager.Snapshot() {
			if l.Manifest.Slug == slug {
				return l
			}
		}
		return nil
	})
	env.Timer = timerHub
	manager.SetOnUnload(timerHub.CancelForPlugin)

	// Fire a once-a-second tick to plugins that subscribe (onTick), and which
	// also drives nothing else — host-scheduled timers run independently. The
	// goroutine is stopped when the host stops.
	tickCtx, tickCancel := context.WithCancel(context.Background())
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			select {
			case <-tickCtx.Done():
				return
			case <-t.C:
				pluginDispatcher.Notify(tickCtx, plugins.EventTick, plugins.TickEvent{Now: time.Now().UnixMilli()})
			}
		}
	}()

	server := plugins.NewLiveServer(manager.Snapshot)
	server.SSE = sseHub
	server.IsAuthenticated = env.IsAuthenticated
	server.RequireAdmin = deps.RequireAdminAuth
	server.GetRequestUser = env.GetRequestUser

	return &Host{
		manager:          manager,
		server:           server,
		sse:              sseHub,
		configRepository: deps.ConfigRepository,
		requireAdminAuth: deps.RequireAdminAuth,
		kv:               env.KV,
		tickCancel:       tickCancel,
	}, nil
}

// AdminHandler returns the HTTP handler for plugin management:
//
//	GET  /api/admin/plugins                       list discovered plugins (admin)
//	POST /api/admin/plugins/<name>/enable|disable|reload  toggle a plugin (admin)
//	GET  /api/admin/plugins/<name>/instructions   bundled INSTRUCTIONS.md markdown (admin)
//	GET  /api/plugins/actions                     merged action-button list (public)
//
// It's mounted by the router on the outer mux so it sits beside, not inside,
// the OpenAPI-generated /api router.
func (p *Host) AdminHandler() http.Handler {
	mux := http.NewServeMux()

	// Admin-protected routes go through the host's RequireAdminAuth
	// middleware — same realm, same CORS handling, same logging as the
	// rest of /api/admin/*.
	requireAdmin := p.requireAdminAuth
	if requireAdmin == nil {
		// Defensive: a caller built Host without wiring the middleware.
		// Treat as "always reject" so a misconfigured boot doesn't expose
		// the management API.
		requireAdmin = func(http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "admin auth not configured", http.StatusUnauthorized)
			}
		}
	}

	mux.Handle("/api/admin/plugins", requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSONResponse(w, http.StatusOK, p.manager.List())
		case http.MethodPost:
			p.handlePluginUpload(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Registry browse + install. Sibling prefix so the catch-all
	// /api/admin/plugins/ trailing-slash matcher below doesn't swallow
	// them. Registered as a prefix (trailing slash) and dispatched
	// internally so /list and /list/ both work — the admin frontend's
	// Next.js dev server 308-redirects every URL to its slash variant
	// because of trailingSlash:true in next.config.js.
	mux.Handle("/api/admin/plugin-registry/", requireAdmin(p.handleRegistryRoute))

	mux.Handle("/api/admin/plugins/", requireAdmin(p.handlePluginAction))

	mux.HandleFunc("/api/plugins/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		actions := make([]plugins.ActionButton, 0)
		for _, loaded := range p.manager.Snapshot() {
			actions = append(actions, loaded.Manifest.Actions...)
		}
		writeJSONResponse(w, http.StatusOK, actions)
	})

	// Plugin icon. Public on purpose: the admin UI renders these in the
	// plugin list and the sidebar via an <img> tag, which doesn't carry
	// the Authorization header, and the bytes themselves aren't
	// sensitive. The host reads icon.png directly from the package on
	// each request, so an admin who swaps the icon gets the new one
	// without reloading the plugin.
	mux.HandleFunc("/api/plugins/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		const prefix = "/api/plugins/"
		const suffix = "/icon"
		rest := strings.TrimPrefix(r.URL.Path, prefix)
		if !strings.HasSuffix(rest, suffix) {
			http.NotFound(w, r)
			return
		}
		name := strings.TrimSuffix(rest, suffix)
		if name == "" || strings.ContainsAny(name, "/") {
			http.NotFound(w, r)
			return
		}
		data, err := p.manager.IconBytes(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-cache")
		// gosec G705: not XSS. We've set Content-Type to image/png and the
		// bytes come from the plugin's bundled icon.png, written by the
		// admin who installed the plugin (not user-controlled input).
		_, _ = w.Write(data) //nolint:gosec // G705 false positive: image/png body
	})

	return mux
}

func writeJSONResponse(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// logPluginSummary writes the boot inventory of plugins to the Owncast log:
// every discovered plugin, its enabled/disabled (and loaded) status, and the
// permissions it declares. Errors from a failed load are surfaced as warnings
// so an admin sees them in the console without checking the admin API.
func logPluginSummary(entries []plugins.DiscoveredEntry) {
	log.Infof("plugins: %d discovered", len(entries))
	for _, e := range entries {
		status := "disabled"
		if e.Enabled {
			status = "enabled"
			if !e.Loaded {
				status = "enabled (load failed)"
			}
		}
		permPart := "with no permissions"
		if len(e.Permissions) > 0 {
			permPart = "with permissions: " + strings.Join(e.Permissions, ", ")
		}
		nameDisplay := e.DisplayName
		if nameDisplay == "" {
			nameDisplay = e.Slug
		}
		log.Infof("  - %s [%s] v%s %s %s", nameDisplay, e.Slug, e.Version, status, permPart)
		if e.LastError != "" {
			log.Warnf("    last error: %s", e.LastError)
		}
	}
}

// handlePluginAction dispatches POST /api/admin/plugins/<slug>/<action>
// to the matching Manager operation (enable, disable, reload, uninstall).
// The URL segment is the plugin's slug, not its display name; Manager
// keys every operation on slug.
func (p *Host) handlePluginAction(w http.ResponseWriter, r *http.Request) {
	slug, action, ok := strings.Cut(strings.TrimPrefix(r.URL.Path, "/api/admin/plugins/"), "/")
	if !ok || slug == "" || action == "" {
		http.Error(w, "expected /<slug>/<action>", http.StatusBadRequest)
		return
	}

	// GET /api/admin/plugins/<slug>/instructions serves the bundled
	// INSTRUCTIONS.md as raw markdown; the admin UI renders it in a details
	// tab. Read fresh from disk on each request so a swapped file shows up
	// without reloading the plugin.
	if action == "instructions" {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		data, err := p.manager.InstructionsBytes(slug)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		// gosec G705 (XSS via taint): the route is admin-authenticated,
		// the response Content-Type is text/markdown (not text/html), and
		// the bytes are the admin's own uploaded INSTRUCTIONS.md content.
		// The admin UI renders them through ReactMarkdown, which sanitizes
		// before insertion.
		_, _ = w.Write(data) //nolint:gosec // G705
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err, known := p.dispatchPluginAction(r.Context(), slug, action)
	if !known {
		http.Error(w, "unknown action; expected enable, disable, reload, or uninstall", http.StatusBadRequest)
		return
	}
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{jsonErrorKey: err.Error()})
		return
	}
	if action == "uninstall" {
		log.Infof("plugin %q uninstalled by admin", slug)
	}
	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "slug": slug, "action": action})
}

// dispatchPluginAction calls the Manager method for one of the admin
// actions. Returns (operationError, true) for a known action, or
// (nil, false) when action is unrecognized.
func (p *Host) dispatchPluginAction(ctx context.Context, slug, action string) (error, bool) {
	switch action {
	case "enable":
		return p.manager.Enable(ctx, slug), true
	case "disable":
		return p.manager.Disable(ctx, slug), true
	case "reload":
		return p.manager.Reload(ctx, slug), true
	case "uninstall":
		return p.manager.Uninstall(ctx, slug), true
	}
	return nil, false
}

// handlePluginUpload accepts a multipart upload of a .ocpkg file from
// the admin UI, validates it, writes it into the plugins directory under
// the manifest's name (not the uploaded filename), and forces a scan
// so the new entry appears in the next list. Caps the request body at
// the manager's MaxUploadBytes to keep a hostile upload from filling
// memory before the validation gate runs.
func (p *Host) handlePluginUpload(w http.ResponseWriter, r *http.Request) {
	// MaxBytesReader returns an error as soon as the limit is exceeded,
	// so we never read more than this into memory.
	r.Body = http.MaxBytesReader(w, r.Body, plugins.MaxUploadBytes)
	// The MaxBytesReader above caps the multipart parse; gosec G120's
	// pattern match doesn't see that wrapping, so it's a false positive.
	if err := r.ParseMultipartForm(plugins.MaxUploadBytes); err != nil { //nolint:gosec // G120: bound via MaxBytesReader above
		writeJSONResponse(w, http.StatusRequestEntityTooLarge, map[string]string{jsonErrorKey: err.Error()})
		return
	}
	file, header, err := r.FormFile("plugin")
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{jsonErrorKey: "expected a file under form field 'plugin': " + err.Error()})
		return
	}
	defer file.Close()
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".ocpkg") {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{jsonErrorKey: "filename must end in .ocpkg"})
		return
	}
	packageBytes, err := io.ReadAll(file)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{jsonErrorKey: "read upload: " + err.Error()})
		return
	}
	entry, err := p.manager.Install(r.Context(), packageBytes)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{jsonErrorKey: err.Error()})
		return
	}
	log.Infof("plugin %q [%s] v%s installed by admin", entry.DisplayName, entry.Slug, entry.Version)
	writeJSONResponse(w, http.StatusOK, entry)
}

// configEnabledStore persists the enabled-plugin set and per-plugin
// admin-approved permission snapshots in Owncast's config datastore
// instead of a .enabled.json file.
type configEnabledStore struct {
	datastore *datastore.Datastore
}

func (s *configEnabledStore) Load() (plugins.StoreData, error) {
	names, err := s.datastore.GetStringSlice(pluginsEnabledConfigKey)
	if err != nil {
		// Unset on a fresh install: start with no plugins enabled.
		return plugins.StoreData{}, nil
	}
	out := plugins.StoreData{Enabled: names}
	raw, err := s.datastore.GetString(pluginsApprovedPermsConfigKey)
	if err == nil && raw != "" {
		var approvals map[string][]string
		if err := json.Unmarshal([]byte(raw), &approvals); err == nil {
			out.ApprovedPermissions = approvals
		}
	}
	return out, nil
}

func (s *configEnabledStore) Save(d plugins.StoreData) error {
	if err := s.datastore.SetStringSlice(pluginsEnabledConfigKey, d.Enabled); err != nil {
		return err
	}
	if len(d.ApprovedPermissions) == 0 {
		return s.datastore.SetString(pluginsApprovedPermsConfigKey, "")
	}
	encoded, err := json.Marshal(d.ApprovedPermissions)
	if err != nil {
		return fmt.Errorf("encode approved permissions: %w", err)
	}
	return s.datastore.SetString(pluginsApprovedPermsConfigKey, string(encoded))
}

// wirePluginHostEnv connects each HostEnv host-function pointer to the
// corresponding Owncast service call. Closures read services lazily so they
// observe current config/state on every call.
func wirePluginHostEnv(env *plugins.HostEnv, deps Deps) {
	chatbots := newPluginChatbotProvisioner(deps.UserRepository, deps.Datastore)
	wireChatSendHostFns(env, deps, chatbots)
	wireChatReadHostFns(env, deps)
	wireChatModerationHostFns(env, deps)
	wireServerReadHostFns(env, deps)
	wireVideoConfigHostFns(env, deps)
	wireUserHostFns(env, deps)
	wireNotificationHostFns(env, deps)
	wireRequestHostFns(env, deps)
	wireFilesystemHostFns(env)
}

// pluginDataRootDirName is the directory under config.DataDirectory that
// holds each plugin's private, sandboxed filesystem (storage.fs). It is
// deliberately separate from the plugins install/scan directory
// (config.DataDirectory/plugins) so a plugin's writable data can never be
// mistaken for, or collide with, an installed package or its assets.
const pluginDataRootDirName = "plugin-data"

// maxPluginFileBytes caps a single storage.fs write. It bounds how much a
// misbehaving plugin can write in one call; the host's overall disk use is
// still the admin's responsibility.
const maxPluginFileBytes = 50 << 20 // 50 MiB

// resolvePluginSandboxPath maps a plugin-supplied relative path to an
// absolute path inside that plugin's sandbox (root/<pluginName>), and
// guarantees the result cannot escape it. rel is treated as rooted before
// cleaning, so "../", absolute paths, and other traversal tricks all
// collapse back inside the sandbox; a defensive prefix check rejects
// anything that still lands outside. pluginName is the plugin slug, which
// the manifest layer has already constrained to [a-z][a-z0-9-]*, so it
// cannot itself contain separators or "..".
func resolvePluginSandboxPath(root, pluginName, rel string) (string, error) {
	sandbox, err := filepath.Abs(filepath.Join(root, pluginName))
	if err != nil {
		return "", err
	}
	// Rooting rel at "/" before Clean neutralizes leading "../" segments:
	// filepath.Clean("/"+"../../etc") == "/etc", which then joins back under
	// the sandbox rather than above it.
	full := filepath.Join(sandbox, filepath.Clean("/"+rel))
	if full != sandbox && !strings.HasPrefix(full, sandbox+string(os.PathSeparator)) {
		return "", fmt.Errorf("path %q escapes the plugin sandbox", rel)
	}
	return full, nil
}

// wireFilesystemHostFns implements the storage.fs host functions against a
// per-plugin sandbox directory under config.DataDirectory/plugin-data.
func wireFilesystemHostFns(env *plugins.HostEnv) {
	wireFilesystemHostFnsWithRoot(env, filepath.Join(config.DataDirectory, pluginDataRootDirName))
}

// wireFilesystemHostFnsWithRoot is wireFilesystemHostFns with the sandbox
// parent directory injected, so tests can point the storage.fs functions at
// a temp directory instead of the real data directory.
func wireFilesystemHostFnsWithRoot(env *plugins.HostEnv, root string) {
	env.FSRead = func(pluginName, path string) ([]byte, error) {
		full, err := resolvePluginSandboxPath(root, pluginName, path)
		if err != nil {
			return nil, err
		}
		// gosec G304: full is confined to the plugin's sandbox by
		// resolvePluginSandboxPath, which rejects any path that escapes it.
		return os.ReadFile(full) //nolint:gosec // G304: path sandboxed above
	}

	env.FSWrite = func(pluginName, path string, data []byte) error {
		if len(data) > maxPluginFileBytes {
			return fmt.Errorf("file is %d bytes; the limit is %d", len(data), maxPluginFileBytes)
		}
		full, err := resolvePluginSandboxPath(root, pluginName, path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(full), 0o700); err != nil {
			return err
		}
		// gosec G304: full is confined to the plugin's sandbox by
		// resolvePluginSandboxPath, which rejects any path that escapes it.
		return os.WriteFile(full, data, 0o600) //nolint:gosec // G304: path sandboxed above
	}

	env.FSList = func(pluginName, dir string) ([]string, error) {
		full, err := resolvePluginSandboxPath(root, pluginName, dir)
		if err != nil {
			return nil, err
		}
		entries, err := os.ReadDir(full)
		if err != nil {
			if os.IsNotExist(err) {
				return []string{}, nil
			}
			return nil, err
		}
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		return names, nil
	}

	env.FSDelete = func(pluginName, path string) error {
		full, err := resolvePluginSandboxPath(root, pluginName, path)
		if err != nil {
			return err
		}
		// os.Remove deletes a file or an empty directory only; it won't
		// recursively wipe a populated subtree, which keeps a single
		// delete call from being more destructive than the plugin asked.
		return os.Remove(full)
	}

	env.FSExists = func(pluginName, path string) (bool, error) {
		full, err := resolvePluginSandboxPath(root, pluginName, path)
		if err != nil {
			return false, err
		}
		if _, err := os.Stat(full); err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
}

// wireVideoConfigHostFns wires the settable video/transcoding configuration:
// videoconfig.read (owncast.videoConfig.read) and videoconfig.write
// (owncast.videoConfig.write). The plugin-facing VideoConfig/StreamVariant
// shapes are a curated, stable wire contract — the mapping to/from Owncast's
// internal models.StreamOutputVariant lives here, at the boundary.
func wireVideoConfigHostFns(env *plugins.HostEnv, deps Deps) {
	cfg := deps.ConfigRepository

	env.VideoConfig = func() plugins.VideoConfig {
		variants := cfg.GetStreamOutputVariants()
		out := plugins.VideoConfig{
			LatencyLevel: cfg.GetStreamLatencyLevel().Level,
			Codec:        cfg.GetVideoCodec(),
			Variants:     make([]plugins.StreamVariant, 0, len(variants)),
		}
		for _, v := range variants {
			out.Variants = append(out.Variants, plugins.StreamVariant{
				Width:         v.ScaledWidth,
				Height:        v.ScaledHeight,
				Framerate:     v.Framerate,
				VideoBitrate:  v.VideoBitrate,
				AudioBitrate:  v.AudioBitrate,
				IsPassthrough: v.IsVideoPassthrough,
			})
		}
		return out
	}

	env.WriteVideoConfig = func(pluginName string, u plugins.VideoConfigUpdate) error {
		if u.LatencyLevel != nil {
			if err := cfg.SetStreamLatencyLevel(float64(*u.LatencyLevel)); err != nil {
				return err
			}
		}
		if u.Codec != nil {
			if err := cfg.SetVideoCodec(*u.Codec); err != nil {
				return err
			}
		}
		if u.Variants != nil {
			variants := make([]models.StreamOutputVariant, 0, len(u.Variants))
			for _, v := range u.Variants {
				variants = append(variants, models.StreamOutputVariant{
					ScaledWidth:        v.Width,
					ScaledHeight:       v.Height,
					Framerate:          v.Framerate,
					VideoBitrate:       v.VideoBitrate,
					AudioBitrate:       v.AudioBitrate,
					IsVideoPassthrough: v.IsPassthrough,
				})
			}
			if err := cfg.SetStreamOutputVariants(variants); err != nil {
				return err
			}
		}
		// Mirror the admin video-config handlers: persist only. The change
		// takes effect on the next stream start (the admin UI already shows
		// this), so we deliberately do NOT restart a live transcoder.
		log.Infof("plugin %q changed video config via videoconfig.write; will take effect on next stream start", pluginName)
		return nil
	}
}

func wireChatSendHostFns(env *plugins.HostEnv, deps Deps, chatbots *pluginChatbotProvisioner) {
	chatSvc := deps.Chat

	env.OnChat = func(req plugins.ChatSendRequest) {
		switch req.Kind {
		case plugins.ChatSendAction:
			if err := chatSvc.SendSystemAction(req.Text, false); err != nil {
				log.Errorln("plugin", req.PluginSlug, "chat action:", err)
			}
		case plugins.ChatSendSystem:
			if err := chatSvc.SendSystemMessage(req.Text, false); err != nil {
				log.Errorln("plugin", req.PluginSlug, "chat system message:", err)
			}
		default: // ChatSendBot — post under the plugin's own chatbot identity.
			// Use slug for the bot's persistent identity (cache + datastore
			// key), and BotDisplayName for the human-readable handle viewers
			// see in chat. BotDisplayName falls back to the plugin's display
			// name in the host fn closure, so it's always non-empty here.
			chatbot, err := chatbots.chatbotUser(req.PluginSlug, req.BotDisplayName)
			if err != nil {
				log.Errorln("plugin", req.PluginSlug, "resolve chatbot user:", err)
				return
			}
			if err := chatSvc.SendMessageAsBot(chatbot, req.Text); err != nil {
				log.Errorln("plugin", req.PluginSlug, "chat send:", err)
			}
		}
	}

	env.SendChatTo = func(pluginSlug string, clientID uint64, text string) {
		chatSvc.SendSystemMessageToClient(uint(clientID), text)
	}
}

func wireChatReadHostFns(env *plugins.HostEnv, deps Deps) {
	chatSvc := deps.Chat

	env.ChatHistory = func(limit int) []plugins.HostChatMessage {
		history := deps.ChatMessageRepository.GetChatHistory()
		out := make([]plugins.HostChatMessage, 0, len(history))
		for _, item := range history {
			msg, ok := item.(events.UserMessageEvent)
			if !ok {
				continue
			}
			hm := plugins.HostChatMessage{
				ID:        msg.ID,
				Body:      msg.Body,
				Timestamp: msg.Timestamp.UTC().Format(time.RFC3339Nano),
			}
			if msg.User != nil {
				hm.User = msg.User.DisplayName
			}
			out = append(out, hm)
		}
		if limit > 0 && len(out) > limit {
			out = out[len(out)-limit:]
		}
		return out
	}

	env.ChatClients = func() []plugins.HostChatClient {
		clients := chatSvc.GetClients()
		out := make([]plugins.HostChatClient, 0, len(clients))
		for _, c := range clients {
			hc := plugins.HostChatClient{
				ID:           uint64(c.Id),
				ConnectedAt:  c.ConnectedAt.UTC().Format(time.RFC3339Nano),
				UserAgent:    c.UserAgent,
				IPAddress:    c.IPAddress,
				MessageCount: c.MessageCount,
			}
			if c.User != nil {
				hc.UserID = c.User.ID
				hc.DisplayName = c.User.DisplayName
			}
			out = append(out, hc)
		}
		return out
	}
}

func wireChatModerationHostFns(env *plugins.HostEnv, deps Deps) {
	chatSvc := deps.Chat

	env.DeleteMessage = func(pluginName, messageID string) {
		if err := chatSvc.SetMessagesVisibility([]string{messageID}, false); err != nil {
			log.Errorln("plugin", pluginName, "delete message:", err)
		}
	}

	env.KickClient = func(pluginName string, clientID uint64) {
		if c, ok := chatSvc.FindClientByID(uint(clientID)); ok {
			chatSvc.DisconnectClients([]*chat.Client{c})
		}
	}
}

func wireServerReadHostFns(env *plugins.HostEnv, deps Deps) {
	cfg := deps.ConfigRepository
	streamSvc := deps.Stream

	env.StreamCurrent = func() plugins.StreamInfo {
		status := streamSvc.GetStatus()
		info := plugins.StreamInfo{
			Online:       status.Online,
			Title:        cfg.GetStreamTitle(),
			Summary:      cfg.GetServerSummary(),
			Viewers:      status.ViewerCount,
			LatencyLevel: cfg.GetStreamLatencyLevel().Level,
		}
		if status.LastConnectTime != nil && status.LastConnectTime.Valid {
			info.StartedAt = status.LastConnectTime.Time.UTC().Format(time.RFC3339Nano)
		}
		return info
	}

	env.ServerInfo = func() plugins.ServerInfo {
		return plugins.ServerInfo{
			Name:           cfg.GetServerName(),
			URL:            cfg.GetServerURL(),
			Summary:        cfg.GetServerSummary(),
			WelcomeMessage: cfg.GetServerWelcomeMessage(),
			Version:        config.VersionNumber,
		}
	}

	env.Socials = func() []plugins.SocialHandle {
		handles := cfg.GetSocialHandles()
		out := make([]plugins.SocialHandle, 0, len(handles))
		for _, h := range handles {
			out = append(out, plugins.SocialHandle{Platform: h.Platform, URL: h.URL, Icon: h.Icon})
		}
		return out
	}

	env.Federation = func() plugins.FederationInfo {
		return plugins.FederationInfo{
			Enabled:   cfg.GetFederationEnabled(),
			Username:  cfg.GetFederationUsername(),
			IsPrivate: cfg.GetFederationIsPrivate(),
		}
	}

	env.Tags = func() []string {
		return cfg.GetServerMetadataTags()
	}

	// Broadcaster is read-only telemetry about the inbound feed (nil between
	// streams), distinct from the settable video config below.
	env.Broadcaster = func() plugins.StreamBroadcaster {
		b := streamSvc.GetBroadcaster()
		if b == nil {
			return plugins.StreamBroadcaster{}
		}
		d := b.StreamDetails
		codecs := make([]string, 0, 2)
		if d.VideoCodec != "" {
			codecs = append(codecs, d.VideoCodec)
		}
		if d.AudioCodec != "" {
			codecs = append(codecs, d.AudioCodec)
		}
		out := plugins.StreamBroadcaster{
			RemoteAddr: b.RemoteAddr,
			Codecs:     codecs,
			Framerate:  int(d.VideoFramerate),
		}
		if d.Width > 0 || d.Height > 0 {
			out.Resolution = fmt.Sprintf("%dx%d", d.Width, d.Height)
		}
		if d.VideoBitrate > 0 {
			out.Bitrates = []int{d.VideoBitrate}
		}
		return out
	}
}

func wireUserHostFns(env *plugins.HostEnv, deps Deps) {
	users := deps.UserRepository
	chatSvc := deps.Chat

	env.Users = func() []plugins.HostUser {
		all := users.GetUsers()
		out := make([]plugins.HostUser, 0, len(all))
		for _, u := range all {
			out = append(out, toHostUser(u))
		}
		return out
	}

	env.UserGet = func(id string) (plugins.HostUser, bool) {
		u := users.GetUserByID(id)
		if u == nil {
			return plugins.HostUser{}, false
		}
		return toHostUser(u), true
	}

	env.SetUserEnabled = func(pluginName, userID string, enabled bool, reason string) {
		if err := users.SetEnabled(userID, enabled); err != nil {
			log.Errorln("plugin", pluginName, "set user enabled:", err)
			return
		}
		if !enabled {
			if clients, err := chatSvc.GetClientsForUser(userID); err == nil {
				chatSvc.DisconnectClients(clients)
			}
		}
	}

	env.BanIP = func(pluginName, ip string) {
		if err := deps.AuthRepository.BanIPAddress(ip, "banned by plugin "+pluginName); err != nil {
			log.Errorln("plugin", pluginName, "ban ip:", err)
		}
	}
}

func wireNotificationHostFns(env *plugins.HostEnv, deps Deps) {
	cfg := deps.ConfigRepository

	env.SendDiscord = func(pluginName, text string) {
		dc := cfg.GetDiscordConfig()
		if !dc.Enabled || dc.Webhook == "" {
			return
		}
		notifier, err := discord.New(cfg.GetServerName(), cfg.GetServerURL()+"/logo", dc.Webhook)
		if err != nil {
			log.Errorln("plugin", pluginName, "discord:", err)
			return
		}
		if err := notifier.Send(text); err != nil {
			log.Errorln("plugin", pluginName, "discord send:", err)
		}
	}

	env.SendBrowserPush = func(pluginName string, p plugins.BrowserPushPayload) {
		publicKey, err := cfg.GetBrowserPushPublicKey()
		if err != nil {
			return
		}
		privateKey, err := cfg.GetBrowserPushPrivateKey()
		if err != nil {
			return
		}
		notifier, err := browser.New(deps.Datastore, publicKey, privateKey)
		if err != nil {
			log.Errorln("plugin", pluginName, "browser push:", err)
			return
		}
		destinations, err := deps.NotificationsRepository.GetNotificationDestinationsForChannel(notificationsrepository.BrowserPushNotification)
		if err != nil {
			return
		}
		for _, destination := range destinations {
			if _, err := notifier.Send(destination, p.Title, p.Body); err != nil {
				log.Debugln("plugin", pluginName, "browser push send:", err)
			}
		}
	}

	env.SendFediverse = func(pluginName string, p plugins.FediversePayload) {
		var image *string
		if p.Image != "" {
			image = &p.Image
		}
		if err := deps.Chat.SendFediverseAction(p.Type, cfg.GetFederationUsername(), image, p.Body, p.Link); err != nil {
			log.Errorln("plugin", pluginName, "fediverse action:", err)
		}
	}

	env.PostFediverse = func(pluginName, text string) (string, error) {
		// Owncast publishes the note but does not return its URL.
		if err := deps.Activitypub.SendPublicFederatedMessage(text); err != nil {
			return "", err
		}
		return "", nil
	}
}

func wireRequestHostFns(env *plugins.HostEnv, deps Deps) {
	cfg := deps.ConfigRepository
	users := deps.UserRepository

	env.UploadStorage = func(pluginName, name string, data []byte) (string, error) {
		return uploadPluginAsset(cfg, pluginName, name, data)
	}

	env.IsAuthenticated = deps.IsAdminRequest

	env.GetRequestUser = func(r *http.Request) *plugins.HostUser {
		token := r.URL.Query().Get("accessToken")
		if token == "" {
			return nil
		}
		u := users.GetUserByToken(token)
		if u == nil {
			return nil
		}
		hu := toHostUser(u)
		return &hu
	}
}

// toHostUser maps an Owncast user model onto the plugin-facing HostUser.
func toHostUser(u *models.User) plugins.HostUser {
	hu := plugins.HostUser{
		ID:              u.ID,
		DisplayName:     u.DisplayName,
		PreviousNames:   u.PreviousNames,
		Scopes:          u.Scopes,
		IsBot:           u.IsBot,
		IsAuthenticated: u.Authenticated,
		CreatedAt:       u.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
	if u.DisabledAt != nil {
		hu.DisabledAt = u.DisabledAt.UTC().Format(time.RFC3339Nano)
	}
	return hu
}

// uploadPluginAsset writes a plugin upload under the public files directory
// and returns the URL it is served at. Names are flattened to their base to
// prevent path traversal. S3-backed storage is a follow-up.
func uploadPluginAsset(cfg configrepository.ConfigRepository, pluginName, name string, data []byte) (string, error) {
	safeName := filepath.Base(filepath.Clean("/" + name))
	if safeName == "." || safeName == "/" || safeName == "" {
		return "", fmt.Errorf("invalid upload name %q", name)
	}
	dir := filepath.Join(config.PublicFilesPath, "plugins", pluginName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, safeName), data, 0o600); err != nil {
		return "", err
	}
	url := strings.TrimSuffix(cfg.GetServerURL(), "/") + "/public/plugins/" + pluginName + "/" + safeName
	return url, nil
}
