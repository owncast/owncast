package handlers

import (
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/cache"
	"github.com/owncast/owncast/services/stream"
	"github.com/owncast/owncast/webserver/handlers/admin"
	"github.com/owncast/owncast/webserver/handlers/auth/fediverse"
)

// Handlers carries the dependencies of HTTP handlers that need injected
// services. Construct one in main() with NewHandlers and pass it to the
// router; methods on *Handlers are registered as routes.
//
// Handlers that don't yet need dependencies remain free functions in this
// package; they migrate to methods as the services they depend on move to
// services/<domain>/ and stop being callable via package-level singletons.
type Handlers struct {
	cache       *cache.Container
	stream      *stream.Service
	admin       *admin.Admin
	activitypub *activitypub.Service
	fediverse   *fediverse.Handler
}

// Deps lists every service a *Handlers consumes. New deps appear here as
// more handlers migrate.
type Deps struct {
	Cache       *cache.Container
	Stream      *stream.Service
	Admin       *admin.Admin
	Activitypub *activitypub.Service
	Fediverse   *fediverse.Handler
}

// NewHandlers constructs the dependency-bearing handler set.
func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		cache:       deps.Cache,
		stream:      deps.Stream,
		admin:       deps.Admin,
		activitypub: deps.Activitypub,
		fediverse:   deps.Fediverse,
	}
}
