package webfinger

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestMain sets the test escape hatches before any test runs: they are read
// once via sync.Once, and httptest servers use loopback + self-signed certs
// that the resolver would otherwise reject.
func TestMain(m *testing.M) {
	if err := os.Setenv("OWNCAST_ALLOW_INTERNAL_FEDERATION", "true"); err != nil {
		panic(err)
	}
	if err := os.Setenv("OWNCAST_INSECURE_SKIP_VERIFY", "true"); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// Regression test for #4696: a webfinger host that redirects the lookup must be
// followed, and the acct resource query must survive the redirect.
func TestGetWebfingerLinksFollowsRedirect(t *testing.T) {
	var gotResource string
	final := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotResource = r.URL.Query().Get("resource")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"links":[{"rel":"self","type":"application/activity+json","href":"https://example.com/actor"}]}`))
	}))
	defer final.Close()

	entry := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, final.URL+r.URL.Path+"?"+r.URL.RawQuery, http.StatusFound)
	}))
	defer entry.Close()

	account := "user@" + strings.TrimPrefix(entry.URL, "https://")

	links, err := GetWebfingerLinks(account)
	if err != nil {
		t.Fatalf("expected redirect to be followed, got error: %v", err)
	}
	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}
	if want := "acct:" + account; gotResource != want {
		t.Errorf("resource query not preserved across redirect: got %q, want %q", gotResource, want)
	}
}

func TestGetWebfingerLinksRedirectLoop(t *testing.T) {
	loop := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}))
	defer loop.Close()

	account := "user@" + strings.TrimPrefix(loop.URL, "https://")
	if _, err := GetWebfingerLinks(account); err == nil {
		t.Error("expected an error from an endless redirect chain, got nil")
	}
}

func TestGetWebfingerLinksTerminalError(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer server.Close()

	account := "user@" + strings.TrimPrefix(server.URL, "https://")
	if _, err := GetWebfingerLinks(account); err == nil {
		t.Error("expected an error for a non-200 webfinger response, got nil")
	}
}
