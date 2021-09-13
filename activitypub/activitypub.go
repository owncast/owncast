package activitypub

import (
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

// Start will initialize and start the federation support.
func Start(datastore *data.Datastore) {
	persistence.Setup(datastore)
	StartRouter()

	// Test
	if data.GetPrivateKey() == "" {
		privateKey, publicKey, err := crypto.GenerateKeys()
		_ = data.SetPrivateKey(string(privateKey))
		_ = data.SetPublicKey(string(publicKey))
		if err != nil {
			log.Errorln("Unable to get private key", err)
		}
	}
}

// SendLive will send a "Go Live" message to followers.
func SendLive() error {
	return outbox.SendLive()
}

// SendPublicFederatedMessage will send an arbitrary provided message to followers.
func SendPublicFederatedMessage(message string) error {
	return outbox.SendPublicMessage(message)
}

// GetFollowerCount will return the local tracked follower count.
func GetFollowerCount() int {
	return persistence.GetFollowerCount()
}
