package webserver

import (
	"net/http"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/storage"
	"github.com/owncast/owncast/utils"
	fediverseauth "github.com/owncast/owncast/webserver/handlers/auth/fediverse"
	"github.com/owncast/owncast/webserver/handlers/auth/indieauth"
	"github.com/owncast/owncast/webserver/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *webServer) setupRoutes() {
	s.setupWebAssetRoutes()
	s.setupInternalAPIRoutes()
	s.setupAdminAPIRoutes()
	s.setupExternalThirdPartyAPIRoutes()
	s.setupModerationAPIRoutes()

	s.router.HandleFunc("/hls/", s.handlers.HandleHLSRequest)

	// websocket
	s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.HandleClientConnection(w, r)
	})

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// This is a hack because Prometheus enables this endpoint by default
		// due to its use of expvar and we do not want this exposed.
		if r.URL.Path == "/debug/vars" {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if r.URL.Path == "/embed/chat/" || r.URL.Path == "/embed/chat" {
			// Redirect /embed/chat
			http.Redirect(w, r, "/embed/chat/readonly", http.StatusTemporaryRedirect)
		} else {
			s.handlers.IndexHandler(w, r)
			// s.ServeHTTP(w, r)
		}
	})

	// ActivityPub has its own router
	activitypub.Start(data.GetDatastore(), s.router)
}

func (s *webServer) setupWebAssetRoutes() {
	c := config.GetConfig()

	// The admin web app.
	s.router.HandleFunc("/admin/", middleware.RequireAdminAuth(s.handlers.IndexHandler))

	// Images
	s.router.HandleFunc("/thumbnail.jpg", s.handlers.GetThumbnail)
	s.router.HandleFunc("/preview.gif", s.handlers.GetPreview)
	s.router.HandleFunc("/logo", s.handlers.GetLogo)

	// Custom Javascript
	s.router.HandleFunc("/customjavascript", s.handlers.ServeCustomJavascript)

	// Return a single emoji image.
	s.router.HandleFunc(config.EmojiDir, s.handlers.GetCustomEmojiImage)

	// return the logo

	// return a logo that's compatible with external social networks
	s.router.HandleFunc("/logo/external", s.handlers.GetCompatibleLogo)

	// robots.txt
	s.router.HandleFunc("/robots.txt", s.handlers.GetRobotsDotTxt)

	// Optional public static files
	s.router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(c.PublicFilesPath))))
}

func (s *webServer) setupInternalAPIRoutes() {
	// Internal APIs

	// status of the system
	s.router.HandleFunc("/api/status", s.handlers.GetStatus)

	// custom emoji supported in the chat
	s.router.HandleFunc("/api/emoji", s.handlers.GetCustomEmojiList)

	// chat history api
	s.router.HandleFunc("/api/chat", middleware.RequireUserAccessToken(s.handlers.GetChatMessages))

	// web config api
	s.router.HandleFunc("/api/config", s.handlers.GetWebConfig)

	// return the YP protocol data
	s.router.HandleFunc("/api/yp", s.handlers.GetYPResponse)

	// list of all social platforms
	s.router.HandleFunc("/api/socialplatforms", s.handlers.GetAllSocialPlatforms)

	// return the list of video variants available
	s.router.HandleFunc("/api/video/variants", s.handlers.GetVideoStreamOutputVariants)

	// tell the backend you're an active viewer
	s.router.HandleFunc("/api/ping", s.handlers.Ping)

	// register a new chat user
	s.router.HandleFunc("/api/chat/register", s.handlers.RegisterAnonymousChatUser)

	// return remote follow details
	s.router.HandleFunc("/api/remotefollow", s.handlers.RemoteFollow)

	// return followers
	s.router.HandleFunc("/api/followers", middleware.HandlePagination(s.handlers.GetFollowers))

	// save client video playback metrics
	s.router.HandleFunc("/api/metrics/playback", s.handlers.ReportPlaybackMetrics)

	// Register for notifications
	s.router.HandleFunc("/api/notifications/register", middleware.RequireUserAccessToken(s.handlers.RegisterForLiveNotifications))

	// Start auth flow
	s.router.HandleFunc("/api/auth/indieauth", middleware.RequireUserAccessToken(indieauth.StartAuthFlow))
	s.router.HandleFunc("/api/auth/indieauth/callback", indieauth.HandleRedirect)
	s.router.HandleFunc("/api/auth/provider/indieauth", indieauth.HandleAuthEndpoint)

	s.router.HandleFunc("/api/auth/fediverse", middleware.RequireUserAccessToken(fediverseauth.RegisterFediverseOTPRequest))
	s.router.HandleFunc("/api/auth/fediverse/verify", fediverseauth.VerifyFediverseOTPRequest)
}

func (s *webServer) setupAdminAPIRoutes() {
	// Current inbound broadcaster
	s.router.HandleFunc("/api/admin/status", middleware.RequireAdminAuth(s.handlers.GetAdminStatus))

	// Disconnect inbound stream
	s.router.HandleFunc("/api/admin/disconnect", middleware.RequireAdminAuth(s.handlers.DisconnectInboundConnection))

	// Server config
	s.router.HandleFunc("/api/admin/serverconfig", middleware.RequireAdminAuth(s.handlers.GetServerConfig))

	// Get viewer count over time
	s.router.HandleFunc("/api/admin/viewersOverTime", middleware.RequireAdminAuth(s.handlers.GetViewersOverTime))

	// Get active viewers
	s.router.HandleFunc("/api/admin/viewers", middleware.RequireAdminAuth(s.handlers.GetActiveViewers))

	// Get hardware stats
	s.router.HandleFunc("/api/admin/hardwarestats", middleware.RequireAdminAuth(s.handlers.GetHardwareStats))

	// Get a a detailed list of currently connected chat clients
	s.router.HandleFunc("/api/admin/chat/clients", middleware.RequireAdminAuth(s.handlers.GetConnectedChatClients))

	// Get all logs
	s.router.HandleFunc("/api/admin/logs", middleware.RequireAdminAuth(s.handlers.GetLogs))

	// Get warning/error logs
	s.router.HandleFunc("/api/admin/logs/warnings", middleware.RequireAdminAuth(s.handlers.GetWarnings))

	// Get all chat messages for the admin, unfiltered.
	s.router.HandleFunc("/api/admin/chat/messages", middleware.RequireAdminAuth(s.handlers.GetAdminChatMessages))

	// Update chat message visibility
	s.router.HandleFunc("/api/admin/chat/messagevisibility", middleware.RequireAdminAuth(s.handlers.UpdateMessageVisibility))

	// Enable/disable a user
	s.router.HandleFunc("/api/admin/chat/users/setenabled", middleware.RequireAdminAuth(s.handlers.UpdateUserEnabled))

	// Ban/unban an IP address
	s.router.HandleFunc("/api/admin/chat/users/ipbans/create", middleware.RequireAdminAuth(s.handlers.BanIPAddress))

	// Remove an IP address ban
	s.router.HandleFunc("/api/admin/chat/users/ipbans/remove", middleware.RequireAdminAuth(s.handlers.UnBanIPAddress))

	// Return all the banned IP addresses
	s.router.HandleFunc("/api/admin/chat/users/ipbans", middleware.RequireAdminAuth(s.handlers.GetIPAddressBans))

	// Get a list of disabled users
	s.router.HandleFunc("/api/admin/chat/users/disabled", middleware.RequireAdminAuth(s.handlers.GetDisabledUsers))

	// Set moderator status for a user
	s.router.HandleFunc("/api/admin/chat/users/setmoderator", middleware.RequireAdminAuth(s.handlers.UpdateUserModerator))

	// Get a list of moderator users
	s.router.HandleFunc("/api/admin/chat/users/moderators", middleware.RequireAdminAuth(s.handlers.GetModerators))

	// return followers
	s.router.HandleFunc("/api/admin/followers", middleware.RequireAdminAuth(middleware.HandlePagination(s.handlers.GetFollowers)))

	// Get a list of pending follow requests
	s.router.HandleFunc("/api/admin/followers/pending", middleware.RequireAdminAuth(s.handlers.GetPendingFollowRequests))

	// Get a list of rejected or blocked follows
	s.router.HandleFunc("/api/admin/followers/blocked", middleware.RequireAdminAuth(s.handlers.GetBlockedAndRejectedFollowers))

	// Set the following state of a follower or follow request.
	s.router.HandleFunc("/api/admin/followers/approve", middleware.RequireAdminAuth(s.handlers.ApproveFollower))

	// Upload custom emoji
	s.router.HandleFunc("/api/admin/emoji/upload", middleware.RequireAdminAuth(s.handlers.UploadCustomEmoji))

	// Delete custom emoji
	s.router.HandleFunc("/api/admin/emoji/delete", middleware.RequireAdminAuth(s.handlers.DeleteCustomEmoji))

	// Update config values

	// Change the current streaming key in memory
	s.router.HandleFunc("/api/admin/config/adminpass", middleware.RequireAdminAuth(s.handlers.SetAdminPassword))

	//  Set an array of valid stream keys
	s.router.HandleFunc("/api/admin/config/streamkeys", middleware.RequireAdminAuth(s.handlers.SetStreamKeys))

	// Change the extra page content in memory
	s.router.HandleFunc("/api/admin/config/pagecontent", middleware.RequireAdminAuth(s.handlers.SetExtraPageContent))

	// Stream title
	s.router.HandleFunc("/api/admin/config/streamtitle", middleware.RequireAdminAuth(s.handlers.SetStreamTitle))

	// Server name
	s.router.HandleFunc("/api/admin/config/name", middleware.RequireAdminAuth(s.handlers.SetServerName))

	// Server summary
	s.router.HandleFunc("/api/admin/config/serversummary", middleware.RequireAdminAuth(s.handlers.SetServerSummary))

	// Offline message
	s.router.HandleFunc("/api/admin/config/offlinemessage", middleware.RequireAdminAuth(s.handlers.SetCustomOfflineMessage))

	// Server welcome message
	s.router.HandleFunc("/api/admin/config/welcomemessage", middleware.RequireAdminAuth(s.handlers.SetServerWelcomeMessage))

	// Disable chat
	s.router.HandleFunc("/api/admin/config/chat/disable", middleware.RequireAdminAuth(s.handlers.SetChatDisabled))

	// Disable chat user join messages
	s.router.HandleFunc("/api/admin/config/chat/joinmessagesenabled", middleware.RequireAdminAuth(s.handlers.SetChatJoinMessagesEnabled))

	// Enable/disable chat established user mode
	s.router.HandleFunc("/api/admin/config/chat/establishedusermode", middleware.RequireAdminAuth(s.handlers.SetEnableEstablishedChatUserMode))

	// Set chat usernames that are not allowed
	s.router.HandleFunc("/api/admin/config/chat/forbiddenusernames", middleware.RequireAdminAuth(s.handlers.SetForbiddenUsernameList))

	// Set the suggested chat usernames that will be assigned automatically
	s.router.HandleFunc("/api/admin/config/chat/suggestedusernames", middleware.RequireAdminAuth(s.handlers.SetSuggestedUsernameList))

	// Set video codec
	s.router.HandleFunc("/api/admin/config/video/codec", middleware.RequireAdminAuth(s.handlers.SetVideoCodec))

	// Set style/color/css values
	s.router.HandleFunc("/api/admin/config/appearance", middleware.RequireAdminAuth(s.handlers.SetCustomColorVariableValues))

	// Return all webhooks
	s.router.HandleFunc("/api/admin/webhooks", middleware.RequireAdminAuth(s.handlers.GetWebhooks))

	// Delete a single webhook
	s.router.HandleFunc("/api/admin/webhooks/delete", middleware.RequireAdminAuth(s.handlers.DeleteWebhook))

	// Create a single webhook
	s.router.HandleFunc("/api/admin/webhooks/create", middleware.RequireAdminAuth(s.handlers.CreateWebhook))

	// Get all access tokens
	s.router.HandleFunc("/api/admin/accesstokens", middleware.RequireAdminAuth(s.handlers.GetExternalAPIUsers))

	// Delete a single access token
	s.router.HandleFunc("/api/admin/accesstokens/delete", middleware.RequireAdminAuth(s.handlers.DeleteExternalAPIUser))

	// Create a single access token
	s.router.HandleFunc("/api/admin/accesstokens/create", middleware.RequireAdminAuth(s.handlers.CreateExternalAPIUser))

	// Return the auto-update features that are supported for this instance.
	s.router.HandleFunc("/api/admin/update/options", middleware.RequireAdminAuth(s.handlers.AutoUpdateOptions))

	// Begin the auto update
	s.router.HandleFunc("/api/admin/update/start", middleware.RequireAdminAuth(s.handlers.AutoUpdateStart))

	// Force quit the service to restart it
	s.router.HandleFunc("/api/admin/update/forcequit", middleware.RequireAdminAuth(s.handlers.AutoUpdateForceQuit))

	// Logo path
	s.router.HandleFunc("/api/admin/config/logo", middleware.RequireAdminAuth(s.handlers.SetLogo))

	// Server tags
	s.router.HandleFunc("/api/admin/config/tags", middleware.RequireAdminAuth(s.handlers.SetTags))

	// ffmpeg
	s.router.HandleFunc("/api/admin/config/ffmpegpath", middleware.RequireAdminAuth(s.handlers.SetFfmpegPath))

	// Server http port
	s.router.HandleFunc("/api/admin/config/webserverport", middleware.RequireAdminAuth(s.handlers.SetWebServerPort))

	// Server http listen address
	s.router.HandleFunc("/api/admin/config/webserverip", middleware.RequireAdminAuth(s.handlers.SetWebServerIP))

	// Server rtmp port
	s.router.HandleFunc("/api/admin/config/rtmpserverport", middleware.RequireAdminAuth(s.handlers.SetRTMPServerPort))

	// Websocket host override
	s.router.HandleFunc("/api/admin/config/sockethostoverride", middleware.RequireAdminAuth(s.handlers.SetSocketHostOverride))

	// Custom video serving endpoint
	s.router.HandleFunc("/api/admin/config/videoservingendpoint", middleware.RequireAdminAuth(s.handlers.SetVideoServingEndpoint))

	// Is server marked as NSFW
	s.router.HandleFunc("/api/admin/config/nsfw", middleware.RequireAdminAuth(s.handlers.SetNSFW))

	// directory enabled
	s.router.HandleFunc("/api/admin/config/directoryenabled", middleware.RequireAdminAuth(s.handlers.SetDirectoryEnabled))

	// social handles
	s.router.HandleFunc("/api/admin/config/socialhandles", middleware.RequireAdminAuth(s.handlers.SetSocialHandles))

	// set the number of video segments and duration per segment in a playlist
	s.router.HandleFunc("/api/admin/config/video/streamlatencylevel", middleware.RequireAdminAuth(s.handlers.SetStreamLatencyLevel))

	// set an array of video output configurations
	s.router.HandleFunc("/api/admin/config/video/streamoutputvariants", middleware.RequireAdminAuth(s.handlers.SetStreamOutputVariants))

	// set s3 configuration
	s.router.HandleFunc("/api/admin/config/s3", middleware.RequireAdminAuth(s.handlers.SetS3Configuration))

	// set server url
	s.router.HandleFunc("/api/admin/config/serverurl", middleware.RequireAdminAuth(s.handlers.SetServerURL))

	// reset the YP registration
	s.router.HandleFunc("/api/admin/yp/reset", middleware.RequireAdminAuth(s.handlers.ResetYPRegistration))

	// set external action links
	s.router.HandleFunc("/api/admin/config/externalactions", middleware.RequireAdminAuth(s.handlers.SetExternalActions))

	// set custom style css
	s.router.HandleFunc("/api/admin/config/customstyles", middleware.RequireAdminAuth(s.handlers.SetCustomStyles))

	// set custom style javascript
	s.router.HandleFunc("/api/admin/config/customjavascript", middleware.RequireAdminAuth(s.handlers.SetCustomJavascript))

	// Video playback metrics
	s.router.HandleFunc("/api/admin/metrics/video", middleware.RequireAdminAuth(s.handlers.GetVideoPlaybackMetrics))

	// Is the viewer count hidden from viewers
	s.router.HandleFunc("/api/admin/config/hideviewercount", middleware.RequireAdminAuth(s.handlers.SetHideViewerCount))

	// set disabling of search indexing
	s.router.HandleFunc("/api/admin/config/disablesearchindexing", middleware.RequireAdminAuth(s.handlers.SetDisableSearchIndexing))

	// enable/disable federation features
	s.router.HandleFunc("/api/admin/config/federation/enable", middleware.RequireAdminAuth(s.handlers.SetFederationEnabled))

	// set if federation activities are private
	s.router.HandleFunc("/api/admin/config/federation/private", middleware.RequireAdminAuth(s.handlers.SetFederationActivityPrivate))

	// set if fediverse engagement appears in chat
	s.router.HandleFunc("/api/admin/config/federation/showengagement", middleware.RequireAdminAuth(s.handlers.SetFederationShowEngagement))

	// set local federated username
	s.router.HandleFunc("/api/admin/config/federation/username", middleware.RequireAdminAuth(s.handlers.SetFederationUsername))

	// set federated go live message
	s.router.HandleFunc("/api/admin/config/federation/livemessage", middleware.RequireAdminAuth(s.handlers.SetFederationGoLiveMessage))

	// Federation blocked domains
	s.router.HandleFunc("/api/admin/config/federation/blockdomains", middleware.RequireAdminAuth(s.handlers.SetFederationBlockDomains))

	// send a public message to the Fediverse from the server's user
	s.router.HandleFunc("/api/admin/federation/send", middleware.RequireAdminAuth(s.handlers.SendFederatedMessage))

	// Return federated activities
	s.router.HandleFunc("/api/admin/federation/actions", middleware.RequireAdminAuth(middleware.HandlePagination(s.handlers.GetFederatedActions)))

	// Prometheus metrics
	s.router.Handle("/api/admin/prometheus", middleware.RequireAdminAuth(func(rw http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(rw, r)
	}))

	// Configure outbound notification channels.
	s.router.HandleFunc("/api/admin/config/notifications/discord", middleware.RequireAdminAuth(s.handlers.SetDiscordNotificationConfiguration))
	s.router.HandleFunc("/api/admin/config/notifications/browser", middleware.RequireAdminAuth(s.handlers.SetBrowserNotificationConfiguration))
}

func (s *webServer) setupExternalThirdPartyAPIRoutes() {
	// Send a system message to chat
	s.router.HandleFunc("/api/integrations/chat/system", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, s.handlers.SendSystemMessage))

	// Send a system message to a single client
	s.router.HandleFunc(utils.RestEndpoint("/api/integrations/chat/system/client/{clientId}", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, s.handlers.SendSystemMessageToConnectedClient)))

	// Send a user message to chat *NO LONGER SUPPORTED
	s.router.HandleFunc("/api/integrations/chat/user", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, s.handlers.SendUserMessage))

	// Send a message to chat as a specific 3rd party bot/integration based on its access token
	s.router.HandleFunc("/api/integrations/chat/send", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, s.handlers.SendIntegrationChatMessage))

	// Send a user action to chat
	s.router.HandleFunc("/api/integrations/chat/action", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, s.handlers.SendChatAction))

	// Hide chat message
	s.router.HandleFunc("/api/integrations/chat/messagevisibility", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, s.handlers.ExternalUpdateMessageVisibility))

	// Stream title
	s.router.HandleFunc("/api/integrations/streamtitle", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, s.handlers.ExternalSetStreamTitle))

	// Get chat history
	s.router.HandleFunc("/api/integrations/chat", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, s.handlers.ExternalGetChatMessages))

	// Connected clients
	s.router.HandleFunc("/api/integrations/clients", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, s.handlers.ExternalGetConnectedChatClients))
}

func (s *webServer) setupModerationAPIRoutes() {
	// Update chat message visibility
	s.router.HandleFunc("/api/chat/messagevisibility", middleware.RequireUserModerationScopeAccesstoken(s.handlers.UpdateMessageVisibility))

	// Enable/disable a user
	s.router.HandleFunc("/api/chat/users/setenabled", middleware.RequireUserModerationScopeAccesstoken(s.handlers.UpdateUserEnabled))

	// Get a user's details
	s.router.HandleFunc("/api/moderation/chat/user/", middleware.RequireUserModerationScopeAccesstoken(s.handlers.GetUserDetails))
}
