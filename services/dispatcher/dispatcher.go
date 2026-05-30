// Package dispatcher is a tiny in-process publish/subscribe + filter mediator.
//
// It lets event producers (chat, webhooks, …) notify and filter through a
// single neutral component that consumers (the plugin host today, other
// subsystems in future) subscribe to. Producers depend on the dispatcher and
// call it; consumers register handlers on it. Neither side has to know about,
// import, or be mutated by the other — which keeps every service on plain
// constructor injection instead of post-construction setters.
//
// Construct one in the composition root and inject it everywhere; it depends
// on nothing.
package dispatcher

import (
	"context"
	"sync"
)

// Event is a named event with an arbitrary payload. Producers and consumers
// agree on the concrete payload type for each event Type (e.g. a
// *events.UserMessageEvent for a chat message).
type Event struct {
	Type    string
	Payload any
}

// Listener is notified of an event that has already happened. Listeners run
// synchronously in the producer's goroutine, so they must not block — offload
// slow work (e.g. wasm plugin calls) to a goroutine.
type Listener func(ctx context.Context, event Event)

// Filter inspects an event before it takes effect and reports whether it
// should proceed; returning false drops the event. A filter may mutate the
// payload in place (it's passed by pointer) to rewrite the event.
type Filter func(ctx context.Context, event Event) (allow bool)

// Dispatcher fans events out to registered listeners and runs them through
// registered filters. Use New to construct one.
type Dispatcher struct {
	mu        sync.RWMutex
	listeners []Listener
	filters   []Filter
}

// New constructs an empty Dispatcher.
func New() *Dispatcher {
	return &Dispatcher{}
}

// AddListener registers l to be notified of every published event. Intended
// to be called during composition, before events start flowing.
func (d *Dispatcher) AddListener(l Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners = append(d.listeners, l)
}

// AddFilter appends f to the filter chain. Filters run in registration order
// and any may drop the event.
func (d *Dispatcher) AddFilter(f Filter) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.filters = append(d.filters, f)
}

// Publish notifies every registered listener of event, in registration order.
// The listener slice is copied under the lock before iteration; the
// outer functions are called without the lock held so a listener that
// re-publishes (or calls AddListener) doesn't deadlock, and a
// concurrent AddListener can't mutate the slice we're iterating.
func (d *Dispatcher) Publish(ctx context.Context, event Event) {
	d.mu.RLock()
	listeners := make([]Listener, len(d.listeners))
	copy(listeners, d.listeners)
	d.mu.RUnlock()

	for _, l := range listeners {
		l(ctx, event)
	}
}

// ApplyFilters runs event through the filter chain and reports whether it
// survived. Filters may mutate the payload in place. With no filters
// registered the event always passes. Same copy-under-lock discipline
// as Publish: AddFilter on another goroutine can't race with the
// iteration here.
func (d *Dispatcher) ApplyFilters(ctx context.Context, event Event) (allow bool) {
	d.mu.RLock()
	filters := make([]Filter, len(d.filters))
	copy(filters, d.filters)
	d.mu.RUnlock()

	for _, f := range filters {
		if !f(ctx, event) {
			return false
		}
	}
	return true
}
