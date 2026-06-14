package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// DefaultMaxTimersPerPlugin caps how many timers a single plugin may have
// pending at once. Like the SSE connection cap, it's per-plugin so one plugin
// can't exhaust host resources at the expense of others.
const DefaultMaxTimersPerPlugin = 64

// DefaultMinTimerDelay is the floor for a timer delay/interval. Requests below
// it are clamped up, so a plugin can't spin the host with a near-zero interval.
const DefaultMinTimerDelay = 100 * time.Millisecond

// maxTimerDelayMs is the ceiling for a requested delay, in milliseconds (24h).
// Longer is pointless for in-memory timers (they don't survive a restart) and
// guards the time.Duration math from overflow. Also keeps the host-function's
// uint64→int64 delay conversion provably in range.
const maxTimerDelayMs = 24 * 60 * 60 * 1000

// timerEntry is one scheduled timer for a plugin.
type timerEntry struct {
	slug      string
	id        uint64
	delay     time.Duration
	repeat    bool
	t         *time.Timer
	cancelled bool
}

// TimerHub schedules host-driven callbacks for plugins. Plugins can't use
// setTimeout (no event loop in the wasm sandbox), so they ask the host to call
// them back: the guest registers a timer via owncast_timer_set, and when it
// fires the host calls the plugin's on_event export with a timer.fire event
// carrying the id. The guest SDK maps the id back to the author's callback.
//
// Timers are in-memory and per-plugin; they're cancelled when the plugin
// unloads and do not survive a host restart.
type TimerHub struct {
	mu           sync.Mutex
	timers       map[string]map[uint64]*timerEntry // slug -> id -> entry
	counts       map[string]int
	maxPerPlugin int
	minDelay     time.Duration

	// resolve maps a plugin slug to its currently-loaded instance so a firing
	// timer can call back into it. Returns nil if the plugin is gone, in which
	// case the fire is dropped. Injected by the host.
	resolve func(slug string) *Loaded
}

// NewTimerHub constructs an empty hub with the default caps. resolve maps a
// plugin slug to its loaded instance (nil-safe) for firing callbacks.
func NewTimerHub(resolve func(slug string) *Loaded) *TimerHub {
	return &TimerHub{
		timers:       make(map[string]map[uint64]*timerEntry),
		counts:       make(map[string]int),
		maxPerPlugin: DefaultMaxTimersPerPlugin,
		minDelay:     DefaultMinTimerDelay,
		resolve:      resolve,
	}
}

// Schedule arms a timer for (slug, id). delayMs is clamped up to the minimum
// delay. repeat=true reschedules after each fire (the next tick is armed only
// once the callback returns, so a slow callback can't cause pile-up).
// Re-using an id replaces the existing timer. Returns false if the plugin is
// already at its pending-timer cap.
func (h *TimerHub) Schedule(slug string, id uint64, delayMs int64, repeat bool) bool {
	delay := time.Duration(delayMs) * time.Millisecond
	if delay < h.minDelay {
		delay = h.minDelay
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if existing := h.timers[slug][id]; existing != nil {
		// Re-set: cancel the old entry first (doesn't count against the cap).
		existing.cancelled = true
		if existing.t != nil {
			existing.t.Stop()
		}
		delete(h.timers[slug], id)
		h.counts[slug]--
	} else if h.counts[slug] >= h.maxPerPlugin {
		return false
	}

	e := &timerEntry{slug: slug, id: id, delay: delay, repeat: repeat}
	if h.timers[slug] == nil {
		h.timers[slug] = make(map[uint64]*timerEntry)
	}
	h.timers[slug][id] = e
	h.counts[slug]++
	e.t = time.AfterFunc(delay, func() { h.fire(e) })
	return true
}

// Clear cancels a single timer. A no-op if the id isn't pending.
func (h *TimerHub) Clear(slug string, id uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e := h.timers[slug][id]
	if e == nil {
		return
	}
	e.cancelled = true
	if e.t != nil {
		e.t.Stop()
	}
	h.removeLocked(e)
}

// CancelForPlugin cancels every pending timer for a plugin. Called when the
// plugin is disabled, reloaded, removed, or the host stops.
func (h *TimerHub) CancelForPlugin(slug string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, e := range h.timers[slug] {
		e.cancelled = true
		if e.t != nil {
			e.t.Stop()
		}
	}
	delete(h.timers, slug)
	delete(h.counts, slug)
}

// removeLocked drops an entry and decrements the plugin's count. Caller holds mu.
func (h *TimerHub) removeLocked(e *timerEntry) {
	if m := h.timers[e.slug]; m != nil {
		if _, ok := m[e.id]; ok {
			delete(m, e.id)
			h.counts[e.slug]--
			if len(m) == 0 {
				delete(h.timers, e.slug)
				delete(h.counts, e.slug)
			}
		}
	}
}

// fire delivers one timer.fire to the plugin, then reschedules (repeat) or
// removes (one-shot) the entry. The callback runs outside the hub lock; the
// reschedule happens only after it returns, so a slow callback can't pile up.
func (h *TimerHub) fire(e *timerEntry) {
	h.mu.Lock()
	cancelled := e.cancelled
	h.mu.Unlock()
	if cancelled {
		return
	}

	if h.resolve != nil {
		if loaded := h.resolve(e.slug); loaded != nil {
			envelope, err := json.Marshal(Envelope{
				EventType: EventTimerFire,
				Payload:   TimerFireEvent{ID: e.id},
			})
			if err == nil {
				if err := callOnEvent(context.Background(), loaded, envelope); err != nil {
					fmt.Fprintf(os.Stderr, "plugin %s: timer %d fire failed: %v\n", e.slug, e.id, err)
				}
			}
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if e.cancelled {
		return
	}
	if e.repeat {
		e.t = time.AfterFunc(e.delay, func() { h.fire(e) })
	} else {
		h.removeLocked(e)
	}
}
