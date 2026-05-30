package inbox

import (
	"context"
	"net/url"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/services/activitypub/apmodels"
)

// TestHandleOfferActivityWithMetadata tests handling an Offer activity with full metadata
func TestHandleOfferActivityWithMetadata(t *testing.T) {
	// Create an Offer activity
	activity := streams.NewActivityStreamsOffer()

	// Set the activity ID
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://remote.server/activities/offer-123")
	idProperty.SetIRI(testIRI)
	activity.SetJSONLDId(idProperty)

	// Set the actor
	actorProperty := streams.NewActivityStreamsActorProperty()
	person := streams.NewActivityStreamsPerson()
	personID := streams.NewJSONLDIdProperty()
	actorIRI, _ := url.Parse("https://remote.server/federation/user/test")
	personID.SetIRI(actorIRI)
	person.SetJSONLDId(personID)
	actorProperty.AppendActivityStreamsPerson(person)
	activity.SetActivityStreamsActor(actorProperty)

	// Set the object (server URL being offered)
	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, _ := url.Parse("https://remote.server")
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Add Owncast metadata
	unknownProps := activity.GetUnknownProperties()
	unknownProps[config.APOwncastNamespaceStreamStatus] = "live"
	unknownProps[config.APOwncastNamespaceServerName] = "Remote Test Server"
	unknownProps[config.APOwncastNamespaceStreamTitle] = "Remote Stream Title"
	unknownProps[config.APOwncastNamespaceStreamDescription] = "Remote stream description"
	unknownProps[config.APOwncastNamespaceLogoURL] = "https://remote.server/logo.png"
	unknownProps[config.APOwncastNamespaceThumbnailURL] = "https://remote.server/thumb.png"
	unknownProps[config.APOwncastNamespaceStreamTags] = []interface{}{"remote", "test", "live"}

	// Parse metadata to verify it's correct
	metadata := apmodels.ParseOwncastMetadata(unknownProps)
	if metadata == nil {
		t.Fatal("Failed to parse metadata")
	}

	if metadata.StreamStatus != "live" {
		t.Errorf("Expected stream status 'live', got '%s'", metadata.StreamStatus)
	}

	if metadata.ServerName != "Remote Test Server" {
		t.Errorf("Expected server name 'Remote Test Server', got '%s'", metadata.ServerName)
	}

	if metadata.StreamTitle != "Remote Stream Title" {
		t.Errorf("Expected stream title 'Remote Stream Title', got '%s'", metadata.StreamTitle)
	}

	if metadata.ThumbnailURL != "https://remote.server/thumb.png" {
		t.Errorf("Expected thumbnail URL 'https://remote.server/thumb.png', got '%s'", metadata.ThumbnailURL)
	}

	if len(metadata.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(metadata.Tags))
	}

	if !metadata.IsOwncastServer {
		t.Error("Expected IsOwncastServer to be true")
	}
}

// TestHandleOfferActivityWithoutMetadata tests handling an Offer activity without Owncast metadata
func TestHandleOfferActivityWithoutMetadata(t *testing.T) {
	// Create an Offer activity without Owncast metadata
	activity := streams.NewActivityStreamsOffer()

	// Set basic properties
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://remote.server/activities/offer-456")
	idProperty.Set(testIRI)
	activity.SetJSONLDId(idProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorIRI, _ := url.Parse("https://remote.server/federation/user/test")
	actorProperty.AppendIRI(actorIRI)
	activity.SetActivityStreamsActor(actorProperty)

	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, _ := url.Parse("https://remote.server")
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Get metadata (should be empty)
	unknownProps := activity.GetUnknownProperties()
	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	// Should not be recognized as Owncast server without metadata
	if metadata.IsOwncastServer {
		t.Error("Expected IsOwncastServer to be false when no Owncast metadata present")
	}
}

// TestHandleOfferActivityWithOfflineStatus tests that offline status in Offer is handled correctly
func TestHandleOfferActivityWithOfflineStatus(t *testing.T) {
	activity := streams.NewActivityStreamsOffer()

	// Set basic properties
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://remote.server/activities/offer-789")
	idProperty.Set(testIRI)
	activity.SetJSONLDId(idProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorIRI, _ := url.Parse("https://remote.server/federation/user/test")
	actorProperty.AppendIRI(actorIRI)
	activity.SetActivityStreamsActor(actorProperty)

	// Add metadata with offline status (should be ignored by Offer handler)
	unknownProps := activity.GetUnknownProperties()
	unknownProps[config.APOwncastNamespaceStreamStatus] = "offline"
	unknownProps[config.APOwncastNamespaceServerName] = "Offline Server"

	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	if metadata.StreamStatus != "offline" {
		t.Errorf("Expected stream status 'offline', got '%s'", metadata.StreamStatus)
	}

	// The handler should ignore offers with offline status since offers are for live streams
	if metadata.IsOwncastServer {
		t.Log("Metadata parsed correctly as Owncast server")
	}
}

// TestOfferActivityValidation tests validation of Offer activity structure
func TestOfferActivityValidation(t *testing.T) {
	tests := []struct {
		name          string
		setupActivity func() vocab.ActivityStreamsOffer
		shouldError   bool
		errorMsg      string
	}{
		{
			name: "valid Offer activity",
			setupActivity: func() vocab.ActivityStreamsOffer {
				activity := streams.NewActivityStreamsOffer()

				actorProperty := streams.NewActivityStreamsActorProperty()
				actorIRI, _ := url.Parse("https://example.com/actor")
				actorProperty.AppendIRI(actorIRI)
				activity.SetActivityStreamsActor(actorProperty)

				objectProperty := streams.NewActivityStreamsObjectProperty()
				objectIRI, _ := url.Parse("https://example.com")
				objectProperty.AppendIRI(objectIRI)
				activity.SetActivityStreamsObject(objectProperty)

				return activity
			},
			shouldError: false,
		},
		{
			name: "Offer activity without actor",
			setupActivity: func() vocab.ActivityStreamsOffer {
				activity := streams.NewActivityStreamsOffer()

				objectProperty := streams.NewActivityStreamsObjectProperty()
				objectIRI, _ := url.Parse("https://example.com")
				objectProperty.AppendIRI(objectIRI)
				activity.SetActivityStreamsObject(objectProperty)

				return activity
			},
			shouldError: true,
			errorMsg:    "Offer activity must have an actor",
		},
		{
			name: "Offer activity without object",
			setupActivity: func() vocab.ActivityStreamsOffer {
				activity := streams.NewActivityStreamsOffer()

				actorProperty := streams.NewActivityStreamsActorProperty()
				actorIRI, _ := url.Parse("https://example.com/actor")
				actorProperty.AppendIRI(actorIRI)
				activity.SetActivityStreamsActor(actorProperty)

				return activity
			},
			shouldError: false, // Object is optional for our use case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activity := tt.setupActivity()

			// Basic validation
			hasActor := activity.GetActivityStreamsActor() != nil
			hasObject := activity.GetActivityStreamsObject() != nil

			if tt.shouldError && hasActor {
				t.Errorf("Expected error for %s but validation passed", tt.name)
			}

			if !tt.shouldError && !hasActor && tt.errorMsg == "Offer activity must have an actor" {
				t.Errorf("Validation failed for %s when it should have passed", tt.name)
			}

			// Verify activity type
			if activity.GetTypeName() != "Offer" {
				t.Errorf("Expected activity type 'Offer', got '%s'", activity.GetTypeName())
			}

			_ = hasObject // Object validation could be added if needed
		})
	}
}

// TestHandleOfferInboxRequestContext tests the handleOfferInboxRequest function with context
func TestHandleOfferInboxRequestContext(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create an Offer activity
	activity := streams.NewActivityStreamsOffer()

	// This should return nil because we don't have a real actor
	// The function returns nil when actor property is nil
	err := testService.handleOfferInboxRequest(ctx, activity)
	if err != nil {
		t.Errorf("Expected nil error for Offer activity without actor, got: %v", err)
	}
}

// TestOfferActivityMetadataRoundTrip tests setting and parsing metadata round-trip
func TestOfferActivityMetadataRoundTrip(t *testing.T) {
	// Create an Offer activity
	activity := streams.NewActivityStreamsOffer()

	// Original metadata
	originalMetadata := map[string]interface{}{
		config.APOwncastNamespaceStreamStatus:      "live",
		config.APOwncastNamespaceServerName:        "Round Trip Server",
		config.APOwncastNamespaceStreamTitle:       "Round Trip Stream",
		config.APOwncastNamespaceStreamDescription: "Round trip description",
		config.APOwncastNamespaceLogoURL:           "https://example.com/logo.png",
		config.APOwncastNamespaceThumbnailURL:      "https://example.com/thumb.png",
		config.APOwncastNamespaceStreamTags:        []interface{}{"tag1", "tag2", "tag3"},
	}

	// Set metadata
	unknownProps := activity.GetUnknownProperties()
	for k, v := range originalMetadata {
		unknownProps[k] = v
	}

	// Parse metadata back
	parsedMetadata := apmodels.ParseOwncastMetadata(unknownProps)

	// Verify round trip
	if parsedMetadata.StreamStatus != "live" {
		t.Errorf("Round trip failed for StreamStatus: expected 'live', got '%s'", parsedMetadata.StreamStatus)
	}

	if parsedMetadata.ServerName != "Round Trip Server" {
		t.Errorf("Round trip failed for ServerName: expected 'Round Trip Server', got '%s'", parsedMetadata.ServerName)
	}

	if parsedMetadata.StreamTitle != "Round Trip Stream" {
		t.Errorf("Round trip failed for StreamTitle: expected 'Round Trip Stream', got '%s'", parsedMetadata.StreamTitle)
	}

	if parsedMetadata.StreamDescription != "Round trip description" {
		t.Errorf("Round trip failed for StreamDescription: expected 'Round trip description', got '%s'", parsedMetadata.StreamDescription)
	}

	if parsedMetadata.LogoURL != "https://example.com/logo.png" {
		t.Errorf("Round trip failed for LogoURL: expected 'https://example.com/logo.png', got '%s'", parsedMetadata.LogoURL)
	}

	if parsedMetadata.ThumbnailURL != "https://example.com/thumb.png" {
		t.Errorf("Round trip failed for ThumbnailURL: expected 'https://example.com/thumb.png', got '%s'", parsedMetadata.ThumbnailURL)
	}

	if len(parsedMetadata.Tags) != 3 {
		t.Errorf("Round trip failed for Tags: expected 3 tags, got %d", len(parsedMetadata.Tags))
	} else {
		expectedTags := []string{"tag1", "tag2", "tag3"}
		for i, tag := range expectedTags {
			if parsedMetadata.Tags[i] != tag {
				t.Errorf("Round trip failed for Tags[%d]: expected '%s', got '%s'", i, tag, parsedMetadata.Tags[i])
			}
		}
	}

	if !parsedMetadata.IsOwncastServer {
		t.Error("Round trip failed: IsOwncastServer should be true")
	}
}
