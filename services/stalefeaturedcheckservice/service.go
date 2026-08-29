package stalefeaturedcheckservice

import (
	"sync"
	"time"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/owncast/owncast/services/activitypub/outbox"
	log "github.com/sirupsen/logrus"
)

const (
	// staleThreshold is the duration after which a server is considered offline
	// if no status update has been received. A live server re-pings its
	// followers every outbox.StreamPingInterval, so this allows two consecutive
	// missed pings (plus a minute of grace for delivery jitter) before marking
	// it offline.
	staleThreshold = 2*outbox.StreamPingInterval + time.Minute

	// checkInterval is how often we check for stale servers. Kept well below
	// staleThreshold so a server is marked offline promptly after it crosses
	// the threshold rather than up to a full check cycle later.
	checkInterval = 1 * time.Minute
)

// Checker periodically marks stale federated servers as offline.
// Dependencies are injected so the logic is testable without globals.
type Checker struct {
	config  configrepository.ConfigRepository
	servers federatedserversrepository.FederatedServersRepository
	now     func() time.Time
	// ticks returns a tick channel plus a func that releases its resources.
	ticks func() (<-chan time.Time, func())

	mu         sync.Mutex
	done       chan bool
	stopTicker func()
	// stopped is closed when the background goroutine exits; used by tests to
	// prove there is no goroutine leak.
	stopped chan struct{}
}

// New returns a Checker using the given dependencies. A nil now defaults to
// time.Now; a nil ticks defaults to a real time.Ticker at checkInterval.
func New(
	config configrepository.ConfigRepository,
	servers federatedserversrepository.FederatedServersRepository,
	now func() time.Time,
	ticks func() (<-chan time.Time, func()),
) *Checker {
	if now == nil {
		now = time.Now
	}
	if ticks == nil {
		ticks = func() (<-chan time.Time, func()) {
			t := time.NewTicker(checkInterval)
			return t.C, t.Stop
		}
	}
	return &Checker{config: config, servers: servers, now: now, ticks: ticks}
}

// Start begins checking for stale federated servers in the background.
func (c *Checker) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.config.GetFederationEnabled() {
		return
	}

	// Don't start if already running
	if c.done != nil {
		log.Debugln("Stale featured server checker already running")
		return
	}

	tickCh, stopTicker := c.ticks()
	c.stopTicker = stopTicker
	c.done = make(chan bool)
	c.stopped = make(chan struct{})

	// Capture in locals so a later Start/Stop cycle can't race this goroutine.
	done := c.done
	stopped := c.stopped

	go func() {
		defer close(stopped)

		// Run immediately on start
		c.check()

		for {
			select {
			case <-tickCh:
				c.check()
			case <-done:
				return
			}
		}
	}()

	log.Infof("Started stale featured server checker (%s interval, %s offline threshold)", checkInterval, staleThreshold)
}

// Stop halts the stale server checker if it is running.
func (c *Checker) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.done == nil {
		return
	}

	c.stopTicker()
	close(c.done)
	c.done = nil
	c.stopTicker = nil
	log.Infoln("Stopped stale featured server checker")
}

// check inspects all online federated servers and marks them as offline
// if they haven't sent a status update within the stale threshold.
func (c *Checker) check() {
	servers, err := c.servers.GetFederatedServers()
	if err != nil {
		log.Errorf("Failed to get federated servers for staleness check: %v", err)
		return
	}

	now := c.now()
	markedOfflineCount := 0

	for _, server := range servers {
		// Only check servers that are currently marked as online
		if !server.IsOnline {
			continue
		}

		// Skip if no last status update (shouldn't happen, but be safe)
		if server.LastStatusUpdate == nil {
			continue
		}

		timeSinceLastUpdate := now.Sub(*server.LastStatusUpdate)

		if timeSinceLastUpdate > staleThreshold {
			log.Infof("Marking federated server %s as offline due to staleness (%v since last update)",
				server.IRI, timeSinceLastUpdate)

			err := c.servers.UpdateServerStatus(server.IRI, false, nil)
			if err != nil {
				log.Errorf("Failed to mark server %s as offline: %v", server.IRI, err)
			} else {
				markedOfflineCount++
			}
		}
	}

	if markedOfflineCount > 0 {
		log.Infof("Marked %d federated server(s) as offline due to staleness", markedOfflineCount)
	}
}

var (
	defaultMu      sync.Mutex
	defaultChecker *Checker
)

// Start begins the default background checker using the global repositories.
func Start() {
	defaultMu.Lock()
	defer defaultMu.Unlock()

	if defaultChecker == nil {
		defaultChecker = New(configrepository.Get(), federatedserversrepository.Get(), nil, nil)
	}
	defaultChecker.Start()
}

// Stop halts the default background checker if it is running.
func Stop() {
	defaultMu.Lock()
	defer defaultMu.Unlock()

	if defaultChecker != nil {
		defaultChecker.Stop()
	}
}
