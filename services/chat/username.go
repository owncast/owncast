package chat

import "strings"

// ForbiddenUsername reports whether a proposed display name contains one of
// the configured reserved names, case-insensitively.
func ForbiddenUsername(proposed string, blocklist []string) (string, bool) {
	proposed = strings.ToLower(strings.TrimSpace(proposed))
	for _, blocked := range blocklist {
		blocked = strings.ToLower(strings.TrimSpace(blocked))
		if blocked != "" && strings.Contains(proposed, blocked) {
			return blocked, true
		}
	}
	return "", false
}
