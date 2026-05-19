package handlers

import (
	"net/http"

	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/cache"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/services/stream"
	"github.com/owncast/owncast/webserver/handlers/admin"
	"github.com/owncast/owncast/webserver/handlers/auth/fediverse"
	"github.com/owncast/owncast/webserver/handlers/auth/indieauth"
	"github.com/owncast/owncast/webserver/handlers/moderation"
)

// Handlers carries the dependencies of HTTP handlers that need injected
// services. Construct one in main() with NewHandlers and pass it to the
// router; methods on *Handlers are registered as routes.
//
// Handlers that don't yet need dependencies remain free functions in this
// package; they migrate to methods as the services they depend on move to
// services/<domain>/ and stop being callable via package-level singletons.
type Handlers struct {
	cache               *cache.Container
	stream              *stream.Service
	chat                *chat.Service
	admin               *admin.Admin
	activitypub         *activitypub.Service
	fediverse           *fediverse.Handler
	indieauth           *indieauth.Handler
	moderation          *moderation.Handler
	configRepository    configrepository.ConfigRepository
	followersRepository followersrepository.FollowersRepository
}

// Deps lists every service a *Handlers consumes. New deps appear here as
// more handlers migrate.
type Deps struct {
	Cache               *cache.Container
	Stream              *stream.Service
	Chat                *chat.Service
	Admin               *admin.Admin
	Activitypub         *activitypub.Service
	Fediverse           *fediverse.Handler
	IndieAuth           *indieauth.Handler
	Moderation          *moderation.Handler
	ConfigRepository    configrepository.ConfigRepository
	FollowersRepository followersrepository.FollowersRepository
}

// HandleWebsocketConnection routes the /ws websocket upgrade to the
// chat service. Lives here so the router can bind a method on
// *Handlers instead of reaching into chat directly.
func (h *Handlers) HandleWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	h.chat.HandleWebsocketConnection(w, r)
}

// NewHandlers constructs the dependency-bearing handler set.
func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		cache:               deps.Cache,
		stream:              deps.Stream,
		chat:                deps.Chat,
		admin:               deps.Admin,
		activitypub:         deps.Activitypub,
		fediverse:           deps.Fediverse,
		indieauth:           deps.IndieAuth,
		moderation:          deps.Moderation,
		configRepository:    deps.ConfigRepository,
		followersRepository: deps.FollowersRepository,
	}
}
