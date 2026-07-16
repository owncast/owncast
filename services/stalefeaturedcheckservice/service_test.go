package stalefeaturedcheckservice

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
)

// fakeConfigRepo implements only the method the checker uses; the embedded
// interface panics if anything else is called.
type fakeConfigRepo struct {
	configrepository.ConfigRepository
	federationEnabled bool
}

func (f *fakeConfigRepo) GetFederationEnabled() bool { return f.federationEnabled }

type fakeServersRepo struct {
	federatedserversrepository.FederatedServersRepository

	mu        sync.Mutex
	servers   []models.FederatedServer
	getErr    error
	updateErr error
	// markedOffline records IRIs passed to UpdateServerStatus with isOnline=false.
	markedOffline []string
	// getCalls, when non-nil, receives a signal per GetFederatedServers call.
	getCalls chan struct{}
}

func (f *fakeServersRepo) GetFederatedServers() ([]models.FederatedServer, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.getCalls != nil {
		f.getCalls <- struct{}{}
	}
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.servers, nil
}

func (f *fakeServersRepo) UpdateServerStatus(iri string, isOnline bool, metadata *models.FederatedStreamUpdate) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if isOnline || metadata != nil {
		panic("checker must only mark servers offline with nil metadata")
	}
	if f.updateErr != nil {
		return f.updateErr
	}
	f.markedOffline = append(f.markedOffline, iri)
	return nil
}

func (f *fakeServersRepo) marked() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.markedOffline...)
}

func TestCheckMarksOnlyStaleServers(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		isOnline   bool
		lastUpdate *time.Time
		wantMarked bool
	}{
		{"stale online server is marked offline", true, new(now.Add(-staleThreshold - time.Hour)), true},
		{"exactly at threshold is not stale", true, new(now.Add(-staleThreshold)), false},
		{"one nanosecond past threshold is stale", true, new(now.Add(-staleThreshold - time.Nanosecond)), true},
		{"fresh online server is untouched", true, new(now.Add(-time.Second)), false},
		{"offline server is untouched even if stale", false, new(now.Add(-staleThreshold - time.Hour)), false},
		{"online server with nil last update is untouched", true, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeServersRepo{servers: []models.FederatedServer{
				{IRI: "https://example.com/actor", IsOnline: tt.isOnline, LastStatusUpdate: tt.lastUpdate},
			}}
			c := New(&fakeConfigRepo{federationEnabled: true}, repo, func() time.Time { return now }, nil)

			c.check()

			marked := repo.marked()
			if tt.wantMarked && len(marked) != 1 {
				t.Fatalf("expected server to be marked offline, marked = %v", marked)
			}
			if !tt.wantMarked && len(marked) != 0 {
				t.Fatalf("expected no servers marked offline, marked = %v", marked)
			}
		})
	}
}

func TestCheckGetErrorMarksNothing(t *testing.T) {
	repo := &fakeServersRepo{getErr: errors.New("db closed")}
	c := New(&fakeConfigRepo{federationEnabled: true}, repo, time.Now, nil)

	c.check()

	if marked := repo.marked(); len(marked) != 0 {
		t.Fatalf("expected no updates after fetch error, marked = %v", marked)
	}
}

func TestCheckUpdateErrorDoesNotPanic(t *testing.T) {
	stale := time.Now().Add(-staleThreshold - time.Hour)
	repo := &fakeServersRepo{
		updateErr: errors.New("db closed"),
		servers: []models.FederatedServer{
			{IRI: "https://a.example.com", IsOnline: true, LastStatusUpdate: &stale},
			{IRI: "https://b.example.com", IsOnline: true, LastStatusUpdate: &stale},
		},
	}
	c := New(&fakeConfigRepo{federationEnabled: true}, repo, time.Now, nil)

	// Must attempt both servers and not panic on per-server failure.
	c.check()
}

func TestStartWithFederationDisabledDoesNothing(t *testing.T) {
	repo := &fakeServersRepo{}
	ticksCalled := false
	c := New(&fakeConfigRepo{federationEnabled: false}, repo, time.Now, func() (<-chan time.Time, func()) {
		ticksCalled = true
		return make(chan time.Time), func() {}
	})

	c.Start()

	if ticksCalled {
		t.Fatal("Start must not create a ticker when federation is disabled")
	}
	if c.done != nil {
		t.Fatal("Start must not begin the background loop when federation is disabled")
	}
	c.Stop() // Stop on a never-started checker must be a no-op.
}

func TestStartStopLifecycle(t *testing.T) {
	tickCh := make(chan time.Time)
	tickerStops := 0
	repo := &fakeServersRepo{getCalls: make(chan struct{}, 16)}
	c := New(&fakeConfigRepo{federationEnabled: true}, repo, time.Now, func() (<-chan time.Time, func()) {
		tickerStops++
		return tickCh, func() {}
	})

	c.Start()
	<-repo.getCalls // immediate check on start

	firstDone := c.done
	c.Start() // second Start is a no-op, not a second goroutine
	if c.done != firstDone {
		t.Fatal("double Start must not restart the checker")
	}
	if tickerStops != 1 {
		t.Fatalf("double Start created %d tickers, want 1", tickerStops)
	}

	tickCh <- time.Now()
	<-repo.getCalls // tick-driven check

	stopped := c.stopped
	c.Stop()
	<-stopped // goroutine exited: no leak
	if c.done != nil {
		t.Fatal("Stop must clear running state")
	}
	c.Stop() // idempotent

	// Restart works after Stop.
	c.Start()
	<-repo.getCalls
	stopped = c.stopped
	c.Stop()
	<-stopped
}
