package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/owncast/owncast/config"
	plugins "github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/utils"
)

// A proposed display name that sanitizes to an empty string (e.g. only HTML
// tags) must fall back to a generated name instead of failing registration.
func TestRegisterAnonymousChatUserEmptyAfterSanitizationFallsBack(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		headers     map[string]string
		wantExactly string
	}{
		{name: "body display name of only HTML", body: `{"displayName": "<b></b>"}`},
		{name: "body display name of HTML around text keeps the text", body: `{"displayName": "<b>gabe</b>"}`, wantExactly: "gabe"},
		{name: "X-Forwarded-User of only HTML", body: `{}`, headers: map[string]string{"X-Forwarded-User": "<b></b>"}},
		{name: "body display name of only whitespace", body: `{"displayName": "   "}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/chat/register", strings.NewReader(tc.body))
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			w := httptest.NewRecorder()

			testHandlers.RegisterAnonymousChatUser(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200, got %d", resp.StatusCode)
			}

			var response struct {
				ID          string `json:"id"`
				AccessToken string `json:"accessToken"`
				DisplayName string `json:"displayName"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if response.DisplayName == "" {
				t.Error("expected a non-empty generated display name")
			}
			if utf8.RuneCountInString(response.DisplayName) > config.MaxChatDisplayNameLength {
				t.Errorf("display name %q exceeds max length %d", response.DisplayName, config.MaxChatDisplayNameLength)
			}
			if tc.wantExactly != "" && response.DisplayName != tc.wantExactly {
				t.Errorf("expected display name %q, got %q", tc.wantExactly, response.DisplayName)
			}
			if response.AccessToken == "" {
				t.Error("expected a non-empty access token")
			}
		})
	}
}

// TestRegisterAnonymousChatUserGateSession verifies that when a verified
// auth-gate session is present in the request context, /api/chat/register
// returns the gate user's existing token instead of creating a new anonymous
// user.
func TestRegisterAnonymousChatUserGateSession(t *testing.T) {
	// Create and authenticate a gate user.
	user, _, err := testHandlers.userRepository.CreateAnonymousUser("GateUser")
	if err != nil {
		t.Fatalf("CreateAnonymousUser: %v", err)
	}
	if err := testHandlers.userRepository.SetUserAsAuthenticated(user.ID); err != nil {
		t.Fatalf("SetUserAsAuthenticated: %v", err)
	}
	// Mint a separate access token simulating what grantPluginSession does.
	gateTok, err := utils.GenerateAccessToken()
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}
	if err := testHandlers.userRepository.AddAccessTokenForUser(gateTok, user.ID); err != nil {
		t.Fatalf("AddAccessTokenForUser: %v", err)
	}
	countBefore := testHandlers.userRepository.GetUsersCount()

	req := httptest.NewRequest(http.MethodPost, "/api/chat/register", strings.NewReader(`{}`))
	req = req.WithContext(plugins.WithSessionToken(req.Context(), gateTok))
	w := httptest.NewRecorder()

	testHandlers.RegisterAnonymousChatUser(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var response struct {
		ID          string `json:"id"`
		AccessToken string `json:"accessToken"`
		DisplayName string `json:"displayName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if response.AccessToken != gateTok {
		t.Errorf("expected gate token %q, got %q", gateTok, response.AccessToken)
	}
	if response.ID != user.ID {
		t.Errorf("expected user ID %q, got %q", user.ID, response.ID)
	}
	if response.DisplayName != "GateUser" {
		t.Errorf("expected display name %q, got %q", "GateUser", response.DisplayName)
	}
	if got := testHandlers.userRepository.GetUsersCount(); got != countBefore {
		t.Fatalf("user count changed from %d to %d", countBefore, got)
	}
}

// TestRegisterAnonymousChatUserGateSessionInvalidToken verifies that a gate
// session carrying a token that doesn't resolve to any user returns 403 and
// does NOT create a new anonymous user.
func TestRegisterAnonymousChatUserGateSessionInvalidToken(t *testing.T) {
	countBefore := testHandlers.userRepository.GetUsersCount()

	req := httptest.NewRequest(http.MethodPost, "/api/chat/register", strings.NewReader(`{}`))
	req = req.WithContext(plugins.WithSessionToken(req.Context(), "does-not-exist"))
	w := httptest.NewRecorder()

	testHandlers.RegisterAnonymousChatUser(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
	if got := testHandlers.userRepository.GetUsersCount(); got != countBefore {
		t.Errorf("user count changed from %d to %d — anonymous user must not be created", countBefore, got)
	}
}

// TestRegisterAnonymousChatUserGateSessionDisabledUser verifies that a gate
// session for a disabled user returns 403 without creating a new anonymous user.
func TestRegisterAnonymousChatUserGateSessionDisabledUser(t *testing.T) {
	user, tok, err := testHandlers.userRepository.CreateAnonymousUser("DisabledGateUser")
	if err != nil {
		t.Fatalf("CreateAnonymousUser: %v", err)
	}
	if err := testHandlers.userRepository.SetEnabled(user.ID, false, "test"); err != nil {
		t.Fatalf("SetEnabled: %v", err)
	}
	countBefore := testHandlers.userRepository.GetUsersCount()

	req := httptest.NewRequest(http.MethodPost, "/api/chat/register", strings.NewReader(`{}`))
	req = req.WithContext(plugins.WithSessionToken(req.Context(), tok))
	w := httptest.NewRecorder()

	testHandlers.RegisterAnonymousChatUser(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
	if got := testHandlers.userRepository.GetUsersCount(); got != countBefore {
		t.Errorf("user count changed from %d to %d — anonymous user must not be created", countBefore, got)
	}
}
