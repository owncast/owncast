package activitypub

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/controllers"
	"github.com/owncast/owncast/router/middleware"
)

func StartRouter() {
	// WebFinger
	http.HandleFunc("/.well-known/webfinger", middleware.RequireActivityPubOrRedirect(controllers.WebfingerHandler))

	// Host Metadata
	http.HandleFunc("/.well-known/host-meta", controllers.HostMetaController)

	// Nodeinfo v1
	http.HandleFunc("/.well-known/nodeinfo", middleware.RequireActivityPubOrRedirect(controllers.NodeInfoController))

	// Nodeinfo v2
	http.HandleFunc("/nodeinfo/2.0", middleware.RequireActivityPubOrRedirect(controllers.NodeInfoV2Controller))

	// Instance details
	http.HandleFunc("/api/v1/instance", controllers.InstanceV1Controller)

	// Single ActivityPub Actor
	http.HandleFunc("/federation/user/", middleware.RequireActivityPubOrRedirect(controllers.ActorHandler))

	// Single AP object
	http.HandleFunc("/federation/", middleware.RequireActivityPubOrRedirect(controllers.ObjectHandler))
}
