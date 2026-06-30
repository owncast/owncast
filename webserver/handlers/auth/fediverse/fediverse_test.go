package fediverse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	fediverseauth "github.com/owncast/owncast/auth/fediverse"
	"github.com/owncast/owncast/models"
)

// When a code is already pending for the same account, RegisterFediverseOTP
// returns success=false with no error (so we don't send a second message).
// The handler must treat that as success and let the client continue to the
// code entry step, rather than dead-ending the user with an error.
func TestRegisterFallsThroughWhenCodeAlreadyPending(t *testing.T) {
	const accessToken = "test-access-token"
	const account = "someone@example.com"

	svc := fediverseauth.New()

	// Seed a pending request directly on the service so the handler's own
	// register call hits the duplicate path. Going through the service here
	// (rather than the handler) avoids sending the first OTP message, which
	// the duplicate path under test never reaches anyway.
	if _, ok, err := svc.RegisterFediverseOTP(accessToken, "uid", "display name", account); err != nil || !ok {
		t.Fatalf("seeding a pending request failed: ok=%v err=%v", ok, err)
	}

	// Only FediverseAuth is needed; the duplicate path returns before touching
	// activitypub, chat, or the config repository.
	h := New(Deps{FediverseAuth: svc})

	body := strings.NewReader(`{"account":"` + account + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/fediverse?accessToken="+accessToken, body)
	rec := httptest.NewRecorder()

	h.RegisterFediverseOTPRequest(models.User{}, rec, req)

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if !resp.Success {
		t.Fatalf("expected success when a code is already pending, got success=false (message=%q)", resp.Message)
	}
}
