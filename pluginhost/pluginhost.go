package pluginhost

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/notifications/browser"
	"github.com/owncast/owncast/notifications/discord"
	"github.com/owncast/owncast/persistence/authrepository"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/services/stream"
	"github.com/owncast/owncast/services/webhooks"
	"github.com/owncast/owncast/utils"
)

// pluginsEnabledConfigKey is the datastore key under which the set of
// admin-enabled plugin names is persisted.
const pluginsEnabledConfigKey = "plugins.enabled"

// adminAuthUsername is the fixed HTTP Basic Auth username for admin requests,
// matching webserver/router/middleware.RequireAdminAuth.
const adminAuthUsername = "admin"

// Deps bundles the Owncast services the plugin runtime adapts into
// HostEnv host functions. Everything here is already constructed in the main
// composition root and passed in by reference.
type Deps struct {
	Datastore               *datastore.Datastore
	Chat                    *chat.Service
	Stream                  *stream.Service
	Activitypub             *activitypub.Service
	Webhooks                *webhooks.Service
	ConfigRepository        configrepository.ConfigRepository
	UserRepository          userrepository.UserRepository
	AuthRepository          *authrepository.SqlAuthRepository
	NotificationsRepository notificationsrepository.NotificationsRepository
	ChatMessageRepository   chatmessagerepository.ChatMessageRepository
}

// Host owns the running plugin runtime: the manager (discovery +
// enable/disable lifecycle), the HTTP handler that serves /plugins/<name>/*,
// and the SSE hub backing host-owned event streams.
type Host struct {
	manager          *plugins.Manager
	server           *plugins.Server
	sse              *plugins.SSEHub
	configRepository configrepository.ConfigRepository
}

// Handler is the http.Handler for /plugins/<name>/* (static assets, dynamic
// on_http_request, and the reserved _sse endpoint).
func (p *Host) Handler() http.Handler { return p.server }

// Stop closes all loaded plugins.
func (p *Host) Stop(ctx context.Context) {
	p.manager.Stop(ctx)
}

// New builds the HostEnv from Owncast services, constructs and
// starts the plugin manager, and returns the assembled host. Plugins are
// optional infrastructure: callers should log and continue on error rather
// than aborting startup.
func New(ctx context.Context, deps Deps) (*Host, error) {
	pluginsDir := filepath.Join(config.DataDirectory, "plugins")
	if err := os.MkdirAll(pluginsDir, 0o700); err != nil {
		return nil, fmt.Errorf("create plugins directory: %w", err)
	}

	env := &plugins.HostEnv{KV: newDatastoreKVStore(deps.Datastore)}
	wirePluginHostEnv(env, deps)

	sseHub := plugins.NewSSEHub()
	env.SSE = sseHub

	enabledStore := &configEnabledStore{datastore: deps.Datastore}
	manager := plugins.NewManagerWithStore(pluginsDir, env, enabledStore)
	if err := manager.Start(ctx); err != nil {
		return nil, fmt.Errorf("start plugin manager: %w", err)
	}

	// Emit delivers plugin-published custom events to other plugins'
	// subscribers. Wired post-Start because it reads the live plugin set.
	dispatcher := plugins.NewLiveDispatcher(manager.Snapshot)
	env.Emit = dispatcher.Dispatch

	// Deliver Owncast's own events (chat, stream lifecycle, moderation, …) to
	// subscribed plugins by listening on the webhooks event source.
	if deps.Webhooks != nil {
		deps.Webhooks.AddEventListener(newPluginEventListener(dispatcher))
	}

	// Run plugin filterChatMessage handlers on each inbound chat message
	// before it's broadcast, so plugins can rewrite or drop messages.
	deps.Chat.SetMessageFilter(newPluginChatFilter(dispatcher))

	server := plugins.NewLiveServer(manager.Snapshot)
	server.SSE = sseHub
	server.IsAuthenticated = env.IsAuthenticated
	server.GetRequestUser = env.GetRequestUser

	return &Host{manager: manager, server: server, sse: sseHub, configRepository: deps.ConfigRepository}, nil
}

// AdminHandler returns the HTTP handler for plugin management:
//
//	GET  /api/admin/plugins                       list discovered plugins (admin)
//	POST /api/admin/plugins/<name>/enable|disable|reload  toggle a plugin (admin)
//	GET  /api/plugins/actions                     merged action-button list (public)
//
// It's mounted by the router on the outer mux so it sits beside, not inside,
// the OpenAPI-generated /api router.
func (p *Host) AdminHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/admin/plugins", func(w http.ResponseWriter, r *http.Request) {
		if !p.requireAdmin(w, r) {
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSONResponse(w, http.StatusOK, p.manager.List())
	})

	mux.HandleFunc("/api/admin/plugins/", func(w http.ResponseWriter, r *http.Request) {
		if !p.requireAdmin(w, r) {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name, action, ok := strings.Cut(strings.TrimPrefix(r.URL.Path, "/api/admin/plugins/"), "/")
		if !ok || name == "" || action == "" {
			http.Error(w, "expected /<name>/<action>", http.StatusBadRequest)
			return
		}
		var err error
		switch action {
		case "enable":
			err = p.manager.Enable(r.Context(), name)
		case "disable":
			err = p.manager.Disable(r.Context(), name)
		case "reload":
			err = p.manager.Reload(r.Context(), name)
		default:
			http.Error(w, "unknown action; expected enable, disable, or reload", http.StatusBadRequest)
			return
		}
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]string{"status": "ok", "name": name, "action": action})
	})

	mux.HandleFunc("/api/plugins/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		actions := make([]plugins.ActionButton, 0)
		for _, loaded := range p.manager.Snapshot() {
			actions = append(actions, loaded.Manifest.Actions...)
		}
		writeJSONResponse(w, http.StatusOK, actions)
	})

	return mux
}

// requireAdmin gates a management endpoint on admin HTTP Basic Auth. It
// writes the 401 response itself and returns false when unauthenticated.
func (p *Host) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	if isAdminRequest(r, p.configRepository) {
		return true
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="Owncast plugin admin"`)
	http.Error(w, "authentication required", http.StatusUnauthorized)
	return false
}

func writeJSONResponse(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// configEnabledStore persists the enabled-plugin set in Owncast's config
// datastore instead of a .enabled.json file.
type configEnabledStore struct {
	datastore *datastore.Datastore
}

func (s *configEnabledStore) LoadEnabled() ([]string, error) {
	names, err := s.datastore.GetStringSlice(pluginsEnabledConfigKey)
	if err != nil {
		// Unset on a fresh install — start with no plugins enabled.
		return nil, nil
	}
	return names, nil
}

func (s *configEnabledStore) SaveEnabled(names []string) error {
	return s.datastore.SetStringSlice(pluginsEnabledConfigKey, names)
}

// wirePluginHostEnv connects each HostEnv host-function pointer to the
// corresponding Owncast service call. Closures read services lazily so they
// observe current config/state on every call.
func wirePluginHostEnv(env *plugins.HostEnv, deps Deps) {
	chatbots := newPluginChatbotProvisioner(deps.UserRepository, deps.Datastore)
	wireChatSendHostFns(env, deps, chatbots)
	wireChatReadHostFns(env, deps)
	wireChatModerationHostFns(env, deps)
	wireServerReadHostFns(env, deps)
	wireUserHostFns(env, deps)
	wireNotificationHostFns(env, deps)
	wireRequestHostFns(env, deps)
}

func wireChatSendHostFns(env *plugins.HostEnv, deps Deps, chatbots *pluginChatbotProvisioner) {
	chatSvc := deps.Chat

	env.OnChat = func(req plugins.ChatSendRequest) {
		switch req.Kind {
		case plugins.ChatSendAction:
			if err := chatSvc.SendSystemAction(req.Text, false); err != nil {
				log.Errorln("plugin", req.PluginName, "chat action:", err)
			}
		case plugins.ChatSendSystem:
			if err := chatSvc.SendSystemMessage(req.Text, false); err != nil {
				log.Errorln("plugin", req.PluginName, "chat system message:", err)
			}
		default: // ChatSendBot — post under the plugin's own chatbot identity.
			chatbot, err := chatbots.chatbotUser(req.PluginName)
			if err != nil {
				log.Errorln("plugin", req.PluginName, "resolve chatbot user:", err)
				return
			}
			if err := chatSvc.SendMessageAsUser(chatbot, req.Text); err != nil {
				log.Errorln("plugin", req.PluginName, "chat send:", err)
			}
		}
	}

	env.SendChatTo = func(pluginName string, clientID uint64, text string) {
		chatSvc.SendSystemMessageToClient(uint(clientID), text)
	}
}

func wireChatReadHostFns(env *plugins.HostEnv, deps Deps) {
	chatSvc := deps.Chat

	env.ChatHistory = func(limit int) []plugins.HostChatMessage {
		history := deps.ChatMessageRepository.GetChatHistory()
		out := make([]plugins.HostChatMessage, 0, len(history))
		for _, item := range history {
			msg, ok := item.(events.UserMessageEvent)
			if !ok {
				continue
			}
			hm := plugins.HostChatMessage{
				ID:        msg.ID,
				Body:      msg.Body,
				Timestamp: msg.Timestamp.UTC().Format(time.RFC3339),
			}
			if msg.User != nil {
				hm.User = msg.User.DisplayName
			}
			out = append(out, hm)
		}
		if limit > 0 && len(out) > limit {
			out = out[len(out)-limit:]
		}
		return out
	}

	env.ChatClients = func() []plugins.HostChatClient {
		clients := chatSvc.GetClients()
		out := make([]plugins.HostChatClient, 0, len(clients))
		for _, c := range clients {
			hc := plugins.HostChatClient{
				ID:           uint64(c.Id),
				ConnectedAt:  c.ConnectedAt.UTC().Format(time.RFC3339),
				UserAgent:    c.UserAgent,
				IPAddress:    c.IPAddress,
				MessageCount: c.MessageCount,
			}
			if c.User != nil {
				hc.UserID = c.User.ID
				hc.DisplayName = c.User.DisplayName
			}
			out = append(out, hc)
		}
		return out
	}
}

func wireChatModerationHostFns(env *plugins.HostEnv, deps Deps) {
	chatSvc := deps.Chat

	env.DeleteMessage = func(pluginName, messageID string) {
		if err := chatSvc.SetMessagesVisibility([]string{messageID}, false); err != nil {
			log.Errorln("plugin", pluginName, "delete message:", err)
		}
	}

	env.KickClient = func(pluginName string, clientID uint64) {
		if c, ok := chatSvc.FindClientByID(uint(clientID)); ok {
			chatSvc.DisconnectClients([]*chat.Client{c})
		}
	}
}

func wireServerReadHostFns(env *plugins.HostEnv, deps Deps) {
	cfg := deps.ConfigRepository
	streamSvc := deps.Stream

	env.StreamCurrent = func() plugins.StreamInfo {
		status := streamSvc.GetStatus()
		info := plugins.StreamInfo{
			Online:       status.Online,
			Title:        cfg.GetStreamTitle(),
			Summary:      cfg.GetServerSummary(),
			Viewers:      status.ViewerCount,
			LatencyLevel: cfg.GetStreamLatencyLevel().Level,
		}
		if status.LastConnectTime != nil && status.LastConnectTime.Valid {
			info.StartedAt = status.LastConnectTime.Time.UTC().Format(time.RFC3339)
		}
		return info
	}

	env.ServerInfo = func() plugins.ServerInfo {
		return plugins.ServerInfo{
			Name:           cfg.GetServerName(),
			URL:            cfg.GetServerURL(),
			Summary:        cfg.GetServerSummary(),
			WelcomeMessage: cfg.GetServerWelcomeMessage(),
			Version:        config.VersionNumber,
		}
	}

	env.Socials = func() []plugins.SocialHandle {
		handles := cfg.GetSocialHandles()
		out := make([]plugins.SocialHandle, 0, len(handles))
		for _, h := range handles {
			out = append(out, plugins.SocialHandle{Platform: h.Platform, URL: h.URL, Icon: h.Icon})
		}
		return out
	}

	env.Federation = func() plugins.FederationInfo {
		return plugins.FederationInfo{
			Enabled:   cfg.GetFederationEnabled(),
			Username:  cfg.GetFederationUsername(),
			IsPrivate: cfg.GetFederationIsPrivate(),
		}
	}
}

func wireUserHostFns(env *plugins.HostEnv, deps Deps) {
	users := deps.UserRepository
	chatSvc := deps.Chat

	env.Users = func() []plugins.HostUser {
		all := users.GetUsers()
		out := make([]plugins.HostUser, 0, len(all))
		for _, u := range all {
			out = append(out, toHostUser(u))
		}
		return out
	}

	env.UserGet = func(id string) (plugins.HostUser, bool) {
		u := users.GetUserByID(id)
		if u == nil {
			return plugins.HostUser{}, false
		}
		return toHostUser(u), true
	}

	env.SetUserEnabled = func(pluginName, userID string, enabled bool, reason string) {
		if err := users.SetEnabled(userID, enabled); err != nil {
			log.Errorln("plugin", pluginName, "set user enabled:", err)
			return
		}
		if !enabled {
			if clients, err := chatSvc.GetClientsForUser(userID); err == nil {
				chatSvc.DisconnectClients(clients)
			}
		}
	}

	env.BanIP = func(pluginName, ip string) {
		if err := deps.AuthRepository.BanIPAddress(ip, "banned by plugin "+pluginName); err != nil {
			log.Errorln("plugin", pluginName, "ban ip:", err)
		}
	}
}

func wireNotificationHostFns(env *plugins.HostEnv, deps Deps) {
	cfg := deps.ConfigRepository

	env.SendDiscord = func(pluginName, text string) {
		dc := cfg.GetDiscordConfig()
		if !dc.Enabled || dc.Webhook == "" {
			return
		}
		notifier, err := discord.New(cfg.GetServerName(), cfg.GetServerURL()+"/logo", dc.Webhook)
		if err != nil {
			log.Errorln("plugin", pluginName, "discord:", err)
			return
		}
		if err := notifier.Send(text); err != nil {
			log.Errorln("plugin", pluginName, "discord send:", err)
		}
	}

	env.SendBrowserPush = func(pluginName string, p plugins.BrowserPushPayload) {
		publicKey, err := cfg.GetBrowserPushPublicKey()
		if err != nil {
			return
		}
		privateKey, err := cfg.GetBrowserPushPrivateKey()
		if err != nil {
			return
		}
		notifier, err := browser.New(deps.Datastore, publicKey, privateKey)
		if err != nil {
			log.Errorln("plugin", pluginName, "browser push:", err)
			return
		}
		destinations, err := deps.NotificationsRepository.GetNotificationDestinationsForChannel(notificationsrepository.BrowserPushNotification)
		if err != nil {
			return
		}
		for _, destination := range destinations {
			if _, err := notifier.Send(destination, p.Title, p.Body); err != nil {
				log.Debugln("plugin", pluginName, "browser push send:", err)
			}
		}
	}

	env.SendFediverse = func(pluginName string, p plugins.FediversePayload) {
		var image *string
		if p.Image != "" {
			image = &p.Image
		}
		if err := deps.Chat.SendFediverseAction(p.Type, cfg.GetFederationUsername(), image, p.Body, p.Link); err != nil {
			log.Errorln("plugin", pluginName, "fediverse action:", err)
		}
	}

	env.PostFediverse = func(pluginName, text string) (string, error) {
		// Owncast publishes the note but does not return its URL.
		if err := deps.Activitypub.SendPublicFederatedMessage(text); err != nil {
			return "", err
		}
		return "", nil
	}
}

func wireRequestHostFns(env *plugins.HostEnv, deps Deps) {
	cfg := deps.ConfigRepository
	users := deps.UserRepository

	env.UploadStorage = func(pluginName, name string, data []byte) (string, error) {
		return uploadPluginAsset(cfg, pluginName, name, data)
	}

	env.IsAuthenticated = func(r *http.Request) bool {
		return isAdminRequest(r, cfg)
	}

	env.GetRequestUser = func(r *http.Request) *plugins.HostUser {
		token := r.URL.Query().Get("accessToken")
		if token == "" {
			return nil
		}
		u := users.GetUserByToken(token)
		if u == nil {
			return nil
		}
		hu := toHostUser(u)
		return &hu
	}
}

// toHostUser maps an Owncast user model onto the plugin-facing HostUser.
func toHostUser(u *models.User) plugins.HostUser {
	hu := plugins.HostUser{
		ID:              u.ID,
		DisplayName:     u.DisplayName,
		PreviousNames:   u.PreviousNames,
		Scopes:          u.Scopes,
		IsBot:           u.IsBot,
		IsAuthenticated: u.Authenticated,
		CreatedAt:       u.CreatedAt.UTC().Format(time.RFC3339),
	}
	if u.DisabledAt != nil {
		hu.DisabledAt = u.DisabledAt.UTC().Format(time.RFC3339)
	}
	return hu
}

// isAdminRequest reports whether r carries valid admin HTTP Basic Auth,
// mirroring the check in webserver/router/middleware.RequireAdminAuth.
func isAdminRequest(r *http.Request, cfg configrepository.ConfigRepository) bool {
	user, pass, ok := r.BasicAuth()
	if !ok {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(user), []byte(adminAuthUsername)) != 1 {
		return false
	}
	return utils.CompareHash(cfg.GetAdminPassword(), pass) == nil
}

// uploadPluginAsset writes a plugin upload under the public files directory
// and returns the URL it is served at. Names are flattened to their base to
// prevent path traversal. S3-backed storage is a follow-up.
func uploadPluginAsset(cfg configrepository.ConfigRepository, pluginName, name string, data []byte) (string, error) {
	safeName := filepath.Base(filepath.Clean("/" + name))
	if safeName == "." || safeName == "/" || safeName == "" {
		return "", fmt.Errorf("invalid upload name %q", name)
	}
	dir := filepath.Join(config.PublicFilesPath, "plugins", pluginName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, safeName), data, 0o600); err != nil {
		return "", err
	}
	url := strings.TrimSuffix(cfg.GetServerURL(), "/") + "/public/plugins/" + pluginName + "/" + safeName
	return url, nil
}
