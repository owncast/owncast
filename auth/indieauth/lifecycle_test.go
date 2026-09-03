package indieauth

import (
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
