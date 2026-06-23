package plugins

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gobwas/glob"
)

// Thin wrappers so the test code reads cleanly. These are only used by the
// path-traversal test, which needs a real on-disk root boundary.
func osMkdirAll(path string) error { return os.MkdirAll(path, 0o755) }
func osDirFS(path string) fs.FS    { return os.DirFS(path) }
func osWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// Build a Server with a single synthetic plugin that has no extism instance
// (we only exercise static-asset paths here). For tests that need the dynamic
// handler, use a real built wasm via the LoadPlugin / LoadPackage path.
func staticOnlyServer(t *testing.T, perms []string, assets fstest.MapFS) *Server {
	t.Helper()
	loaded := &Loaded{
		Manifest: &Manifest{
			API: "1", DisplayName: "demo", Slug: "demo", Version: "1.0.0", Permissions: perms,
		},
		PublicFS: assets,
	}
	return NewServer([]*Loaded{loaded})
}

func TestServer_NotFound_UnknownPlugin(t *testing.T) {
	s := staticOnlyServer(t, []string{"http.serve"}, fstest.MapFS{})
	req := httptest.NewRequest("GET", "/plugins/nonexistent/anything", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown plugin: status %d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestServer_NotFound_WithoutHttpServePermission(t *testing.T) {
	// Plugin exists but never declared http.serve — should appear as not-found.
	s := staticOnlyServer(t, []string{} /* no http.serve */, fstest.MapFS{
		"index.html": {Data: []byte("hi")},
	})
	req := httptest.NewRequest("GET", "/plugins/demo/index.html", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("no http.serve permission: status %d want %d", rec.Code, http.StatusNotFound)
	}
}

func TestServer_StaticFile(t *testing.T) {
	s := staticOnlyServer(t, []string{"http.serve"}, fstest.MapFS{
		"style.css": {Data: []byte("body { color: red; }")},
	})
	req := httptest.NewRequest("GET", "/plugins/demo/style.css", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status: got %d want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "color: red") {
		t.Errorf("body: got %q", rec.Body.String())
	}
}

func TestServer_DirectoryServesIndexHTML(t *testing.T) {
	s := staticOnlyServer(t, []string{"http.serve"}, fstest.MapFS{
		"index.html":     {Data: []byte("<h1>root</h1>")},
		"sub/index.html": {Data: []byte("<h1>sub</h1>")},
	})

	// Slash-terminated directory URLs serve the index.html directly. Without
	// a trailing slash, the server issues a 301 to the canonical slash form
	// so the browser's base URL is right for any relative links on the page.
	for _, tc := range []struct {
		path         string
		wantStatus   int
		wantBody     string
		wantLocation string
	}{
		{path: "/plugins/demo/", wantStatus: http.StatusOK, wantBody: "<h1>root</h1>"},
		{path: "/plugins/demo/sub/", wantStatus: http.StatusOK, wantBody: "<h1>sub</h1>"},
		{path: "/plugins/demo/sub", wantStatus: http.StatusMovedPermanently, wantLocation: "/plugins/demo/sub/"},
	} {
		req := httptest.NewRequest("GET", tc.path, nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
		if rec.Code != tc.wantStatus {
			t.Errorf("%s: status %d want %d", tc.path, rec.Code, tc.wantStatus)
			continue
		}
		if tc.wantBody != "" && rec.Body.String() != tc.wantBody {
			t.Errorf("%s: body %q want %q", tc.path, rec.Body.String(), tc.wantBody)
		}
		if tc.wantLocation != "" && rec.Header().Get("Location") != tc.wantLocation {
			t.Errorf("%s: Location %q want %q", tc.path, rec.Header().Get("Location"), tc.wantLocation)
		}
	}
}

func TestServer_InjectsAdminStylesOnAdminPath(t *testing.T) {
	// HTML served on a manifest-declared admin path gets the host's admin
	// stylesheet link injected before </head>, so the plugin iframe matches
	// the surrounding admin visually without plugin-author opt-in.
	loaded := &Loaded{
		Manifest: &Manifest{
			API: "1", DisplayName: "demo", Slug: "demo", Version: "1.0.0",
			Permissions: []string{"http.serve"},
		},
		PublicFS: fstest.MapFS{
			"admin/index.html": {Data: []byte("<html><head><title>x</title></head><body>hi</body></html>")},
			"index.html":       {Data: []byte("<html><head></head><body>public</body></html>")},
		},
		adminGlobs: []glob.Glob{glob.MustCompile("/admin/*")},
		adminPaths: []string{"/admin/*"},
	}
	s := NewServer([]*Loaded{loaded})
	// IsAuthenticated returns true so the admin-path gate lets the request through.
	s.IsAuthenticated = func(*http.Request) bool { return true }
	s.RequireAdmin = func(h http.HandlerFunc) http.HandlerFunc { return h }

	// Admin path: snippet present.
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest("GET", "/plugins/demo/admin/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("admin path status %d want 200; body=%q", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), adminStyleMarker) {
		t.Errorf("admin HTML did not contain injected stylesheet marker; body=%q", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "/styles/plugin.css") {
		t.Errorf("admin HTML missing plugin.css link; body=%q", rec.Body.String())
	}

	// Public path: snippet absent (we don't theme non-admin pages).
	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, httptest.NewRequest("GET", "/plugins/demo/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("public path status %d want 200; body=%q", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), adminStyleMarker) {
		t.Errorf("public HTML should not have injected stylesheet marker; body=%q", rec.Body.String())
	}
}

func TestInjectAdminStyles_Idempotent(t *testing.T) {
	html := []byte("<html><head></head><body></body></html>")
	once := injectAdminStyles(html)
	twice := injectAdminStyles(once)
	if string(once) != string(twice) {
		t.Errorf("second pass changed the document\nfirst:  %q\nsecond: %q", once, twice)
	}
}

func TestInjectAdminStyles_NoHeadStillInjects(t *testing.T) {
	// A body-only fragment still gets the link prepended so the iframe is
	// themed even when the plugin emits a bare HTML snippet.
	html := []byte("<body>just a fragment</body>")
	out := injectAdminStyles(html)
	if !strings.Contains(string(out), adminStyleMarker) {
		t.Errorf("expected marker in output: %q", out)
	}
	if !strings.HasPrefix(string(out), "<link") {
		t.Errorf("expected link prepended, got: %q", out)
	}
}

// TestServer_PathTraversal_CannotEscapePublicRoot verifies that traversal
// attempts ("../") in the URL cannot reach files outside the plugin's
// PublicFS root. This is the actual security property — that the on-disk
// boundary of the public directory is respected.
func TestServer_PathTraversal_CannotEscapePublicRoot(t *testing.T) {
	// Sibling directory layout:
	//   <tmp>/
	//     plugin-public/      ← this becomes the plugin's PublicFS
	//       safe.txt
	//     forbidden/
	//       secret.txt        ← must NOT be reachable via traversal
	root := t.TempDir()
	pluginPublic := root + "/plugin-public"
	forbidden := root + "/forbidden"
	for _, d := range []string{pluginPublic, forbidden} {
		if err := osMkdirAll(d); err != nil {
			t.Fatal(err)
		}
	}
	osWriteFile(t, pluginPublic+"/safe.txt", []byte("safe content"))
	osWriteFile(t, forbidden+"/secret.txt", []byte("FORBIDDEN_CONTENT"))

	loaded := &Loaded{
		Manifest: &Manifest{
			API: "1", DisplayName: "demo", Slug: "demo", Version: "1.0.0",
			Permissions: []string{"http.serve"},
		},
		PublicFS: osDirFS(pluginPublic),
	}
	s := NewServer([]*Loaded{loaded})

	// Sanity: in-root file is reachable
	{
		req := httptest.NewRequest("GET", "/plugins/demo/safe.txt", nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "safe content") {
			t.Fatalf("baseline reachability broken: status=%d body=%q", rec.Code, rec.Body.String())
		}
	}

	// Real test: every traversal attempt must NOT return forbidden content.
	traversalPaths := []string{
		"/plugins/demo/../forbidden/secret.txt",
		"/plugins/demo/../../forbidden/secret.txt",
		"/plugins/demo/./../forbidden/secret.txt",
		"/plugins/demo/safe.txt/../../forbidden/secret.txt",
	}
	for _, p := range traversalPaths {
		req := httptest.NewRequest("GET", p, nil)
		req.URL.Path = p
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
		if strings.Contains(rec.Body.String(), "FORBIDDEN_CONTENT") {
			t.Errorf("%s: traversal escaped PublicFS root: %q", p, rec.Body.String())
		}
	}
}

// TestServer_AdminGate_NotBypassableByTraversal verifies that a path-normalization
// trick cannot reach an admin-gated static asset unauthenticated. The bug: the
// admin gate checked the raw rest path while the static server resolved a
// separately path.Clean'd path, so a request whose decoded path contained "..",
// "." or "//" (e.g. arriving percent-encoded as %2e%2e) evaded IsAdminPath yet
// still resolved to the protected file. ServeHTTP now normalizes once up front
// so the gate and the file server can't disagree.
func TestServer_AdminGate_NotBypassableByTraversal(t *testing.T) {
	loaded := &Loaded{
		Manifest: &Manifest{
			API: "1", DisplayName: "demo", Slug: "demo", Version: "1.0.0",
			Permissions: []string{"http.serve"},
		},
		PublicFS: fstest.MapFS{
			"admin/secret.txt": {Data: []byte("ADMIN_SECRET")},
			"public.txt":       {Data: []byte("public ok")},
		},
		adminGlobs: []glob.Glob{glob.MustCompile("/admin/*")},
		adminPaths: []string{"/admin/*"},
	}
	s := NewServer([]*Loaded{loaded})
	s.IsAuthenticated = func(*http.Request) bool { return false }
	// An unauthenticated admin gate: anything routed through RequireAdmin is
	// blocked, so a 401 (or anything that isn't the secret) means "gated".
	s.RequireAdmin = func(http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "admin required", http.StatusUnauthorized)
		}
	}

	// Baselines: a public file is served, the admin file is gated when hit directly.
	{
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, httptest.NewRequest("GET", "/plugins/demo/public.txt", nil))
		if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "public ok") {
			t.Fatalf("baseline public file broken: status=%d body=%q", rec.Code, rec.Body.String())
		}
	}
	{
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, httptest.NewRequest("GET", "/plugins/demo/admin/secret.txt", nil))
		if strings.Contains(rec.Body.String(), "ADMIN_SECRET") {
			t.Fatalf("baseline admin gate broken: secret served directly: %q", rec.Body.String())
		}
	}

	// Every normalization trick must resolve to the admin path (and be gated),
	// never slip the secret out unauthenticated. URL.Path is set directly to the
	// decoded form net/http hands ServeHTTP for the matching %xx-encoded request.
	for _, p := range []string{
		"/plugins/demo/x/../admin/secret.txt",
		"/plugins/demo/./admin/secret.txt",
		"/plugins/demo//admin/secret.txt",
		"/plugins/demo/admin/../admin/secret.txt",
	} {
		req := httptest.NewRequest("GET", "/plugins/demo/", nil)
		req.URL.Path = p
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
		if strings.Contains(rec.Body.String(), "ADMIN_SECRET") {
			t.Errorf("%s: admin gate bypassed, secret leaked: status=%d body=%q", p, rec.Code, rec.Body.String())
		}
	}
}

// TestWritePluginHTTPResponse_StripsCoreSessionCookie is the A4 regression: a
// plugin (even one without auth.gate, so sessionCookies is empty) must not be
// able to set or clobber the core-owned owncast_session cookie via a Set-Cookie
// response header.
func TestWritePluginHTTPResponse_StripsCoreSessionCookie(t *testing.T) {
	out := []byte(`{"status":200,"headers":{"Set-Cookie":"` + SessionCookieName + `=forged; Path=/"},"body":"ok"}`)
	rec := httptest.NewRecorder()

	writePluginHTTPResponse(rec, out, false, nil /* no host-minted session cookies */)

	for _, c := range rec.Result().Cookies() {
		if c.Name == SessionCookieName {
			t.Fatalf("plugin-set %s cookie was not stripped: %q", SessionCookieName, c.Value)
		}
	}
}

// TestIsAdminPath_CaseInsensitive is the A5 regression: the admin gate must
// fire regardless of the request path's case, so a case-insensitively-routing
// plugin can't be tricked into serving admin content via "/Admin/...".
func TestIsAdminPath_CaseInsensitive(t *testing.T) {
	// adminGlobs/adminPaths are stored lowered by loadFromBytes; mirror that.
	l := &Loaded{
		adminGlobs: []glob.Glob{glob.MustCompile("/admin/*")},
		adminPaths: []string{"/admin/*"},
	}
	for _, p := range []string{"/admin/secret", "/Admin/secret", "/ADMIN/secret"} {
		if !l.IsAdminPath(p) {
			t.Errorf("%q should be gated as an admin path", p)
		}
	}
	if l.IsAdminPath("/public/x") {
		t.Error("/public/x should not be gated")
	}
}

func TestServer_NonGetFallsThroughStatic(t *testing.T) {
	// POST to a static asset path shouldn't serve the file. (Static is
	// read-only; non-GET/HEAD requests fall through to the dynamic handler.)
	s := staticOnlyServer(t, []string{"http.serve"}, fstest.MapFS{
		"data.txt": {Data: []byte("static data")},
	})
	req := httptest.NewRequest("POST", "/plugins/demo/data.txt", strings.NewReader(""))
	rec := httptest.NewRecorder()
	func() {
		defer func() { _ = recover() }() // would dereference nil plugin
		s.ServeHTTP(rec, req)
	}()
	if strings.Contains(rec.Body.String(), "static data") {
		t.Errorf("POST returned static content: %q", rec.Body.String())
	}
}

func TestIsAllowedResponseHeader(t *testing.T) {
	allowed := []string{
		"Content-Type", "content-type", "CONTENT-TYPE", // case-insensitive
		"Cache-Control", "ETag", "Last-Modified", "Location", "Vary",
		"Access-Control-Allow-Origin", "Access-Control-Allow-Methods",
		"Content-Encoding", "Content-Language", "Link",
		"Set-Cookie", "set-cookie", // plugins can set cookies in their own namespace
	}
	for _, h := range allowed {
		if !isAllowedResponseHeader(h) {
			t.Errorf("%q should be allowed", h)
		}
	}

	denied := []string{
		"Authorization",
		"WWW-Authenticate",
		"X-Custom-Plugin-Header",
		"Content-Security-Policy",   // host owns CSP, not plugins
		"Strict-Transport-Security", // host owns transport security
		"Server",                    // host identifies itself
		"X-Frame-Options",
	}
	for _, h := range denied {
		if isAllowedResponseHeader(h) {
			t.Errorf("%q should be DENIED but was allowed", h)
		}
	}
}
