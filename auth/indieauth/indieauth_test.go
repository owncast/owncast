package indieauth

import (
	"testing"

	"github.com/owncast/owncast/utils"
)

func TestLimitGlobalPendingRequests(t *testing.T) {
	// client_id and redirect_uri must be valid same-host absolute URLs to
	// pass StartServerAuth's validation; the slug just keeps them unique.
	clientURLs := func() (string, string) {
		slug, _ := utils.GenerateRandomString(10)
		return "https://client.example/" + slug, "https://client.example/" + slug + "/callback"
	}

	// Simulate maxPendingRequests-1 pending requests.
	for i := 0; i < maxPendingRequests-1; i++ {
		cid, redirectURL := clientURLs()
		cc, _ := utils.GenerateRandomString(10)
		state, _ := utils.GenerateRandomString(10)
		me, _ := utils.GenerateRandomString(10)

		_, err := StartServerAuth(cid, redirectURL, cc, state, me)
		if err != nil {
			t.Error("Registration should be permitted.", i, " of ", len(pendingServerAuthRequests), err)
		}
	}

	// This should throw an error
	cid, redirectURL := clientURLs()
	cc, _ := utils.GenerateRandomString(10)
	state, _ := utils.GenerateRandomString(10)
	me, _ := utils.GenerateRandomString(10)

	_, err := StartServerAuth(cid, redirectURL, cc, state, me)
	if err == nil {
		t.Error("Registration should not be permitted.")
	}
}

func TestRejectMismatchedRedirectURI(t *testing.T) {
	// Reset the package-level pending-request map; an earlier test in this
	// package fills it, and StartServerAuth rejects once it's near capacity.
	pendingServerAuthRequests = map[string]ServerAuthRequest{}

	// A redirect_uri on a different host than client_id must be rejected so
	// the auth endpoint can't be used as an open redirect.
	if _, err := StartServerAuth("https://client.example", "https://attacker.example/callback", "cc", "state", "me"); err == nil {
		t.Error("redirect_uri on a foreign host should be rejected")
	}

	// A same-host redirect_uri is accepted.
	if _, err := StartServerAuth("https://client.example", "https://client.example/callback", "cc", "state", "me"); err != nil {
		t.Error("same-host redirect_uri should be permitted:", err)
	}
}
