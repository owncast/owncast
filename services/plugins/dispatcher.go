package plugins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// Envelope is the JSON payload passed across the plugin boundary for every
// event. Plugins receive this serialized via Host.inputString().
type Envelope struct {
	EventType string `json:"eventType"`
	Payload   any    `json:"payload"`
}

type FilterAction string

const (
	FilterPass   FilterAction = "pass"
	FilterModify FilterAction = "modify"
	FilterDrop   FilterAction = "drop"
)

type FilterResult struct {
	Action  FilterAction `json:"action"`
	Payload any          `json:"payload,omitempty"`
	Reason  string       `json:"reason,omitempty"`
}

// Dispatcher fans events to subscribed plugins. It holds a snapshot function
// rather than a static slice so the Manager can enable/disable plugins at
// runtime and the dispatcher will always see the current set.
type Dispatcher struct {
	snapshot func() []*Loaded
}

// NewDispatcher builds a dispatcher over a fixed plugin set. Use this in
// tests and any context where the plugin set doesn't change after construction.
func NewDispatcher(loaded []*Loaded) *Dispatcher {
	snap := loaded
	return &Dispatcher{snapshot: func() []*Loaded { return snap }}
}

// NewLiveDispatcher builds a dispatcher backed by a snapshot function — the
// Manager passes its Snapshot method here so admin enable/disable shows up
// immediately without restarting the host.
func NewLiveDispatcher(snapshot func() []*Loaded) *Dispatcher {
	return &Dispatcher{snapshot: snapshot}
}

// MaxEmitDepth caps how deep host_emit_event recursion is allowed to go.
// Above this, the dispatcher silently drops the event and logs once. This
// stops a plugin that emits an event it itself subscribes to from blowing
// the stack.
const MaxEmitDepth = 8

// FilterTimeout caps how long a single filter call can run before the host
// cancels it and fails-open (same semantics as a thrown error: skip this
// filter, continue with the previous payload, count as a strike toward
// auto-disable). The chain is unbounded; per-filter caps are sufficient
// because a runaway plugin in slot 1 doesn't block slot 2.
//
// Declared as a var (not const) so tests that need to exercise post-call
// safeguards — e.g. the output-size check — can stretch this temporarily.
// Production callers should treat it as read-only.
var FilterTimeout = 50 * time.Millisecond

type depthKey struct{}

func emitDepth(ctx context.Context) int {
	if v, ok := ctx.Value(depthKey{}).(int); ok {
		return v
	}
	return 0
}

// Dispatch runs the filter chain for an event, then fans the (possibly
// modified) payload out as notifications. If the filter chain drops the
// event, no notifications fire.
func (d *Dispatcher) Dispatch(ctx context.Context, eventType string, payload any) {
	depth := emitDepth(ctx)
	if depth >= MaxEmitDepth {
		fmt.Fprintf(os.Stderr, "dispatcher: max emit depth %d reached for %q — dropping\n", MaxEmitDepth, eventType)
		return
	}
	ctx = context.WithValue(ctx, depthKey{}, depth+1)

	final, allowed, reason := d.Filter(ctx, eventType, payload)
	if !allowed {
		fmt.Fprintf(os.Stderr, "[filtered] %s dropped: %s\n", eventType, reason)
		return
	}
	d.Notify(ctx, eventType, final)
}

// Filter runs the ordered filter chain for an event. Plugin errors are
// fail-open — the event passes through that plugin unchanged. Returns the
// final payload, whether the event survived, and a drop reason if not.
func (d *Dispatcher) Filter(ctx context.Context, eventType string, payload any) (any, bool, string) {
	type subbed struct {
		plugin   *Loaded
		priority int
	}
	var chain []subbed
	for _, p := range d.snapshot() {
		for _, s := range p.Manifest.Subscriptions.Filter {
			if s.Event == eventType {
				chain = append(chain, subbed{p, s.Priority})
			}
		}
	}
	sort.SliceStable(chain, func(i, j int) bool {
		if chain[i].priority != chain[j].priority {
			return chain[i].priority < chain[j].priority
		}
		return chain[i].plugin.Manifest.Slug < chain[j].plugin.Manifest.Slug
	})

	current := payload
	for _, c := range chain {
		if c.plugin.IsDisabled() {
			// Auto-disabled by the strike system; skip silently.
			continue
		}
		result, err := callOnFilter(ctx, c.plugin, eventType, current)
		if err != nil {
			// A plugin torn down mid-chain (or one exporting no filter) is a
			// benign no-op, not a fault: don't log it or count it as a strike.
			if errors.Is(err, errPluginNotLoaded) || errors.Is(err, errPluginNoSuchExport) {
				continue
			}
			fmt.Fprintf(os.Stderr, "plugin %s: on_filter(%s) failed (fail-open): %v\n", c.plugin.Manifest.Slug, eventType, err)
			if c.plugin.recordFilterFailure() {
				fmt.Fprintf(os.Stderr, "plugin %s: auto-disabled after %d consecutive filter failures\n",
					c.plugin.Manifest.Slug, FilterStrikeThreshold)
			}
			continue
		}
		c.plugin.recordFilterSuccess()
		switch result.Action {
		case FilterPass, "":
			// payload unchanged
		case FilterModify:
			current = result.Payload
		case FilterDrop:
			return nil, false, result.Reason
		default:
			fmt.Fprintf(os.Stderr, "plugin %s: unknown filter action %q (fail-open)\n", c.plugin.Manifest.Slug, result.Action)
		}
	}
	return current, true, ""
}

func callOnFilter(ctx context.Context, p *Loaded, eventType string, payload any) (*FilterResult, error) {
	envelope, err := json.Marshal(Envelope{EventType: eventType, Payload: payload})
	if err != nil {
		return nil, fmt.Errorf("marshal envelope: %w", err)
	}
	// Per-call deadline. FilterTimeout has to cover BOTH waiting for the
	// per-plugin mutex AND the wasm call itself: an HTTP request or
	// on_event in flight on the same plugin holds p.mu, so a naive
	// deadline started after Lock() could let the chat hot path stall
	// for the full duration of that other call before the filter clock
	// even starts. Bounding both phases with one context keeps the
	// chat path responsive even when the plugin is otherwise busy.
	callCtx, cancel := context.WithTimeout(ctx, FilterTimeout)
	defer cancel()
	if err := lockWithContext(callCtx, &p.mu); err != nil {
		return nil, fmt.Errorf("on_filter timed out after %s waiting for plugin mutex", FilterTimeout)
	}
	defer p.mu.Unlock()
	// Snapshot the instance under the lock: a concurrent Loaded.Close may have
	// torn it down. Reading p.plugin before the lock raced with that write and
	// could nil-deref.
	pl := p.plugin
	if pl == nil {
		return nil, errPluginNotLoaded
	}
	if !pl.FunctionExists("on_filter") {
		return nil, errPluginNoSuchExport
	}
	_, out, err := pl.CallWithContext(callCtx, "on_filter", envelope)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || callCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("on_filter timed out after %s", FilterTimeout)
		}
		return nil, fmt.Errorf("call on_filter: %w", err)
	}
	if len(out) > MaxFilterOutputBytes {
		return nil, fmt.Errorf("on_filter output too large: %d bytes (max %d)", len(out), MaxFilterOutputBytes)
	}
	var result FilterResult
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("parse on_filter result: %w", err)
	}
	return &result, nil
}

// Notify fans an event out to every plugin subscribed to it via
// `subscriptions.notify`. Plugins run in parallel; errors are logged but
// don't affect other plugins or the caller.
func (d *Dispatcher) Notify(ctx context.Context, eventType string, payload any) {
	envelope := Envelope{EventType: eventType, Payload: payload}
	encoded, err := json.Marshal(envelope)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dispatcher: marshal envelope for %s: %v\n", eventType, err)
		return
	}

	var wg sync.WaitGroup
	for _, p := range d.snapshot() {
		if !subscribed(p.Manifest.Subscriptions.Notify, eventType) {
			continue
		}
		wg.Add(1)
		go func(p *Loaded) {
			defer wg.Done()
			// A plugin must never take down the host. A panic in this goroutine
			// can't be recovered by the caller, so recover it here.
			defer func() {
				if rec := recover(); rec != nil {
					fmt.Fprintf(os.Stderr, "plugin %s: on_event(%s) panicked: %v\n", p.Manifest.Slug, eventType, rec)
				}
			}()
			if err := callOnEvent(ctx, p, encoded); err != nil &&
				!errors.Is(err, errPluginNotLoaded) && !errors.Is(err, errPluginNoSuchExport) {
				fmt.Fprintf(os.Stderr, "plugin %s: on_event(%s) failed: %v\n", p.Manifest.Slug, eventType, err)
			}
		}(p)
	}
	wg.Wait()
}

func subscribed(subs []Subscription, eventType string) bool {
	for _, s := range subs {
		if s.Event == eventType {
			return true
		}
	}
	return false
}

// errPluginNotLoaded / errPluginNoSuchExport are benign "nothing to do"
// sentinels from the per-plugin call helpers: the instance was torn down by a
// concurrent Close, or the plugin doesn't export the requested function. Best-
// effort callers (Notify, SSE connect/disconnect) skip logging them.
var (
	errPluginNotLoaded    = errors.New("plugin is not loaded")
	errPluginNoSuchExport = errors.New("plugin does not export the requested function")
)

func callOnEvent(ctx context.Context, p *Loaded, input []byte) error {
	// Extism plugins are not safe for concurrent calls; serialize per-plugin.
	p.mu.Lock()
	defer p.mu.Unlock()
	// Snapshot the instance under the lock: a concurrent Loaded.Close may have
	// torn it down. Reading p.plugin before the lock raced with that write and
	// could nil-deref.
	pl := p.plugin
	if pl == nil {
		return errPluginNotLoaded
	}
	if !pl.FunctionExists("on_event") {
		return errPluginNoSuchExport
	}
	callCtx, cancel := context.WithTimeout(ctx, NotifyTimeout)
	defer cancel()
	_, _, err := pl.CallWithContext(callCtx, "on_event", input)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || callCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("on_event timed out after %s", NotifyTimeout)
		}
	}
	return err
}

// lockWithContext attempts to acquire mu, but gives up if ctx fires
// first. Go's sync.Mutex can't be cancelled mid-Lock, so when the
// context expires before the lock is acquired we spawn a throwaway
// goroutine that drops the lock as soon as the kernel eventually
// grants it. That goroutine is bounded by the original Lock duration
// and exits cleanly; the caller treats the timeout as an error.
//
// This matters for the chat filter hot path: another call on the
// same plugin (an HTTP handler, an on_event) can hold p.mu for many
// seconds, and a plain Lock() inside Filter would let the chat path
// stall for that whole window before the filter's own deadline even
// started.
func lockWithContext(ctx context.Context, mu *sync.Mutex) error {
	locked := make(chan struct{})
	go func() {
		mu.Lock()
		close(locked)
	}()
	select {
	case <-locked:
		return nil
	case <-ctx.Done():
		go func() {
			<-locked
			mu.Unlock()
		}()
		return ctx.Err()
	}
}
