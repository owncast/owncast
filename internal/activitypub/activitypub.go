package activitypub

import (
	"fmt"

	"github.com/owncast/owncast/internal/activitypub/apmodels"
	"github.com/owncast/owncast/internal/activitypub/controllers"
	"github.com/owncast/owncast/internal/activitypub/crypto"
	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/internal/activitypub/inbox"
	"github.com/owncast/owncast/internal/activitypub/outbox"
	"github.com/owncast/owncast/internal/activitypub/persistence"
	"github.com/owncast/owncast/internal/activitypub/workerpool"
	"github.com/owncast/owncast/models"
)

// New will initialize and start the ActivityPub service
func New(p *persistence.Service) (s *Service, err error) {
	s = &Service{}
	s.Persistence = p

	if s.Crypto, err = crypto.New(p.Data); err != nil {
		return nil, fmt.Errorf("initializing crypto sub-service: %v", err)
	}

	if s.Models, err = apmodels.New(p.Data, s.Crypto); err != nil {
		return nil, fmt.Errorf("initializing activity pub models sub-service: %v", err)
	}

	if s.Controller, err = controllers.New(p, s.Models, s.Follower); err != nil {
		return nil, fmt.Errorf("initializing activity pub controllers: %v", err)
	}

	if s.Inbox, err = inbox.New(p); err != nil {
		return nil, fmt.Errorf("initializing inbox sub-service: %v", err)
	}

	if s.WorkerPool, err = workerpool.New(); err != nil {
		return nil, fmt.Errorf("initializing worker pool sub-service: %v", err)
	}

	if s.Follower, err = follower.New(p.Data, s.Crypto, s.Models, s.WorkerPool); err != nil {
		return nil, fmt.Errorf("initializing follower sub-service: %v", err)
	}

	if s.Outbox, err = outbox.New(s.Persistence, s.Crypto, s.Models, s.Follower, s.WorkerPool); err != nil {
		return nil, fmt.Errorf("initializing outbox sub-service: %v", err)
	}

	s.WorkerPool.InitOutboundWorkerPool()
	s.Inbox.InitInboxWorkerPool()

	// Generate the keys for signing federated activity if needed.
	if p.Data.GetPrivateKey() == "" {
		privateKey, publicKey, err := s.Crypto.GenerateKeys()
		_ = p.Data.SetPrivateKey(string(privateKey))
		_ = p.Data.SetPublicKey(string(publicKey))
		if err != nil {
			return nil, fmt.Errorf("getting private key: %v", err)
		}
	}

	return s, nil
}

type Service struct {
	Controller *controllers.Service
	subservices
}

type subservices struct {
	Crypto      *crypto.Service
	Models      *apmodels.Service
	Follower    *follower.Service
	Inbox       *inbox.Service
	Outbox      *outbox.Service
	Persistence *persistence.Service
	WorkerPool  *workerpool.Service
}

// SendLive will send a "Go Live" message to followers.
func (s *Service) SendLive() error {
	return s.Outbox.SendLive()
}

// SendPublicFederatedMessage will send an arbitrary provided message to followers.
func (s *Service) SendPublicFederatedMessage(message string) error {
	return s.Outbox.SendPublicMessage(message, s.Persistence)
}

// SendDirectFederatedMessage will send a direct message to a single account.
func (s *Service) SendDirectFederatedMessage(message, account string) error {
	return s.Outbox.SendDirectMessageToAccount(message, account)
}

// GetFollowerCount will return the local tracked follower count.
func (s *Service) GetFollowerCount() (int64, error) {
	return s.Follower.GetFollowerCount()
}

// GetPendingFollowRequests will return the pending follow requests.
func (s *Service) GetPendingFollowRequests() ([]models.Follower, error) {
	return s.Follower.GetPendingFollowRequests()
}
