package activitypub

import (
	"fmt"

	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/core/data"
)

func Start(datastore *data.Datastore) {
	persistence.Setup(datastore)
	StartRouter()

	// Test
	if data.GetPrivateKey() == "" {
		privateKey, publicKey, err := crypto.GenerateKeys()
		fmt.Println(string(privateKey), string(publicKey))

		data.SetPrivateKey(string(privateKey))
		data.SetPublicKey(string(publicKey))
		if err != nil {
			panic(err)
		}
	}

}
