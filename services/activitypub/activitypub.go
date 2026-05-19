// Package activitypub is the composition root for the federation
// subsystem. It wires together the persistence, workerpool, outbox,
// inbox, controllers, and jobs sub-services, exposes a small
// stream/admin-facing API (SendLive, SendPublic, GetFollowerCount, …),
// and hands the controllers set to the router.
package activitypub

import (
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub/controllers"
	"github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/inbox"
	"github.com/owncast/owncast/services/activitypub/jobs"
	"github.com/owncast/owncast/services/activitypub/outbox"
	"github.com/owncast/owncast/services/activitypub/persistence"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/activitypub/workerpool"
)

// Service is the composed federation subsystem. Construct one in
// main.go with New(Deps) and call Start to spin up the worker pools and
// background jobs. Hand .Controllers() to the router for HTTP routing.
type Service struct {
	persistence *persistence.Service
	workerpool  *workerpool.Service
	outbox      *outbox.Service
	inbox       *inbox.Service
	jobs        *jobs.Service
	controllers *controllers.Controllers
	followers   followersrepository.FollowersRepository
}

// Deps lists the explicit construction inputs for the federation
// subsystem.
type Deps struct {
	Datastore *data.Datastore
}

// New constructs the federation subsystem in dependency order. It does
// not spawn any goroutines or bind ports; call Start for that.
func New(deps Deps) *Service {
	persistenceSvc := persistence.New(deps.Datastore)
	followers := followersrepository.New(deps.Datastore)

	wpSvc := workerpool.New(outboundWorkerPoolSize(followers))

	outboxSvc := outbox.New(outbox.Deps{
		Persistence: persistenceSvc,
		Workerpool:  wpSvc,
		Followers:   followers,
	})

	inboxSvc := inbox.New(inbox.Deps{
		Persistence: persistenceSvc,
		Workerpool:  wpSvc,
		Followers:   followers,
	})

	jobsSvc := jobs.New(jobs.Deps{
		Followers: followers,
	})

	ctrls := controllers.New(controllers.Deps{
		Persistence: persistenceSvc,
		Outbox:      outboxSvc,
		Inbox:       inboxSvc,
		Followers:   followers,
	})

	return &Service{
		persistence: persistenceSvc,
		workerpool:  wpSvc,
		outbox:      outboxSvc,
		inbox:       inboxSvc,
		jobs:        jobsSvc,
		controllers: ctrls,
		followers:   followers,
	}
}

// Start brings up the worker pools and recurring jobs. Also generates
// the signing keypair on first run.
func (s *Service) Start() {
	configRepository := configrepository.Get()

	s.workerpool.Start()
	s.inbox.Start()

	// Generate the keys for signing federated activity if needed.
	if configRepository.GetPrivateKey() == "" {
		privateKey, publicKey, err := crypto.GenerateKeys()
		_ = configRepository.SetPrivateKey(string(privateKey))
		_ = configRepository.SetPublicKey(string(publicKey))
		if err != nil {
			log.Errorln("Unable to get private key", err)
		}
	}

	s.jobs.Start()
}

// Controllers returns the HTTP controller set so the router can mount
// federation routes against it.
func (s *Service) Controllers() *controllers.Controllers {
	return s.controllers
}

// Outbox returns the outbox service. Callers that send federated
// messages on behalf of the stream/admin layer go through this.
func (s *Service) Outbox() *outbox.Service {
	return s.outbox
}

// SendLive sends a "Go Live" message to followers.
func (s *Service) SendLive() error {
	return s.outbox.SendLive()
}

// SendPublicFederatedMessage sends an arbitrary message to all followers.
func (s *Service) SendPublicFederatedMessage(message string) error {
	return s.outbox.SendPublicMessage(message)
}

// SendDirectFederatedMessage sends a direct message to a single account.
func (s *Service) SendDirectFederatedMessage(message, account string) error {
	return s.outbox.SendDirectMessageToAccount(message, account)
}

// GetFollowerCount returns the local tracked follower count.
func (s *Service) GetFollowerCount() (int64, error) {
	return s.followers.GetCount()
}

// GetPendingFollowRequests returns the pending follow requests.
func (s *Service) GetPendingFollowRequests() ([]models.Follower, error) {
	return s.followers.GetPendingFollowRequests()
}

// UpdateFollowersWithAccountUpdates broadcasts a profile-update activity to followers.
func (s *Service) UpdateFollowersWithAccountUpdates() error {
	return s.outbox.UpdateFollowersWithAccountUpdates()
}

// GetInboundActivities returns saved inbound federated activities (paginated).
func (s *Service) GetInboundActivities(limit int, offset int) ([]models.FederatedActivity, int, error) {
	return s.persistence.GetInboundActivities(limit, offset)
}

// Followers returns the underlying followers repository for callers
// that need direct repo access (admin handlers paginating follower
// lists, etc.). Prefer the typed methods above where they exist.
func (s *Service) Followers() followersrepository.FollowersRepository {
	return s.followers
}

// Workerpool returns the outbound delivery worker pool. Exposed for the
// activitypub/requests helper functions that build and queue signed
// outbound requests directly.
func (s *Service) Workerpool() *workerpool.Service {
	return s.workerpool
}

// outboundWorkerPoolSize sizes the outbound delivery pool from the
// current follower count: base workers + 1 per 100 followers, clamped.
// Prevents excessive resource usage on instances with many followers.
func outboundWorkerPoolSize(followers followersrepository.FollowersRepository) int {
	const (
		minWorkers     = 10
		maxWorkers     = 50
		defaultWorkers = 20
	)

	fc, err := followers.GetCount()
	if err != nil {
		log.Errorln("Unable to get follower count", err)
		return defaultWorkers
	}

	workers := minWorkers + int(fc/100)
	if workers > maxWorkers {
		workers = maxWorkers
	}

	log.Debugf("Initializing ActivityPub outbound worker pool with %d workers for %d followers", workers, fc)
	return workers
}
