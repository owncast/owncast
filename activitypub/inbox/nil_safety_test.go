package inbox

import (
	"context"
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

// These tests verify that handler functions don't panic when given
// ActivityPub activities with nil or missing properties that could
// cause nil pointer dereferences.

func TestHandleFollowWithNilObject(t *testing.T) {
	activity := streams.NewActivityStreamsFollow()
	// Don't set object or actor - they will be nil

	// This should return an error, not panic
	err := handleFollowInboxRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleFollowInboxRequest with nil object should return error")
	}
}

func TestHandleFollowWithEmptyObject(t *testing.T) {
	activity := streams.NewActivityStreamsFollow()
	object := streams.NewActivityStreamsObjectProperty()
	activity.SetActivityStreamsObject(object)
	// Object is set but empty (no items)

	err := handleFollowInboxRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleFollowInboxRequest with empty object should return error")
	}
}

func TestHandleFollowWithNilActorIRI(t *testing.T) {
	activity := streams.NewActivityStreamsFollow()

	// Set a valid object with IRI
	object := streams.NewActivityStreamsObjectProperty()
	objectNote := streams.NewActivityStreamsNote()
	objectID := streams.NewJSONLDIdProperty()
	objectID.SetIRI(mustParseURL("https://example.com/note/1"))
	objectNote.SetJSONLDId(objectID)
	object.AppendActivityStreamsNote(objectNote)
	activity.SetActivityStreamsObject(object)

	// Actor is nil
	err := handleFollowInboxRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleFollowInboxRequest with nil actor should return error")
	}
}

func TestHandleAnnounceWithNilObject(t *testing.T) {
	activity := streams.NewActivityStreamsAnnounce()
	// Don't set object or actor

	err := handleAnnounceRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleAnnounceRequest with nil object should return error")
	}
}

func TestHandleAnnounceWithEmptyObject(t *testing.T) {
	activity := streams.NewActivityStreamsAnnounce()
	object := streams.NewActivityStreamsObjectProperty()
	activity.SetActivityStreamsObject(object)

	err := handleAnnounceRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleAnnounceRequest with empty object should return error")
	}
}

func TestHandleAnnounceWithNilActorIRI(t *testing.T) {
	activity := streams.NewActivityStreamsAnnounce()

	// Set object with IRI
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendIRI(mustParseURL("https://example.com/note/1"))
	activity.SetActivityStreamsObject(object)

	// Actor is nil
	err := handleAnnounceRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleAnnounceRequest with nil actor should return error")
	}
}

func TestHandleLikeWithNilObject(t *testing.T) {
	activity := streams.NewActivityStreamsLike()
	// Don't set object or actor

	err := handleLikeRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleLikeRequest with nil object should return error")
	}
}

func TestHandleLikeWithEmptyObject(t *testing.T) {
	activity := streams.NewActivityStreamsLike()
	object := streams.NewActivityStreamsObjectProperty()
	activity.SetActivityStreamsObject(object)

	err := handleLikeRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleLikeRequest with empty object should return error")
	}
}

func TestHandleLikeWithNilActorIRI(t *testing.T) {
	activity := streams.NewActivityStreamsLike()

	// Set object with IRI
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendIRI(mustParseURL("https://example.com/note/1"))
	activity.SetActivityStreamsObject(object)

	// Actor is nil
	err := handleLikeRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleLikeRequest with nil actor should return error")
	}
}

func TestHandleCreateWithNilId(t *testing.T) {
	activity := streams.NewActivityStreamsCreate()
	// Don't set JSONLD ID

	err := handleCreateRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleCreateRequest with nil ID should return error")
	}
}

func TestHandleCreateWithIdButNilIRI(t *testing.T) {
	activity := streams.NewActivityStreamsCreate()
	id := streams.NewJSONLDIdProperty()
	// Set the ID property but don't set an IRI on it
	activity.SetJSONLDId(id)

	err := handleCreateRequest(context.Background(), activity)
	if err == nil {
		t.Error("handleCreateRequest with ID but nil IRI should return error")
	}
}

func TestHandleUpdateWithNilObject(t *testing.T) {
	activity := streams.NewActivityStreamsUpdate()
	// Don't set object - should return nil (not an error, just skip)

	// This should not panic and should return nil since we only care about Person updates
	err := handleUpdateRequest(context.Background(), activity)
	if err != nil {
		t.Errorf("handleUpdateRequest with nil object should return nil (skip), got %v", err)
	}
}

func TestHandleUpdateWithEmptyObject(t *testing.T) {
	activity := streams.NewActivityStreamsUpdate()
	object := streams.NewActivityStreamsObjectProperty()
	activity.SetActivityStreamsObject(object)

	// Should return nil since empty object means it's not a Person update
	err := handleUpdateRequest(context.Background(), activity)
	if err != nil {
		t.Errorf("handleUpdateRequest with empty object should return nil (skip), got %v", err)
	}
}

func TestHandleUpdateWithNonPersonObject(t *testing.T) {
	activity := streams.NewActivityStreamsUpdate()
	object := streams.NewActivityStreamsObjectProperty()
	note := streams.NewActivityStreamsNote()
	object.AppendActivityStreamsNote(note)
	activity.SetActivityStreamsObject(object)

	// Should return nil since it's not a Person update
	err := handleUpdateRequest(context.Background(), activity)
	if err != nil {
		t.Errorf("handleUpdateRequest with non-Person object should return nil (skip), got %v", err)
	}
}

// TestNilSafetyNoPanic verifies that none of the handlers panic when given
// completely empty activities. This is the most important test - we want to
// ensure that malformed ActivityPub payloads don't crash the server.
func TestNilSafetyNoPanic(t *testing.T) {
	ctx := context.Background()

	t.Run("Follow with nil properties", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("handleFollowInboxRequest panicked: %v", r)
			}
		}()
		activity := streams.NewActivityStreamsFollow()
		_ = handleFollowInboxRequest(ctx, activity)
	})

	t.Run("Announce with nil properties", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("handleAnnounceRequest panicked: %v", r)
			}
		}()
		activity := streams.NewActivityStreamsAnnounce()
		_ = handleAnnounceRequest(ctx, activity)
	})

	t.Run("Like with nil properties", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("handleLikeRequest panicked: %v", r)
			}
		}()
		activity := streams.NewActivityStreamsLike()
		_ = handleLikeRequest(ctx, activity)
	})

	t.Run("Create with nil properties", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("handleCreateRequest panicked: %v", r)
			}
		}()
		activity := streams.NewActivityStreamsCreate()
		_ = handleCreateRequest(ctx, activity)
	})

	t.Run("Update with nil properties", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("handleUpdateRequest panicked: %v", r)
			}
		}()
		activity := streams.NewActivityStreamsUpdate()
		_ = handleUpdateRequest(ctx, activity)
	})
}
