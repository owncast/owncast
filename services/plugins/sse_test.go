package plugins

import (
	"strings"
	"testing"
)

func TestSSEHubPublishFansOutToSubscribers(t *testing.T) {
	h := NewSSEHub()

	a, unsubA, _, ok := h.Subscribe("p1", "overlay")
	if !ok {
		t.Fatal("first subscribe should succeed")
	}
	defer unsubA()
	b, unsubB, _, ok := h.Subscribe("p1", "overlay")
	if !ok {
		t.Fatal("second subscribe should succeed")
	}
	defer unsubB()

	if n := h.Publish("p1", "overlay", "emoji", []byte("🦉")); n != 2 {
		t.Fatalf("expected delivery to 2 clients, got %d", n)
	}

	for _, ch := range []<-chan []byte{a, b} {
		frame := string(<-ch)
		if !strings.Contains(frame, "event: emoji\n") || !strings.Contains(frame, "data: 🦉\n") {
			t.Errorf("unexpected SSE frame: %q", frame)
		}
		if !strings.HasSuffix(frame, "\n\n") {
			t.Errorf("frame must end with a blank line: %q", frame)
		}
	}
}

func TestSSEHubIsolatesByPluginAndChannel(t *testing.T) {
	h := NewSSEHub()

	other, unsub, _, _ := h.Subscribe("p2", "overlay")
	defer unsub()
	wrongChannel, unsub2, _, _ := h.Subscribe("p1", "stats")
	defer unsub2()

	// No subscribers on p1/overlay → zero deliveries, and neither of the
	// other subscriptions should receive anything.
	if n := h.Publish("p1", "overlay", "x", []byte("y")); n != 0 {
		t.Fatalf("expected 0 deliveries, got %d", n)
	}
	select {
	case f := <-other:
		t.Errorf("p2 subscriber should not receive p1 events, got %q", f)
	case f := <-wrongChannel:
		t.Errorf("p1/stats subscriber should not receive p1/overlay events, got %q", f)
	default:
	}
}

func TestSSEHubUnsubscribeStopsDelivery(t *testing.T) {
	h := NewSSEHub()

	ch, unsub, _, _ := h.Subscribe("p1", "c")
	unsub()

	if n := h.Publish("p1", "c", "e", []byte("d")); n != 0 {
		t.Fatalf("expected 0 deliveries after unsubscribe, got %d", n)
	}
	// Unsubscribe must be idempotent (the stream handler defers it, and it
	// may also run on context cancel).
	unsub()
	select {
	case f := <-ch:
		t.Errorf("no frame expected after unsubscribe, got %q", f)
	default:
	}
}

func TestSSEHubPerPluginConnectionCap(t *testing.T) {
	h := NewSSEHub()
	h.maxPerPlugin = 2

	_, u1, _, ok1 := h.Subscribe("p1", "c")
	_, u2, _, ok2 := h.Subscribe("p1", "c")
	_, _, _, ok3 := h.Subscribe("p1", "c")
	if !ok1 || !ok2 {
		t.Fatal("first two subscribes should succeed under the cap")
	}
	if ok3 {
		t.Fatal("third subscribe should be rejected by the per-plugin cap")
	}

	// A different plugin is unaffected by p1's cap.
	if _, u, _, ok := h.Subscribe("p2", "c"); !ok {
		t.Fatal("other plugin should not be affected by p1's cap")
	} else {
		u()
	}

	// Freeing a slot lets a new connection in.
	u1()
	if _, u, _, ok := h.Subscribe("p1", "c"); !ok {
		t.Fatal("subscribe should succeed after a slot is freed")
	} else {
		u()
	}
	u2()
}

func TestSSEHubDropsFramesForSlowClient(t *testing.T) {
	h := NewSSEHub()
	_, unsub, _, _ := h.Subscribe("p1", "c") // never drained
	defer unsub()

	// The client buffer is 16; publishing well beyond it must not block and
	// must report fewer deliveries than sends once the buffer fills.
	delivered := 0
	for i := 0; i < 100; i++ {
		delivered += h.Publish("p1", "c", "e", []byte("d"))
	}
	if delivered == 0 {
		t.Fatal("expected some deliveries before the buffer filled")
	}
	if delivered >= 100 {
		t.Fatalf("expected dropped frames for the undrained client, delivered=%d", delivered)
	}
}

func TestFrameSSEMultilineData(t *testing.T) {
	got := string(frameSSE("", []byte("line1\nline2")))
	want := "data: line1\ndata: line2\n\n"
	if got != want {
		t.Errorf("frameSSE multiline = %q, want %q", got, want)
	}
}
