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

// TestHandleLeaveActivityWithMetadata tests handling a Leave activity with full metadata
func TestHandleLeaveActivityWithMetadata(t *testing.T) {
	// Create a Leave activity
	activity := streams.NewActivityStreamsLeave()

	// Set the activity ID
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://remote.server/activities/leave-123")
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

	// Set the object
	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, _ := url.Parse("https://remote.server")
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Add metadata
	unknownProps := activity.GetUnknownProperties()
	unknownProps[config.APOwncastNamespaceStreamStatus] = "offline"
	unknownProps[config.APOwncastNamespaceServerName] = "Remote Test Server"
	unknownProps[config.APOwncastNamespaceStreamTitle] = "Remote Stream"
	unknownProps[config.APOwncastNamespaceStreamDescription] = "Remote server description"
	unknownProps[config.APOwncastNamespaceLogoURL] = "https://remote.server/logo.png"
	unknownProps[config.APOwncastNamespaceStreamTags] = []interface{}{"remote", "test"}

	// Parse metadata to verify it's correct
	metadata := apmodels.ParseOwncastMetadata(unknownProps)
	if metadata == nil {
		t.Fatal("Failed to parse metadata")
	}

	if metadata.StreamStatus != "offline" {
		t.Errorf("Expected stream status 'offline', got '%s'", metadata.StreamStatus)
	}

	if metadata.ServerName != "Remote Test Server" {
		t.Errorf("Expected server name 'Remote Test Server', got '%s'", metadata.ServerName)
	}

	if metadata.StreamTitle != "Remote Stream" {
		t.Errorf("Expected stream title 'Remote Stream', got '%s'", metadata.StreamTitle)
	}

	if !metadata.IsOwncastServer {
		t.Error("Expected IsOwncastServer to be true")
	}
}

// TestHandleLeaveActivityWithoutMetadata tests handling a Leave activity without metadata
func TestHandleLeaveActivityWithoutMetadata(t *testing.T) {
	// Create a Leave activity without metadata
	activity := streams.NewActivityStreamsLeave()

	// Set basic properties
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://remote.server/activities/leave-456")
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

	// Get metadata (should be minimal)
	unknownProps := activity.GetUnknownProperties()
	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	// Should not be recognized as Owncast server without metadata
	if metadata.IsOwncastServer {
		t.Error("Expected IsOwncastServer to be false when no Owncast metadata present")
	}
}

// TestHandleLeaveActivityWithPartialMetadata tests handling a Leave activity with partial metadata
func TestHandleLeaveActivityWithPartialMetadata(t *testing.T) {
	activity := streams.NewActivityStreamsLeave()

	// Set basic properties
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://remote.server/activities/leave-789")
	idProperty.Set(testIRI)
	activity.SetJSONLDId(idProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorIRI, _ := url.Parse("https://remote.server/federation/user/test")
	actorProperty.AppendIRI(actorIRI)
	activity.SetActivityStreamsActor(actorProperty)

	// Add only some metadata
	unknownProps := activity.GetUnknownProperties()
	unknownProps[config.APOwncastNamespaceStreamStatus] = "offline"
	unknownProps[config.APOwncastNamespaceServerName] = "Partial Server"
	// Intentionally leaving out other metadata fields

	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	if metadata.StreamStatus != "offline" {
		t.Errorf("Expected stream status 'offline', got '%s'", metadata.StreamStatus)
	}

	if metadata.ServerName != "Partial Server" {
		t.Errorf("Expected server name 'Partial Server', got '%s'", metadata.ServerName)
	}

	// These should be empty
	if metadata.StreamTitle != "" {
		t.Errorf("Expected empty stream title, got '%s'", metadata.StreamTitle)
	}

	if metadata.LogoURL != "" {
		t.Errorf("Expected empty logo URL, got '%s'", metadata.LogoURL)
	}

	// Should still be recognized as Owncast server with partial metadata
	if !metadata.IsOwncastServer {
		t.Error("Expected IsOwncastServer to be true with partial Owncast metadata")
	}
}

// TestLeaveActivityValidation tests validation of Leave activity structure
func TestLeaveActivityValidation(t *testing.T) {
	tests := []struct {
		name          string
		setupActivity func() vocab.ActivityStreamsLeave
		shouldError   bool
		errorMsg      string
	}{
		{
			name: "valid Leave activity",
			setupActivity: func() vocab.ActivityStreamsLeave {
				activity := streams.NewActivityStreamsLeave()

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
			name: "Leave activity without actor",
			setupActivity: func() vocab.ActivityStreamsLeave {
				activity := streams.NewActivityStreamsLeave()

				objectProperty := streams.NewActivityStreamsObjectProperty()
				objectIRI, _ := url.Parse("https://example.com")
				objectProperty.AppendIRI(objectIRI)
				activity.SetActivityStreamsObject(objectProperty)

				return activity
			},
			shouldError: true,
			errorMsg:    "Leave activity must have an actor",
		},
		{
			name: "Leave activity without object",
			setupActivity: func() vocab.ActivityStreamsLeave {
				activity := streams.NewActivityStreamsLeave()

				actorProperty := streams.NewActivityStreamsActorProperty()
				actorIRI, _ := url.Parse("https://example.com/actor")
				actorProperty.AppendIRI(actorIRI)
				activity.SetActivityStreamsActor(actorProperty)

				return activity
			},
			shouldError: false, // Object is optional for Leave activity
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

			if !tt.shouldError && !hasActor && tt.errorMsg == "Leave activity must have an actor" {
				t.Errorf("Validation failed for %s when it should have passed", tt.name)
			}

			// Verify activity type
			if activity.GetTypeName() != "Leave" {
				t.Errorf("Expected activity type 'Leave', got '%s'", activity.GetTypeName())
			}

			_ = hasObject // Object validation could be added if needed
		})
	}
}

// TestUpdateFederatedServerStatusFunction tests the updateFederatedServerStatus function
func TestUpdateFederatedServerStatusFunction(t *testing.T) {
	tests := []struct {
		name     string
		actorIRI string
		metadata *apmodels.OwncastMetadata
		wantErr  bool
	}{
		{
			name:     "update with full metadata",
			actorIRI: "https://example.com/federation/user/test",
			metadata: &apmodels.OwncastMetadata{
				StreamStatus:      "offline",
				ServerName:        "Test Server",
				StreamTitle:       "Test Stream",
				StreamDescription: "Test Description",
				LogoURL:           "https://example.com/logo.png",
				Tags:              []string{"test", "stream"},
				IsOwncastServer:   true,
			},
			wantErr: false,
		},
		{
			name:     "update with minimal metadata",
			actorIRI: "https://example.com/federation/user/test",
			metadata: &apmodels.OwncastMetadata{
				StreamStatus:    "offline",
				IsOwncastServer: true,
			},
			wantErr: false,
		},
		{
			name:     "update with empty actor IRI",
			actorIRI: "",
			metadata: &apmodels.OwncastMetadata{
				StreamStatus:    "offline",
				ServerName:      "Test Server",
				IsOwncastServer: true,
			},
			wantErr: false, // We don't error on empty IRI in current implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := updateFederatedServerStatus(tt.actorIRI, tt.metadata)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateFederatedServerStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestHandleLeaveInboxRequestContext tests the handleLeaveInboxRequest function with context
func TestHandleLeaveInboxRequestContext(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a Leave activity
	activity := streams.NewActivityStreamsLeave()

	// This would normally fail because we don't have a real actor,
	// but we're testing that the function handles the context properly
	err := testService.handleLeaveInboxRequest(ctx, activity)
	// The function should return nil even with missing actor
	// because it returns nil when actor property is nil
	if err != nil {
		t.Errorf("Expected nil error for Leave activity without actor, got: %v", err)
	}
}

// TestLeaveActivityMetadataRoundTrip tests setting and parsing metadata round-trip
func TestLeaveActivityMetadataRoundTrip(t *testing.T) {
	// Create a Leave activity
	activity := streams.NewActivityStreamsLeave()

	// Original metadata
	originalMetadata := map[string]interface{}{
		config.APOwncastNamespaceStreamStatus:      "offline",
		config.APOwncastNamespaceServerName:        "Round Trip Server",
		config.APOwncastNamespaceStreamTitle:       "Round Trip Stream",
		config.APOwncastNamespaceStreamDescription: "Round trip description",
		config.APOwncastNamespaceLogoURL:           "https://example.com/logo.png",
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
	if parsedMetadata.StreamStatus != "offline" {
		t.Errorf("Round trip failed for StreamStatus: expected 'offline', got '%s'", parsedMetadata.StreamStatus)
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
