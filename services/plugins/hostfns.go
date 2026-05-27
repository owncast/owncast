package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	extism "github.com/extism/go-sdk"
	"github.com/owncast/owncast/services/plugins/kv"
)

var _ = (*http.Request)(nil) // ensure import retained for the HostEnv signature

// Permission identifiers. The manifest declares which a plugin needs; only
// the corresponding wasm imports are wired into the plugin instance.
const (
	PermStorageKV         = "storage.kv"
	PermStorageUpload     = "storage.upload"
	PermChatSend          = "chat.send"
	PermChatHistory       = "chat.history"
	PermChatModerate      = "chat.moderate"
	PermNetworkFetch      = "network.fetch"
	PermEmitEvent         = "events.emit"
	PermHttpServe         = "http.serve"
	PermHttpSSE           = "http.sse"
	PermServerRead        = "server.read"
	PermNotificationsSend = "notifications.send"
	PermUsersRead         = "users.read"
	PermUsersModerate     = "users.moderate"
	PermFediversePost     = "fediverse.post"
)

// ChatSendKind distinguishes how a plugin asked to post to chat. All sends
// post under the plugin's own chat identity — provisioned by the host at
// install time as a chat user with IsBot=true and DisplayName=plugin name.
// Plugins cannot post under arbitrary or other users' identities; to use a
// different chat name, ship as a different plugin.
type ChatSendKind int

const (
	ChatSendBot    ChatSendKind = iota // owncast.chat.send — posts as the plugin's bot
	ChatSendAction                     // owncast.chat.sendAction — italic, "/me" style
	ChatSendSystem                     // owncast.chat.system — server-announcement style, no user identity, body is HTML
)

// ChatSendRequest is everything the host needs to dispatch a chat post made
// by a plugin. The host looks up the plugin's bot access token and posts
// through Owncast's normal chat pipeline using that token.
type ChatSendRequest struct {
	PluginName string
	Kind       ChatSendKind
	Text       string
}

// StreamInfo is what owncast.stream.current() returns to a plugin. Wired to
// real Owncast state in production; in the PoC the demo binary fills it in.
type StreamInfo struct {
	Online       bool   `json:"online"`
	Title        string `json:"title,omitempty"`
	Summary      string `json:"summary,omitempty"`
	Viewers      int    `json:"viewers"`
	StartedAt    string `json:"startedAt,omitempty"` // ISO-8601, empty when offline
	LatencyLevel int    `json:"latencyLevel,omitempty"`
}

// ServerInfo is what owncast.server.info() returns to a plugin.
type ServerInfo struct {
	Name           string `json:"name,omitempty"`
	URL            string `json:"url,omitempty"`
	Summary        string `json:"summary,omitempty"`
	WelcomeMessage string `json:"welcomeMessage,omitempty"`
	Version        string `json:"version,omitempty"`
}

// SocialHandle is one entry returned by owncast.server.socials().
type SocialHandle struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
	Icon     string `json:"icon,omitempty"`
}

// FederationInfo is what owncast.server.federation() returns.
type FederationInfo struct {
	Enabled   bool   `json:"enabled"`
	Username  string `json:"username,omitempty"`
	IsPrivate bool   `json:"isPrivate,omitempty"`
}

// FediversePayload is what a plugin passes to owncast.notifications.fediverse.
type FediversePayload struct {
	Type  string `json:"type"` // "follow", "like", "repost", or a custom string
	Body  string `json:"body"`
	Image string `json:"image,omitempty"`
	Link  string `json:"link,omitempty"`
}

// HostChatMessage is the shape returned by ChatHistory. Wider than the
// onChatMessage event payload — production wires this to whatever the chat
// repository hands back; tests construct it directly.
type HostChatMessage struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}

// BrowserPushPayload is what a plugin asks Owncast to send via the
// configured browser push channel.
type BrowserPushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	URL   string `json:"url,omitempty"`
}

// HostUser is the shape returned by Users() / UserGet().
type HostUser struct {
	ID              string   `json:"id"`
	DisplayName     string   `json:"displayName"`
	PreviousNames   []string `json:"previousNames,omitempty"`
	CreatedAt       string   `json:"createdAt,omitempty"`
	DisabledAt      string   `json:"disabledAt,omitempty"` // ISO-8601 if banned, empty otherwise
	Scopes          []string `json:"scopes,omitempty"`
	IsBot           bool     `json:"isBot,omitempty"`
	IsAuthenticated bool     `json:"isAuthenticated,omitempty"`
}

// HostChatClient is the shape returned by ChatClients() — a connected chat
// session, not a User (one user may have several clients).
type HostChatClient struct {
	ID           uint64 `json:"id"`
	UserID       string `json:"userId,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	ConnectedAt  string `json:"connectedAt,omitempty"`
	UserAgent    string `json:"userAgent,omitempty"`
	IPAddress    string `json:"ipAddress,omitempty"`
	MessageCount int    `json:"messageCount"`
}

// UploadResult is what storage.upload returns to the plugin.
type UploadResult struct {
	URL string `json:"url"`
}

// HostEnv is everything host functions need to do their job. Function-pointer
// fields are wired by the host (the production Owncast binary, the demo
// binary, or the test runner); each host function reads them lazily at call
// time so that fields wired post-load (Emit) work correctly.
type HostEnv struct {
	KV              kv.Store
	OnChat          func(ChatSendRequest)
	Emit            func(ctx context.Context, eventType string, payload any)
	StreamCurrent   func() StreamInfo
	ServerInfo      func() ServerInfo
	ChatHistory     func(limit int) []HostChatMessage
	ChatClients     func() []HostChatClient                  // chat.history
	DeleteMessage   func(pluginName, messageID string)       // chat.moderate
	KickClient      func(pluginName string, clientID uint64) // chat.moderate
	SendDiscord     func(pluginName, text string)            // notifications.send
	SendBrowserPush func(pluginName string, p BrowserPushPayload)
	Users           func() []HostUser                                            // users.read
	UserGet         func(id string) (HostUser, bool)                             // users.read
	SetUserEnabled  func(pluginName, userID string, enabled bool, reason string) // users.moderate
	BanIP           func(pluginName, ip string)                                  // users.moderate
	UploadStorage   func(pluginName, name string, data []byte) (string, error)   // storage.upload
	Socials         func() []SocialHandle                                        // server.read
	Federation      func() FederationInfo                                        // server.read
	SendFediverse   func(pluginName string, p FediversePayload)                  // notifications.send
	SendChatTo      func(pluginName string, clientID uint64, text string)        // chat.send
	// PostFediverse publishes a public, text-only note to the fediverse
	// on the streamer's behalf. Returns the resulting post URL. The host is
	// responsible for rate-limiting (max 5/hour per plugin by default) and
	// for honoring the admin's "disable plugin fediverse posting" toggle.
	// fediverse.post permission required.
	PostFediverse func(pluginName, text string) (url string, err error)
	// IsAuthenticated is forwarded to plugin.Server (which uses it both to
	// gate admin paths and to populate req.authenticated). Optional; nil
	// means "no auth available" — admin paths always return 401.
	IsAuthenticated func(r *http.Request) bool
	// GetRequestUser returns the User the request came from when the request
	// carries a user-token (not admin Basic Auth). Plugins see this in
	// req.user. Optional; nil → req.user is always omitted.
	GetRequestUser func(r *http.Request) *HostUser
	// SSE fans plugin-published events out to browser clients connected to
	// the plugin's host-owned event stream. The plugin only publishes (via
	// the owncast_sse_send host function, gated by http.sse); the host owns
	// the long-lived connections. Optional; nil → owncast.sse.send is a
	// no-op even if the plugin declared http.sse.
	SSE *SSEHub
}

// BuildHostFunctions returns the list of extism host functions a single
// plugin should be granted, based on its declared permissions. A plugin
// only sees imports for permissions it declared; importing anything else
// will fail to link at instantiation time.
func BuildHostFunctions(env *HostEnv, manifest *Manifest) []extism.HostFunction {
	var fns []extism.HostFunction
	granted := stringSet(manifest.Permissions)

	if granted[PermStorageKV] {
		ns := env.KV.Namespace(manifest.Name)
		fns = append(fns, hostKVGet(ns), hostKVSet(ns))
	}
	if granted[PermChatSend] {
		fns = append(fns,
			hostSendChat(env.OnChat, manifest.Name),
			hostSendChatAction(env.OnChat, manifest.Name),
			hostSendChatSystem(env.OnChat, manifest.Name),
			hostSendChatTo(env, manifest.Name),
		)
	}
	if granted[PermEmitEvent] {
		fns = append(fns, hostEmitEvent(env, manifest.Name))
	}
	if granted[PermServerRead] {
		fns = append(fns,
			hostStreamCurrent(env),
			hostServerInfo(env),
			hostServerSocials(env),
			hostServerFederation(env),
		)
	}
	if granted[PermChatHistory] {
		fns = append(fns, hostChatHistory(env))
	}
	if granted[PermChatModerate] {
		fns = append(fns,
			hostDeleteMessage(env, manifest.Name),
			hostKickClient(env, manifest.Name),
		)
	}
	if granted[PermNotificationsSend] {
		fns = append(fns,
			hostSendDiscord(env, manifest.Name),
			hostSendBrowserPush(env, manifest.Name),
			hostSendFediverse(env, manifest.Name),
		)
	}
	if granted[PermChatHistory] {
		fns = append(fns, hostChatClients(env))
	}
	if granted[PermUsersRead] {
		fns = append(fns, hostUsersList(env), hostUserGet(env))
	}
	if granted[PermUsersModerate] {
		fns = append(fns,
			hostUserSetEnabled(env, manifest.Name),
			hostBanIP(env, manifest.Name),
		)
	}
	if granted[PermStorageUpload] {
		fns = append(fns, hostStorageUpload(env, manifest.Name))
	}
	if granted[PermFediversePost] {
		fns = append(fns, hostFediversePost(env, manifest.Name))
	}
	if granted[PermHttpSSE] {
		fns = append(fns, hostSSESend(env, manifest.Name))
	}
	return fns
}

// hostSSESend backs owncast.sse.send(channel, event, data). It publishes a
// single Server-Sent-Event to every browser currently connected to this
// plugin's <channel> stream. Fire-and-forget: the call returns as soon as
// the frame is queued to each client, so it never blocks the plugin on a
// slow browser. Requires the http.sse permission.
func hostSSESend(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_sse_send",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			channel, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			event, err := p.ReadString(stack[1])
			if err != nil {
				return
			}
			data, err := p.ReadBytes(stack[2])
			if err != nil {
				return
			}
			if env.SSE == nil {
				return
			}
			env.SSE.Publish(pluginName, channel, event, data)
		},
		[]extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR, extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostFediversePost(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_fediverse_post",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			text, err := p.ReadString(stack[0])
			if err != nil {
				stack[0] = 0
				return
			}
			if env.PostFediverse == nil {
				stack[0] = 0
				return
			}
			url, err := env.PostFediverse(pluginName, text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_fediverse_post from %s: %v\n", pluginName, err)
				stack[0] = 0
				return
			}
			result, err := json.Marshal(map[string]string{"url": url})
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(result)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostChatClients(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_chat_clients",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var clients []HostChatClient
			if env.ChatClients != nil {
				clients = env.ChatClients()
			}
			if clients == nil {
				clients = []HostChatClient{}
			}
			data, err := json.Marshal(clients)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(data)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostUsersList(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_users_list",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var users []HostUser
			if env.Users != nil {
				users = env.Users()
			}
			if users == nil {
				users = []HostUser{}
			}
			data, err := json.Marshal(users)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(data)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostUserGet(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_user_get",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			id, err := p.ReadString(stack[0])
			if err != nil {
				stack[0] = 0
				return
			}
			if env.UserGet == nil {
				stack[0] = 0
				return
			}
			user, ok := env.UserGet(id)
			if !ok {
				stack[0] = 0
				return
			}
			data, err := json.Marshal(user)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(data)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostUserSetEnabled(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_user_set_enabled",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			id, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			enabled := stack[1] != 0
			reason, _ := p.ReadString(stack[2])
			if env.SetUserEnabled != nil {
				env.SetUserEnabled(pluginName, id, enabled, reason)
			}
		},
		[]extism.ValueType{extism.ValueTypePTR, extism.ValueTypeI32, extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostBanIP(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_ban_ip",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			ip, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			if env.BanIP != nil {
				env.BanIP(pluginName, ip)
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostServerSocials(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_server_socials",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var socials []SocialHandle
			if env.Socials != nil {
				socials = env.Socials()
			}
			if socials == nil {
				socials = []SocialHandle{}
			}
			data, err := json.Marshal(socials)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(data)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostServerFederation(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_server_federation",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var info FederationInfo
			if env.Federation != nil {
				info = env.Federation()
			}
			data, err := json.Marshal(info)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(data)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendFediverse(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_notify_fediverse",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			payloadBytes, err := p.ReadBytes(stack[0])
			if err != nil {
				return
			}
			var payload FediversePayload
			if err := json.Unmarshal(payloadBytes, &payload); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_notify_fediverse from %s: invalid JSON: %v\n", pluginName, err)
				return
			}
			if env.SendFediverse != nil {
				env.SendFediverse(pluginName, payload)
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendChatTo(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_send_chat_to",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			clientID := stack[0]
			text, err := p.ReadString(stack[1])
			if err != nil {
				return
			}
			if env.SendChatTo != nil {
				env.SendChatTo(pluginName, clientID, text)
			}
		},
		[]extism.ValueType{extism.ValueTypeI64, extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostStorageUpload(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_storage_upload",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			name, err := p.ReadString(stack[0])
			if err != nil {
				stack[0] = 0
				return
			}
			data, err := p.ReadBytes(stack[1])
			if err != nil {
				stack[0] = 0
				return
			}
			if env.UploadStorage == nil {
				stack[0] = 0
				return
			}
			url, err := env.UploadStorage(pluginName, name, data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_storage_upload from %s: %v\n", pluginName, err)
				stack[0] = 0
				return
			}
			result, err := json.Marshal(UploadResult{URL: url})
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(result)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostDeleteMessage(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_delete_message",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			id, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			if env.DeleteMessage != nil {
				env.DeleteMessage(pluginName, id)
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostKickClient(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_kick_client",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			clientID := stack[0]
			if env.KickClient != nil {
				env.KickClient(pluginName, clientID)
			}
		},
		[]extism.ValueType{extism.ValueTypeI64},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendDiscord(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_notify_discord",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			text, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			if env.SendDiscord != nil {
				env.SendDiscord(pluginName, text)
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendBrowserPush(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_notify_browser_push",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			payloadBytes, err := p.ReadBytes(stack[0])
			if err != nil {
				return
			}
			var payload BrowserPushPayload
			if err := json.Unmarshal(payloadBytes, &payload); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_notify_browser_push from %s: invalid JSON: %v\n", pluginName, err)
				return
			}
			if env.SendBrowserPush != nil {
				env.SendBrowserPush(pluginName, payload)
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostChatHistory(env *HostEnv) extism.HostFunction {
	const defaultLimit = 50
	fn := extism.NewHostFunctionWithStack(
		"owncast_chat_history",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			limit := int(int32(stack[0])) //nolint:gosec // G115: truncation is intentional; non-positive results fall back to defaultLimit
			if limit <= 0 {
				limit = defaultLimit
			}
			var msgs []HostChatMessage
			if env.ChatHistory != nil {
				msgs = env.ChatHistory(limit)
			}
			if msgs == nil {
				msgs = []HostChatMessage{}
			}
			data, err := json.Marshal(msgs)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(data)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{extism.ValueTypeI32},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostEmitEvent(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_emit_event",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			eventType, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			payloadBytes, err := p.ReadBytes(stack[1])
			if err != nil {
				return
			}
			var payload any
			if len(payloadBytes) > 0 {
				if err := json.Unmarshal(payloadBytes, &payload); err != nil {
					fmt.Fprintf(os.Stderr, "owncast_emit_event from %s: invalid JSON payload: %v\n", pluginName, err)
					return
				}
			}
			if env.Emit == nil {
				return
			}
			env.Emit(ctx, eventType, payload)
		},
		[]extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendChat(sink func(ChatSendRequest), pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_send_chat",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			text, err := p.ReadString(stack[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_send_chat: read string: %v\n", err)
				return
			}
			if sink != nil {
				sink(ChatSendRequest{PluginName: pluginName, Kind: ChatSendBot, Text: text})
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendChatSystem(sink func(ChatSendRequest), pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_send_chat_system",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			body, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			if sink != nil {
				sink(ChatSendRequest{PluginName: pluginName, Kind: ChatSendSystem, Text: body})
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendChatAction(sink func(ChatSendRequest), pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_send_chat_action",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			text, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			if sink != nil {
				sink(ChatSendRequest{PluginName: pluginName, Kind: ChatSendAction, Text: text})
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostStreamCurrent(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_stream_current",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var info StreamInfo
			if env.StreamCurrent != nil {
				info = env.StreamCurrent()
			}
			data, err := json.Marshal(info)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(data)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostServerInfo(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_server_info",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var info ServerInfo
			if env.ServerInfo != nil {
				info = env.ServerInfo()
			}
			data, err := json.Marshal(info)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(data)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostKVGet(ns kv.Namespace) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_kv_get",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			key, err := p.ReadString(stack[0])
			if err != nil {
				stack[0] = 0
				return
			}
			val, err := ns.Get(key)
			if err != nil || val == nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(val)
			if err != nil {
				stack[0] = 0
				return
			}
			stack[0] = offset
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{extism.ValueTypePTR},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostKVSet(ns kv.Namespace) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_kv_set",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			key, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			val, err := p.ReadBytes(stack[1])
			if err != nil {
				return
			}
			if err := ns.Set(key, val); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_kv_set: %v\n", err)
			}
		},
		[]extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}
