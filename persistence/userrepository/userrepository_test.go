package userrepository

import (
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
// is the durable ownership check behind that, matching the "<slug>:<authId>"
// plugin.auth identity RegisterUser stores.
func TestUserRegisteredByPlugin(t *testing.T) {
	r := newTestRepo(t)

	user, _, err := r.CreateAnonymousUser("gated viewer")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	// Mirror what RegisterUser stores: "<slug>:<authId>" under PluginAuth.
	if err := r.AddAuth(user.ID, "basic-auth:shared", models.PluginAuth); err != nil {
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
		// Slug prefix must respect the ":" separator, not match a partial slug.
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
