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

var (
	stalenessChecker      *time.Ticker
	stalenessCheckerDone  chan bool
	stalenessCheckerMutex sync.Mutex
)

// Start begins checking for stale federated servers in the background.
func Start() {
	stalenessCheckerMutex.Lock()
	defer stalenessCheckerMutex.Unlock()

	configRepository := configrepository.Get()
	if !configRepository.GetFederationEnabled() {
		return
	}

	// Don't start if already running
	if stalenessChecker != nil {
		log.Debugln("Stale featured server checker already running")
		return
	}

	stalenessChecker = time.NewTicker(checkInterval)
	stalenessCheckerDone = make(chan bool)

	// Capture the done channel in a local variable to avoid race conditions
	done := stalenessCheckerDone
	ticker := stalenessChecker

	go func() {
		// Run immediately on start
		checkAndMarkStaleServers()

		for {
			select {
			case <-ticker.C:
				checkAndMarkStaleServers()
			case <-done:
				return
			}
		}
	}()

	log.Infof("Started stale featured server checker (%s interval, %s offline threshold)", checkInterval, staleThreshold)
}

// Stop halts the stale server checker if it is running.
func Stop() {
	stalenessCheckerMutex.Lock()
	defer stalenessCheckerMutex.Unlock()

	if stalenessChecker != nil {
		stalenessChecker.Stop()
		close(stalenessCheckerDone)
		stalenessChecker = nil
		stalenessCheckerDone = nil
		log.Infoln("Stopped stale featured server checker")
	}
}

// checkAndMarkStaleServers checks all online federated servers and marks them as offline
// if they haven't sent a status update within the stale threshold.
func checkAndMarkStaleServers() {
	repo := federatedserversrepository.Get()

	servers, err := repo.GetFederatedServers()
	if err != nil {
		log.Errorf("Failed to get federated servers for staleness check: %v", err)
		return
	}

	now := time.Now()
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

			err := repo.UpdateServerStatus(server.IRI, false, nil)
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
