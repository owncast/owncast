package admin

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/authrepository"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	"github.com/owncast/owncast/services/chat"
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
	stream                *stream.Service
	rtmp                  *rtmp.Service
	activitypub           *activitypub.Service
	webhooks              *webhooks.Service
	chat                  *chat.Service
	metrics               *metrics.Service
	configRepository      configrepository.ConfigRepository
	authRepository        authrepository.AuthRepository
	followersRepository   followersrepository.FollowersRepository
	webhookRepository     webhookrepository.WebhookRepository
	chatMessageRepository chatmessagerepository.ChatMessageRepository
	userRepository        userrepository.UserRepository
	apBuilder             *apmodels.Builder
	apSigner              *apcrypto.Signer
	cfg                   *config.Config

	// pluginStyleContributors, when non-nil, returns the per-plugin
	// page-styling report (which enabled plugins emit CSS and which
	// appearance tokens they set). Wired to the plugin host's
	// StyleContributors() method; nil when the plugin host is disabled.
	// GetServerConfig emits it as `styleContributors` so the admin
	// Appearance UI can flag plugin styling that combines with the
	// admin's own colors.
	pluginStyleContributors func() []models.PluginStyleInfo
}

// Deps lists every service a *Admin consumes. New deps appear here as
// more admin handlers migrate.
type Deps struct {
	Stream                *stream.Service
	Rtmp                  *rtmp.Service
	Activitypub           *activitypub.Service
	Webhooks              *webhooks.Service
	Chat                  *chat.Service
	Metrics               *metrics.Service
	ConfigRepository      configrepository.ConfigRepository
	AuthRepository        authrepository.AuthRepository
	FollowersRepository   followersrepository.FollowersRepository
	WebhookRepository     webhookrepository.WebhookRepository
	ChatMessageRepository chatmessagerepository.ChatMessageRepository
	UserRepository        userrepository.UserRepository
	APBuilder             *apmodels.Builder
	APSigner              *apcrypto.Signer
	Config                *config.Config
	// PluginStyleContributors is an optional getter that returns the
	// per-plugin page-styling report. Wired by main.go to the plugin
	// host's StyleContributors() method; nil when the plugin host is
	// disabled.
	PluginStyleContributors func() []models.PluginStyleInfo
}

// New constructs the dependency-bearing admin handler set.
func New(deps Deps) *Admin {
	return &Admin{
		stream:                  deps.Stream,
		rtmp:                    deps.Rtmp,
		activitypub:             deps.Activitypub,
		webhooks:                deps.Webhooks,
		chat:                    deps.Chat,
		metrics:                 deps.Metrics,
		configRepository:        deps.ConfigRepository,
		authRepository:          deps.AuthRepository,
		followersRepository:     deps.FollowersRepository,
		webhookRepository:       deps.WebhookRepository,
		chatMessageRepository:   deps.ChatMessageRepository,
		userRepository:          deps.UserRepository,
		apBuilder:               deps.APBuilder,
		apSigner:                deps.APSigner,
		cfg:                     deps.Config,
		pluginStyleContributors: deps.PluginStyleContributors,
	}
}
