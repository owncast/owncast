package activitypub

import (
	"net/http"

	"github.com/owncast/owncast/app/middleware"
)

// StartRouter will start the federation specific http router.
func (s *Service) StartRouter(mux *http.ServeMux) {
	// WebFinger
	mux.HandleFunc("/.well-known/webfinger", s.WebfingerHandler)

	// Host Metadata
	mux.HandleFunc("/.well-known/host-meta", s.HostMetaController)

	// Nodeinfo v1
	mux.HandleFunc("/.well-known/nodeinfo", s.NodeInfoController)

	// x-nodeinfo v2
	mux.HandleFunc("/.well-known/x-nodeinfo2", s.XNodeInfo2Controller)

	// Nodeinfo v2
	mux.HandleFunc("/nodeinfo/2.0", s.NodeInfoV2Controller)

	// Instance details
	mux.HandleFunc("/api/v1/instance", s.InstanceV1Controller)

	// Single ActivityPub Actor
	mux.HandleFunc("/federation/user/", middleware.RequireActivityPubOrRedirect(s.ActorHandler, s.Persistence.Data))

	// Single AP object
	mux.HandleFunc("/federation/", middleware.RequireActivityPubOrRedirect(s.ObjectHandler, s.Persistence.Data))
}
