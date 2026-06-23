package userrepository

import (
	"database/sql"
	"os"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/datastore"
)

func newTestRepo(t *testing.T) *SqlUserRepository {
	t.Helper()
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		t.Fatalf("setup persistence: %v", err)
	}
	return New(ds).(*SqlUserRepository)
}

// The scoped external API is for true third-party integrations
// (users.type = 'API'). A regular user that happens to carry an admin scope on
// an access token — e.g. one a viewer-auth plugin created and granted a session
// for — must NOT be able to authenticate to it.
func TestGetExternalAPIUserForAccessTokenAndScope_OnlyAPIUsers(t *testing.T) {
	r := newTestRepo(t)

	// A real API integration with the admin scope: authorized.
	const apiToken = "api-token-aaa"
	if err := r.InsertExternalAPIUser(apiToken, "integration", 0, []string{models.ScopeHasAdminAccess}); err != nil {
		t.Fatalf("insert API user: %v", err)
	}
	if u, err := r.GetExternalAPIUserForAccessTokenAndScope(apiToken, models.ScopeHasAdminAccess); err != nil || u == nil {
		t.Fatalf("API user with the scope should be authorized: user=%v err=%v", u, err)
	}

	// A STANDARD user carrying the same scope on a plain access token: rejected.
	user, _, err := r.CreateAnonymousUser("viewer")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := r.SetUserScopes(user.ID, []string{models.ScopeHasAdminAccess}); err != nil {
		t.Fatalf("set scopes: %v", err)
	}
	const stdToken = "std-token-bbb"
	if err := r.AddAccessTokenForUser(stdToken, user.ID); err != nil {
		t.Fatalf("add token: %v", err)
	}
	// No match returns a nil user (and sql.ErrNoRows); the security property is
	// that no authorized user comes back for a non-API token.
	if u, _ := r.GetExternalAPIUserForAccessTokenAndScope(stdToken, models.ScopeHasAdminAccess); u != nil {
		t.Fatalf("non-API user must NOT be authorized to the scoped external API: got %v", u)
	}
}

// grantSession is confined to users a plugin registered. UserRegisteredByPlugin
// is the durable ownership check behind that, matching the plugin.auth identity
// RegisterUser stores with provider=<slug>.
func TestUserRegisteredByPlugin(t *testing.T) {
	r := newTestRepo(t)

	user, _, err := r.CreateAnonymousUser("gated viewer")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	// Mirror what RegisterUser stores: auth_key=authId, provider=<slug>, under
	// the PluginAuth type.
	if err := r.AddAuth(user.ID, "shared", models.PluginAuth, &models.LinkedIdentityFields{Provider: "basic-auth"}); err != nil {
		t.Fatalf("add auth: %v", err)
	}

	cases := []struct {
		name   string
		plugin string
		userID string
		want   bool
	}{
		{"owning plugin", "basic-auth", user.ID, true},
		{"different plugin", "github-auth", user.ID, false},
		{"unknown user", "basic-auth", "no-such-user", false},
		// Provider is matched on exact equality, so a partial slug is no match.
		{"partial slug is not a match", "basic", user.ID, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := r.UserRegisteredByPlugin(tc.plugin, tc.userID); got != tc.want {
				t.Fatalf("UserRegisteredByPlugin(%q, %q) = %v, want %v", tc.plugin, tc.userID, got, tc.want)
			}
		})
	}
}

// A plugin identity (provider=slug, auth_key) must map to exactly one user. The
// partial unique index backs RegisterUser's find-or-create so a concurrent
// double-registration can't mint two users for the same external identity.
func TestPluginIdentityIsUnique(t *testing.T) {
	r := newTestRepo(t)

	a, _, err := r.CreateAnonymousUser("first")
	if err != nil {
		t.Fatalf("create user a: %v", err)
	}
	b, _, err := r.CreateAnonymousUser("second")
	if err != nil {
		t.Fatalf("create user b: %v", err)
	}

	githubID := &models.LinkedIdentityFields{Provider: "github-auth"}
	if err := r.AddAuth(a.ID, "583231", models.PluginAuth, githubID); err != nil {
		t.Fatalf("first link: %v", err)
	}
	// The same (provider, auth_key) for a different user must be rejected.
	if err := r.AddAuth(b.ID, "583231", models.PluginAuth, githubID); err == nil {
		t.Fatal("expected duplicate plugin identity to be rejected, got nil")
	}
	// The same raw id under a different provider is a distinct identity: allowed.
	if err := r.AddAuth(b.ID, "583231", models.PluginAuth, &models.LinkedIdentityFields{Provider: "discord-auth"}); err != nil {
		t.Fatalf("distinct provider link: %v", err)
	}
}

// TestGetUserByAuthBuiltinRoundTrip is the core regression guard for the
// token->auth_key rename: a returning IndieAuth/Fediverse user must still
// resolve by (auth_key, type), and type must keep the two providers distinct.
func TestGetUserByAuthBuiltinRoundTrip(t *testing.T) {
	r := newTestRepo(t)

	indie, _, err := r.CreateAnonymousUser("indie user")
	if err != nil {
		t.Fatalf("create indie user: %v", err)
	}
	const me = "https://me.example.com"
	if err := r.AddAuth(indie.ID, me, models.IndieAuth, &models.LinkedIdentityFields{ProfileURL: me}); err != nil {
		t.Fatalf("add indieauth: %v", err)
	}

	fedi, _, err := r.CreateAnonymousUser("fedi user")
	if err != nil {
		t.Fatalf("create fedi user: %v", err)
	}
	const account = "@me@host.example"
	if err := r.AddAuth(fedi.ID, account, models.Fediverse, &models.LinkedIdentityFields{Handle: account}); err != nil {
		t.Fatalf("add fediverse: %v", err)
	}

	if u := r.GetUserByAuth(me, models.IndieAuth); u == nil || u.ID != indie.ID {
		t.Errorf("GetUserByAuth(me, IndieAuth) = %v, want %s", u, indie.ID)
	}
	if u := r.GetUserByAuth(account, models.Fediverse); u == nil || u.ID != fedi.ID {
		t.Errorf("GetUserByAuth(account, Fediverse) = %v, want %s", u, fedi.ID)
	}
	// type discriminates: the right key under the wrong provider type matches no one.
	if u := r.GetUserByAuth(me, models.Fediverse); u != nil {
		t.Errorf("GetUserByAuth(me, Fediverse) = %v, want nil", u)
	}
	if u := r.GetUserByAuth("https://nobody.example", models.IndieAuth); u != nil {
		t.Errorf("GetUserByAuth(unknown) = %v, want nil", u)
	}
}

// TestPluginAuthNamespacing locks in the security invariant behind the
// type+provider split: a plugin identity resolves only within its own slug, and
// a plugin whose slug equals a built-in provider name ("fediverse") can neither
// resolve nor claim ownership of the built-in identity — that row's type is
// 'fediverse', not 'plugin.auth'.
func TestPluginAuthNamespacing(t *testing.T) {
	r := newTestRepo(t)

	// A real Fediverse user.
	fedi, _, err := r.CreateAnonymousUser("real fedi")
	if err != nil {
		t.Fatalf("create fedi: %v", err)
	}
	const account = "@victim@host.example"
	if err := r.AddAuth(fedi.ID, account, models.Fediverse, &models.LinkedIdentityFields{Handle: account}); err != nil {
		t.Fatalf("add fediverse: %v", err)
	}

	// A plugin whose slug is "fediverse" registers a user with the SAME raw key.
	pluginUser, _, err := r.CreateAnonymousUser("plugin user")
	if err != nil {
		t.Fatalf("create plugin user: %v", err)
	}
	if err := r.AddAuth(pluginUser.ID, account, models.PluginAuth, &models.LinkedIdentityFields{Provider: "fediverse"}); err != nil {
		t.Fatalf("add plugin auth: %v", err)
	}

	// The plugin lookup resolves to the plugin's own user, never the built-in one.
	if u := r.GetUserByPluginAuth("fediverse", account); u == nil || u.ID != pluginUser.ID {
		t.Errorf("GetUserByPluginAuth(fediverse, account) = %v, want %s", u, pluginUser.ID)
	}
	// A different slug with the same key resolves to nobody.
	if u := r.GetUserByPluginAuth("github-auth", account); u != nil {
		t.Errorf("GetUserByPluginAuth(github-auth, account) = %v, want nil", u)
	}
	// Ownership: the plugin owns the user it registered, but NOT the built-in one.
	if !r.UserRegisteredByPlugin("fediverse", pluginUser.ID) {
		t.Error("plugin should own the user it registered")
	}
	if r.UserRegisteredByPlugin("fediverse", fedi.ID) {
		t.Error("plugin slug 'fediverse' must NOT own the built-in Fediverse user")
	}
}

// TestAddAuthPersistsLinkedIdentityFields verifies AddAuth actually writes the
// provider/profile_url/handle/is_public columns. There is no Go reader for them
// yet, so assert against the row directly.
func TestAddAuthPersistsLinkedIdentityFields(t *testing.T) {
	r := newTestRepo(t)
	u, _, err := r.CreateAnonymousUser("plugin user")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	fields := &models.LinkedIdentityFields{
		Provider:   "github-auth",
		ProfileURL: "https://github.com/octocat",
		Handle:     "octocat",
		Public:     true,
	}
	if err := r.AddAuth(u.ID, "583231", models.PluginAuth, fields); err != nil {
		t.Fatalf("add auth: %v", err)
	}

	var authKey, typ, provider string
	var profileURL, handle sql.NullString
	var isPublic bool
	if err := r.datastore.DB.QueryRow(
		`SELECT auth_key, type, provider, profile_url, handle, is_public FROM auth WHERE user_id = ?`, u.ID,
	).Scan(&authKey, &typ, &provider, &profileURL, &handle, &isPublic); err != nil {
		t.Fatalf("read back: %v", err)
	}
	if authKey != "583231" || typ != string(models.PluginAuth) || provider != "github-auth" {
		t.Errorf("got auth_key=%q type=%q provider=%q", authKey, typ, provider)
	}
	if profileURL.String != "https://github.com/octocat" || handle.String != "octocat" || !isPublic {
		t.Errorf("profile fields not persisted: url=%v handle=%v public=%v", profileURL, handle, isPublic)
	}
}

// TestAuthProviderLabel covers the admin user-list label after the switch from
// token-prefix parsing to the provider column.
func TestAuthProviderLabel(t *testing.T) {
	cases := []struct {
		authType string
		provider string
		want     string
	}{
		{"indieauth", "indieauth", "IndieAuth"},
		{"fediverse", "fediverse", "Fediverse"},
		{"plugin.auth", "github-auth", "github-auth"},
		{"plugin.auth", "", "Plugin"},
		{"something-else", "something-else", "something-else"},
	}
	for _, tc := range cases {
		if got := authProviderLabel(tc.authType, tc.provider); got != tc.want {
			t.Errorf("authProviderLabel(%q, %q) = %q, want %q", tc.authType, tc.provider, got, tc.want)
		}
	}
}
