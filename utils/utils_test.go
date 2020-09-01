package utils

import "testing"

func TestUserAgent(t *testing.T) {
	testAgents := []string{
		"Pleroma 1.0.0-1168-ge18c7866-pleroma-dot-site; https://pleroma.site info@pleroma.site",
		"Mastodon 1.2.3 Bot",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)",
	}

	for _, agent := range testAgents {
		if !IsUserAgentABot(agent) {
			t.Error("Incorrect parsing of useragent", agent)
		}
	}
}
