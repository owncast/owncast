package activitypub

import (
	"fmt"
	"net/http"

	"github.com/owncast/owncast/activitypub/controllers"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/core/data"
)

func Start() {
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

	// WebFinger
	http.HandleFunc("/.well-known/webfinger", controllers.WebfingerHandler)

	// Single ActivityPub Actor
	http.HandleFunc("/federation/user/", controllers.ActorHandler)
}
