package resolvers

import (
	"net/url"
	"testing"

	"github.com/go-fed/activity/streams"
)

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

// TestGetResolvedActorFromActorPropertyWithNil verifies that the function
// doesn't panic when given a nil actor property.
func TestGetResolvedActorFromActorPropertyWithNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetResolvedActorFromActorProperty panicked with nil: %v", r)
		}
	}()

	_, err := GetResolvedActorFromActorProperty(nil)
	if err == nil {
		t.Error("GetResolvedActorFromActorProperty(nil) should return error")
	}
}

// TestGetResolvedActorFromActorPropertyWithEmpty verifies that the function
// doesn't panic when given an empty actor property.
func TestGetResolvedActorFromActorPropertyWithEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetResolvedActorFromActorProperty panicked with empty: %v", r)
		}
	}()

	actor := streams.NewActivityStreamsActorProperty()
	_, err := GetResolvedActorFromActorProperty(actor)
	if err == nil {
		t.Error("GetResolvedActorFromActorProperty with empty actor should return error")
	}
}

// TestGetResolvedActorFromActorPropertyWithIRI verifies that the function
// handles the case where an IRI needs to be resolved.
// Note: This test involves network calls and config repository access which may
// panic in a test environment. This is expected since we're testing nil safety
// of ActivityPub property handling, not network behavior.
func TestGetResolvedActorFromActorPropertyWithIRI(t *testing.T) {
	t.Skip("Skipping IRI resolution test - requires network access and proper config setup")
}

// TestGetResolvedActorFromActorPropertyWithPersonButMissingFields verifies that
// the function handles Person objects that are missing required fields.
func TestGetResolvedActorFromActorPropertyWithPersonMissingFields(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetResolvedActorFromActorProperty panicked: %v", r)
		}
	}()

	actor := streams.NewActivityStreamsActorProperty()
	person := streams.NewActivityStreamsPerson()
	// Person has no ID, inbox, username, or public key
	actor.AppendActivityStreamsPerson(person)

	_, err := GetResolvedActorFromActorProperty(actor)
	if err == nil {
		t.Error("GetResolvedActorFromActorProperty with incomplete Person should return error")
	}
}

// TestGetResolvedActorFromActorPropertyWithServiceMissingFields verifies that
// the function handles Service objects that are missing required fields.
func TestGetResolvedActorFromActorPropertyWithServiceMissingFields(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetResolvedActorFromActorProperty panicked: %v", r)
		}
	}()

	actor := streams.NewActivityStreamsActorProperty()
	service := streams.NewActivityStreamsService()
	// Service has no ID, inbox, username, or public key
	actor.AppendActivityStreamsService(service)

	_, err := GetResolvedActorFromActorProperty(actor)
	if err == nil {
		t.Error("GetResolvedActorFromActorProperty with incomplete Service should return error")
	}
}

// TestGetResolvedActorFromActorPropertyWithApplicationMissingFields verifies that
// the function handles Application objects that are missing required fields.
func TestGetResolvedActorFromActorPropertyWithApplicationMissingFields(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetResolvedActorFromActorProperty panicked: %v", r)
		}
	}()

	actor := streams.NewActivityStreamsActorProperty()
	app := streams.NewActivityStreamsApplication()
	// Application has no ID, inbox, username, or public key
	actor.AppendActivityStreamsApplication(app)

	_, err := GetResolvedActorFromActorProperty(actor)
	if err == nil {
		t.Error("GetResolvedActorFromActorProperty with incomplete Application should return error")
	}
}

// TestNilSafetyNoPanic is a comprehensive test that ensures none of the
// resolver functions panic with nil/empty inputs.
func TestNilSafetyNoPanic(t *testing.T) {
	t.Run("GetResolvedActorFromActorProperty with nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panicked: %v", r)
			}
		}()
		_, _ = GetResolvedActorFromActorProperty(nil)
	})

	t.Run("GetResolvedActorFromActorProperty with empty", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panicked: %v", r)
			}
		}()
		_, _ = GetResolvedActorFromActorProperty(streams.NewActivityStreamsActorProperty())
	})

	t.Run("GetResolvedActorFromActorProperty with Person no fields", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panicked: %v", r)
			}
		}()
		actor := streams.NewActivityStreamsActorProperty()
		actor.AppendActivityStreamsPerson(streams.NewActivityStreamsPerson())
		_, _ = GetResolvedActorFromActorProperty(actor)
	})

	t.Run("GetResolvedActorFromActorProperty with Service no fields", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panicked: %v", r)
			}
		}()
		actor := streams.NewActivityStreamsActorProperty()
		actor.AppendActivityStreamsService(streams.NewActivityStreamsService())
		_, _ = GetResolvedActorFromActorProperty(actor)
	})

	t.Run("GetResolvedActorFromActorProperty with Application no fields", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panicked: %v", r)
			}
		}()
		actor := streams.NewActivityStreamsActorProperty()
		actor.AppendActivityStreamsApplication(streams.NewActivityStreamsApplication())
		_, _ = GetResolvedActorFromActorProperty(actor)
	})
}
