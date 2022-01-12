package activitypub

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/controllers"
	"github.com/owncast/owncast/router/middleware"
)

// StartRouter will start the federation specific http router.
func StartRouter() {
	// WebFinger
	http.HandleFunc("/.well-known/webfinger", controllers.WebfingerHandler)

	// Host Metadata
	http.HandleFunc("/.well-known/host-meta", controllers.HostMetaController)

	// Nodeinfo v1
	http.HandleFunc("/.well-known/nodeinfo", controllers.NodeInfoController)

	// x-nodeinfo v2
	http.HandleFunc("/.well-known/x-nodeinfo2", controllers.XNodeInfo2Controller)

	// Nodeinfo v2
	http.HandleFunc("/nodeinfo/2.0", controllers.NodeInfoV2Controller)

	// Instance details
	http.HandleFunc("/api/v1/instance", controllers.InstanceV1Controller)

	// Single ActivityPub Actor
	http.HandleFunc("/federation/user/", middleware.RequireActivityPubOrRedirect(controllers.ActorHandler))

	// Single AP object
	http.HandleFunc("/federation/", middleware.RequireActivityPubOrRedirect(controllers.ObjectHandler))
}
