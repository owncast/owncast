// Package inbox is the inbound side of the ActivityPub federation
// subsystem: parses verified inbound activities, dispatches them to
// per-type handlers (Follow, Like, Announce, Undo, Create, Update), and
// records accepted activities. Construct via New(Deps) and call
// Start(ctx) to spin up the request worker pool.
package inbox

import (
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	apresolvers "github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/activitypub/workerpool"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/services/webhooks"
)

// Job bundles a single inbound HTTP request for a worker.
type Job struct {
	request apmodels.InboxRequest
}

// Service owns the inbound inbox worker pool and routes verified
// activities to their per-type handlers. It composes the persistence
// service (to record accepted activities) and the followers repository
// (to track follow/unfollow state).
type Service struct {
	workerPoolSize int
	queue          chan Job

	persistence      *persistence.Service
	workerpool       *workerpool.Service
	webhooks         *webhooks.Service
	chat             *chat.Service
	followers        followersrepository.FollowersRepository
	configRepository configrepository.ConfigRepository
	builder          *apmodels.Builder
	signer           *apcrypto.Signer
	resolver         *apresolvers.Resolver
}

// Deps is the explicit dependency contract for inbox.
type Deps struct {
	Persistence      *persistence.Service
	Workerpool       *workerpool.Service
	Webhooks         *webhooks.Service
	Chat             *chat.Service
	Followers        followersrepository.FollowersRepository
	ConfigRepository configrepository.ConfigRepository
	Builder          *apmodels.Builder
	Signer           *apcrypto.Signer
	Resolver         *apresolvers.Resolver
}

// New constructs an idle inbox Service. Call Start to launch the worker
// pool.
func New(deps Deps) *Service {
	return &Service{
		workerPoolSize:   runtime.GOMAXPROCS(0),
		persistence:      deps.Persistence,
		workerpool:       deps.Workerpool,
		webhooks:         deps.Webhooks,
		chat:             deps.Chat,
		followers:        deps.Followers,
		configRepository: deps.ConfigRepository,
		builder:          deps.Builder,
		signer:           deps.Signer,
		resolver:         deps.Resolver,
	}
}

// Start launches the inbox worker pool. Safe to call once.
func (s *Service) Start() {
	s.queue = make(chan Job)
	for i := 1; i <= s.workerPoolSize; i++ {
		go s.worker(i)
	}
}

// AddToQueue queues an inbound request for the worker pool.
func (s *Service) AddToQueue(req apmodels.InboxRequest) {
	log.Tracef("Queued request for ActivityPub inbox handler")
	s.queue <- Job{req}
}

func (s *Service) worker(workerID int) {
	log.Debugf("Started ActivityPub worker %d", workerID)

	for job := range s.queue {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("Recovered from panic in ActivityPub worker %d: %v", workerID, r)
				}
			}()
			s.handle(job.request)
		}()

		log.Tracef("Done with ActivityPub inbox handler using worker %d", workerID)
	}
}
