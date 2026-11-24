package indieauth

import (
	"sync"
	"testing"
)

func TestRandStringFormatAndUniqueness(t *testing.T) {
	const workers, perWorker = 8, 2000
	results := make(chan string, workers*perWorker)

	// Concurrent callers (simultaneous logins) must be safe and never collide.
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perWorker {
				results <- randString(20)
			}
		}()
	}
	wg.Wait()
	close(results)

	seen := make(map[string]struct{}, workers*perWorker)
	for s := range results {
		if _, dup := seen[s]; dup {
			t.Fatalf("duplicate value %q", s)
		}
		seen[s] = struct{}{}
	}

	for _, n := range []int{1, 20, 50} {
		s := randString(n)
		if len(s) != n {
			t.Fatalf("randString(%d) returned %d characters", n, len(s))
		}
		for _, c := range s {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
				t.Fatalf("randString(%d) = %q contains non-letter %q", n, s, c)
			}
		}
	}
}
