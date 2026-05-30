package dispatcher

import (
	"context"
	"sync"
	"testing"
)

func TestPublishFansOutToListeners(t *testing.T) {
	d := New()
	var got []string
	d.AddListener(func(_ context.Context, e Event) { got = append(got, "a:"+e.Type) })
	d.AddListener(func(_ context.Context, e Event) { got = append(got, "b:"+e.Type) })

	d.Publish(context.Background(), Event{Type: "x"})

	if len(got) != 2 || got[0] != "a:x" || got[1] != "b:x" {
		t.Fatalf("listeners not invoked in order: %v", got)
	}
}

func TestApplyFiltersPassesWhenNoneRegistered(t *testing.T) {
	d := New()
	if !d.ApplyFilters(context.Background(), Event{Type: "x"}) {
		t.Fatal("expected pass with no filters registered")
	}
}

func TestApplyFiltersRunsChainAndMutatesInPlace(t *testing.T) {
	d := New()
	d.AddFilter(func(_ context.Context, e Event) bool {
		*e.Payload.(*string) += "-1"
		return true
	})
	d.AddFilter(func(_ context.Context, e Event) bool {
		*e.Payload.(*string) += "-2"
		return true
	})

	body := "body"
	if !d.ApplyFilters(context.Background(), Event{Type: "x", Payload: &body}) {
		t.Fatal("expected event to pass")
	}
	if body != "body-1-2" {
		t.Errorf("in-place mutation chain = %q, want %q", body, "body-1-2")
	}
}

func TestApplyFiltersDropStopsChain(t *testing.T) {
	d := New()
	ran := 0
	d.AddFilter(func(_ context.Context, _ Event) bool { ran++; return false })
	d.AddFilter(func(_ context.Context, _ Event) bool { ran++; return true })

	if d.ApplyFilters(context.Background(), Event{Type: "x"}) {
		t.Fatal("expected event to be dropped")
	}
	if ran != 1 {
		t.Errorf("filters run = %d, want 1 (chain should stop at the drop)", ran)
	}
}

func TestPublishWithNoListeners(t *testing.T) {
	// Publishing with nothing registered must be a safe no-op.
	New().Publish(context.Background(), Event{Type: "x"})
}

func TestConcurrentRegisterAndDispatch(t *testing.T) {
	// Registration and dispatch must be race-free (run with -race).
	d := New()
	var wg sync.WaitGroup
	for i := 0; i < 25; i++ {
		wg.Add(3)
		go func() { defer wg.Done(); d.AddListener(func(context.Context, Event) {}) }()
		go func() { defer wg.Done(); d.Publish(context.Background(), Event{Type: "x"}) }()
		go func() { defer wg.Done(); d.ApplyFilters(context.Background(), Event{Type: "x"}) }()
	}
	wg.Wait()
}
