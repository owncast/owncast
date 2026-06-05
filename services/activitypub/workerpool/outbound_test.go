package workerpool

import (
	"testing"
	"time"
)

func newTestService() *Service {
	return New(1)
}

func TestCircuitBreaker(t *testing.T) {
	s := newTestService()

	testDomain := "failing.example.com"

	// Initially, domain should not be skipped.
	if s.ShouldSkipDomain(testDomain) {
		t.Error("Domain should not be skipped initially")
	}

	// Record failures.
	s.recordDomainFailure(testDomain)
	s.recordDomainFailure(testDomain)
	s.recordDomainFailure(testDomain)

	// Domain should now be skipped.
	if !s.ShouldSkipDomain(testDomain) {
		t.Error("Domain should be skipped after failures")
	}

	// After successful delivery, domain should be reset.
	s.resetDomainFailure(testDomain)
	if s.ShouldSkipDomain(testDomain) {
		t.Error("Domain should not be skipped after reset")
	}
}

func TestHTTPTimeouts(t *testing.T) {
	s := New(1)
	s.Start()

	if s.httpClient == nil {
		t.Error("HTTP client should be initialized")
	}

	if s.httpClient.Timeout != 8*time.Second {
		t.Errorf("HTTP client should have 8 second timeout, got %v", s.httpClient.Timeout)
	}
}

func TestWorkerPoolSizing(t *testing.T) {
	// Queue buffer should be at least the 500-item minimum even for
	// small worker pools.
	s := New(5)
	s.Start()
	if cap(s.queue) < 500 {
		t.Errorf("Queue capacity should be at least 500, got %d", cap(s.queue))
	}

	// Larger worker pools get proportionally larger buffers.
	s2 := New(100)
	s2.Start()
	if cap(s2.queue) != 1000 {
		t.Errorf("Queue capacity should be 1000 for 100 workers, got %d", cap(s2.queue))
	}
}

func TestBackoffDurations(t *testing.T) {
	expectedDurations := []time.Duration{
		1 * time.Minute,
		5 * time.Minute,
		15 * time.Minute,
		30 * time.Minute,
		60 * time.Minute,
	}

	if len(circuitBreakerBackoffDurations) != len(expectedDurations) {
		t.Errorf("Expected %d backoff durations, got %d", len(expectedDurations), len(circuitBreakerBackoffDurations))
	}

	for i, expected := range expectedDurations {
		if circuitBreakerBackoffDurations[i] != expected {
			t.Errorf("Backoff duration at index %d: expected %v, got %v", i, expected, circuitBreakerBackoffDurations[i])
		}
	}
}

func TestCircuitBreakerIsolation(t *testing.T) {
	s := newTestService()

	domain1 := "test1.example.com"
	domain2 := "test2.example.com"

	// Neither domain should be blocked initially.
	if s.ShouldSkipDomain(domain1) || s.ShouldSkipDomain(domain2) {
		t.Error("Domains should not be blocked initially")
	}

	// Record failures for domain1 only.
	s.recordDomainFailure(domain1)
	s.recordDomainFailure(domain1)
	s.recordDomainFailure(domain1)

	// Only domain1 should be blocked.
	if !s.ShouldSkipDomain(domain1) {
		t.Error("Domain1 should be blocked after failures")
	}
	if s.ShouldSkipDomain(domain2) {
		t.Error("Domain2 should not be blocked")
	}
}
