package activitypub

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/controllers"
	"github.com/owncast/owncast/webserver/middleware"
)

// StartRouter will start the federation specific http router.
func StartRouter(router *http.ServeMux) {
	// WebFinger
	router.HandleFunc("/.well-known/webfinger", controllers.WebfingerHandler)

	// Host Metadata
	router.HandleFunc("/.well-known/host-meta", controllers.HostMetaController)

	// Nodeinfo v1
	router.HandleFunc("/.well-known/nodeinfo", controllers.NodeInfoController)

	// x-nodeinfo v2
	router.HandleFunc("/.well-known/x-nodeinfo2", controllers.XNodeInfo2Controller)

	// Nodeinfo v2
	router.HandleFunc("/nodeinfo/2.0", controllers.NodeInfoV2Controller)

	// Instance details
	router.HandleFunc("/api/v1/instance", controllers.InstanceV1Controller)

	// Single ActivityPub Actor
	router.HandleFunc("/federation/user/", middleware.RequireActivityPubOrRedirect(controllers.ActorHandler))

	// Single AP object
	router.HandleFunc("/federation/", middleware.RequireActivityPubOrRedirect(controllers.ObjectHandler))
}
