package activitypub

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/controllers"
)

func StartRouter() {
	// WebFinger
	http.HandleFunc("/.well-known/webfinger", controllers.WebfingerHandler)

	// Nodeinfo v1
	http.HandleFunc("/.well-known/nodeinfo", controllers.NodeInfoController)

	// Nodeinfo v2
	http.HandleFunc("/nodeinfo/2.0", controllers.NodeInfoV2Controller)

	// Instance details
	http.HandleFunc("/api/v1/instance", controllers.InstanceV1Controller)

	// Single ActivityPub Actor
	http.HandleFunc("/federation/user/", controllers.ActorHandler)

	// Single AP object
	http.HandleFunc("/federation/", controllers.ObjectHandler)
}
