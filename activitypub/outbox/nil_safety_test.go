package outbox

import (
	"testing"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
)

// TestAddWithNilId verifies that Add doesn't panic when given an item with nil JSONLD ID.
func TestAddWithNilId(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Add panicked with nil ID: %v", r)
		}
	}()

	note := streams.NewActivityStreamsNote()
	// Don't set JSONLD ID

	err := Add(note, "test-id", false)
	if err == nil {
		t.Error("Add with nil ID should return error")
	}
}

// TestAddWithIdButNilIRI verifies that Add doesn't panic when given an item
// with a JSONLD ID property that has no IRI set.
func TestAddWithIdButNilIRI(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Add panicked with ID but nil IRI: %v", r)
		}
	}()

	note := streams.NewActivityStreamsNote()
	id := streams.NewJSONLDIdProperty()
	// Set the ID property but don't set an IRI on it
	note.SetJSONLDId(id)

	err := Add(note, "test-id", false)
	if err == nil {
		t.Error("Add with ID but nil IRI should return error")
	}
}

// TestAddNoPanic verifies that Add doesn't panic with various nil/empty inputs.
func TestAddNoPanic(t *testing.T) {
	testCases := []struct {
		name string
		item func() vocab.ActivityStreamsNote
	}{
		{
			name: "note with no properties",
			item: func() vocab.ActivityStreamsNote {
				return streams.NewActivityStreamsNote()
			},
		},
		{
			name: "note with empty ID property",
			item: func() vocab.ActivityStreamsNote {
				note := streams.NewActivityStreamsNote()
				note.SetJSONLDId(streams.NewJSONLDIdProperty())
				return note
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Add panicked: %v", r)
				}
			}()
			_ = Add(tc.item(), "test-id", false)
		})
	}
}
