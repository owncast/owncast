package chat

import (
	"sync"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"golang.org/x/time/rate"
)

type fakeSpamConfig struct {
	fakeJoinPartConfig
	spamProtection bool
}

func (f fakeSpamConfig) GetChatSpamProtectionEnabled() bool {
	return f.spamProtection
}

// TestClient_TimeoutTimerRace ensures that concurrent access to chat rejection
// timeout methods (startChatRejectionTimeout, passesRateLimit, isInTimeout, close)
// is data-race free under Go's race detector.
func TestClient_TimeoutTimerRace(t *testing.T) {
	chatSvc := &Service{
		configRepository: fakeSpamConfig{spamProtection: true},
		unregister:       make(chan uint, 100),
		rejectionTimeout: 1 * time.Millisecond,
	}

	client := &Client{
		Id:          1,
		User:        &models.User{ID: "test-user", DisplayName: "Test User"},
		server:      chatSvc,
		send:        make(chan []byte, 1000),
		rateLimiter: rate.NewLimiter(rate.Inf, 0),
	}

	const iterations = 500
	var wg sync.WaitGroup

	// Goroutine 1: repeatedly triggers chat rejection timeout
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			client.startChatRejectionTimeout()
		}
	}()

	// Goroutine 2: repeatedly checks rate limit which inspects inTimeout
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = client.passesRateLimit()
		}
	}()

	// Goroutine 3: repeatedly checks isInTimeout directly
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = client.isInTimeout()
		}
	}()

	wg.Wait()
}

// TestClient_PassesRateLimit validates the rate limiting rules across
// different user scopes, spam protection settings, and timeout states.
func TestClient_PassesRateLimit(t *testing.T) {
	for _, tc := range []struct {
		name             string
		spamProtection   bool
		isModerator      bool
		inTimeout        bool
		rateLimiterAllow bool
		wantPass         bool
	}{
		{
			name:             "spam protection disabled allows message",
			spamProtection:   false,
			isModerator:      false,
			inTimeout:        true,
			rateLimiterAllow: false,
			wantPass:         true,
		},
		{
			name:             "moderator bypasses rate limits and timeouts",
			spamProtection:   true,
			isModerator:      true,
			inTimeout:        true,
			rateLimiterAllow: false,
			wantPass:         true,
		},
		{
			name:             "normal user within rate limit passes",
			spamProtection:   true,
			isModerator:      false,
			inTimeout:        false,
			rateLimiterAllow: true,
			wantPass:         true,
		},
		{
			name:             "normal user in timeout is rejected",
			spamProtection:   true,
			isModerator:      false,
			inTimeout:        true,
			rateLimiterAllow: true,
			wantPass:         false,
		},
		{
			name:             "normal user exceeding rate limiter is rejected",
			spamProtection:   true,
			isModerator:      false,
			inTimeout:        false,
			rateLimiterAllow: false,
			wantPass:         false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chatSvc := &Service{
				configRepository: fakeSpamConfig{spamProtection: tc.spamProtection},
			}
			user := &models.User{
				ID: "test-user",
			}
			if tc.isModerator {
				user.Scopes = []string{"MODERATOR"}
			}

			limiterLimit := rate.Inf
			if !tc.rateLimiterAllow {
				limiterLimit = rate.Limit(0)
			}
			client := &Client{
				server:      chatSvc,
				User:        user,
				rateLimiter: rate.NewLimiter(limiterLimit, 0),
			}
			if tc.inTimeout {
				client.timeoutTimer = time.NewTimer(time.Hour)
				t.Cleanup(func() { client.timeoutTimer.Stop() })
			}

			got := client.passesRateLimit()
			if got != tc.wantPass {
				t.Fatalf("passesRateLimit() = %v, want %v", got, tc.wantPass)
			}
		})
	}
}

func waitForCondition(timeout time.Duration, condition func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(500 * time.Microsecond)
	}
	return condition()
}

// TestClient_TimeoutLifecycle tests client timeout transitions, including
// triggering, duplicate suppression, background expiration, and cancellation.
func TestClient_TimeoutLifecycle(t *testing.T) {
	for _, tc := range []struct {
		name           string
		initialTimeout time.Duration
		action         string
		wantInTimeout  bool
		wantCanSend    bool
	}{
		{
			name:           "initial start triggers timeout",
			initialTimeout: 50 * time.Millisecond,
			action:         "none",
			wantInTimeout:  true,
			wantCanSend:    false,
		},
		{
			name:           "duplicate trigger while active is a no-op",
			initialTimeout: 50 * time.Millisecond,
			action:         "re-trigger",
			wantInTimeout:  true,
			wantCanSend:    false,
		},
		{
			name:           "timer expiration clears timeout and restores sending",
			initialTimeout: 10 * time.Millisecond,
			action:         "wait-expire",
			wantInTimeout:  false,
			wantCanSend:    true,
		},
		{
			name:           "close cancels active timeout and resets timeout status",
			initialTimeout: 50 * time.Millisecond,
			action:         "close",
			wantInTimeout:  false,
			wantCanSend:    true,
		},
		{
			name:           "can re-trigger timeout after previous expired",
			initialTimeout: 10 * time.Millisecond,
			action:         "re-trigger-after-expire",
			wantInTimeout:  true,
			wantCanSend:    false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chatSvc := &Service{
				configRepository: fakeSpamConfig{spamProtection: true},
				unregister:       make(chan uint, 1),
				rejectionTimeout: tc.initialTimeout,
			}
			client := &Client{
				Id:          1,
				User:        &models.User{ID: "test-user"},
				server:      chatSvc,
				send:        make(chan []byte, 10),
				rateLimiter: rate.NewLimiter(rate.Inf, 0),
			}

			client.startChatRejectionTimeout()

			switch tc.action {
			case "re-trigger":
				client.startChatRejectionTimeout()
			case "wait-expire":
				if !waitForCondition(500*time.Millisecond, func() bool { return !client.isInTimeout() }) {
					t.Fatal("timed out waiting for rejection timeout to expire")
				}
			case "close":
				client.close()
			case "re-trigger-after-expire":
				if !waitForCondition(500*time.Millisecond, func() bool { return !client.isInTimeout() }) {
					t.Fatal("timed out waiting for rejection timeout to expire")
				}
				client.startChatRejectionTimeout()
			}

			if got := client.isInTimeout(); got != tc.wantInTimeout {
				t.Fatalf("client.isInTimeout() = %v, want %v", got, tc.wantInTimeout)
			}
			if got := client.passesRateLimit(); got != tc.wantCanSend {
				t.Fatalf("client.passesRateLimit() = %v, want %v", got, tc.wantCanSend)
			}
		})
	}
}

// TestService_RejectionTimeoutResolution validates fallback and override behavior
// for rejection timeout duration configured via chat.Deps.
func TestService_RejectionTimeoutResolution(t *testing.T) {
	for _, tc := range []struct {
		name        string
		depsTimeout time.Duration
		wantTimeout time.Duration
	}{
		{
			name:        "default fallback when unset",
			depsTimeout: 0,
			wantTimeout: defaultChatRejectionTimeout,
		},
		{
			name:        "service override via Deps",
			depsTimeout: 75 * time.Millisecond,
			wantTimeout: 75 * time.Millisecond,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := New(Deps{RejectionTimeout: tc.depsTimeout})
			if svc.rejectionTimeout != tc.wantTimeout {
				t.Fatalf("svc.rejectionTimeout = %v, want %v", svc.rejectionTimeout, tc.wantTimeout)
			}
		})
	}

	t.Run("standalone client without server uses default timeout", func(t *testing.T) {
		client := &Client{
			send: make(chan []byte, 1),
		}
		client.startChatRejectionTimeout()
		defer client.close()

		if !client.isInTimeout() {
			t.Fatal("expected standalone client to be in timeout")
		}
	})
}
