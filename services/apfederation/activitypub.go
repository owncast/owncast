package apfederation

import (
	"github.com/owncast/owncast/services/apfederation/crypto"
	"github.com/owncast/owncast/services/apfederation/outbox"

	"github.com/owncast/owncast/services/apfederation/workerpool"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/data"
	"github.com/owncast/owncast/storage/federationrepository"

	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

type APFederation struct {
	workers *workerpool.WorkerPool
	outbox  *outbox.APOutbox
}

func New() *APFederation {
	ds := data.GetDatastore()
	apf := &APFederation{
		outbox: outbox.Get(),
	}
	apf.Start(ds)
	return apf
}

var temporaryGlobalInstance *APFederation

func Get() *APFederation {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}
	return temporaryGlobalInstance
}

// Start will initialize and start the federation support.
func (ap *APFederation) Start(datastore *data.Store) {
	configRepository := configrepository.Get()

	// workerpool.InitOutboundWorkerPool()
	// ap.InitInboxWorkerPool()

	// Generate the keys for signing federated activity if needed.
	if configRepository.GetPrivateKey() == "" {
		privateKey, publicKey, err := crypto.GenerateKeys()
		_ = configRepository.SetPrivateKey(string(privateKey))
		_ = configRepository.SetPublicKey(string(publicKey))
		if err != nil {
			log.Errorln("Unable to get private key", err)
		}
	}
}

// SendLive will send a "Go Live" message to followers.
func (ap *APFederation) SendLive() error {
	return ap.SendLive()
}

// SendPublicFederatedMessage will send an arbitrary provided message to followers.
func (ap *APFederation) SendPublicFederatedMessage(message string) error {
	return ap.outbox.SendPublicMessage(message)
}

// SendDirectFederatedMessage will send a direct message to a single account.
func (ap *APFederation) SendDirectFederatedMessage(message, account string) error {
	return ap.outbox.SendDirectMessageToAccount(message, account)
}

// GetFollowerCount will return the local tracked follower count.
func (ap *APFederation) GetFollowerCount() (int64, error) {
	federationRepository := federationrepository.Get()
	return federationRepository.GetFollowerCount()
}

// GetPendingFollowRequests will return the pending follow requests.
func (ap *APFederation) GetPendingFollowRequests() ([]models.Follower, error) {
	federationRepository := federationrepository.Get()
	return federationRepository.GetPendingFollowRequests()
}
