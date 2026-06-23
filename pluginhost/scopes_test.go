package pluginhost

import (
	"testing"

	"github.com/owncast/owncast/models"
)

// F1 regression: a viewer-auth plugin may grant MODERATOR (the documented
// scope) and the chat-send scopes via users.register, but never
// HAS_ADMIN_ACCESS — and an unknown scope is rejected rather than silently
// dropped or applied.
func TestValidatePluginGrantedScopes(t *testing.T) {
	cases := []struct {
		name    string
		scopes  []string
		wantErr bool
	}{
		{"none", nil, false},
		{"moderator", []string{models.ModeratorScopeKey}, false},
		{"chat-send scopes", []string{models.ScopeCanSendChatMessages, models.ScopeCanSendSystemMessages}, false},
		{"moderator + chat", []string{models.ModeratorScopeKey, models.ScopeCanSendChatMessages}, false},
		{"admin access rejected", []string{models.ScopeHasAdminAccess}, true},
		{"admin mixed in is rejected", []string{models.ModeratorScopeKey, models.ScopeHasAdminAccess}, true},
		{"unknown scope rejected", []string{"BOGUS"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePluginGrantedScopes(tc.scopes)
			if tc.wantErr && err == nil {
				t.Fatalf("scopes %v should be rejected", tc.scopes)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("scopes %v should be allowed, got: %v", tc.scopes, err)
			}
		})
	}
}
