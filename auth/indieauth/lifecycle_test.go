package indieauth

import (
	"net/url"
	"testing"
	"time"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
)

func newLifecycleTestService(t *testing.T) *Service {
	t.Helper()
	ds, err := datastore.SetupPersistence(":memory:", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ds.DB.Close() })
	config := configrepository.New(ds)
	if err := config.SetServerURL("https://owncast.example"); err != nil {
		t.Fatal(err)
	}
	return New(Deps{ConfigRepository: config})
}

func TestServerAuthCodeIsSingleUse(t *testing.T) {
	svc := newLifecycleTestService(t)
	request, err := svc.StartServerAuth("https://client.example", "https://client.example/callback", createCodeChallenge("verifier"), "state", "me")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CompleteServerAuth(request.Code, request.RedirectURI, request.ClientID, "verifier"); err != nil {
		t.Fatalf("first completion: %v", err)
	}
	if _, err := svc.CompleteServerAuth(request.Code, request.RedirectURI, request.ClientID, "verifier"); err == nil {
		t.Fatal("replayed authorization code was accepted")
	}
}

func TestServerAuthRejectsMismatchedCodeVerifier(t *testing.T) {
	svc := newLifecycleTestService(t)
	verifier := randString(50)
	request, err := svc.StartServerAuth("https://client.example", "https://client.example/callback", createCodeChallenge(verifier), randString(20), "me")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CompleteServerAuth(request.Code, request.RedirectURI, request.ClientID, randString(50)); err == nil {
		t.Fatal("mismatched code verifier was accepted")
	}
	// A failed attempt consumes the code, so the real verifier cannot retry it.
	if _, err := svc.CompleteServerAuth(request.Code, request.RedirectURI, request.ClientID, verifier); err == nil {
		t.Fatal("authorization code was accepted after a failed verification")
	}
}

func TestClientCallbackStateMustMatchAndIsSingleUse(t *testing.T) {
	svc := newLifecycleTestService(t)
	request, err := createAuthRequestForTest()
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.storePendingAuthRequest(request); err != nil {
		t.Fatal(err)
	}

	if _, _, err := svc.HandleCallbackCode("code", randString(20)); err == nil {
		t.Fatal("mismatched state was accepted")
	}

	// The matching state is consumed. The token exchange fails because the
	// endpoint is unreachable, which is fine: only replay matters here.
	_, _, _ = svc.HandleCallbackCode("code", request.State)
	if _, _, err := svc.HandleCallbackCode("code", request.State); err == nil || err.Error() != "no auth requests pending" {
		t.Fatalf("reused state was not rejected as unknown: %v", err)
	}
}

// createAuthRequestForTest runs the real request builder, skipping only the
// network discovery of the authorization endpoint.
func createAuthRequestForTest() (*Request, error) {
	me, err := url.Parse("https://user.example/")
	if err != nil {
		return nil, err
	}
	endpoint, err := url.Parse("https://127.0.0.1:1/auth")
	if err != nil {
		return nil, err
	}
	return buildAuthRequest(me, endpoint, "user", "name", "token", "https://owncast.example")
}

func TestAuthRequestGeneratesFreshStateAndPKCE(t *testing.T) {
	isLetters := func(s string) bool {
		for _, c := range s {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
				return false
			}
		}
		return true
	}

	a, err := createAuthRequestForTest()
	if err != nil {
		t.Fatal(err)
	}
	b, err := createAuthRequestForTest()
	if err != nil {
		t.Fatal(err)
	}

	if len(a.CodeVerifier) != 50 || !isLetters(a.CodeVerifier) {
		t.Fatalf("verifier %q is not 50 letters", a.CodeVerifier)
	}
	if len(a.State) != 20 || !isLetters(a.State) {
		t.Fatalf("state %q is not 20 letters", a.State)
	}
	if a.CodeChallenge != createCodeChallenge(a.CodeVerifier) {
		t.Fatal("code challenge is not derived from the verifier")
	}

	q := a.Redirect.Query()
	if q.Get("state") != a.State || q.Get("code_challenge") != a.CodeChallenge || q.Get("code_challenge_method") != "S256" {
		t.Fatalf("redirect query does not carry the generated state and S256 challenge: %s", a.Redirect.RawQuery)
	}

	if a.State == b.State || a.CodeVerifier == b.CodeVerifier {
		t.Fatal("two auth requests share a state or verifier")
	}
}

func TestExpiredAuthRequestsAreRejectedAndPruned(t *testing.T) {
	svc := newLifecycleTestService(t)
	expired := time.Now().Add(-registrationTimeout - time.Second)
	svc.pendingAuthRequests["client"] = &Request{State: "client", Timestamp: expired}
	svc.pendingServerAuthRequests["server"] = ServerAuthRequest{Code: "server", Timestamp: expired}

	if _, _, err := svc.HandleCallbackCode("code", "client"); err == nil {
		t.Fatal("expired client request was accepted")
	}
	if _, err := svc.CompleteServerAuth("server", "https://client.example/callback", "https://client.example", "verifier"); err == nil {
		t.Fatal("expired server request was accepted")
	}

	svc.pendingAuthRequests["client"] = &Request{State: "client", Timestamp: expired}
	svc.pendingServerAuthRequests["server"] = ServerAuthRequest{Code: "server", Timestamp: expired}
	svc.pruneExpiredRequests(time.Now())
	if len(svc.pendingAuthRequests) != 0 || len(svc.pendingServerAuthRequests) != 0 {
		t.Fatal("expired auth requests were not pruned")
	}
}
