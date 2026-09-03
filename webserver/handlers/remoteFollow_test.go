package handlers

import "testing"

func TestRemoteFollowTemplateIgnoresNonStringTemplate(t *testing.T) {
	links := []map[string]interface{}{
		{
			"rel":      "http://ostatus.org/schema/1.0/subscribe",
			"template": 42,
		},
		{
			"rel":      "http://ostatus.org/schema/1.0/subscribe",
			"template": "https://remote.example/follow?uri={uri}",
		},
	}

	if got := remoteFollowTemplate(links); got != "https://remote.example/follow?uri={uri}" {
		t.Fatalf("remoteFollowTemplate() = %q", got)
	}
}
