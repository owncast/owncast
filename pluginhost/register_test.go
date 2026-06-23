package pluginhost

import (
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/plugins"
)

// validateProfileURL is a trust boundary: a plugin-supplied profile URL is
// stored as a verified identity and later rendered as a clickable public link,
// so only plain http(s) links may pass. Non-http schemes (javascript:, data:),
// hostless values, and unparseable input must be rejected.
func TestValidateProfileURL(t *testing.T) {
	cases := []struct {
		url string
		ok  bool
	}{
		{"", true}, // a gate-only plugin surfaces nothing publicly
		{"https://github.com/octocat", true},
		{"http://example.com/u/1", true},
		{"javascript:alert(1)", false},
		{"data:text/html,<script>", false},
		{"ftp://example.com", false},
		{"https://", false},       // no host
		{"/relative/path", false}, // no scheme or host
		{"not a url", false},
		{"ht tp://x", false}, // unparseable (space)
	}
	for _, tc := range cases {
		err := validateProfileURL(tc.url)
		if tc.ok && err != nil {
			t.Errorf("validateProfileURL(%q) = %v, want nil", tc.url, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("validateProfileURL(%q) = nil, want error", tc.url)
		}
	}
}

// newRegisterEnv wires the auth host functions against a fresh in-memory
// datastore, returning the live RegisterUser closure and the repository behind
// it so tests can assert against stored state.
func newRegisterEnv(t *testing.T) (*plugins.HostEnv, userrepository.UserRepository) {
	t.Helper()
	ds, err := datastore.SetupPersistence(":memory:", t.TempDir())
	if err != nil {
		t.Fatalf("setup persistence: %v", err)
	}
	users := userrepository.New(ds)
	env := &plugins.HostEnv{}
	wireAuthHostFns(env, Deps{Datastore: ds, UserRepository: users})
	if env.RegisterUser == nil {
		t.Fatal("RegisterUser was not wired (signing secret failed to establish)")
	}
	return env, users
}

// TestRegisterUser exercises the whole users.register flow: find-or-create
// idempotency, per-plugin namespacing, and that scope/URL rejections fail
// before any user is created (so a rejected call leaves no orphan account).
func TestRegisterUser(t *testing.T) {
	env, users := newRegisterEnv(t)

	// Happy path creates an authenticated user.
	id, err := env.RegisterUser("github-auth", plugins.UserRegisterRequest{
		AuthID:      "583231",
		DisplayName: "octocat",
		ProfileURL:  "https://github.com/octocat",
		Handle:      "octocat",
		Public:      true,
	})
	if err != nil || id == "" {
		t.Fatalf("register: id=%q err=%v", id, err)
	}

	// Idempotent: the same identity resolves to the same user, not a duplicate.
	again, err := env.RegisterUser("github-auth", plugins.UserRegisterRequest{AuthID: "583231"})
	if err != nil || again != id {
		t.Fatalf("re-register: got %q err=%v, want %q", again, err, id)
	}

	// A different plugin with the same raw authId is a distinct user. (A new-user
	// registration must carry a display name; CreateAnonymousUser rejects empty.)
	other, err := env.RegisterUser("discord-auth", plugins.UserRegisterRequest{AuthID: "583231", DisplayName: "octocat"})
	if err != nil {
		t.Fatalf("other plugin register: %v", err)
	}
	if other == id {
		t.Error("a different plugin must not resolve to the same user")
	}

	// A disallowed scope is rejected, and no orphan user is left behind.
	if _, err := env.RegisterUser("github-auth", plugins.UserRegisterRequest{
		AuthID: "scope-reject", Scopes: []string{models.ScopeHasAdminAccess},
	}); err == nil {
		t.Error("expected admin scope to be rejected")
	}
	if u := users.GetUserByPluginAuth("github-auth", "scope-reject"); u != nil {
		t.Error("scope-rejected registration left an orphan user")
	}

	// A bad profile URL is rejected, and no orphan user is left behind.
	if _, err := env.RegisterUser("github-auth", plugins.UserRegisterRequest{
		AuthID: "url-reject", ProfileURL: "javascript:alert(1)",
	}); err == nil {
		t.Error("expected bad profile URL to be rejected")
	}
	if u := users.GetUserByPluginAuth("github-auth", "url-reject"); u != nil {
		t.Error("url-rejected registration left an orphan user")
	}
}
