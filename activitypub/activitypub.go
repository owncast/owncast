package activitypub

import (
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

var _datastore *data.Datastore

func Start(datastore *data.Datastore) {
	_datastore = datastore

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

func SendLive() {
	outbox.SendLive()
}

func SendPublicFederatedMessage(message string) {
	outbox.SendPublicMessage(message)
}

func GetFollowerCount() int {
	return persistence.GetFollowerCount()
}
