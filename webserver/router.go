package webserver

import (
	"net/http"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/admin"
	fediverseauth "github.com/owncast/owncast/controllers/auth/fediverse"
	"github.com/owncast/owncast/controllers/auth/indieauth"
	"github.com/owncast/owncast/controllers/moderation"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/middleware"
	"github.com/owncast/owncast/yp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *webServer) setupRoutes() {
	s.router.HandleFunc("/test", s.handlers.HandleTesting)

	s.setupWebAssetRoutes()
	s.setupInternalAPIRoutes()
	s.setupAdminAPIRoutes()
	s.setupExternalThirdPartyAPIRoutes()
	s.setupModerationAPIRoutes()

	s.router.HandleFunc("/hls/", controllers.HandleHLSRequest)

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
			controllers.IndexHandler(w, r)
			// s.ServeHTTP(w, r)
		}
	})

	// ActivityPub has its own router
	activitypub.Start(data.GetDatastore(), s.router)
}

func (s *webServer) setupWebAssetRoutes() {
	// The admin web app.
	s.router.HandleFunc("/admin/", middleware.RequireAdminAuth(controllers.IndexHandler))

	// Images
	s.router.HandleFunc("/thumbnail.jpg", controllers.GetThumbnail)
	s.router.HandleFunc("/preview.gif", controllers.GetPreview)
	s.router.HandleFunc("/logo", controllers.GetLogo)

	// Custom Javascript
	s.router.HandleFunc("/customjavascript", controllers.ServeCustomJavascript)

	// Return a single emoji image.
	s.router.HandleFunc(config.EmojiDir, controllers.GetCustomEmojiImage)

	// return the logo

	// return a logo that's compatible with external social networks
	s.router.HandleFunc("/logo/external", controllers.GetCompatibleLogo)

	// robots.txt
	s.router.HandleFunc("/robots.txt", controllers.GetRobotsDotTxt)

	// Optional public static files
	s.router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(config.PublicFilesPath))))
}

func (s *webServer) setupInternalAPIRoutes() {
	// Internal APIs

	// status of the system
	s.router.HandleFunc("/api/status", controllers.GetStatus)

	// custom emoji supported in the chat
	s.router.HandleFunc("/api/emoji", controllers.GetCustomEmojiList)

	// chat history api
	s.router.HandleFunc("/api/chat", middleware.RequireUserAccessToken(controllers.GetChatMessages))

	// web config api
	s.router.HandleFunc("/api/config", controllers.GetWebConfig)

	// return the YP protocol data
	s.router.HandleFunc("/api/yp", yp.GetYPResponse)

	// list of all social platforms
	s.router.HandleFunc("/api/socialplatforms", controllers.GetAllSocialPlatforms)

	// return the list of video variants available
	s.router.HandleFunc("/api/video/variants", controllers.GetVideoStreamOutputVariants)

	// tell the backend you're an active viewer
	s.router.HandleFunc("/api/ping", controllers.Ping)

	// register a new chat user
	s.router.HandleFunc("/api/chat/register", controllers.RegisterAnonymousChatUser)

	// return remote follow details
	s.router.HandleFunc("/api/remotefollow", controllers.RemoteFollow)

	// return followers
	s.router.HandleFunc("/api/followers", middleware.HandlePagination(controllers.GetFollowers))

	// save client video playback metrics
	s.router.HandleFunc("/api/metrics/playback", controllers.ReportPlaybackMetrics)

	// Register for notifications
	s.router.HandleFunc("/api/notifications/register", middleware.RequireUserAccessToken(controllers.RegisterForLiveNotifications))

	// Start auth flow
	http.HandleFunc("/api/auth/indieauth", middleware.RequireUserAccessToken(indieauth.StartAuthFlow))
	http.HandleFunc("/api/auth/indieauth/callback", indieauth.HandleRedirect)
	http.HandleFunc("/api/auth/provider/indieauth", indieauth.HandleAuthEndpoint)

	http.HandleFunc("/api/auth/fediverse", middleware.RequireUserAccessToken(fediverseauth.RegisterFediverseOTPRequest))
	http.HandleFunc("/api/auth/fediverse/verify", fediverseauth.VerifyFediverseOTPRequest)
}

func (s *webServer) setupAdminAPIRoutes() {
	// Current inbound broadcaster
	s.router.HandleFunc("/api/admin/status", middleware.RequireAdminAuth(admin.Status))

	// Disconnect inbound stream
	s.router.HandleFunc("/api/admin/disconnect", middleware.RequireAdminAuth(admin.DisconnectInboundConnection))

	// Server config
	s.router.HandleFunc("/api/admin/serverconfig", middleware.RequireAdminAuth(admin.GetServerConfig))

	// Get viewer count over time
	s.router.HandleFunc("/api/admin/viewersOverTime", middleware.RequireAdminAuth(admin.GetViewersOverTime))

	// Get active viewers
	s.router.HandleFunc("/api/admin/viewers", middleware.RequireAdminAuth(admin.GetActiveViewers))

	// Get hardware stats
	s.router.HandleFunc("/api/admin/hardwarestats", middleware.RequireAdminAuth(admin.GetHardwareStats))

	// Get a a detailed list of currently connected chat clients
	s.router.HandleFunc("/api/admin/chat/clients", middleware.RequireAdminAuth(admin.GetConnectedChatClients))

	// Get all logs
	s.router.HandleFunc("/api/admin/logs", middleware.RequireAdminAuth(admin.GetLogs))

	// Get warning/error logs
	s.router.HandleFunc("/api/admin/logs/warnings", middleware.RequireAdminAuth(admin.GetWarnings))

	// Get all chat messages for the admin, unfiltered.
	s.router.HandleFunc("/api/admin/chat/messages", middleware.RequireAdminAuth(admin.GetChatMessages))

	// Update chat message visibility
	s.router.HandleFunc("/api/admin/chat/messagevisibility", middleware.RequireAdminAuth(admin.UpdateMessageVisibility))

	// Enable/disable a user
	s.router.HandleFunc("/api/admin/chat/users/setenabled", middleware.RequireAdminAuth(admin.UpdateUserEnabled))

	// Ban/unban an IP address
	s.router.HandleFunc("/api/admin/chat/users/ipbans/create", middleware.RequireAdminAuth(admin.BanIPAddress))

	// Remove an IP address ban
	s.router.HandleFunc("/api/admin/chat/users/ipbans/remove", middleware.RequireAdminAuth(admin.UnBanIPAddress))

	// Return all the banned IP addresses
	s.router.HandleFunc("/api/admin/chat/users/ipbans", middleware.RequireAdminAuth(admin.GetIPAddressBans))

	// Get a list of disabled users
	s.router.HandleFunc("/api/admin/chat/users/disabled", middleware.RequireAdminAuth(admin.GetDisabledUsers))

	// Set moderator status for a user
	s.router.HandleFunc("/api/admin/chat/users/setmoderator", middleware.RequireAdminAuth(admin.UpdateUserModerator))

	// Get a list of moderator users
	s.router.HandleFunc("/api/admin/chat/users/moderators", middleware.RequireAdminAuth(admin.GetModerators))

	// return followers
	s.router.HandleFunc("/api/admin/followers", middleware.RequireAdminAuth(middleware.HandlePagination(controllers.GetFollowers)))

	// Get a list of pending follow requests
	s.router.HandleFunc("/api/admin/followers/pending", middleware.RequireAdminAuth(admin.GetPendingFollowRequests))

	// Get a list of rejected or blocked follows
	s.router.HandleFunc("/api/admin/followers/blocked", middleware.RequireAdminAuth(admin.GetBlockedAndRejectedFollowers))

	// Set the following state of a follower or follow request.
	s.router.HandleFunc("/api/admin/followers/approve", middleware.RequireAdminAuth(admin.ApproveFollower))

	// Upload custom emoji
	s.router.HandleFunc("/api/admin/emoji/upload", middleware.RequireAdminAuth(admin.UploadCustomEmoji))

	// Delete custom emoji
	s.router.HandleFunc("/api/admin/emoji/delete", middleware.RequireAdminAuth(admin.DeleteCustomEmoji))

	// Update config values

	// Change the current streaming key in memory
	s.router.HandleFunc("/api/admin/config/adminpass", middleware.RequireAdminAuth(admin.SetAdminPassword))

	//  Set an array of valid stream keys
	s.router.HandleFunc("/api/admin/config/streamkeys", middleware.RequireAdminAuth(admin.SetStreamKeys))

	// Change the extra page content in memory
	s.router.HandleFunc("/api/admin/config/pagecontent", middleware.RequireAdminAuth(admin.SetExtraPageContent))

	// Stream title
	s.router.HandleFunc("/api/admin/config/streamtitle", middleware.RequireAdminAuth(admin.SetStreamTitle))

	// Server name
	s.router.HandleFunc("/api/admin/config/name", middleware.RequireAdminAuth(admin.SetServerName))

	// Server summary
	s.router.HandleFunc("/api/admin/config/serversummary", middleware.RequireAdminAuth(admin.SetServerSummary))

	// Offline message
	s.router.HandleFunc("/api/admin/config/offlinemessage", middleware.RequireAdminAuth(admin.SetCustomOfflineMessage))

	// Server welcome message
	s.router.HandleFunc("/api/admin/config/welcomemessage", middleware.RequireAdminAuth(admin.SetServerWelcomeMessage))

	// Disable chat
	s.router.HandleFunc("/api/admin/config/chat/disable", middleware.RequireAdminAuth(admin.SetChatDisabled))

	// Disable chat user join messages
	s.router.HandleFunc("/api/admin/config/chat/joinmessagesenabled", middleware.RequireAdminAuth(admin.SetChatJoinMessagesEnabled))

	// Enable/disable chat established user mode
	s.router.HandleFunc("/api/admin/config/chat/establishedusermode", middleware.RequireAdminAuth(admin.SetEnableEstablishedChatUserMode))

	// Set chat usernames that are not allowed
	s.router.HandleFunc("/api/admin/config/chat/forbiddenusernames", middleware.RequireAdminAuth(admin.SetForbiddenUsernameList))

	// Set the suggested chat usernames that will be assigned automatically
	s.router.HandleFunc("/api/admin/config/chat/suggestedusernames", middleware.RequireAdminAuth(admin.SetSuggestedUsernameList))

	// Set video codec
	s.router.HandleFunc("/api/admin/config/video/codec", middleware.RequireAdminAuth(admin.SetVideoCodec))

	// Set style/color/css values
	s.router.HandleFunc("/api/admin/config/appearance", middleware.RequireAdminAuth(admin.SetCustomColorVariableValues))

	// Return all webhooks
	s.router.HandleFunc("/api/admin/webhooks", middleware.RequireAdminAuth(admin.GetWebhooks))

	// Delete a single webhook
	s.router.HandleFunc("/api/admin/webhooks/delete", middleware.RequireAdminAuth(admin.DeleteWebhook))

	// Create a single webhook
	s.router.HandleFunc("/api/admin/webhooks/create", middleware.RequireAdminAuth(admin.CreateWebhook))

	// Get all access tokens
	s.router.HandleFunc("/api/admin/accesstokens", middleware.RequireAdminAuth(admin.GetExternalAPIUsers))

	// Delete a single access token
	s.router.HandleFunc("/api/admin/accesstokens/delete", middleware.RequireAdminAuth(admin.DeleteExternalAPIUser))

	// Create a single access token
	s.router.HandleFunc("/api/admin/accesstokens/create", middleware.RequireAdminAuth(admin.CreateExternalAPIUser))

	// Return the auto-update features that are supported for this instance.
	s.router.HandleFunc("/api/admin/update/options", middleware.RequireAdminAuth(admin.AutoUpdateOptions))

	// Begin the auto update
	s.router.HandleFunc("/api/admin/update/start", middleware.RequireAdminAuth(admin.AutoUpdateStart))

	// Force quit the service to restart it
	s.router.HandleFunc("/api/admin/update/forcequit", middleware.RequireAdminAuth(admin.AutoUpdateForceQuit))

	// Logo path
	s.router.HandleFunc("/api/admin/config/logo", middleware.RequireAdminAuth(admin.SetLogo))

	// Server tags
	s.router.HandleFunc("/api/admin/config/tags", middleware.RequireAdminAuth(admin.SetTags))

	// ffmpeg
	s.router.HandleFunc("/api/admin/config/ffmpegpath", middleware.RequireAdminAuth(admin.SetFfmpegPath))

	// Server http port
	s.router.HandleFunc("/api/admin/config/webserverport", middleware.RequireAdminAuth(admin.SetWebServerPort))

	// Server http listen address
	s.router.HandleFunc("/api/admin/config/webserverip", middleware.RequireAdminAuth(admin.SetWebServerIP))

	// Server rtmp port
	s.router.HandleFunc("/api/admin/config/rtmpserverport", middleware.RequireAdminAuth(admin.SetRTMPServerPort))

	// Websocket host override
	s.router.HandleFunc("/api/admin/config/sockethostoverride", middleware.RequireAdminAuth(admin.SetSocketHostOverride))

	// Custom video serving endpoint
	s.router.HandleFunc("/api/admin/config/videoservingendpoint", middleware.RequireAdminAuth(admin.SetVideoServingEndpoint))

	// Is server marked as NSFW
	s.router.HandleFunc("/api/admin/config/nsfw", middleware.RequireAdminAuth(admin.SetNSFW))

	// directory enabled
	s.router.HandleFunc("/api/admin/config/directoryenabled", middleware.RequireAdminAuth(admin.SetDirectoryEnabled))

	// social handles
	s.router.HandleFunc("/api/admin/config/socialhandles", middleware.RequireAdminAuth(admin.SetSocialHandles))

	// set the number of video segments and duration per segment in a playlist
	s.router.HandleFunc("/api/admin/config/video/streamlatencylevel", middleware.RequireAdminAuth(admin.SetStreamLatencyLevel))

	// set an array of video output configurations
	s.router.HandleFunc("/api/admin/config/video/streamoutputvariants", middleware.RequireAdminAuth(admin.SetStreamOutputVariants))

	// set s3 configuration
	s.router.HandleFunc("/api/admin/config/s3", middleware.RequireAdminAuth(admin.SetS3Configuration))

	// set server url
	s.router.HandleFunc("/api/admin/config/serverurl", middleware.RequireAdminAuth(admin.SetServerURL))

	// reset the YP registration
	s.router.HandleFunc("/api/admin/yp/reset", middleware.RequireAdminAuth(admin.ResetYPRegistration))

	// set external action links
	s.router.HandleFunc("/api/admin/config/externalactions", middleware.RequireAdminAuth(admin.SetExternalActions))

	// set custom style css
	s.router.HandleFunc("/api/admin/config/customstyles", middleware.RequireAdminAuth(admin.SetCustomStyles))

	// set custom style javascript
	s.router.HandleFunc("/api/admin/config/customjavascript", middleware.RequireAdminAuth(admin.SetCustomJavascript))

	// Video playback metrics
	s.router.HandleFunc("/api/admin/metrics/video", middleware.RequireAdminAuth(admin.GetVideoPlaybackMetrics))

	// Is the viewer count hidden from viewers
	s.router.HandleFunc("/api/admin/config/hideviewercount", middleware.RequireAdminAuth(admin.SetHideViewerCount))

	// set disabling of search indexing
	s.router.HandleFunc("/api/admin/config/disablesearchindexing", middleware.RequireAdminAuth(admin.SetDisableSearchIndexing))

	// enable/disable federation features
	s.router.HandleFunc("/api/admin/config/federation/enable", middleware.RequireAdminAuth(admin.SetFederationEnabled))

	// set if federation activities are private
	s.router.HandleFunc("/api/admin/config/federation/private", middleware.RequireAdminAuth(admin.SetFederationActivityPrivate))

	// set if fediverse engagement appears in chat
	s.router.HandleFunc("/api/admin/config/federation/showengagement", middleware.RequireAdminAuth(admin.SetFederationShowEngagement))

	// set local federated username
	s.router.HandleFunc("/api/admin/config/federation/username", middleware.RequireAdminAuth(admin.SetFederationUsername))

	// set federated go live message
	s.router.HandleFunc("/api/admin/config/federation/livemessage", middleware.RequireAdminAuth(admin.SetFederationGoLiveMessage))

	// Federation blocked domains
	s.router.HandleFunc("/api/admin/config/federation/blockdomains", middleware.RequireAdminAuth(admin.SetFederationBlockDomains))

	// send a public message to the Fediverse from the server's user
	s.router.HandleFunc("/api/admin/federation/send", middleware.RequireAdminAuth(admin.SendFederatedMessage))

	// Return federated activities
	s.router.HandleFunc("/api/admin/federation/actions", middleware.RequireAdminAuth(middleware.HandlePagination(admin.GetFederatedActions)))

	// Prometheus metrics
	s.router.Handle("/api/admin/prometheus", middleware.RequireAdminAuth(func(rw http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(rw, r)
	}))

	// Configure outbound notification channels.
	http.HandleFunc("/api/admin/config/notifications/discord", middleware.RequireAdminAuth(admin.SetDiscordNotificationConfiguration))
	http.HandleFunc("/api/admin/config/notifications/browser", middleware.RequireAdminAuth(admin.SetBrowserNotificationConfiguration))
}

func (s *webServer) setupExternalThirdPartyAPIRoutes() {
	// Send a system message to chat
	s.router.HandleFunc("/api/integrations/chat/system", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessage))

	// Send a system message to a single client
	s.router.HandleFunc(utils.RestEndpoint("/api/integrations/chat/system/client/{clientId}", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessageToConnectedClient)))

	// Send a user message to chat *NO LONGER SUPPORTED
	s.router.HandleFunc("/api/integrations/chat/user", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, admin.SendUserMessage))

	// Send a message to chat as a specific 3rd party bot/integration based on its access token
	s.router.HandleFunc("/api/integrations/chat/send", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, admin.SendIntegrationChatMessage))

	// Send a user action to chat
	s.router.HandleFunc("/api/integrations/chat/action", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendChatAction))

	// Hide chat message
	s.router.HandleFunc("/api/integrations/chat/messagevisibility", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalUpdateMessageVisibility))

	// Stream title
	s.router.HandleFunc("/api/integrations/streamtitle", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalSetStreamTitle))

	// Get chat history
	s.router.HandleFunc("/api/integrations/chat", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, controllers.ExternalGetChatMessages))

	// Connected clients
	s.router.HandleFunc("/api/integrations/clients", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalGetConnectedChatClients))
}

func (s *webServer) setupModerationAPIRoutes() {
	// Update chat message visibility
	s.router.HandleFunc("/api/chat/messagevisibility", middleware.RequireUserModerationScopeAccesstoken(admin.UpdateMessageVisibility))

	// Enable/disable a user
	s.router.HandleFunc("/api/chat/users/setenabled", middleware.RequireUserModerationScopeAccesstoken(admin.UpdateUserEnabled))

	// Get a user's details
	s.router.HandleFunc("/api/moderation/chat/user/", middleware.RequireUserModerationScopeAccesstoken(moderation.GetUserDetails))
}
