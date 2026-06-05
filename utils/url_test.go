package utils

import "testing"

func TestCanonicalizeURLHostname(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unicode hostname",
			input:    "https://live.retrospection.みんな",
			expected: "https://live.retrospection.xn--q9jyb4c",
		},
		{
			name:     "punycode hostname",
			input:    "https://live.retrospection.xn--q9jyb4c",
			expected: "https://live.retrospection.xn--q9jyb4c",
		},
		{
			name:     "port",
			input:    "https://live.retrospection.みんな:8443",
			expected: "https://live.retrospection.xn--q9jyb4c:8443",
		},
		{
			name:     "path and query",
			input:    "https://live.retrospection.みんな/federation/user/retrots3m?page=1",
			expected: "https://live.retrospection.xn--q9jyb4c/federation/user/retrots3m?page=1",
		},
		{
			name:     "localhost",
			input:    "http://localhost:8080",
			expected: "http://localhost:8080",
		},
		{
			name:     "ipv4",
			input:    "http://203.0.113.10:8080",
			expected: "http://203.0.113.10:8080",
		},
		{
			name:     "ipv6",
			input:    "http://[2001:db8::1]:8080",
			expected: "http://[2001:db8::1]:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CanonicalizeURLHostname(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			if result != tt.expected {
				t.Errorf("CanonicalizeURLHostname() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCanonicalizeURLHostnameWithInvalidURL(t *testing.T) {
	if _, err := CanonicalizeURLHostname("https://"); err == nil {
		t.Error("CanonicalizeURLHostname() should return an error for a URL without a hostname")
	}
}

func TestCanonicalizeHost(t *testing.T) {
	result, err := CanonicalizeHost("live.retrospection.みんな:8443")
	if err != nil {
		t.Fatal(err)
	}

	if result != "live.retrospection.xn--q9jyb4c:8443" {
		t.Errorf("CanonicalizeHost() = %v, want %v", result, "live.retrospection.xn--q9jyb4c:8443")
	}
}

func TestCanonicalizeHostRejectsURLParts(t *testing.T) {
	tests := []string{
		"",
		"live.retrospection.みんな/path",
		"live.retrospection.みんな?query=true",
		"live.retrospection.みんな#fragment",
		"user@live.retrospection.みんな",
		"https://live.retrospection.みんな",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if _, err := CanonicalizeHost(tt); err == nil {
				t.Errorf("CanonicalizeHost(%q) should return an error", tt)
			}
		})
	}
}
