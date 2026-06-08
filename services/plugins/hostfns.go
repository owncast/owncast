package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	extism "github.com/extism/go-sdk"
	"github.com/owncast/owncast/services/plugins/kv"
)

var _ = (*http.Request)(nil) // ensure import retained for the HostEnv signature

// Permission identifiers. The manifest declares which a plugin needs; only
// the corresponding wasm imports are wired into the plugin instance.
const (
	PermStorageKV     = "storage.kv"
	PermStorageUpload = "storage.upload"
	// PermStorageFS grants a plugin a private, sandboxed area on disk
	// (data/plugin-data/<slug>/) it can read, write, list, and delete
	// within. Unlike storage.upload (which publishes browser-accessible
	// files under public/), this storage is server-side only and never
	// exposed over HTTP. A plugin cannot reach outside its own directory
	// or read another plugin's files.
	PermStorageFS    = "storage.fs"
	PermChatSend     = "chat.send"
	PermChatHistory  = "chat.history"
	PermChatModerate = "chat.moderate"
	// PermChatFilter is required for any plugin that subscribes to
	// filterChatMessage. The host refuses to load a plugin whose
	// runtime-declared filter subscriptions include the chat-message
	// event without this permission, so a plugin can't silently start
	// dropping or rewriting chat without the admin's consent.
	PermChatFilter        = "chat.filter"
	PermNetworkFetch      = "network.fetch"
	PermEmitEvent         = "events.emit"
	PermHttpServe         = "http.serve"
	PermHttpSSE           = "http.sse"
	PermServerRead        = "server.read"
	PermNotificationsSend = "notifications.send"
	PermUsersRead         = "users.read"
	PermUsersModerate     = "users.moderate"
	PermFediversePost     = "fediverse.post"
	PermVideoConfigRead   = "videoconfig.read"
	PermVideoConfigWrite  = "videoconfig.write"
	// PermUIModify lets a plugin add UI surfaces to Owncast: admin pages
	// (manifest.admin.pages) and viewer action buttons (manifest.actions).
	// A plugin can serve HTTP on /plugins/<name>/ without this permission
	// (e.g. for headless APIs), but cannot publish anything that shows up
	// in the admin UI or viewer chrome without opting in explicitly.
	PermUIModify = "ui.modify"
)

// resultErrorKey is the JSON key host functions use to return an error
// string to the plugin in their {ok, error?} result envelope.
const resultErrorKey = "error"

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
//
// Two identity fields, not one: PluginSlug is the stable identifier used to
// look up the bot's persistent access token (it's how the chatbot user is
// indexed in the datastore, and how the cache keys), while BotDisplayName is
// the human-readable name shown to chat viewers. The split exists so a plugin
// authored as "Awesome Echo Bot" (display) with slug "awesome-echo-bot"
// (identifier) shows the friendly name in chat instead of the slug.
type ChatSendRequest struct {
	PluginSlug     string
	BotDisplayName string
	Kind           ChatSendKind
	Text           string
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

// StreamBroadcaster is what owncast.stream.broadcaster() returns — details
// about the currently-connected inbound broadcast. Zero-valued when offline.
type StreamBroadcaster struct {
	RemoteAddr string   `json:"remoteAddr,omitempty"`
	Codecs     []string `json:"codecs,omitempty"`
	Resolution string   `json:"resolution,omitempty"`
	Framerate  int      `json:"framerate,omitempty"`
	Bitrates   []int    `json:"bitrates,omitempty"`
}

// StreamVariant is one configured output rendition, part of the VideoConfig
// returned/accepted by owncast.videoConfig.read()/write().
type StreamVariant struct {
	Width         int  `json:"width"`
	Height        int  `json:"height"`
	Framerate     int  `json:"framerate"`
	VideoBitrate  int  `json:"videoBitrate"`
	AudioBitrate  int  `json:"audioBitrate"`
	IsPassthrough bool `json:"isPassthrough"`
}

// VideoConfig is the current output/transcoding configuration returned by
// owncast.videoConfig.read(). These are settable knobs (see VideoConfigUpdate),
// as opposed to read-only inbound-broadcast telemetry (StreamBroadcaster).
// Requires the videoconfig.read permission.
type VideoConfig struct {
	LatencyLevel int             `json:"latencyLevel"`
	Codec        string          `json:"codec"`
	Variants     []StreamVariant `json:"variants"`
}

// VideoConfigUpdate is a partial change passed to owncast.videoConfig.write().
// Nil/omitted fields are left unchanged, so a plugin can set just the latency
// level without disturbing the configured output variants. Requires the
// videoconfig.write permission.
type VideoConfigUpdate struct {
	LatencyLevel *int            `json:"latencyLevel,omitempty"`
	Codec        *string         `json:"codec,omitempty"`
	Variants     []StreamVariant `json:"variants,omitempty"`
}

// SocialHandle is one entry returned by owncast.server.socials().
type SocialHandle struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
	Icon     string `json:"icon,omitempty"`
}

// Emote is one custom chat emote returned by owncast.server.emotes(): the
// `:code:` chat clients substitute and the URL of the image it renders to.
type Emote struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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

// HostChatUser is the sender identity carried by a chat message. It mirrors
// the SDK's ChatUser TypeScript interface (and pluginhost.pluginChatUser on
// the event path) so chat.history() hands a plugin the same nested object
// shape its onChatMessage handler already receives.
type HostChatUser struct {
	ID              string   `json:"id"`
	DisplayName     string   `json:"displayName"`
	IsBot           bool     `json:"isBot,omitempty"`
	IsAuthenticated bool     `json:"isAuthenticated,omitempty"`
	Scopes          []string `json:"scopes,omitempty"`
}

// HostChatMessage is the shape returned by ChatHistory. It matches the
// onChatMessage / chat.message.received event payload: User is the full
// ChatUser object (nil for the rare message with no associated account),
// not a bare display-name string. Production wires this to whatever the chat
// repository hands back; tests construct it directly.
type HostChatMessage struct {
	ID        string        `json:"id"`
	User      *HostChatUser `json:"user,omitempty"`
	Body      string        `json:"body"`
	Timestamp string        `json:"timestamp"`
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

// SSEConnectionEvent is the payload for the sse.connect / sse.disconnect
// events, fired when a browser opens or closes one of the plugin's
// Server-Sent-Events streams. ConnectionID is unique for the life of the host
// process, so a plugin can pair a disconnect with its connect and count
// distinct connections (e.g. one user open in several tabs). User is the
// resolved chat user when the connection carried a chat identity, and is
// omitted for anonymous viewers.
type SSEConnectionEvent struct {
	Channel      string    `json:"channel"`
	ConnectionID uint64    `json:"connectionId"`
	User         *HostUser `json:"user,omitempty"`
}

// TickEvent is the payload for the once-a-second tick event, delivered to
// plugins that subscribe (define onTick). Now is the host's wall-clock time in
// unix milliseconds at the moment the tick fired.
type TickEvent struct {
	Now int64 `json:"now"`
}

// TimerFireEvent is the payload delivered (via the timer.fire event) when a
// host-scheduled timer elapses. ID is the timer's id, which the guest SDK maps
// back to the author's callback.
type TimerFireEvent struct {
	ID uint64 `json:"id"`
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
	Broadcaster     func() StreamBroadcaster // server.read (read-only telemetry)
	Tags            func() []string          // server.read
	VideoConfig     func() VideoConfig       // videoconfig.read
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
	// Sandboxed per-plugin filesystem (storage.fs). Each plugin sees only
	// its own directory under data/plugin-data/<slug>/; the host rejects any
	// path that escapes it. Paths are relative to that root.
	FSRead     func(pluginName, path string) ([]byte, error) // storage.fs
	FSWrite    func(pluginName, path string, data []byte) error
	FSList     func(pluginName, dir string) ([]string, error)
	FSDelete   func(pluginName, path string) error
	FSExists   func(pluginName, path string) (bool, error)
	Socials    func() []SocialHandle // server.read
	Emotes     func() []Emote        // server.read
	Federation func() FederationInfo // server.read
	// ConfigValue resolves an admin-set override for one of the plugin's
	// manifest-declared config keys (owncast.config.get). Returns the override
	// value and true when the admin has set one; false to fall back to the
	// manifest's declared default. Optional; nil → defaults only (the common
	// case until an admin edits the value).
	ConfigValue func(pluginName, key string) (any, bool)
	// WriteVideoConfig applies a partial video/transcoding configuration
	// change. Returns an error the plugin can see if the host rejects the
	// config (e.g. an invalid variant). videoconfig.write permission required.
	WriteVideoConfig func(pluginName string, u VideoConfigUpdate) error
	SendFediverse    func(pluginName string, p FediversePayload)           // notifications.send
	SendChatTo       func(pluginName string, clientID uint64, text string) // chat.send
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
	// OnSSESend, when set, is invoked for every owncast.sse.send in addition to
	// (and independently of) SSE delivery. It exists so the test harness can
	// observe SSE output, which otherwise vanishes when no browser client is
	// subscribed. Production leaves it nil. Optional.
	OnSSESend func(pluginName, channel, event string, data []byte)
	// Timer schedules host-driven callbacks (owncast.timer.*). Ambient: every
	// plugin gets the host functions, since a plugin can't setTimeout in the
	// sandbox. Optional; nil → owncast_timer_set reports success but never
	// fires (used by the test harness, which simulates fires via events).
	Timer *TimerHub
}

// BuildHostFunctions returns the list of extism host functions a single
// plugin should be granted, based on its declared permissions. A plugin
// only sees imports for permissions it declared; importing anything else
// will fail to link at instantiation time.
func BuildHostFunctions(env *HostEnv, manifest *Manifest, assetsFS fs.FS) []extism.HostFunction {
	var fns []extism.HostFunction
	granted := stringSet(manifest.Permissions)

	if granted[PermStorageKV] {
		ns := env.KV.Namespace(manifest.Slug)
		fns = append(fns, hostKVGet(ns), hostKVSet(ns))
	}
	if granted[PermChatSend] {
		// Chat send fns capture both slug and chat display name in
		// their closure: slug routes to the right bot user, display
		// name is what chat viewers see.
		chatDisplay := manifest.ChatDisplayName()
		fns = append(fns,
			hostSendChat(env.OnChat, manifest.Slug, chatDisplay),
			hostSendChatAction(env.OnChat, manifest.Slug, chatDisplay),
			hostSendChatSystem(env.OnChat, manifest.Slug, chatDisplay),
			hostSendChatTo(env, manifest.Slug),
		)
	}
	if granted[PermEmitEvent] {
		fns = append(fns, hostEmitEvent(env, manifest.Slug))
	}
	if granted[PermServerRead] {
		fns = append(fns,
			hostStreamCurrent(env),
			hostServerInfo(env),
			hostServerSocials(env),
			hostServerEmotes(env),
			hostServerFederation(env),
			hostStreamBroadcaster(env),
			hostServerTags(env),
		)
	}
	if granted[PermChatHistory] {
		fns = append(fns, hostChatHistory(env))
	}
	if granted[PermChatModerate] {
		fns = append(fns,
			hostDeleteMessage(env, manifest.Slug),
			hostKickClient(env, manifest.Slug),
		)
	}
	if granted[PermNotificationsSend] {
		fns = append(fns,
			hostSendDiscord(env, manifest.Slug),
			hostSendBrowserPush(env, manifest.Slug),
			hostSendFediverse(env, manifest.Slug),
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
			hostUserSetEnabled(env, manifest.Slug),
			hostBanIP(env, manifest.Slug),
		)
	}
	fns = append(fns, storageHostFunctions(env, manifest, granted)...)
	if granted[PermFediversePost] {
		fns = append(fns, hostFediversePost(env, manifest.Slug))
	}
	if granted[PermHttpSSE] {
		fns = append(fns, hostSSESend(env, manifest.Slug))
	}

	// Timers are ambient (no permission): a plugin can't setTimeout in the
	// sandbox, so scheduling is a baseline capability. The act of scheduling
	// is benign — whatever the callback does still needs its own permissions —
	// and TimerHub's per-plugin caps bound abuse.
	fns = append(fns, hostTimerSet(env, manifest.Slug), hostTimerClear(env, manifest.Slug))

	// Config is ambient too: reading the plugin's own manifest-declared config
	// (admin override, else declared default) is benign and needs no grant.
	fns = append(fns, hostConfigGet(env, manifest))
	// Asset reading is ambient: a plugin reads only files it shipped itself.
	fns = append(fns, hostAssetRead(assetsFS))
	if granted[PermUIModify] {
		fns = append(fns,
			hostAddActions(env, manifest),
			hostClearActions(env, manifest.Slug),
		)
	}
	fns = append(fns, videoConfigHostFunctions(env, manifest, granted)...)
	return fns
}

// storageHostFunctions returns the storage.upload / storage.fs host functions
// a plugin is granted. Split out of BuildHostFunctions to keep that function's
// cyclomatic complexity in check.
func storageHostFunctions(env *HostEnv, manifest *Manifest, granted map[string]bool) []extism.HostFunction {
	var fns []extism.HostFunction
	if granted[PermStorageUpload] {
		fns = append(fns, hostStorageUpload(env, manifest.Slug))
	}
	if granted[PermStorageFS] {
		fns = append(fns,
			hostFSRead(env, manifest.Slug),
			hostFSWrite(env, manifest.Slug),
			hostFSList(env, manifest.Slug),
			hostFSDelete(env, manifest.Slug),
			hostFSExists(env, manifest.Slug),
		)
	}
	return fns
}

// videoConfigHostFunctions returns the videoconfig.read / videoconfig.write
// host functions a plugin is granted. Split out of BuildHostFunctions to keep
// that function's cyclomatic complexity in check.
func videoConfigHostFunctions(env *HostEnv, manifest *Manifest, granted map[string]bool) []extism.HostFunction {
	var fns []extism.HostFunction
	if granted[PermVideoConfigRead] {
		fns = append(fns, hostVideoConfigRead(env))
	}
	if granted[PermVideoConfigWrite] {
		fns = append(fns, hostVideoConfigWrite(env, manifest.Slug))
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
			if env.OnSSESend != nil {
				env.OnSSESend(pluginName, channel, event, data)
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

// hostTimerSet backs owncast.timer.setTimeout / setInterval. The guest passes a
// guest-allocated id, the delay in milliseconds, and whether it repeats. The
// host arms a timer that, on fire, calls the plugin's on_event with a
// timer.fire event carrying the id. Returns 1 on success, 0 if the plugin is
// at its pending-timer cap. A nil Timer (test harness) reports success so the
// guest keeps its callback; fires are then simulated via events.
func hostTimerSet(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_timer_set",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			id := stack[0]
			// Clamp the requested delay to a sane ceiling before narrowing to
			// int64 — bounds the duration math and keeps the conversion in range.
			rawDelayMs := stack[1]
			if rawDelayMs > maxTimerDelayMs {
				rawDelayMs = maxTimerDelayMs
			}
			delayMs := int64(rawDelayMs)
			repeat := stack[2] == 1
			if env.Timer == nil {
				stack[0] = 1
				return
			}
			if env.Timer.Schedule(pluginName, id, delayMs, repeat) {
				stack[0] = 1
			} else {
				stack[0] = 0
			}
		},
		[]extism.ValueType{extism.ValueTypeI64, extism.ValueTypeI64, extism.ValueTypeI32},
		[]extism.ValueType{extism.ValueTypeI32},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

// hostTimerClear backs owncast.timer.clear(id), cancelling a pending timer.
func hostTimerClear(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_timer_clear",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			id := stack[0]
			if env.Timer != nil {
				env.Timer.Clear(pluginName, id)
			}
		},
		[]extism.ValueType{extism.ValueTypeI64},
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

// hostConfigGet returns the effective value of a manifest-declared config key
// as JSON: the admin-set override when present, otherwise the manifest's
// declared default. Returns 0 (→ undefined in the guest) for an unknown key or
// a declared key with no override and no default.
func hostConfigGet(env *HostEnv, manifest *Manifest) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_config_get",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			key, err := p.ReadString(stack[0])
			if err != nil {
				stack[0] = 0
				return
			}
			field, declared := manifest.Config[key]
			if !declared {
				stack[0] = 0
				return
			}
			value := field.Default
			if env.ConfigValue != nil {
				if override, ok := env.ConfigValue(manifest.Slug, key); ok {
					value = override
				}
			}
			if value == nil {
				stack[0] = 0
				return
			}
			data, err := json.Marshal(value)
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

// hostAssetRead backs owncast.assets.read/readText. Reads a file from the
// plugin's bundled assets/ directory. The path must be relative — no ".."
// segments, no leading "/". Returns 0 when assetsFS is nil or the file
// doesn't exist. Ambient — no permission required.
func hostAssetRead(assetsFS fs.FS) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_asset_read",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			path, err := p.ReadString(stack[0])
			if err != nil || assetsFS == nil {
				stack[0] = 0
				return
			}
			if strings.Contains(path, "..") || strings.HasPrefix(path, "/") {
				stack[0] = 0
				return
			}
			data, err := fs.ReadFile(assetsFS, path)
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

func hostServerEmotes(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_server_emotes",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var emotes []Emote
			if env.Emotes != nil {
				emotes = env.Emotes()
			}
			if emotes == nil {
				emotes = []Emote{}
			}
			data, err := json.Marshal(emotes)
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

// hostFSRead backs owncast.fs.read(path). Returns the file's raw bytes,
// or 0 (null to the plugin) when the path is missing, escapes the
// sandbox, or can't be read. Requires the storage.fs permission.
func hostFSRead(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_fs_read",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			path, err := p.ReadString(stack[0])
			if err != nil || env.FSRead == nil {
				stack[0] = 0
				return
			}
			data, err := env.FSRead(pluginName, path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_fs_read from %s: %v\n", pluginName, err)
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

// hostFSWrite backs owncast.fs.write(path, data). Creates parent
// directories as needed and writes the bytes, returning a JSON
// {ok, error?} result so the plugin can react to a rejected write
// (sandbox escape, oversized payload, disk error). Requires storage.fs.
func hostFSWrite(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_fs_write",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			path, err := p.ReadString(stack[0])
			if err != nil {
				stack[0] = 0
				return
			}
			data, err := p.ReadBytes(stack[1])
			if err != nil {
				stack[0] = 0
				return
			}
			result := map[string]any{"ok": true}
			if env.FSWrite == nil {
				result = map[string]any{"ok": false, resultErrorKey: "filesystem unavailable"}
			} else if err := env.FSWrite(pluginName, path, data); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_fs_write from %s: %v\n", pluginName, err)
				result = map[string]any{"ok": false, resultErrorKey: err.Error()}
			}
			out, err := json.Marshal(result)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(out)
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

// hostFSList backs owncast.fs.list(dir). Returns a JSON array of the
// entry names (files and subdirectories) directly inside dir. A missing
// directory lists as empty rather than erroring. Requires storage.fs.
func hostFSList(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_fs_list",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			dir, err := p.ReadString(stack[0])
			if err != nil || env.FSList == nil {
				stack[0] = 0
				return
			}
			names, err := env.FSList(pluginName, dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_fs_list from %s: %v\n", pluginName, err)
				stack[0] = 0
				return
			}
			if names == nil {
				names = []string{}
			}
			data, err := json.Marshal(names)
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

// hostFSDelete backs owncast.fs.delete(path). Removes a single file or
// empty directory, returning a JSON {ok, error?} result. Requires
// storage.fs.
func hostFSDelete(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_fs_delete",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			path, err := p.ReadString(stack[0])
			if err != nil {
				stack[0] = 0
				return
			}
			result := map[string]any{"ok": true}
			if env.FSDelete == nil {
				result = map[string]any{"ok": false, resultErrorKey: "filesystem unavailable"}
			} else if err := env.FSDelete(pluginName, path); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_fs_delete from %s: %v\n", pluginName, err)
				result = map[string]any{"ok": false, resultErrorKey: err.Error()}
			}
			out, err := json.Marshal(result)
			if err != nil {
				stack[0] = 0
				return
			}
			offset, err := p.WriteBytes(out)
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

// hostFSExists backs owncast.fs.exists(path). Returns 1 if the path
// exists inside the sandbox, 0 otherwise (including on a sandbox-escape
// attempt or stat error). Requires storage.fs.
func hostFSExists(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_fs_exists",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			path, err := p.ReadString(stack[0])
			if err != nil || env.FSExists == nil {
				stack[0] = 0
				return
			}
			exists, err := env.FSExists(pluginName, path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_fs_exists from %s: %v\n", pluginName, err)
				stack[0] = 0
				return
			}
			if exists {
				stack[0] = 1
			} else {
				stack[0] = 0
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{extism.ValueTypeI32},
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

func hostSendChat(sink func(ChatSendRequest), pluginSlug, botDisplayName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_send_chat",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			text, err := p.ReadString(stack[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_send_chat: read string: %v\n", err)
				return
			}
			if sink != nil {
				sink(ChatSendRequest{
					PluginSlug:     pluginSlug,
					BotDisplayName: botDisplayName,
					Kind:           ChatSendBot,
					Text:           text,
				})
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendChatSystem(sink func(ChatSendRequest), pluginSlug, botDisplayName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_send_chat_system",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			body, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			if sink != nil {
				sink(ChatSendRequest{
					PluginSlug:     pluginSlug,
					BotDisplayName: botDisplayName,
					Kind:           ChatSendSystem,
					Text:           body,
				})
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

func hostSendChatAction(sink func(ChatSendRequest), pluginSlug, botDisplayName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_send_chat_action",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			text, err := p.ReadString(stack[0])
			if err != nil {
				return
			}
			if sink != nil {
				sink(ChatSendRequest{
					PluginSlug:     pluginSlug,
					BotDisplayName: botDisplayName,
					Kind:           ChatSendAction,
					Text:           text,
				})
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

func hostStreamBroadcaster(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_stream_broadcaster",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var info StreamBroadcaster
			if env.Broadcaster != nil {
				info = env.Broadcaster()
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

func hostServerTags(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_server_tags",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var tags []string
			if env.Tags != nil {
				tags = env.Tags()
			}
			if tags == nil {
				tags = []string{}
			}
			data, err := json.Marshal(tags)
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

func hostVideoConfigRead(env *HostEnv) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_video_config_read",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			var cfg VideoConfig
			if env.VideoConfig != nil {
				cfg = env.VideoConfig()
			}
			if cfg.Variants == nil {
				cfg.Variants = []StreamVariant{}
			}
			data, err := json.Marshal(cfg)
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

// hostVideoConfigWrite backs owncast.videoConfig.write(config). It applies a
// partial video/transcoding configuration change via the host. Returns a
// JSON {ok, error?} result so the plugin can react to a rejected config.
// Requires the videoconfig.write permission.
func hostVideoConfigWrite(env *HostEnv, pluginName string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_video_config_write",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			payloadBytes, err := p.ReadBytes(stack[0])
			if err != nil {
				stack[0] = 0
				return
			}
			var update VideoConfigUpdate
			if err := json.Unmarshal(payloadBytes, &update); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_video_config_write from %s: invalid JSON: %v\n", pluginName, err)
				stack[0] = 0
				return
			}
			result := map[string]any{"ok": true}
			if env.WriteVideoConfig != nil {
				if err := env.WriteVideoConfig(pluginName, update); err != nil {
					result = map[string]any{"ok": false, resultErrorKey: err.Error()}
				}
			}
			data, err := json.Marshal(result)
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

// RuntimeActionsConfigKey is the reserved key inside a plugin's own
// config namespace that holds action buttons the plugin has added at
// runtime via owncast.actions.add. The plugin host reads this key on
// every /api/config request and returns manifest.actions ++ this list
// to viewers.
const RuntimeActionsConfigKey = "owncast.actions"

// hostAddActions backs owncast.actions.add(actions). Takes a JSON array
// of ActionButton entries, validates each (title required, exactly one
// of url/html, relative URLs rewritten into this plugin's namespace,
// cross-plugin URLs rejected), and appends to the runtime list in the
// plugin's config.
//
// Requires the ui.modify permission; invalid input is logged but not
// surfaced back to the plugin.
func hostAddActions(env *HostEnv, manifest *Manifest) extism.HostFunction {
	pluginSlug := manifest.Slug
	hasHTTPServe := false
	for _, perm := range manifest.Permissions {
		if perm == PermHttpServe {
			hasHTTPServe = true
			break
		}
	}
	fn := extism.NewHostFunctionWithStack(
		"owncast_add_actions",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			payloadBytes, err := p.ReadBytes(stack[0])
			if err != nil {
				return
			}
			var incoming []ActionButton
			if err := json.Unmarshal(payloadBytes, &incoming); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_add_actions from %s: invalid JSON: %v\n", pluginSlug, err)
				return
			}
			normalized, err := validateRuntimeActions(pluginSlug, hasHTTPServe, incoming)
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_add_actions from %s: %v\n", pluginSlug, err)
				return
			}
			if env.KV == nil {
				return
			}
			ns := env.KV.Namespace(pluginSlug)
			var existing []ActionButton
			if raw, err := ns.Get(RuntimeActionsConfigKey); err == nil && len(raw) > 0 {
				_ = json.Unmarshal(raw, &existing)
			}
			combined := make([]ActionButton, 0, len(existing)+len(normalized))
			combined = append(combined, existing...)
			combined = append(combined, normalized...)
			out, err := json.Marshal(combined)
			if err != nil {
				fmt.Fprintf(os.Stderr, "owncast_add_actions from %s: marshal: %v\n", pluginSlug, err)
				return
			}
			if err := ns.Set(RuntimeActionsConfigKey, out); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_add_actions from %s: kv write: %v\n", pluginSlug, err)
			}
		},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{},
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

// hostClearActions backs owncast.actions.clear(). Removes the runtime
// list from the plugin's config so only manifest.actions remain in the
// effective set returned by /api/config.
func hostClearActions(env *HostEnv, pluginSlug string) extism.HostFunction {
	fn := extism.NewHostFunctionWithStack(
		"owncast_clear_actions",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			if env.KV == nil {
				return
			}
			if err := env.KV.Namespace(pluginSlug).Delete(RuntimeActionsConfigKey); err != nil {
				fmt.Fprintf(os.Stderr, "owncast_clear_actions from %s: %v\n", pluginSlug, err)
			}
		},
		nil,
		nil,
	)
	fn.SetNamespace("extism:host/user")
	return fn
}

// validateRuntimeActions checks and rewrites plugin-supplied action
// entries using the same rules as manifest validation, so the runtime
// path can't accept a malformed entry or a cross-plugin URL.
func validateRuntimeActions(pluginSlug string, hasHTTPServe bool, actions []ActionButton) ([]ActionButton, error) {
	pluginPrefix := "/plugins/" + pluginSlug + "/"
	for i := range actions {
		a := &actions[i]
		if a.Title == "" {
			return nil, fmt.Errorf("actions[%d].title is required", i)
		}
		hasURL, hasHTML := a.Url != "", a.Html != ""
		if hasURL == hasHTML {
			return nil, fmt.Errorf("actions[%d]: exactly one of url or html is required", i)
		}
		if !hasURL {
			continue
		}
		if strings.HasPrefix(a.Url, "/") && !strings.HasPrefix(a.Url, "/plugins/") {
			a.Url = pluginPrefix + strings.TrimPrefix(a.Url, "/")
		}
		if strings.HasPrefix(a.Url, pluginPrefix) && !hasHTTPServe {
			return nil, fmt.Errorf("actions[%d].url targets this plugin (%s) but http.serve permission is not declared", i, a.Url)
		}
		if strings.HasPrefix(a.Url, "/plugins/") && !strings.HasPrefix(a.Url, pluginPrefix) {
			return nil, fmt.Errorf("actions[%d].url points at another plugin's namespace: %s", i, a.Url)
		}
	}
	for i := range actions {
		rewritten, err := rewriteActionIcon(pluginPrefix, hasHTTPServe, actions[i].Icon)
		if err != nil {
			return nil, fmt.Errorf("actions[%d].icon: %w", i, err)
		}
		actions[i].Icon = rewritten
	}
	return actions, nil
}
