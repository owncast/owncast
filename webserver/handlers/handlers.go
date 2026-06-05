package handlers

import (
	"net/http"
	"sync"

	"github.com/jellydator/ttlcache/v3"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/models"
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

	// pluginActions, when non-nil, returns the current set of action
	// buttons contributed by loaded plugins. Merged into the
	// externalActions list returned by GetWebConfig so the viewer sees
	// plugin actions alongside admin-defined ones. nil = no plugin host
	// (boot disabled or failed).
	pluginActions func() []models.ExternalAction

	// pluginCSSContent, when non-nil, returns the concatenated CSS
	// bytes contributed by loaded plugins that declared `styles` in
	// their manifest. /api/config appends these bytes to the admin's
	// customStyles so the viewer renders one inline <style> block
	// covering both sources. nil = no plugin host.
	pluginCSSContent func() []byte

	// pluginJSContent mirrors pluginCSSContent for JavaScript: the
	// concatenated JS bytes contributed by loaded plugins. The
	// /customjavascript handler appends these to the admin's
	// customJavascript so the viewer loads one <script> tag covering
	// both sources.
	pluginJSContent func() []byte

	// pluginPageContent returns the concatenated HTML bytes
	// contributed by loaded plugins via manifest.extraPageContent.
	// /api/config prepends these bytes to the admin's rendered
	// extraPageContent.
	pluginPageContent func(*http.Request) []byte

	// pluginTabs returns the list of viewer-page tabs contributed by
	// loaded plugins via manifest.tabs. /api/config emits this list
	// as `pluginTabs`; the viewer page renders one tab per entry.
	pluginTabs func(*http.Request) []models.PluginTab

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
	// PluginActions is an optional getter that returns action buttons
	// contributed by loaded plugins. Wired by main.go to the plugin host's
	// Actions() method; nil when the plugin host is disabled.
	PluginActions func() []models.ExternalAction
	// PluginCSSContent is an optional getter that returns the
	// concatenated CSS bytes contributed by loaded plugins. Wired by
	// main.go to the plugin host's StylesContent() method; nil when
	// the plugin host is disabled.
	PluginCSSContent func() []byte
	// PluginJSContent mirrors PluginCSSContent for JavaScript: the
	// concatenated JS bytes contributed by loaded plugins. Wired to
	// the plugin host's ScriptsContent() method.
	PluginJSContent func() []byte
	// PluginPageContent returns the concatenated HTML bytes from
	// each loaded plugin's manifest.extraPageContent. Wired to the
	// plugin host's PageContent() method.
	PluginPageContent func(*http.Request) []byte
	// PluginTabs returns the list of viewer-page tabs contributed by
	// loaded plugins via manifest.tabs. Wired to the plugin host's
	// Tabs() method.
	PluginTabs func(*http.Request) []models.PluginTab
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
		pluginActions:           deps.PluginActions,
		pluginCSSContent:        deps.PluginCSSContent,
		pluginJSContent:         deps.PluginJSContent,
		pluginPageContent:       deps.PluginPageContent,
		pluginTabs:              deps.PluginTabs,
		previewThumbCache: ttlcache.New(
			ttlcache.WithTTL[string, []byte](15),
			ttlcache.WithCapacity[string, []byte](1),
			ttlcache.WithDisableTouchOnHit[string, []byte](),
		),
	}
}
