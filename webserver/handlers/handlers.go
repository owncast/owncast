package handlers

import (
	"net/http"
	"sync"

	"github.com/jellydator/ttlcache/v3"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/cache"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/services/stream"
	"github.com/owncast/owncast/webserver/handlers/admin"
	"github.com/owncast/owncast/webserver/handlers/auth/fediverse"
	"github.com/owncast/owncast/webserver/handlers/auth/indieauth"
	"github.com/owncast/owncast/webserver/handlers/moderation"
	"github.com/owncast/owncast/webserver/router/middleware"
	"github.com/owncast/owncast/yp"
)

// Handlers carries the dependencies of HTTP handlers that need injected
// services. Construct one in main() with NewHandlers and pass it to the
// router; methods on *Handlers are registered as routes.
//
// Handlers that don't yet need dependencies remain free functions in this
// package; they migrate to methods as the services they depend on move to
// services/<domain>/ and stop being callable via package-level singletons.
type Handlers struct {
	cache                   *cache.Container
	stream                  *stream.Service
	chat                    *chat.Service
	admin                   *admin.Admin
	activitypub             *activitypub.Service
	fediverse               *fediverse.Handler
	indieauth               *indieauth.Handler
	moderation              *moderation.Handler
	middleware              *middleware.Middleware
	yp                      *yp.YP
	metrics                 *metrics.Service
	configRepository        configrepository.ConfigRepository
	followersRepository     followersrepository.FollowersRepository
	chatMessageRepository   chatmessagerepository.ChatMessageRepository
	userRepository          userrepository.UserRepository
	notificationsRepository notificationsrepository.NotificationsRepository
	apBuilder               *apmodels.Builder
	cfg                     *config.Config

	// previewThumbCache caches thumbnail/preview bytes for a short window
	// so frequent polling from chat clients doesn't re-read the file
	// every request.
	previewThumbCache *ttlcache.Cache[string, []byte]

	// hasWarnedSVGLogo gates the one-time warning logged when an
	// external site requests an SVG logo via /logo/external.
	hasWarnedSVGLogoLock sync.Mutex
	hasWarnedSVGLogo     bool
}

// Deps lists every service a *Handlers consumes. New deps appear here as
// more handlers migrate.
type Deps struct {
	Cache                   *cache.Container
	Stream                  *stream.Service
	Chat                    *chat.Service
	Admin                   *admin.Admin
	Activitypub             *activitypub.Service
	Fediverse               *fediverse.Handler
	IndieAuth               *indieauth.Handler
	Moderation              *moderation.Handler
	Middleware              *middleware.Middleware
	YP                      *yp.YP
	Metrics                 *metrics.Service
	ConfigRepository        configrepository.ConfigRepository
	FollowersRepository     followersrepository.FollowersRepository
	ChatMessageRepository   chatmessagerepository.ChatMessageRepository
	UserRepository          userrepository.UserRepository
	NotificationsRepository notificationsrepository.NotificationsRepository
	APBuilder               *apmodels.Builder
	Config                  *config.Config
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
		cache:                   deps.Cache,
		stream:                  deps.Stream,
		chat:                    deps.Chat,
		admin:                   deps.Admin,
		activitypub:             deps.Activitypub,
		fediverse:               deps.Fediverse,
		indieauth:               deps.IndieAuth,
		moderation:              deps.Moderation,
		middleware:              deps.Middleware,
		yp:                      deps.YP,
		metrics:                 deps.Metrics,
		configRepository:        deps.ConfigRepository,
		followersRepository:     deps.FollowersRepository,
		chatMessageRepository:   deps.ChatMessageRepository,
		userRepository:          deps.UserRepository,
		notificationsRepository: deps.NotificationsRepository,
		apBuilder:               deps.APBuilder,
		cfg:                     deps.Config,
		previewThumbCache: ttlcache.New(
			ttlcache.WithTTL[string, []byte](15),
			ttlcache.WithCapacity[string, []byte](1),
			ttlcache.WithDisableTouchOnHit[string, []byte](),
		),
	}
}
