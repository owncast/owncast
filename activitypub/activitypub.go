package activitypub

import (
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/outbox"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/core/data"
)

func Start(datastore *data.Datastore) {
	persistence.Setup(datastore)
	StartRouter()

	// Test
	if data.GetPrivateKey() == "" {
		privateKey, publicKey, err := crypto.GenerateKeys()
		_ = data.SetPrivateKey(string(privateKey))
		_ = data.SetPublicKey(string(publicKey))
		if err != nil {
			panic(err)
		}
	}
}

func SendLive() {
	outbox.SendLive()
}
