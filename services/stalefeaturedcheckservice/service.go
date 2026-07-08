// Package stalefeaturedcheckservice marks featured federated servers
// offline when they stop sending status updates, keeping the
// featured-streams directory honest. It runs as a job on the central
// scheduler.
package stalefeaturedcheckservice

import (
	"context"
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

	// CheckInterval is how often Run should execute, registered on the
	// central scheduler by main.go. Kept well below staleThreshold so a
	// server is marked offline promptly after it crosses the threshold
	// rather than up to a full check cycle later.
	CheckInterval = 1 * time.Minute
)

// Service sweeps featured federated servers for staleness. Construct with
// New in main.go and register Run on the scheduler.
type Service struct {
	configRepository           configrepository.ConfigRepository
	federatedServersRepository federatedserversrepository.FederatedServersRepository
}

// Deps lists everything the service consumes.
type Deps struct {
	ConfigRepository           configrepository.ConfigRepository
	FederatedServersRepository federatedserversrepository.FederatedServersRepository
}

// New constructs the stale featured server checker.
func New(deps Deps) *Service {
	return &Service{
		configRepository:           deps.ConfigRepository,
		federatedServersRepository: deps.FederatedServersRepository,
	}
}

// Run is the scheduler job body. It no-ops while federation is disabled,
// which also means enabling federation at runtime starts sweeping on the
// next tick instead of requiring a restart like the old self-gating ticker
// did.
func (s *Service) Run(ctx context.Context) {
	if !s.configRepository.GetFederationEnabled() {
		return
	}

	servers, err := s.federatedServersRepository.GetFederatedServers()
	if err != nil {
		log.Errorf("Failed to get federated servers for staleness check: %v", err)
		return
	}

	now := time.Now()
	markedOfflineCount := 0

	for _, server := range servers {
		// Shutdown cancels the context; stop sweeping promptly.
		if ctx.Err() != nil {
			return
		}

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

			err := s.federatedServersRepository.UpdateServerStatus(server.IRI, false, nil)
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
