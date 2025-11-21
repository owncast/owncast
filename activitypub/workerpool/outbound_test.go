package workerpool

import (
	"testing"
	"time"
)

// resetCircuitBreakerForTesting clears all failed domain state for testing purposes.
// This function provides test isolation by resetting the global circuit breaker state.
func resetCircuitBreakerForTesting() {
	failedDomainsMutex.Lock()
	defer failedDomainsMutex.Unlock()
	failedDomains = make(map[string]*domainFailure)
}

func TestCircuitBreaker(t *testing.T) {
	// Ensure clean state before test
	resetCircuitBreakerForTesting()
	defer resetCircuitBreakerForTesting() // Clean up after test

	testDomain := "failing.example.com"

	// Initially, domain should not be skipped
	if shouldSkipDomain(testDomain) {
		t.Error("Domain should not be skipped initially")
	}

	// Record failures
	recordDomainFailure(testDomain)
	recordDomainFailure(testDomain)
	recordDomainFailure(testDomain)

	// Domain should now be skipped
	if !shouldSkipDomain(testDomain) {
		t.Error("Domain should be skipped after failures")
	}

	// After successful delivery, domain should be reset
	resetDomainFailure(testDomain)
	if shouldSkipDomain(testDomain) {
		t.Error("Domain should not be skipped after reset")
	}
}

func TestHTTPTimeouts(t *testing.T) {
	// Ensure clean state before test
	resetCircuitBreakerForTesting()
	defer resetCircuitBreakerForTesting() // Clean up after test

	// Initialize HTTP client
	InitOutboundWorkerPool(1)

	if httpClient == nil {
		t.Error("HTTP client should be initialized")
	}

	if httpClient.Timeout != 15*time.Second {
		t.Error("HTTP client should have 15 second timeout")
	}
}

func TestWorkerPoolSizing(t *testing.T) {
	// Ensure clean state before test
	resetCircuitBreakerForTesting()
	defer resetCircuitBreakerForTesting() // Clean up after test

	// Test the queue sizing
	InitOutboundWorkerPool(5)

	if cap(queue) != 5 {
		t.Errorf("Queue capacity should be 5, got %d", cap(queue))
	}
}

func TestBackoffDurations(t *testing.T) {
	// Test that backoff durations are properly configured
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
	// Test that multiple tests don't interfere with each other
	resetCircuitBreakerForTesting()
	defer resetCircuitBreakerForTesting()

	domain1 := "test1.example.com"
	domain2 := "test2.example.com"

	// Neither domain should be blocked initially
	if shouldSkipDomain(domain1) || shouldSkipDomain(domain2) {
		t.Error("Domains should not be blocked initially")
	}

	// Record failures for domain1 only
	recordDomainFailure(domain1)
	recordDomainFailure(domain1)
	recordDomainFailure(domain1)

	// Only domain1 should be blocked
	if !shouldSkipDomain(domain1) {
		t.Error("Domain1 should be blocked after failures")
	}
	if shouldSkipDomain(domain2) {
		t.Error("Domain2 should not be blocked")
	}

	// Reset and verify clean state
	resetCircuitBreakerForTesting()
	if shouldSkipDomain(domain1) || shouldSkipDomain(domain2) {
		t.Error("Both domains should be unblocked after reset")
	}
}
