package admin

import (
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/rtmp"
	"github.com/owncast/owncast/services/stream"
	"github.com/owncast/owncast/services/webhooks"
)

// Admin carries the dependencies of admin HTTP handlers that need
// injected services. Construct one in main() and pass it to the OpenAPI
// shim layer (webserver/handlers.NewHandlers) so the matching
// ServerInterface methods can delegate to it.
//
// Admin handlers without injected dependencies remain free functions in
// this package; they migrate to methods on *Admin as the services they
// need to consume move to services/<domain>/.
type Admin struct {
	stream      *stream.Service
	rtmp        *rtmp.Service
	activitypub *activitypub.Service
	webhooks    *webhooks.Service
}

// Deps lists every service a *Admin consumes. New deps appear here as
// more admin handlers migrate.
type Deps struct {
	Stream      *stream.Service
	Rtmp        *rtmp.Service
	Activitypub *activitypub.Service
	Webhooks    *webhooks.Service
}

// New constructs the dependency-bearing admin handler set.
func New(deps Deps) *Admin {
	return &Admin{
		stream:      deps.Stream,
		rtmp:        deps.Rtmp,
		activitypub: deps.Activitypub,
		webhooks:    deps.Webhooks,
	}
}
