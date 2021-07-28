package router

import (
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/yp"
)

// Start starts the router for the http, ws, and rtmp.
func Start() error {
	// static files
	http.HandleFunc("/", controllers.IndexHandler)

	// admin static files
	http.HandleFunc("/admin/", middleware.RequireAdminAuth(admin.ServeAdmin))

	// status of the system
	http.HandleFunc("/api/status", controllers.GetStatus)

	// custom emoji supported in the chat
	http.HandleFunc("/api/emoji", controllers.GetCustomEmoji)

	// chat rest api
	http.HandleFunc("/api/chat", middleware.RequireUserAccessToken(controllers.GetChatMessages))

	// web config api
	http.HandleFunc("/api/config", controllers.GetWebConfig)

	// pre v0.0.8 chat embed
	http.HandleFunc("/embed/chat", controllers.GetChatEmbedreadonly)

	// readonly chat embed
	http.HandleFunc("/embed/chat/readonly", controllers.GetChatEmbedreadonly)

	// readwrite chat embed
	http.HandleFunc("/embed/chat/readwrite", controllers.GetChatEmbedreadwrite)

	// video embed
	http.HandleFunc("/embed/video", controllers.GetVideoEmbed)

	// return the YP protocol data
	http.HandleFunc("/api/yp", yp.GetYPResponse)

	// list of all social platforms
	http.HandleFunc("/api/socialplatforms", controllers.GetAllSocialPlatforms)

	// return the logo
	http.HandleFunc("/logo", controllers.GetLogo)

	// return a logo that's compatible with external social networks
	http.HandleFunc("/logo/external", controllers.GetCompatibleLogo)

	// return the list of video variants available
	http.HandleFunc("/api/video/variants", controllers.GetVideoStreamOutputVariants)

	// tell the backend you're an active viewer
	http.HandleFunc("/api/ping", controllers.Ping)

	// register a new chat user
	http.HandleFunc("/api/chat/register", controllers.RegisterAnonymousChatUser)

	// Authenticated admin requests

	// Current inbound broadcaster
	http.HandleFunc("/api/admin/status", middleware.RequireAdminAuth(admin.Status))

	// Return HLS video
	http.HandleFunc("/hls/", controllers.HandleHLSRequest)

	// Disconnect inbound stream
	http.HandleFunc("/api/admin/disconnect", middleware.RequireAdminAuth(admin.DisconnectInboundConnection))

	// Server config
	http.HandleFunc("/api/admin/serverconfig", middleware.RequireAdminAuth(admin.GetServerConfig))

	// Get viewer count over time
	http.HandleFunc("/api/admin/viewersOverTime", middleware.RequireAdminAuth(admin.GetViewersOverTime))

	// Get hardware stats
	http.HandleFunc("/api/admin/hardwarestats", middleware.RequireAdminAuth(admin.GetHardwareStats))

	// Get a a detailed list of currently connected chat clients
	http.HandleFunc("/api/admin/chat/clients", middleware.RequireAdminAuth(admin.GetConnectedChatClients))

	// Get all logs
	http.HandleFunc("/api/admin/logs", middleware.RequireAdminAuth(admin.GetLogs))

	// Get warning/error logs
	http.HandleFunc("/api/admin/logs/warnings", middleware.RequireAdminAuth(admin.GetWarnings))

	// Get all chat messages for the admin, unfiltered.
	http.HandleFunc("/api/admin/chat/messages", middleware.RequireAdminAuth(admin.GetChatMessages))

	// Update chat message visibility
	http.HandleFunc("/api/admin/chat/updatemessagevisibility", middleware.RequireAdminAuth(admin.UpdateMessageVisibility))

	// Enable/disable a user
	http.HandleFunc("/api/admin/chat/users/setenabled", middleware.RequireAdminAuth(admin.UpdateUserEnabled))

	// Get a list of disabled users
	http.HandleFunc("/api/admin/chat/users/disabled", middleware.RequireAdminAuth(admin.GetDisabledUsers))

	// Update config values

	// Change the current streaming key in memory
	http.HandleFunc("/api/admin/config/key", middleware.RequireAdminAuth(admin.SetStreamKey))

	// Change the extra page content in memory
	http.HandleFunc("/api/admin/config/pagecontent", middleware.RequireAdminAuth(admin.SetExtraPageContent))

	// Stream title
	http.HandleFunc("/api/admin/config/streamtitle", middleware.RequireAdminAuth(admin.SetStreamTitle))

	// Server name
	http.HandleFunc("/api/admin/config/name", middleware.RequireAdminAuth(admin.SetServerName))

	// Server summary
	http.HandleFunc("/api/admin/config/serversummary", middleware.RequireAdminAuth(admin.SetServerSummary))

	// Server welcome message
	http.HandleFunc("/api/admin/config/welcomemessage", middleware.RequireAdminAuth(admin.SetServerWelcomeMessage))

	// Disable chat
	http.HandleFunc("/api/admin/config/chat/disable", middleware.RequireAdminAuth(admin.SetChatDisabled))

	// Set chat usernames that are not allowed
	http.HandleFunc("/api/admin/config/chat/forbiddenusernames", middleware.RequireAdminAuth(admin.SetForbiddenUsernameList))

	// Set video codec
	http.HandleFunc("/api/admin/config/video/codec", middleware.RequireAdminAuth(admin.SetVideoCodec))

	// Return all webhooks
	http.HandleFunc("/api/admin/webhooks", middleware.RequireAdminAuth(admin.GetWebhooks))

	// Delete a single webhook
	http.HandleFunc("/api/admin/webhooks/delete", middleware.RequireAdminAuth(admin.DeleteWebhook))

	// Create a single webhook
	http.HandleFunc("/api/admin/webhooks/create", middleware.RequireAdminAuth(admin.CreateWebhook))

	// Get all access tokens
	http.HandleFunc("/api/admin/accesstokens", middleware.RequireAdminAuth(admin.GetExternalAPIUsers))

	// Delete a single access token
	http.HandleFunc("/api/admin/accesstokens/delete", middleware.RequireAdminAuth(admin.DeleteExternalAPIUser))

	// Create a single access token
	http.HandleFunc("/api/admin/accesstokens/create", middleware.RequireAdminAuth(admin.CreateExternalAPIUser))

	// Send a system message to chat
	http.HandleFunc("/api/integrations/chat/system", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessage))

	// Send a system message to a single client
	http.HandleFunc(utils.RestEndpoint("/api/integrations/chat/system/client/{clientId}", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessageToConnectedClient)))

	// Send a user message to chat *NO LONGER SUPPORTED
	http.HandleFunc("/api/integrations/chat/user", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, admin.SendUserMessage))

	// Send a message to chat as a specific 3rd party bot/integration based on its access token
	http.HandleFunc("/api/integrations/chat/send", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendChatMessages, admin.SendIntegrationChatMessage))

	// Send a user action to chat
	http.HandleFunc("/api/integrations/chat/action", middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendChatAction))

	// Hide chat message
	http.HandleFunc("/api/integrations/chat/messagevisibility", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalUpdateMessageVisibility))

	// Stream title
	http.HandleFunc("/api/integrations/streamtitle", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalSetStreamTitle))

	// Get chat history
	http.HandleFunc("/api/integrations/chat", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, controllers.ExternalGetChatMessages))

	// Connected clients
	http.HandleFunc("/api/integrations/clients", middleware.RequireExternalAPIAccessToken(user.ScopeHasAdminAccess, admin.ExternalGetConnectedChatClients))

	// Logo path
	http.HandleFunc("/api/admin/config/logo", middleware.RequireAdminAuth(admin.SetLogo))

	// Server tags
	http.HandleFunc("/api/admin/config/tags", middleware.RequireAdminAuth(admin.SetTags))

	// ffmpeg
	http.HandleFunc("/api/admin/config/ffmpegpath", middleware.RequireAdminAuth(admin.SetFfmpegPath))

	// Server http port
	http.HandleFunc("/api/admin/config/webserverport", middleware.RequireAdminAuth(admin.SetWebServerPort))

	// Server http listen address
	http.HandleFunc("/api/admin/config/webserverip", middleware.RequireAdminAuth(admin.SetWebServerIP))

	// Server rtmp port
	http.HandleFunc("/api/admin/config/rtmpserverport", middleware.RequireAdminAuth(admin.SetRTMPServerPort))

	// Is server marked as NSFW
	http.HandleFunc("/api/admin/config/nsfw", middleware.RequireAdminAuth(admin.SetNSFW))

	// directory enabled
	http.HandleFunc("/api/admin/config/directoryenabled", middleware.RequireAdminAuth(admin.SetDirectoryEnabled))

	// social handles
	http.HandleFunc("/api/admin/config/socialhandles", middleware.RequireAdminAuth(admin.SetSocialHandles))

	// set the number of video segments and duration per segment in a playlist
	http.HandleFunc("/api/admin/config/video/streamlatencylevel", middleware.RequireAdminAuth(admin.SetStreamLatencyLevel))

	// set an array of video output configurations
	http.HandleFunc("/api/admin/config/video/streamoutputvariants", middleware.RequireAdminAuth(admin.SetStreamOutputVariants))

	// set s3 configuration
	http.HandleFunc("/api/admin/config/s3", middleware.RequireAdminAuth(admin.SetS3Configuration))

	// set server url
	http.HandleFunc("/api/admin/config/serverurl", middleware.RequireAdminAuth(admin.SetServerURL))

	// reset the YP registration
	http.HandleFunc("/api/admin/yp/reset", middleware.RequireAdminAuth(admin.ResetYPRegistration))

	// set external action links
	http.HandleFunc("/api/admin/config/externalactions", middleware.RequireAdminAuth(admin.SetExternalActions))

	// set custom style css
	http.HandleFunc("/api/admin/config/customstyles", middleware.RequireAdminAuth(admin.SetCustomStyles))

	// ActivityPub has its own router
	activitypub.Start()

	// websocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.HandleClientConnection(w, r)
	})

	port := config.WebServerPort
	ip := config.WebServerIP

	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		log.Fatalln("Invalid IP address", ip)
	}
	log.Infof("Web server is listening on IP %s port %d.", ipAddr.String(), port)
	log.Infoln("The web admin interface is available at /admin.")

	return http.ListenAndServe(fmt.Sprintf("%s:%d", ipAddr.String(), port), nil)
}
