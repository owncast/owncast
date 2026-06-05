package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// Server is an http.Handler that serves /plugins/<name>/* — first by
// looking for a static asset on the plugin's PublicFS, then falling through
// to the plugin's on_http_request wasm export. Plugins without the
// http.serve permission produce 404 regardless.
//
// Asset storage is abstracted via io/fs so loose-files plugins (os.DirFS),
// packaged plugins (.ocpkg via zip.Reader), and any future layout share the
// same serving code.
//
// Mount this on the parent router at /plugins/, e.g.
//
//	mux.Handle("/plugins/", server)
type Server struct {
	// snapshot returns the currently-loaded plugins. Called per-request so
	// admin enable/disable takes effect immediately.
	snapshot func() []*Loaded
	// IsAuthenticated reports whether an incoming HTTP request carries
	// authenticated-admin credentials. Used to populate req.authenticated
	// for the plugin's dynamic handler (so a plugin page can render
	// differently for authenticated viewers without owning auth itself).
	// Admin-path gating is done via RequireAdmin, not by checking this and
	// writing a 401 here — that's the host's responsibility.
	// nil = always false (no auth available).
	IsAuthenticated func(*http.Request) bool
	// RequireAdmin wraps a handler in the host's admin Basic Auth
	// middleware. The plugin static server calls it when a request lands
	// on a manifest-declared admin path so the 401 (realm, CORS, log line)
	// comes from one place in the codebase instead of being duplicated
	// here. nil = admin paths cannot be served (404).
	RequireAdmin func(http.HandlerFunc) http.HandlerFunc
	// GetRequestUser returns the user identity attached to the request
	// (when the request came with a user-token, not admin auth). nil →
	// req.user is always omitted from the envelope.
	GetRequestUser func(*http.Request) *HostUser
	// SSE is the host-owned hub that backs the reserved
	// /plugins/<name>/_sse/<channel> endpoint. The host holds these
	// long-lived connections; plugins push to them via owncast.sse.send.
	// nil → the _sse endpoint returns 503.
	SSE *SSEHub
}

// SSEReservedPrefix is the request path (relative to the plugin namespace)
// the host reserves for Server-Sent-Events streams. A request to
// /plugins/<name>/_sse/<channel> connects to the plugin's <channel> stream;
// the segment after the prefix is the channel name (empty = default).
// Plugins cannot serve their own routes under this prefix.
const SSEReservedPrefix = "/_sse"

// SSEKeepAliveInterval is how often the host writes an SSE comment line to
// an idle stream so proxies don't close it for inactivity.
const SSEKeepAliveInterval = 15 * time.Second

// HTTP enforcement limits. Per-plugin, per-request.
const (
	MaxHTTPRequestBodyBytes  = 1 << 20  // 1 MB
	MaxHTTPResponseBodyBytes = 10 << 20 // 10 MB
)

// adminStyleHrefs are the host's stylesheet URLs auto-injected into HTML
// responses on a plugin's manifest-declared admin paths. It's the shared,
// generated plugin stylesheet (web/style-definitions builds it from the
// design tokens, so values stay in sync) that styles plain HTML elements
// (input, button, table, label, …) with the Owncast theme tokens — not the
// surrounding admin's AntD-specific stylesheets, which only apply to .ant-*
// class selectors the plugin pages don't use. The viewer-tab iframes inject
// the same sheet, so admin and tab plugin content share one baseline.
//
// The file is served by the Owncast HTTP server (from web/public/styles/
// in dev, from the bundled assets in prod), so the same origin-relative URL
// works in both modes.
var adminStyleHrefs = []string{
	"/styles/plugin.css",
}

// adminStyleSnippet is the bytes injected into HTML responses on admin
// paths. Built once at startup so the per-request cost is just a single
// substring search + concatenation.
var adminStyleSnippet = buildAdminStyleSnippet(adminStyleHrefs)

// adminStyleMarker is a sentinel attribute on the injected <link> tags
// so the injector is idempotent: if a plugin's HTML already includes our
// snippet (e.g. baked in by their build), we don't add it again.
const adminStyleMarker = `data-owncast-admin-style="1"`

func buildAdminStyleSnippet(hrefs []string) []byte {
	var b bytes.Buffer
	for _, h := range hrefs {
		fmt.Fprintf(&b, `<link rel="stylesheet" href=%q %s>`, h, adminStyleMarker)
	}
	return b.Bytes()
}

// injectAdminStyles returns html with the admin-stylesheet snippet inserted
// before </head>, or — when no </head> is present — prepended. Idempotent:
// returns the input unchanged when the snippet's marker is already in the
// document. text/html responses on a manifest-declared admin path go
// through this; nothing else.
func injectAdminStyles(html []byte) []byte {
	if bytes.Contains(html, []byte(adminStyleMarker)) {
		return html
	}
	if i := bytes.Index(bytes.ToLower(html), []byte("</head>")); i >= 0 {
		out := make([]byte, 0, len(html)+len(adminStyleSnippet))
		out = append(out, html[:i]...)
		out = append(out, adminStyleSnippet...)
		out = append(out, html[i:]...)
		return out
	}
	// No <head> — prepend so styles still apply. Browsers tolerate a
	// <link> at the top of a body-only document.
	out := make([]byte, 0, len(html)+len(adminStyleSnippet))
	out = append(out, adminStyleSnippet...)
	out = append(out, html...)
	return out
}

// allowedResponseHeaders is the set of headers a plugin response is allowed
// to set. We block headers that would let a plugin override Owncast's own
// transport-security, CSP, or server-identification headers. Set-Cookie
// is allowed: a plugin's response defaults to a Path scoped to its own
// /plugins/<name>/ namespace, so cookies don't leak into the rest of
// Owncast unless the plugin explicitly broadens them. CORS headers
// (Access-Control-*) are matched via prefix below.
var allowedResponseHeaders = map[string]bool{
	"content-type":     true,
	"content-encoding": true,
	"content-language": true,
	"cache-control":    true,
	"last-modified":    true,
	"etag":             true,
	"location":         true,
	"vary":             true,
	"link":             true,
	"set-cookie":       true,
}

func isAllowedResponseHeader(name string) bool {
	lower := strings.ToLower(name)
	if allowedResponseHeaders[lower] {
		return true
	}
	if strings.HasPrefix(lower, "access-control-") {
		return true
	}
	return false
}

// NewServer constructs an HTTP handler over a fixed plugin set. Used in
// tests and any context where the plugin set doesn't change after
// construction. Each plugin's PublicFS is used for static asset serving;
// plugins with nil PublicFS just don't serve static files.
func NewServer(loaded []*Loaded) *Server {
	snap := loaded
	return &Server{snapshot: func() []*Loaded { return snap }}
}

// NewLiveServer constructs an HTTP handler backed by a snapshot function —
// the Manager passes its Snapshot method here so admin enable/disable takes
// effect on subsequent requests without restarting the host.
func NewLiveServer(snapshot func() []*Loaded) *Server {
	return &Server{snapshot: snapshot}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Path: /plugins/<name>/<rest>. We're mounted at /plugins/, so
	// r.URL.Path starts with /plugins/. Strip and split.
	rel := strings.TrimPrefix(r.URL.Path, "/plugins/")
	parts := strings.SplitN(rel, "/", 2)
	if parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	name := parts[0]
	p := s.lookup(name)
	if p == nil {
		http.NotFound(w, r)
		return
	}

	rest := "/"
	if len(parts) > 1 {
		rest = "/" + parts[1]
	}

	// Admin-only routes (declared in the plugin's manifest) are gated by
	// the host's admin auth middleware before the plugin sees the request.
	// We invoke RequireAdmin directly so the 401 — realm, CORS, log line —
	// comes from middleware/auth.go rather than being duplicated here.
	if p.IsAdminPath(rest) {
		if s.RequireAdmin == nil {
			http.NotFound(w, r)
			return
		}
		s.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
			s.serveAuthorized(w, r, p, rest, true)
		})(w, r)
		return
	}

	authenticated := false
	if s.IsAuthenticated != nil {
		authenticated = s.IsAuthenticated(r)
	}
	s.serveAuthorized(w, r, p, rest, authenticated)
}

// serveAuthorized runs the rest of the plugin request pipeline (SSE,
// permission gate, static asset, dynamic handler) after the admin-path
// check has been resolved. authenticated is the boolean handed to the
// plugin's dynamic handler — true on admin paths (callers reach this via
// RequireAdmin's success continuation), the IsAuthenticated predicate
// otherwise.
//
// An HTML response served on a manifest-declared admin path also gets the
// host's admin stylesheet injected so the plugin iframe matches the
// surrounding admin visually with no plugin-author opt-in.
func (s *Server) serveAuthorized(w http.ResponseWriter, r *http.Request, p *Loaded, rest string, authenticated bool) {
	injectStyles := p.IsAdminPath(rest)

	// Host-reserved Server-Sent-Events endpoint. Gated on http.sse
	// (independent of http.serve) and handled entirely by the host — the
	// plugin never sees these requests.
	if rest == SSEReservedPrefix || strings.HasPrefix(rest, SSEReservedPrefix+"/") {
		if !pluginHasPermission(p.Manifest, PermHttpSSE) {
			http.NotFound(w, r)
			return
		}
		channel := strings.TrimPrefix(strings.TrimPrefix(rest, SSEReservedPrefix), "/")
		s.serveSSE(w, r, p, channel)
		return
	}

	if !pluginHasPermission(p.Manifest, PermHttpServe) {
		http.NotFound(w, r)
		return
	}

	if s.tryStatic(w, r, p, rest, injectStyles) {
		return
	}
	// Strip the trailing slash before handing the path to the plugin's
	// dynamic handler so plugin code can match canonical no-slash form
	// regardless of how the URL was canonicalized upstream. Next.js dev's
	// trailingSlash:true rewrites every URL to a slash-terminated form,
	// which would otherwise force every plugin to accept both /foo and
	// /foo/ in its on_http_request switch. Static asset serving keeps the
	// original path so directory-vs-file resolution is unaffected.
	dynamicPath := rest
	if len(dynamicPath) > 1 && strings.HasSuffix(dynamicPath, "/") {
		dynamicPath = strings.TrimRight(dynamicPath, "/")
	}
	s.serveDynamic(w, r, p, dynamicPath, authenticated, injectStyles)
}

// serveSSE holds a long-lived Server-Sent-Events connection open and streams
// frames the plugin publishes to (plugin, channel). The connection is owned
// entirely by the host: the per-plugin mutex is never taken during streaming,
// so an idle stream costs only a goroutine. The loop exits when the client
// disconnects (request context cancelled).
//
// The plugin is notified of the connection opening and closing via the
// sse.connect / sse.disconnect events (a single short wasm call each, around
// the stream, not during it) so it can track who is connected.
func (s *Server) serveSSE(w http.ResponseWriter, r *http.Request, p *Loaded, channel string) {
	if s.SSE == nil {
		http.Error(w, "server-sent events not available", http.StatusServiceUnavailable)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	stream, unsubscribe, connID, ok := s.SSE.Subscribe(p.Manifest.Slug, channel)
	if !ok {
		http.Error(w, "too many event-stream connections for this plugin", http.StatusServiceUnavailable)
		return
	}
	defer unsubscribe()

	// Resolve the connecting chat user (if any) once, and tell the plugin who
	// connected. The matching disconnect fires when this handler returns.
	var user *HostUser
	if s.GetRequestUser != nil {
		user = s.GetRequestUser(r)
	}
	s.notifySSEConnection(p, EventSSEConnect, channel, connID, user)
	defer s.notifySSEConnection(p, EventSSEDisconnect, channel, connID, user)

	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	keepAlive := time.NewTicker(SSEKeepAliveInterval)
	defer keepAlive.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case frame := <-stream:
			if _, err := w.Write(frame); err != nil {
				return
			}
			flusher.Flush()
		case <-keepAlive.C:
			if _, err := io.WriteString(w, ": keep-alive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// sseConnectionEnvelope marshals the on_event payload delivered to a plugin
// for an SSE connect/disconnect.
func sseConnectionEnvelope(eventType, channel string, connID uint64, user *HostUser) ([]byte, error) {
	return json.Marshal(Envelope{
		EventType: eventType,
		Payload:   SSEConnectionEvent{Channel: channel, ConnectionID: connID, User: user},
	})
}

// notifySSEConnection fires an sse.connect / sse.disconnect event to the
// plugin that owns the channel, carrying the channel, connection id, and the
// resolved chat user (if any). It is best-effort: a plugin that doesn't export
// on_event or doesn't handle the event is a silent no-op, and any error is
// logged but never affects the stream. A fresh context is used because the
// disconnect path fires after the request context has already been cancelled.
func (s *Server) notifySSEConnection(p *Loaded, eventType, channel string, connID uint64, user *HostUser) {
	if p.plugin == nil || !p.plugin.FunctionExists("on_event") {
		return
	}
	envelope, err := sseConnectionEnvelope(eventType, channel, connID, user)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), NotifyTimeout)
	defer cancel()
	if err := callOnEvent(ctx, p, envelope); err != nil {
		fmt.Fprintf(os.Stderr, "plugin %s: %s notify failed: %v\n", p.Manifest.Slug, eventType, err)
	}
}

// lookup returns the currently-loaded plugin with the given slug, or nil
// if there isn't one. Called per-request against the live snapshot.
// The lookup key is the manifest slug (also the URL segment under
// /plugins/<slug>/), not the human-readable display name.
func (s *Server) lookup(slug string) *Loaded {
	for _, p := range s.snapshot() {
		if p.Manifest.Slug == slug {
			return p
		}
	}
	return nil
}

func pluginHasPermission(m *Manifest, perm string) bool {
	for _, p := range m.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

func (s *Server) tryStatic(w http.ResponseWriter, r *http.Request, loaded *Loaded, requestPath string, injectStyles bool) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	if loaded.PublicFS == nil {
		return false
	}

	// fs.FS paths must be slash-separated, without a leading slash, and
	// can't contain ".." segments — path.Clean handles the first two; the
	// fs.ValidPath check rejects traversal.
	cleaned := strings.TrimPrefix(path.Clean("/"+requestPath), "/")
	if cleaned == "" {
		cleaned = "."
	}
	if !fs.ValidPath(cleaned) {
		return false
	}

	info, err := fs.Stat(loaded.PublicFS, cleaned)
	if err != nil {
		return false
	}
	if info.IsDir() {
		indexPath := path.Join(cleaned, "index.html")
		idx, err := fs.Stat(loaded.PublicFS, indexPath)
		if err != nil || idx.IsDir() {
			return false
		}
		// Standard directory-index behavior: if the request URL doesn't end
		// with a slash, redirect to the canonical slash form before serving
		// the index.html. Without this the browser's base URL stays at the
		// parent path and any relative links in the page (`./api/settings`,
		// `<img src="logo.png">`, etc.) resolve to the wrong place.
		if !strings.HasSuffix(r.URL.Path, "/") {
			// Not an open redirect — destination is the request's own path
			// with a trailing slash, not an attacker-controlled URL.
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently) //nolint:gosec // G710
			return true
		}
		serveAssetFile(w, r, loaded.PublicFS, indexPath, idx, injectStyles)
		return true
	}
	serveAssetFile(w, r, loaded.PublicFS, cleaned, info, injectStyles)
	return true
}

// serveAssetFile reads a file from the plugin's AssetsFS into memory and
// hands it to http.ServeContent. Reading into memory avoids the seekability
// problems with non-file-backed fs.FS implementations (zip entries aren't
// seekable as ReadClosers); plugin assets are small enough that this is
// fine in practice. http.ServeContent gives us correct content-type
// sniffing, range support, ETag/conditional-GET handling — without
// net/http.ServeFile's path-canonicalization redirects.
//
// injectStyles=true rewrites HTML responses to include the host's admin
// stylesheet links before serving — only set by serveAuthorized when the
// path is on a manifest-declared admin route.
func serveAssetFile(w http.ResponseWriter, r *http.Request, root fs.FS, name string, info fs.FileInfo, injectStyles bool) {
	data, err := fs.ReadFile(root, name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if injectStyles && isHTMLName(name) {
		data = injectAdminStyles(data)
	}
	modtime := info.ModTime()
	if modtime.IsZero() {
		modtime = time.Time{} // ServeContent skips Last-Modified if zero
	}
	http.ServeContent(w, r, path.Base(name), modtime, bytes.NewReader(data))
}

func isHTMLName(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".html") || strings.HasSuffix(lower, ".htm")
}

func (s *Server) serveDynamic(w http.ResponseWriter, r *http.Request, p *Loaded, requestPath string, authenticated bool, injectStyles bool) {
	// p.plugin can be nil during shutdown (Loaded.Close clears it) or in
	// tests that only exercise the static path. Either way, no plugin
	// instance means no dynamic handler.
	if p.plugin == nil || !p.plugin.FunctionExists("on_http_request") {
		http.NotFound(w, r)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, MaxHTTPRequestBodyBytes+1))
	if err != nil {
		http.Error(w, "request body read error", http.StatusBadRequest)
		return
	}
	if int64(len(body)) > MaxHTTPRequestBodyBytes {
		http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	envelopeJSON, err := s.buildRequestEnvelope(r, requestPath, authenticated, body)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	callCtx, cancel := context.WithTimeout(r.Context(), HTTPHandlerTimeout)
	defer cancel()
	p.mu.Lock()
	_, out, err := p.plugin.CallWithContext(callCtx, "on_http_request", envelopeJSON)
	p.mu.Unlock()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || callCtx.Err() == context.DeadlineExceeded {
			http.Error(w, "plugin timed out", http.StatusGatewayTimeout)
			fmt.Fprintf(os.Stderr, "plugin %s: on_http_request timed out after %s\n", p.Manifest.Slug, HTTPHandlerTimeout)
			return
		}
		http.Error(w, "plugin error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "plugin %s: on_http_request failed: %v\n", p.Manifest.Slug, err)
		return
	}
	if len(out) > MaxHTTPHandlerOutputBytes {
		http.Error(w, "plugin response too large", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "plugin %s: on_http_request output too large: %d bytes (max %d)\n",
			p.Manifest.Slug, len(out), MaxHTTPHandlerOutputBytes)
		return
	}

	writePluginHTTPResponse(w, out, injectStyles)
}

// buildRequestEnvelope marshals the JSON envelope passed to a plugin's
// on_http_request export.
func (s *Server) buildRequestEnvelope(r *http.Request, requestPath string, authenticated bool, body []byte) ([]byte, error) {
	// Redact the credentials a request may carry so a sandboxed plugin never
	// sees the raw access token. Identity reaches the plugin only via the
	// host-resolved, trusted req.user below. The Cookie header carries the
	// chat identity cookie (and admin session); the accessToken query param
	// is the legacy token source.
	headers := flattenValues(r.Header)
	delete(headers, "Cookie")
	query := flattenValues(r.URL.Query())
	delete(query, "accessToken")

	envelope := map[string]any{
		"method":        r.Method,
		"path":          requestPath,
		"query":         query,
		"headers":       headers,
		"body":          string(body),
		"remoteAddr":    r.RemoteAddr,
		"authenticated": authenticated,
	}
	if s.GetRequestUser != nil {
		if user := s.GetRequestUser(r); user != nil {
			envelope["user"] = user
		}
	}
	return json.Marshal(envelope)
}

// writePluginHTTPResponse parses a plugin's on_http_request output envelope
// and writes it to the client, filtering disallowed headers.
//
// injectStyles=true rewrites an HTML response (Content-Type starting with
// text/html) to include the host's admin stylesheet links — same behavior
// as static HTML assets, so a plugin returning admin HTML from
// on_http_request gets the iframe theme automatically.
func writePluginHTTPResponse(w http.ResponseWriter, out []byte, injectStyles bool) {
	var resp struct {
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`
		Body    string            `json:"body"`
	}
	if err := json.Unmarshal(out, &resp); err != nil {
		http.Error(w, "plugin returned invalid response", http.StatusInternalServerError)
		return
	}
	if resp.Status == 0 {
		resp.Status = http.StatusOK
	}
	if len(resp.Body) > MaxHTTPResponseBodyBytes {
		http.Error(w, "plugin response too large", http.StatusInternalServerError)
		return
	}

	body := resp.Body
	if injectStyles && responseIsHTML(resp.Headers) {
		body = string(injectAdminStyles([]byte(body)))
	}

	for k, v := range resp.Headers {
		if !isAllowedResponseHeader(k) {
			continue
		}
		w.Header().Set(k, v)
	}
	w.WriteHeader(resp.Status)
	_, _ = io.WriteString(w, body)
}

// responseIsHTML reports whether the plugin-declared headers identify the
// response body as HTML. Case-insensitive on the header name; only the
// media type prefix is inspected so charset/boundary parameters don't
// matter.
func responseIsHTML(headers map[string]string) bool {
	for k, v := range headers {
		if strings.EqualFold(k, "content-type") {
			return strings.HasPrefix(strings.ToLower(strings.TrimSpace(v)), "text/html")
		}
	}
	return false
}

func flattenValues(v map[string][]string) map[string]string {
	out := make(map[string]string, len(v))
	for k, vs := range v {
		if len(vs) > 0 {
			out[k] = vs[0]
		}
	}
	return out
}
