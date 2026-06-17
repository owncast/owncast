package models

import (
	"encoding/json"
	"testing"
	"time"
)

// TestFederatedServerJSONContract pins the JSON field names the API emits for
// a federated server. The web consumes these names directly (see
// web/hooks/useFederatedServers.tsx and web/tests/useFederatedServers.test.tsx),
// so a rename here without a matching web change silently breaks the
// featured-streams admin table and public listing. This test fails loudly if
// the contract drifts -- including reverting to the legacy url/logo/thumbnail/
// lastChecked names that previously broke the feature.
func TestFederatedServerJSONContract(t *testing.T) {
	name := "goodnight"
	displayName := "Goodnight TV"
	logoURL := "https://goodnight.example.com/logo.png"
	streamTitle := "Late night coding"
	streamDescription := "Chill coding vibes"
	thumbnailURL := "https://goodnight.example.com/thumb.jpg"
	summary := "Goodnight TV"
	username := "goodnight"
	now := time.Date(2026, 6, 16, 21, 55, 0, 0, time.UTC)

	server := FederatedServer{
		ID:                7,
		IRI:               "https://goodnight.example.com",
		Name:              &name,
		DisplayName:       &displayName,
		Summary:           &summary,
		Username:          &username,
		LogoURL:           &logoURL,
		IsOnline:          true,
		StreamTitle:       &streamTitle,
		StreamDescription: &streamDescription,
		Tags:              []string{"coding", "chill"},
		ThumbnailURL:      &thumbnailURL,
		LastSeenOnline:    &now,
		LastStatusUpdate:  &now,
		AddedAt:           now,
		FollowedAt:        &now,
		Pending:           false,
		FollowStatus:      "accepted",
	}

	data, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("failed to marshal FederatedServer: %v", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatalf("failed to unmarshal FederatedServer JSON: %v", err)
	}

	// The web reads each of these keys by name.
	required := []string{
		"id",
		"iri",
		"name",
		"displayName",
		"logoUrl",
		"isOnline",
		"streamTitle",
		"streamDescription",
		"tags",
		"thumbnailUrl",
		"lastStatusUpdate",
		"addedAt",
		"followStatus",
	}
	for _, key := range required {
		if _, ok := fields[key]; !ok {
			t.Errorf("FederatedServer JSON is missing required field %q; the web depends on it", key)
		}
	}

	// These are the legacy names the web mistakenly used before the contract
	// was unified. They must never reappear -- their presence would mean the
	// two sides have diverged again.
	forbidden := []string{"url", "logo", "thumbnail", "lastChecked"}
	for _, key := range forbidden {
		if _, ok := fields[key]; ok {
			t.Errorf("FederatedServer JSON contains unexpected legacy field %q", key)
		}
	}
}
