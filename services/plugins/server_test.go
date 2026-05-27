package plugins

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"testing/fstest"
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
			API: "1", Name: "demo", Version: "1.0.0", Permissions: perms,
		},
		AssetsFS: assets,
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

	for _, tc := range []struct {
		path string
		want string
	}{
		{"/plugins/demo/", "<h1>root</h1>"},
		{"/plugins/demo/sub", "<h1>sub</h1>"},
		{"/plugins/demo/sub/", "<h1>sub</h1>"},
	} {
		req := httptest.NewRequest("GET", tc.path, nil)
		rec := httptest.NewRecorder()
		s.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("%s: status %d want 200", tc.path, rec.Code)
			continue
		}
		if rec.Body.String() != tc.want {
			t.Errorf("%s: body %q want %q", tc.path, rec.Body.String(), tc.want)
		}
	}
}

// TestServer_PathTraversal_CannotEscapeAssetsRoot verifies that traversal
// attempts ("../") in the URL cannot reach files outside the plugin's
// AssetsFS root. This is the actual security property — that the on-disk
// boundary of the asset directory is respected.
func TestServer_PathTraversal_CannotEscapeAssetsRoot(t *testing.T) {
	// Sibling directory layout:
	//   <tmp>/
	//     plugin-assets/      ← this becomes the plugin's AssetsFS
	//       safe.txt
	//     forbidden/
	//       secret.txt        ← must NOT be reachable via traversal
	root := t.TempDir()
	pluginAssets := root + "/plugin-assets"
	forbidden := root + "/forbidden"
	for _, d := range []string{pluginAssets, forbidden} {
		if err := osMkdirAll(d); err != nil {
			t.Fatal(err)
		}
	}
	osWriteFile(t, pluginAssets+"/safe.txt", []byte("safe content"))
	osWriteFile(t, forbidden+"/secret.txt", []byte("FORBIDDEN_CONTENT"))

	loaded := &Loaded{
		Manifest: &Manifest{
			API: "1", Name: "demo", Version: "1.0.0",
			Permissions: []string{"http.serve"},
		},
		AssetsFS: osDirFS(pluginAssets),
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
			t.Errorf("%s: traversal escaped AssetsFS root: %q", p, rec.Body.String())
		}
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
	}
	for _, h := range allowed {
		if !isAllowedResponseHeader(h) {
			t.Errorf("%q should be allowed", h)
		}
	}

	denied := []string{
		"Set-Cookie",
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
