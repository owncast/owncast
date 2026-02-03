package activitypub

import (
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/inbox"
	"github.com/owncast/owncast/activitypub/jobs"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/activitypub/workerpool"
	"github.com/owncast/owncast/persistence/configrepository"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// Start will initialize and start the federation support.
func Start(datastore *data.Datastore) {
	configRepository := configrepository.Get()
	persistence.Setup(datastore)

	outboundWorkerPoolSize := getOutboundWorkerPoolSize()
	workerpool.InitOutboundWorkerPool(outboundWorkerPoolSize)
	inbox.InitInboxWorkerPool()

	// Generate the keys for signing federated activity if needed.
	if configRepository.GetPrivateKey() == "" {
		privateKey, publicKey, err := crypto.GenerateKeys()
		_ = configRepository.SetPrivateKey(string(privateKey))
		_ = configRepository.SetPublicKey(string(publicKey))
		if err != nil {
			log.Errorln("Unable to get private key", err)
		}
	}

	// Start follower validation background job.
	jobs.StartFollowerValidationJob()
}

func getOutboundWorkerPoolSize() int {
	// Use a reasonable fixed worker pool size instead of scaling with followers
	// This prevents excessive resource usage when streamers have many followers
	const (
		minWorkers     = 10 // Minimum workers for small instances
		maxWorkers     = 50 // Maximum workers to prevent resource exhaustion
		defaultWorkers = 20 // Default for most instances
	)

	followersRepo := followersrepository.Get()
	var followerCount int64
	fc, err := followersRepo.GetCount()
	if err != nil {
		log.Errorln("Unable to get follower count", err)
		return defaultWorkers
	}
	followerCount = fc

	// Scale more conservatively: start with base workers, add 1 worker per 100 followers
	// This gives a much more reasonable scaling than the previous followerCount * 5
	workers := minWorkers + int(followerCount/100)

	if workers > maxWorkers {
		workers = maxWorkers
	}

	log.Debugf("Initializing ActivityPub outbound worker pool with %d workers for %d followers", workers, followerCount)
	return workers
}

// SendLive will send a "Go Live" message to followers.
func SendLive() error {
	return outbox.SendLive()
}

// SendPublicFederatedMessage will send an arbitrary provided message to followers.
func SendPublicFederatedMessage(message string) error {
	return outbox.SendPublicMessage(message)
}

// SendDirectFederatedMessage will send a direct message to a single account.
func SendDirectFederatedMessage(message, account string) error {
	return outbox.SendDirectMessageToAccount(message, account)
}

// GetFollowerCount will return the local tracked follower count.
func GetFollowerCount() (int64, error) {
	followersRepo := followersrepository.Get()
	return followersRepo.GetCount()
}

// GetPendingFollowRequests will return the pending follow requests.
func GetPendingFollowRequests() ([]models.Follower, error) {
	followersRepo := followersrepository.Get()
	return followersRepo.GetPendingFollowRequests()
}
