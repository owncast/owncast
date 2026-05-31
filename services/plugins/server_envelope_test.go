package plugins

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The host hands a plugin a trusted, resolved req.user but must never leak the
// raw access token the request carried (cookie or accessToken query param).
func TestBuildRequestEnvelope_RedactsTokenAndForwardsUser(t *testing.T) {
	s := NewServer(nil)
	s.GetRequestUser = func(r *http.Request) *HostUser {
		return &HostUser{ID: "u-1", DisplayName: "Ann", IsAuthenticated: true}
	}

	req := httptest.NewRequest(http.MethodGet, "/plugins/demo/api/me?accessToken=secret&keep=yes", nil)
	req.Header.Set("Cookie", "owncast_chat_token=secret; theme=dark")
	req.Header.Set("X-Demo", "hi")

	raw, err := s.buildRequestEnvelope(req, "/api/me", false, nil)
	if err != nil {
		t.Fatalf("buildRequestEnvelope: %v", err)
	}

	var env struct {
		Query   map[string]string `json:"query"`
		Headers map[string]string `json:"headers"`
		User    *HostUser         `json:"user"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	// The raw access token must never reach the plugin.
	if _, ok := env.Query["accessToken"]; ok {
		t.Error("accessToken query param was not redacted")
	}
	if _, ok := env.Headers["Cookie"]; ok {
		t.Error("Cookie header was not redacted")
	}

	// Non-sensitive query params and headers still pass through.
	if env.Query["keep"] != "yes" {
		t.Errorf("benign query param dropped: got %q", env.Query["keep"])
	}
	if env.Headers["X-Demo"] != "hi" {
		t.Errorf("benign header dropped: got %q", env.Headers["X-Demo"])
	}

	// Identity still reaches the plugin via the trusted, host-resolved user.
	if env.User == nil || env.User.ID != "u-1" {
		t.Errorf("expected req.user forwarded, got %+v", env.User)
	}
}

// With no GetRequestUser wired (or no identity on the request), the envelope
// simply omits user, which plugins treat as optional.
func TestBuildRequestEnvelope_OmitsUserWhenAbsent(t *testing.T) {
	s := NewServer(nil) // GetRequestUser left nil

	req := httptest.NewRequest(http.MethodGet, "/plugins/demo/api/me", nil)
	raw, err := s.buildRequestEnvelope(req, "/api/me", false, nil)
	if err != nil {
		t.Fatalf("buildRequestEnvelope: %v", err)
	}

	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if _, ok := env["user"]; ok {
		t.Error("expected user to be omitted when no identity is present")
	}
}
