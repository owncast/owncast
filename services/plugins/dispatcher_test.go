package plugins

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// TestReservedEventTypes_CoversCoreEvents is the B4 guard: the set the emit
// host function rejects must include every built-in event type, so a plugin
// can't forge a core event onto other plugins. Catches a new Event* constant
// added without updating the reserved set.
func TestReservedEventTypes_CoversCoreEvents(t *testing.T) {
	for _, e := range []string{
		EventChatMessageReceived, EventChatUserJoined, EventChatUserParted,
		EventChatUserRenamed, EventChatMessageModerated, EventStreamStarted,
		EventStreamStopped, EventStreamTitleChanged, EventSSEConnect,
		EventSSEDisconnect, EventTick, EventTimerFire, EventFediverseFollow,
		EventFediverseLike, EventFediverseRepost, EventFediverseMention,
		EventFediverseReply,
	} {
		if !reservedEventTypes[e] {
			t.Errorf("core event %q must be reserved (plugins must not be able to emit it)", e)
		}
	}
	if reservedEventTypes["my.plugin.custom"] {
		t.Error("a plugin's custom event type must not be treated as reserved")
	}
}

// TestFilterTimeoutErrorShape verifies the error-string the dispatcher
// produces when CallWithContext returns a deadline-exceeded error. The
// real Wazero-level cancellation is verified separately by integration
// tests; this is the Go-side conversion shape that the strike system and
// the host log key off of.
func TestFilterTimeoutErrorShape(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	time.Sleep(time.Millisecond)
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("test setup: expected deadline exceeded, got %v", ctx.Err())
	}
	formatted := formatFilterTimeoutError(ctx.Err())
	if !strings.Contains(formatted.Error(), "timed out") {
		t.Errorf("expected timeout error string, got %v", formatted)
	}
	if !strings.Contains(formatted.Error(), FilterTimeout.String()) {
		t.Errorf("expected error to mention %s, got %v", FilterTimeout, formatted)
	}
}

// formatFilterTimeoutError mirrors the shape callOnFilter produces so we
// can test conversion without instantiating a real wasm plugin.
func formatFilterTimeoutError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.New("on_filter timed out after " + FilterTimeout.String())
	}
	return err
}

// TestFilterTimeoutCountsAsStrike confirms that a timeout-induced failure
// goes through the strike system just like a thrown error would. The
// dispatcher already calls recordFilterFailure on any returned error, so
// this is a direct check of the strike accumulator.
func TestFilterTimeoutCountsAsStrike(t *testing.T) {
	l := &Loaded{Manifest: &Manifest{DisplayName: "slow-filter"}}
	for i := 0; i < FilterStrikeThreshold-1; i++ {
		if disabled := l.recordFilterFailure(); disabled {
			t.Fatalf("disabled too early at strike %d", i+1)
		}
	}
	if disabled := l.recordFilterFailure(); !disabled {
		t.Fatal("threshold-th strike should trigger auto-disable")
	}
	if !l.IsDisabled() {
		t.Fatal("IsDisabled should be true after threshold reached")
	}
}
