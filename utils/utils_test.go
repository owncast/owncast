package utils

import (
	"testing"
)

func TestUserAgent(t *testing.T) {
	testAgents := []string{
		"Pleroma 1.0.0-1168-ge18c7866-pleroma-dot-site; https://pleroma.site info@pleroma.site",
		"Mastodon 1.2.3 Bot",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)",
		"WhatsApp",
	}

	for _, agent := range testAgents {
		if !IsUserAgentABot(agent) {
			t.Error("Incorrect parsing of useragent", agent)
		}
	}
}

func TestGetHashtagsFromText(t *testing.T) {
	text := `Some text with a #hashtag goes here.\n\n
	Another #secondhashtag, goes here.\n\n
	#thirdhashtag`

	hashtags := GetHashtagsFromText(text)

	if hashtags[0] != "#hashtag" || hashtags[1] != "#secondhashtag" || hashtags[2] != "#thirdhashtag" {
		t.Error("Incorrect hashtags fetched from text.")
	}
}

func TestPercentageUtilsTest(t *testing.T) {
	total := 42
	number := 18

	percent := IntPercentage(number, total)

	if percent != 42 {
		t.Error("Incorrect percentage calculation.")
	}
}
