package utils

import "testing"

func TestUserAgent(t *testing.T) {
	testAgents := []string{
		"Pleroma 1.0.0-1168-ge18c7866-pleroma-dot-site; https://pleroma.site info@pleroma.site",
		"Mastodon 1.2.3 Bot",
	}

	for _, agent := range testAgents {
		if !IsUserAgentABot(agent) {
			t.Error("Incorrect parsing of useragent", agent)
		}
	}
}
