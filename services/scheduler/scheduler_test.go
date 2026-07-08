package scheduler

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// waitFor polls cond until it returns true or the timeout elapses. Timing
// assertions in this file are eventually-style so the suite stays
// deterministic under CI load: delays can only make a condition take
// longer, never flip a passing assertion.
func waitFor(t *testing.T, timeout time.Duration, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out after %v waiting for %s", timeout, msg)
}

func newService(t *testing.T) *Service {
	t.Helper()
	s, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
	return s
}

func TestRegisterValidation(t *testing.T) {
	noop := func(ctx context.Context) {}

	cases := []struct {
		name    string
		job     Job
		wantErr bool
	}{
		{"empty name", Job{Name: "", Interval: time.Second, Run: noop}, true},
		{"zero interval", Job{Name: "zero", Interval: 0, Run: noop}, true},
		{"negative interval", Job{Name: "negative", Interval: -time.Second, Run: noop}, true},
		{"nil run", Job{Name: "nilrun", Interval: time.Second, Run: nil}, true},
		{"valid job", Job{Name: "valid", Interval: time.Second, Run: noop}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newService(t)
			err := s.Register(tc.job)
			if tc.wantErr && err == nil {
				t.Errorf("Register(%+v) returned nil error, want error", tc.job)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Register(%+v) returned error %v, want nil", tc.job, err)
			}
		})
	}

	t.Run("duplicate name", func(t *testing.T) {
		s := newService(t)
		if err := s.Register(Job{Name: "dupe", Interval: time.Second, Run: noop}); err != nil {
			t.Fatalf("first Register returned error: %v", err)
		}
		if err := s.Register(Job{Name: "dupe", Interval: time.Minute, Run: noop}); err == nil {
			t.Errorf("second Register with duplicate name returned nil error, want error")
		}
	})
}

func TestJobFiresRepeatedly(t *testing.T) {
	s := newService(t)
	var runs atomic.Int32
	err := s.Register(Job{
		Name:     "repeat",
		Interval: 50 * time.Millisecond,
		Run:      func(ctx context.Context) { runs.Add(1) },
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	s.Start()
	defer s.Stop()

	waitFor(t, 3*time.Second, func() bool { return runs.Load() >= 2 },
		"job to run at least twice")
}

func TestRunAtStartFiresImmediately(t *testing.T) {
	s := newService(t)
	var runs atomic.Int32
	// Interval far beyond the poll deadline: the only way this job fires
	// within the deadline is the immediate run-at-start execution.
	err := s.Register(Job{
		Name:       "immediate",
		Interval:   time.Hour,
		RunAtStart: true,
		Run:        func(ctx context.Context) { runs.Add(1) },
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	s.Start()
	defer s.Stop()

	waitFor(t, 2*time.Second, func() bool { return runs.Load() >= 1 },
		"run-at-start job to fire without waiting an interval")
}

func TestNoRunAtStartWaitsAnInterval(t *testing.T) {
	s := newService(t)
	const interval = 600 * time.Millisecond
	var firstRunElapsed atomic.Int64
	start := time.Now()
	err := s.Register(Job{
		Name:     "patient",
		Interval: interval,
		Run: func(ctx context.Context) {
			// Record only the first run's elapsed time.
			firstRunElapsed.CompareAndSwap(0, int64(time.Since(start)))
		},
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	s.Start()
	defer s.Stop()

	waitFor(t, 3*time.Second, func() bool { return firstRunElapsed.Load() != 0 },
		"job to fire once")

	// The scheduler never fires early, so load delays cannot flip this:
	// the elapsed time only grows. Margin below the interval absorbs
	// timer granularity.
	if elapsed := time.Duration(firstRunElapsed.Load()); elapsed < 500*time.Millisecond {
		t.Errorf("job without RunAtStart fired after %v, want at least ~one interval (%v)", elapsed, interval)
	}
}

func TestOverlapProtection(t *testing.T) {
	s := newService(t)
	var current, maxConcurrent, completed atomic.Int32
	err := s.Register(Job{
		Name:       "slow",
		Interval:   50 * time.Millisecond,
		RunAtStart: true,
		Run: func(ctx context.Context) {
			c := current.Add(1)
			for {
				m := maxConcurrent.Load()
				if c <= m || maxConcurrent.CompareAndSwap(m, c) {
					break
				}
			}
			// Outlive the interval 3x so ticks arrive mid-run.
			time.Sleep(150 * time.Millisecond)
			current.Add(-1)
			completed.Add(1)
		},
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	s.Start()
	defer s.Stop()

	// Two completions prove later ticks still fire after a long run.
	waitFor(t, 3*time.Second, func() bool { return completed.Load() >= 2 },
		"slow job to complete at least twice")

	if got := maxConcurrent.Load(); got != 1 {
		t.Errorf("max concurrent executions = %d, want 1 (overlap protection)", got)
	}
}

func TestPanicRecovery(t *testing.T) {
	s := newService(t)
	var panicRuns, siblingRuns atomic.Int32
	err := s.Register(Job{
		Name:     "panics",
		Interval: 50 * time.Millisecond,
		Run: func(ctx context.Context) {
			panicRuns.Add(1)
			panic("job blew up")
		},
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	err = s.Register(Job{
		Name:     "sibling",
		Interval: 50 * time.Millisecond,
		Run:      func(ctx context.Context) { siblingRuns.Add(1) },
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	s.Start()
	defer s.Stop()

	// The panicking job keeps getting rescheduled and the sibling keeps
	// firing: a panic never takes the scheduler down.
	waitFor(t, 3*time.Second, func() bool {
		return panicRuns.Load() >= 2 && siblingRuns.Load() >= 2
	}, "both the panicking job and its sibling to run at least twice")
}

func TestStopCancelsContextAndWaitsForInflight(t *testing.T) {
	s := newService(t)
	started := make(chan struct{})
	var ctxCancelled, returned atomic.Bool
	err := s.Register(Job{
		Name:       "longrunner",
		Interval:   time.Hour,
		RunAtStart: true,
		Run: func(ctx context.Context) {
			close(started)
			select {
			case <-ctx.Done():
				ctxCancelled.Store(true)
			case <-time.After(2 * time.Second):
				// Escape hatch so a broken Stop fails the test instead
				// of hanging it.
			}
			// Simulate in-flight work after cancellation; Stop must wait
			// for this to finish before returning.
			time.Sleep(100 * time.Millisecond)
			returned.Store(true)
		},
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	s.Start()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("job never started running")
	}

	s.Stop()

	if !ctxCancelled.Load() {
		t.Error("Stop did not cancel the context passed to Run")
	}
	if !returned.Load() {
		t.Error("Stop returned before the in-flight run finished")
	}
}

func TestRegisterAfterStart(t *testing.T) {
	s := newService(t)
	s.Start()
	defer s.Stop()

	var runs atomic.Int32
	err := s.Register(Job{
		Name:     "latecomer",
		Interval: 50 * time.Millisecond,
		Run:      func(ctx context.Context) { runs.Add(1) },
	})
	if err != nil {
		t.Fatalf("Register after Start returned error: %v", err)
	}

	waitFor(t, 3*time.Second, func() bool { return runs.Load() >= 2 },
		"job registered after Start to run at least twice")
}

func TestRegisterAfterStopErrors(t *testing.T) {
	s := newService(t)
	s.Start()
	s.Stop()

	var runs atomic.Int32
	err := s.Register(Job{
		Name:       "toolate",
		Interval:   20 * time.Millisecond,
		RunAtStart: true,
		Run:        func(ctx context.Context) { runs.Add(1) },
	})
	if err == nil {
		t.Fatal("Register after Stop returned nil error, want error")
	}
	if !strings.Contains(err.Error(), "toolate") {
		t.Errorf("Register after Stop error %q does not mention the job name", err)
	}

	// The rejected job must never actually run. RunAtStart plus a tiny
	// interval means any incorrect scheduling would fire well within
	// this window.
	time.Sleep(200 * time.Millisecond)
	if got := runs.Load(); got != 0 {
		t.Errorf("job registered after Stop ran %d time(s), want 0", got)
	}
}

func TestStatuses(t *testing.T) {
	s := newService(t)
	var runs atomic.Int32
	err := s.Register(Job{
		Name:       "worker",
		Interval:   250 * time.Millisecond,
		RunAtStart: true,
		Run:        func(ctx context.Context) { runs.Add(1) },
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	err = s.Register(Job{
		Name:     "idler",
		Interval: time.Hour,
		Run:      func(ctx context.Context) {},
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	byName := func() map[string]JobStatus {
		m := make(map[string]JobStatus)
		for _, st := range s.Statuses() {
			m[st.Name] = st
		}
		return m
	}

	statuses := byName()
	if len(statuses) != 2 {
		t.Fatalf("Statuses() returned %d entries before Start, want 2: %v", len(statuses), statuses)
	}
	for _, name := range []string{"worker", "idler"} {
		if _, ok := statuses[name]; !ok {
			t.Errorf("Statuses() missing job %q", name)
		}
	}

	s.Start()
	defer s.Stop()

	waitFor(t, 3*time.Second, func() bool {
		m := byName()
		return !m["worker"].NextRunAt.IsZero() &&
			!m["idler"].NextRunAt.IsZero() &&
			!m["worker"].LastRunAt.IsZero()
	}, "Statuses to report NextRunAt for started jobs and LastRunAt for a job that ran")
}
