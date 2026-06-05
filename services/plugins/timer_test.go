package plugins

import (
	"sync/atomic"
	"testing"
	"time"
)

// fireCounter builds a TimerHub whose resolve records each fire and returns nil
// (so the callOnEvent path is skipped — we only assert that timers fire on
// schedule). minDelay is shrunk so timing tests run fast.
func fireCounter() (*TimerHub, *int64) {
	var n int64
	h := NewTimerHub(func(slug string) *Loaded {
		atomic.AddInt64(&n, 1)
		return nil
	})
	h.minDelay = time.Millisecond
	return h, &n
}

func TestTimerHubFiresOnceThenRemoves(t *testing.T) {
	h, n := fireCounter()
	if !h.Schedule("p", 1, 5, false) {
		t.Fatal("schedule should succeed")
	}
	time.Sleep(60 * time.Millisecond)
	if got := atomic.LoadInt64(n); got != 1 {
		t.Fatalf("expected exactly 1 fire, got %d", got)
	}
	h.mu.Lock()
	pending := h.counts["p"]
	h.mu.Unlock()
	if pending != 0 {
		t.Errorf("one-shot should be removed after firing, %d still pending", pending)
	}
}

func TestTimerHubRepeatsUntilCleared(t *testing.T) {
	h, n := fireCounter()
	h.Schedule("p", 1, 5, true)
	time.Sleep(70 * time.Millisecond)
	h.Clear("p", 1)
	time.Sleep(30 * time.Millisecond)
	after := atomic.LoadInt64(n)
	if after < 2 {
		t.Fatalf("interval should fire repeatedly, got %d", after)
	}
	time.Sleep(40 * time.Millisecond)
	if again := atomic.LoadInt64(n); again != after {
		t.Errorf("interval kept firing after Clear: %d -> %d", after, again)
	}
}

func TestTimerHubClearBeforeFire(t *testing.T) {
	var n int64
	h := NewTimerHub(func(string) *Loaded { atomic.AddInt64(&n, 1); return nil })
	// Keep the default 100ms floor so we can cancel before it fires.
	h.Schedule("p", 1, 100, false)
	h.Clear("p", 1)
	time.Sleep(160 * time.Millisecond)
	if atomic.LoadInt64(&n) != 0 {
		t.Error("a cleared timer must not fire")
	}
}

func TestTimerHubClampsDelayToMinimum(t *testing.T) {
	var n int64
	h := NewTimerHub(func(string) *Loaded { atomic.AddInt64(&n, 1); return nil })
	// Default minDelay is 100ms; a 0ms request must be clamped up, not fire now.
	h.Schedule("p", 1, 0, false)
	time.Sleep(40 * time.Millisecond)
	if atomic.LoadInt64(&n) != 0 {
		t.Error("0ms delay should be clamped up to the minimum, not fire immediately")
	}
	time.Sleep(120 * time.Millisecond)
	if atomic.LoadInt64(&n) != 1 {
		t.Error("timer should fire once after the clamped minimum delay")
	}
}

func TestTimerHubPerPluginCap(t *testing.T) {
	h := NewTimerHub(func(string) *Loaded { return nil })
	h.maxPerPlugin = 2
	// Long delays so nothing fires during the test.
	if !h.Schedule("p", 1, 100000, false) || !h.Schedule("p", 2, 100000, false) {
		t.Fatal("first two timers should fit under the cap")
	}
	if h.Schedule("p", 3, 100000, false) {
		t.Error("third timer should be rejected by the per-plugin cap")
	}
	// Re-using an existing id replaces it and doesn't count against the cap.
	if !h.Schedule("p", 1, 100000, false) {
		t.Error("re-setting an existing id should succeed")
	}
	// A different plugin is unaffected by p's cap.
	if !h.Schedule("q", 1, 100000, false) {
		t.Error("another plugin should not be limited by p's cap")
	}
	h.CancelForPlugin("p")
	h.CancelForPlugin("q")
}

func TestTimerHubCancelForPlugin(t *testing.T) {
	h, _ := fireCounter()
	h.Schedule("p", 1, 100000, true)
	h.Schedule("p", 2, 100000, true)
	h.CancelForPlugin("p")
	h.mu.Lock()
	remaining := len(h.timers["p"])
	count := h.counts["p"]
	h.mu.Unlock()
	if remaining != 0 || count != 0 {
		t.Errorf("CancelForPlugin should drop all timers, got %d entries / count %d", remaining, count)
	}
}
