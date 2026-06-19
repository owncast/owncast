package outbox

import (
	"net/url"
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/services/activitypub/apmodels"
)

// TestLeaveActivityStructure tests that the Leave activity is created with the correct structure
func TestLeaveActivityStructure(t *testing.T) {
	// Create a Leave activity manually to test structure
	activity := streams.NewActivityStreamsLeave()

	// Set the activity ID
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://example.com/activities/leave-123")
	idProperty.SetIRI(testIRI)
	activity.SetJSONLDId(idProperty)

	// Set the actor
	actorProperty := streams.NewActivityStreamsActorProperty()
	actorIRI, _ := url.Parse("https://example.com/federation/user/test")
	actorProperty.AppendIRI(actorIRI)
	activity.SetActivityStreamsActor(actorProperty)

	// Set the object
	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, _ := url.Parse("https://example.com")
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Verify the activity type
	if activity.GetTypeName() != "Leave" {
		t.Errorf("Expected activity type 'Leave', got '%s'", activity.GetTypeName())
	}

	// Verify the actor
	if actorProperty.Len() != 1 {
		t.Errorf("Expected 1 actor, got %d", actorProperty.Len())
	}

	// Verify the object
	if objectProperty.Len() != 1 {
		t.Errorf("Expected 1 object, got %d", objectProperty.Len())
	}

	// Serialize and check that it doesn't error
	_, err := apmodels.Serialize(activity)
	if err != nil {
		t.Errorf("Failed to serialize Leave activity: %v", err)
	}
}

// TestLeaveActivityWithMetadata tests that metadata is properly added to Leave activity
func TestLeaveActivityWithMetadata(t *testing.T) {
	activity := streams.NewActivityStreamsLeave()

	// Get unknown properties and add metadata
	unknownProps := activity.GetUnknownProperties()
	unknownProps[config.APOwncastNamespaceStreamStatus] = "offline"
	unknownProps[config.APOwncastNamespaceServerName] = "Test Server"
	unknownProps[config.APOwncastNamespaceStreamTitle] = "Test Stream"
	unknownProps[config.APOwncastNamespaceStreamDescription] = "Test Description"
	unknownProps[config.APOwncastNamespaceLogoURL] = "https://example.com/logo.png"
	unknownProps[config.APOwncastNamespaceStreamTags] = []interface{}{"gaming", "tech"}

	// Parse the metadata back
	metadata := apmodels.ParseOwncastMetadata(unknownProps)

	// Verify metadata was properly set
	if metadata.StreamStatus != "offline" {
		t.Errorf("Expected stream status 'offline', got '%s'", metadata.StreamStatus)
	}

	if metadata.ServerName != "Test Server" {
		t.Errorf("Expected server name 'Test Server', got '%s'", metadata.ServerName)
	}

	if metadata.StreamTitle != "Test Stream" {
		t.Errorf("Expected stream title 'Test Stream', got '%s'", metadata.StreamTitle)
	}

	if metadata.StreamDescription != "Test Description" {
		t.Errorf("Expected stream description 'Test Description', got '%s'", metadata.StreamDescription)
	}

	if metadata.LogoURL != "https://example.com/logo.png" {
		t.Errorf("Expected logo URL 'https://example.com/logo.png', got '%s'", metadata.LogoURL)
	}

	if len(metadata.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(metadata.Tags))
	} else {
		if metadata.Tags[0] != "gaming" || metadata.Tags[1] != "tech" {
			t.Errorf("Expected tags ['gaming', 'tech'], got %v", metadata.Tags)
		}
	}

	// A Leave carries stream metadata but is not a directory marker.
	if metadata.IsDirectory {
		t.Error("Expected IsDirectory to be false for a stream Leave")
	}
}

// TestLeaveActivityAddressing tests that the Leave activity has proper addressing
func TestLeaveActivityAddressing(t *testing.T) {
	activity := streams.NewActivityStreamsLeave()

	// Create addressing properties
	toProperty := streams.NewActivityStreamsToProperty()
	publicIRI, _ := url.Parse("https://www.w3.org/ns/activitystreams#Public")
	toProperty.AppendIRI(publicIRI)
	activity.SetActivityStreamsTo(toProperty)

	ccProperty := streams.NewActivityStreamsCcProperty()
	followersIRI, _ := url.Parse("https://example.com/federation/user/test/followers")
	ccProperty.AppendIRI(followersIRI)
	activity.SetActivityStreamsCc(ccProperty)

	// Verify addressing
	to := activity.GetActivityStreamsTo()
	if to == nil || to.Len() != 1 {
		t.Error("Expected 'to' property to have 1 recipient")
	}

	cc := activity.GetActivityStreamsCc()
	if cc == nil || cc.Len() != 1 {
		t.Error("Expected 'cc' property to have 1 recipient")
	}
}

// TestParseLeaveActivity tests parsing an incoming Leave activity
func TestParseLeaveActivity(t *testing.T) {
	// Create a Leave activity
	activity := streams.NewActivityStreamsLeave()

	// Set basic properties
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://remote.server/activities/leave-456")
	idProperty.SetIRI(testIRI)
	activity.SetJSONLDId(idProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorIRI, _ := url.Parse("https://remote.server/federation/user/remote")
	actorProperty.AppendIRI(actorIRI)
	activity.SetActivityStreamsActor(actorProperty)

	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, _ := url.Parse("https://remote.server")
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Add metadata
	unknownProps := activity.GetUnknownProperties()
	unknownProps[config.APOwncastNamespaceStreamStatus] = "offline"
	unknownProps[config.APOwncastNamespaceServerName] = "Remote Server"

	// Verify we can access all properties
	if activity.GetActivityStreamsActor() == nil {
		t.Error("Failed to get actor from Leave activity")
	}

	if activity.GetActivityStreamsObject() == nil {
		t.Error("Failed to get object from Leave activity")
	}

	metadata := apmodels.ParseOwncastMetadata(unknownProps)
	if metadata.StreamStatus != "offline" {
		t.Errorf("Failed to parse stream status from Leave activity, got '%s'", metadata.StreamStatus)
	}

	if metadata.ServerName != "Remote Server" {
		t.Errorf("Failed to parse server name from Leave activity, got '%s'", metadata.ServerName)
	}
}

// TestLeaveActivitySerialization tests that Leave activities can be serialized/deserialized
func TestLeaveActivitySerialization(t *testing.T) {
	// Create a Leave activity
	activity := streams.NewActivityStreamsLeave()

	// Set properties
	idProperty := streams.NewJSONLDIdProperty()
	testIRI, _ := url.Parse("https://example.com/activities/leave-789")
	idProperty.SetIRI(testIRI)
	activity.SetJSONLDId(idProperty)

	actorProperty := streams.NewActivityStreamsActorProperty()
	actorIRI, _ := url.Parse("https://example.com/federation/user/test")
	actorProperty.AppendIRI(actorIRI)
	activity.SetActivityStreamsActor(actorProperty)

	objectProperty := streams.NewActivityStreamsObjectProperty()
	serverIRI, _ := url.Parse("https://example.com")
	objectProperty.AppendIRI(serverIRI)
	activity.SetActivityStreamsObject(objectProperty)

	// Add metadata
	unknownProps := activity.GetUnknownProperties()
	unknownProps[config.APOwncastNamespaceStreamStatus] = "offline"
	unknownProps[config.APOwncastNamespaceServerName] = "Test Server"
	unknownProps[config.APOwncastNamespaceStreamTitle] = "Test Stream"

	// Serialize
	data, err := apmodels.Serialize(activity)
	if err != nil {
		t.Fatalf("Failed to serialize Leave activity: %v", err)
	}

	// Verify the serialized data contains expected elements
	serialized := string(data)
	if !contains(serialized, "\"type\":\"Leave\"") {
		t.Error("Serialized Leave activity doesn't contain expected type")
	}

	if !contains(serialized, "https://example.com/activities/leave-789") {
		t.Error("Serialized Leave activity doesn't contain expected ID")
	}

	if !contains(serialized, "https://example.com/federation/user/test") {
		t.Error("Serialized Leave activity doesn't contain expected actor")
	}

	if !contains(serialized, config.APOwncastNamespaceStreamStatus) {
		t.Error("Serialized Leave activity doesn't contain stream status metadata")
	}
}

// TestLeaveActivityType tests that the activity is recognized as a Leave type
func TestLeaveActivityType(t *testing.T) {
	activity := streams.NewActivityStreamsLeave()

	// Check if it's recognized as a Leave activity
	if _, ok := activity.(vocab.ActivityStreamsLeave); !ok {
		t.Error("Activity is not recognized as ActivityStreamsLeave")
	}

	// Check type name
	if activity.GetTypeName() != "Leave" {
		t.Errorf("Expected type name 'Leave', got '%s'", activity.GetTypeName())
	}
}

// Helper function to check if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(str) > 0 && len(substr) == 0 || (len(substr) > 0 && str != "" && str[0:] != "" && (str == substr || (len(str) > len(substr) && str[:len(substr)] == substr) || (len(str) > len(substr) && str[len(str)-len(substr):] == substr) || (len(str) > len(substr) && findSubstring(str, substr)))))
}

func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
