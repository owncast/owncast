package indieauth

import (
	"testing"

	"github.com/owncast/owncast/utils"
)

func TestLimitGlobalPendingRequests(t *testing.T) {
	// Construct an isolated Service for this test. CompleteServerAuth is
	// not called here, so the ConfigRepository can be nil. Each Service owns
	// its own pending-request map, so no global reset is needed.
	svc := New(Deps{ConfigRepository: nil})

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

		_, err := svc.StartServerAuth(cid, redirectURL, cc, state, me)
		if err != nil {
			t.Error("Registration should be permitted.", i, " of ", len(svc.pendingServerAuthRequests), err)
		}
	}

	// This should throw an error.
	cid, redirectURL := clientURLs()
	cc, _ := utils.GenerateRandomString(10)
	state, _ := utils.GenerateRandomString(10)
	me, _ := utils.GenerateRandomString(10)

	_, err := svc.StartServerAuth(cid, redirectURL, cc, state, me)
	if err == nil {
		t.Error("Registration should not be permitted.")
	}
}

func TestRejectMismatchedRedirectURI(t *testing.T) {
	svc := New(Deps{ConfigRepository: nil})

	// A redirect_uri on a different host than client_id must be rejected so
	// the auth endpoint can't be used as an open redirect.
	if _, err := svc.StartServerAuth("https://client.example", "https://attacker.example/callback", "cc", "state", "me"); err == nil {
		t.Error("redirect_uri on a foreign host should be rejected")
	}

	// A same-host redirect_uri is accepted.
	if _, err := svc.StartServerAuth("https://client.example", "https://client.example/callback", "cc", "state", "me"); err != nil {
		t.Error("same-host redirect_uri should be permitted:", err)
	}
}

func TestRejectNonWebRedirectURI(t *testing.T) {
	svc := New(Deps{ConfigRepository: nil})

	// Opaque/non-http(s) URIs (javascript:, data:, mailto:) have an empty
	// hostname; without a scheme+host check, two of them would pass the
	// host-match comparison and smuggle a hostile target into the redirect.
	cases := []struct {
		name        string
		clientID    string
		redirectURI string
	}{
		{"javascript scheme", "javascript:alert(1)", "javascript:alert(2)"},
		{"data scheme", "data:text/html,a", "data:text/html,b"},
		{"non-web client_id", "https://client.example", "mailto:someone@client.example"},
	}

	for _, tc := range cases {
		if _, err := svc.StartServerAuth(tc.clientID, tc.redirectURI, "cc", "state", "me"); err == nil {
			t.Errorf("%s: should be rejected", tc.name)
		}
	}
}
