package indieauth

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/persistence/configrepository"
)

const registrationTimeout = time.Minute * 10

// Service bundles the dependencies and per-instance state needed by the
// IndieAuth flows (both this Owncast server acting as a client against a
// remote IndieAuth endpoint, and this Owncast server acting as its own
// IndieAuth server). Construct once in main() and inject into HTTP
// handlers that drive auth flows.
type Service struct {
	configRepository configrepository.ConfigRepository

	pendingAuthRequestsLock sync.Mutex
	pendingAuthRequests     map[string]*Request

	pendingServerAuthRequestsLock sync.Mutex
	pendingServerAuthRequests     map[string]ServerAuthRequest
}

// Deps lists the services this Service needs.
type Deps struct {
	ConfigRepository configrepository.ConfigRepository
}

// New constructs a Service. It also starts a background goroutine that
// periodically prunes expired client- and server-side auth requests.
func New(deps Deps) *Service {
	s := &Service{
		configRepository:          deps.ConfigRepository,
		pendingAuthRequests:       make(map[string]*Request),
		pendingServerAuthRequests: make(map[string]ServerAuthRequest),
	}
	go s.runExpiredRequestPruner()
	return s
}

func (s *Service) runExpiredRequestPruner() {
	ticker := time.NewTicker(registrationTimeout)
	defer ticker.Stop()
	for now := range ticker.C {
		log.Debugln("Pruning expired IndieAuth requests.")
		s.pruneExpiredRequests(now)
	}
}

func (s *Service) pruneExpiredRequests(now time.Time) {
	s.pendingAuthRequestsLock.Lock()
	for k, v := range s.pendingAuthRequests {
		if now.Sub(v.Timestamp) > registrationTimeout {
			delete(s.pendingAuthRequests, k)
		}
	}
	s.pendingAuthRequestsLock.Unlock()

	s.pendingServerAuthRequestsLock.Lock()
	for k, v := range s.pendingServerAuthRequests {
		if now.Sub(v.Timestamp) > registrationTimeout {
			delete(s.pendingServerAuthRequests, k)
		}
	}
	s.pendingServerAuthRequestsLock.Unlock()
}
