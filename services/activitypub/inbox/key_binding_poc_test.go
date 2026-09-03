package inbox

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-fed/httpsig"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/datastore"
)

func TestVerifyRejectsCrossHostKeyOwner(t *testing.T) {
	attackerKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(&attackerKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	publicPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER}))

	const victim = "https://victim.example/users/alice"
	var keyID string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/activity+json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"@context":          []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"},
			"id":                victim,
			"type":              "Person",
			"preferredUsername": "alice",
			"inbox":             victim + "/inbox",
			"publicKey": map[string]any{
				"id":           keyID,
				"owner":        victim,
				"publicKeyPem": publicPEM,
			},
		})
	}))
	defer server.Close()
	keyID = server.URL + "/key"

	ds, err := datastore.SetupPersistence(":memory:", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ds.DB.Close() })
	cfg := configrepository.New(ds)
	localPrivatePEM, localPublicPEM, err := apcrypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}
	for _, set := range []func() error{
		func() error { return cfg.SetServerURL("https://local.example") },
		func() error { return cfg.SetFederationUsername("streamer") },
		func() error { return cfg.SetPrivateKey(string(localPrivatePEM)) },
		func() error { return cfg.SetPublicKey(string(localPublicPEM)) },
	} {
		if err := set(); err != nil {
			t.Fatal(err)
		}
	}
	localSigner := apcrypto.New(apcrypto.Deps{ConfigRepository: cfg})
	builder := apmodels.New(apmodels.Deps{ConfigRepository: cfg, Signer: localSigner})
	resolver := resolvers.New(resolvers.Deps{ConfigRepository: cfg, Builder: builder, Signer: localSigner})
	service := New(Deps{
		ConfigRepository: cfg,
		Followers:        followersrepository.New(ds),
		Resolver:         resolver,
	})

	body := []byte(`{"type":"Like","actor":"` + victim + `"}`)
	request, err := http.NewRequest(http.MethodPost, "https://local.example/inbox", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	request.Header.Set("Host", request.URL.Host)
	signer, _, err := httpsig.NewSigner(
		[]httpsig.Algorithm{httpsig.RSA_SHA256},
		httpsig.DigestSha256,
		[]string{httpsig.RequestTarget, "host", "date", "digest"},
		httpsig.Signature,
		0,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := signer.SignRequest(attackerKey, keyID, request, body); err != nil {
		t.Fatal(err)
	}

	owner, err := service.Verify(request, body)
	if err == nil {
		t.Fatalf("expected cross-host key ownership to be rejected, got owner %v", owner)
	}
}
