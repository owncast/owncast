package persistence

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClaimInboundFediverseActivityConcurrent(t *testing.T) {
	service := New(_datastore, nil)
	const workers = 32

	start := make(chan struct{})
	errs := make(chan error, workers)
	var claimed atomic.Int32
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			<-start
			ok, err := service.ClaimInboundFediverseActivity(
				"https://remote.example/activities/concurrent-claim",
				"https://remote.example/users/alice",
				"Create",
				time.Now(),
			)
			if err != nil {
				errs <- err
				return
			}
			if ok {
				claimed.Add(1)
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatalf("claim failed: %v", err)
	}
	if got := claimed.Load(); got != 1 {
		t.Fatalf("successful claims = %d, want 1", got)
	}

	ok, err := service.ClaimInboundFediverseActivity(
		"https://remote.example/activities/concurrent-claim",
		"https://remote.example/users/alice",
		"Create",
		time.Now(),
	)
	if err != nil {
		t.Fatalf("retry claim failed: %v", err)
	}
	if ok {
		t.Fatal("retry unexpectedly claimed an already handled activity")
	}
}
