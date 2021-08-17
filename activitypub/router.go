package activitypub

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/controllers"
)

func StartRouter() {
	// WebFinger
	http.HandleFunc("/.well-known/webfinger", controllers.WebfingerHandler)

	// Single ActivityPub Actor
	http.HandleFunc("/federation/user/", controllers.ActorHandler)

	// Single AP object
	http.HandleFunc("/federation/", controllers.ObjectHandler)
}
