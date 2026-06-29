package models

import (
	"encoding/json"
	"testing"
)

func strptr(s string) *string { return &s }

// TestStreamKeyJSONStability locks the on-the-wire / on-disk JSON shape of a
// StreamKey. Stream keys are persisted in the datastore and returned over the
// admin API with this exact shape; a tag change here would make existing
// installs' stored keys unreadable and break the admin client.
func TestStreamKeyJSONStability(t *testing.T) {
	cases := []struct {
		name string
		in   StreamKey
		want string
	}{
		{
			name: "key and comment",
			in:   StreamKey{Key: strptr("abc123"), Comment: strptr("main key")},
			want: `{"comment":"main key","key":"abc123"}`,
		},
		{
			name: "key only omits comment",
			in:   StreamKey{Key: strptr("abc123")},
			want: `{"key":"abc123"}`,
		},
		{
			name: "empty omits all",
			in:   StreamKey{},
			want: `{}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.in)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("StreamKey JSON changed.\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}

// TestStreamKeyRoundTrip proves a persisted blob decodes back into the same
// fields, i.e. keys written by an older build still load.
func TestStreamKeyRoundTrip(t *testing.T) {
	stored := `[{"key":"k1","comment":"first"},{"key":"k2"}]`

	var keys []StreamKey
	if err := json.Unmarshal([]byte(stored), &keys); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0].Key == nil || *keys[0].Key != "k1" {
		t.Errorf("keys[0].Key = %v, want k1", keys[0].Key)
	}
	if keys[0].Comment == nil || *keys[0].Comment != "first" {
		t.Errorf("keys[0].Comment = %v, want first", keys[0].Comment)
	}
	if keys[1].Comment != nil {
		t.Errorf("keys[1].Comment = %v, want nil", keys[1].Comment)
	}
}
