package chat

import "testing"

func TestForbiddenUsernameMatchesReservedNameInsideCandidate(t *testing.T) {
	for _, name := range []string{"owncast", "OwnCast_Admin", "SYSTEM-status"} {
		if blocked, ok := ForbiddenUsername(name, []string{"owncast", "admin", "system"}); !ok || blocked == "" {
			t.Errorf("ForbiddenUsername(%q) = %q, %t; want a blocklist match", name, blocked, ok)
		}
	}

	for _, name := range []string{"own", "moderator", "viewer"} {
		if blocked, ok := ForbiddenUsername(name, []string{"owncast", "admin", "system"}); ok {
			t.Errorf("ForbiddenUsername(%q) = %q, true; want allowed", name, blocked)
		}
	}
}
