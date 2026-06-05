package plugins

import (
	"encoding/json"
	"testing"
)

// Each subscription gets a distinct, non-zero connection id so a plugin can
// pair a disconnect with its connect and tell concurrent connections apart.
func TestSSEHubAssignsUniqueConnectionIDs(t *testing.T) {
	h := NewSSEHub()

	_, u1, id1, ok1 := h.Subscribe("p1", "c")
	_, u2, id2, ok2 := h.Subscribe("p1", "c")
	_, u3, id3, ok3 := h.Subscribe("p2", "other")
	if !ok1 || !ok2 || !ok3 {
		t.Fatal("subscribes should succeed")
	}
	defer u1()
	defer u2()
	defer u3()

	if id1 == 0 || id2 == 0 || id3 == 0 {
		t.Errorf("connection ids must be non-zero, got %d %d %d", id1, id2, id3)
	}
	if id1 == id2 || id1 == id3 || id2 == id3 {
		t.Errorf("connection ids must be unique, got %d %d %d", id1, id2, id3)
	}
}

// The sse.connect / sse.disconnect payload carries the channel, connection id,
// and the resolved user (omitted when anonymous).
func TestSSEConnectionEnvelope(t *testing.T) {
	user := &HostUser{ID: "u-1", DisplayName: "Ann", IsAuthenticated: true}
	raw, err := sseConnectionEnvelope(EventSSEConnect, "presence", 7, user)
	if err != nil {
		t.Fatalf("sseConnectionEnvelope: %v", err)
	}

	var env struct {
		EventType string             `json:"eventType"`
		Payload   SSEConnectionEvent `json:"payload"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.EventType != EventSSEConnect {
		t.Errorf("eventType: got %q, want %q", env.EventType, EventSSEConnect)
	}
	if env.Payload.Channel != "presence" {
		t.Errorf("channel: got %q", env.Payload.Channel)
	}
	if env.Payload.ConnectionID != 7 {
		t.Errorf("connectionId: got %d, want 7", env.Payload.ConnectionID)
	}
	if env.Payload.User == nil || env.Payload.User.ID != "u-1" {
		t.Errorf("user not carried through: %+v", env.Payload.User)
	}

	// Anonymous connection: user is omitted entirely from the JSON.
	raw, err = sseConnectionEnvelope(EventSSEDisconnect, "presence", 8, nil)
	if err != nil {
		t.Fatalf("sseConnectionEnvelope (anon): %v", err)
	}
	var generic struct {
		Payload map[string]any `json:"payload"`
	}
	if err := json.Unmarshal(raw, &generic); err != nil {
		t.Fatalf("unmarshal anon: %v", err)
	}
	if _, ok := generic.Payload["user"]; ok {
		t.Error("user should be omitted for an anonymous connection")
	}
}
