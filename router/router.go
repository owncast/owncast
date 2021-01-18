package router

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/router/middleware"
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

	// websocket chat server
	go func() {
		err := chat.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// chat rest api
	http.HandleFunc("/api/chat", controllers.GetChatMessages)

	// web config api
	http.HandleFunc("/api/config", controllers.GetWebConfig)

	// chat embed
	http.HandleFunc("/embed/chat", controllers.GetChatEmbed)

	// video embed
	http.HandleFunc("/embed/video", controllers.GetVideoEmbed)

	http.HandleFunc("/api/yp", yp.GetYPResponse)

	// Authenticated admin requests

	// Current inbound broadcaster
	http.HandleFunc("/api/admin/status", middleware.RequireAdminAuth(admin.Status))

	// Disconnect inbound stream
	http.HandleFunc("/api/admin/disconnect", middleware.RequireAdminAuth(admin.DisconnectInboundConnection))

	// Server config
	http.HandleFunc("/api/admin/serverconfig", middleware.RequireAdminAuth(admin.GetServerConfig))

	// Get viewer count over time
	http.HandleFunc("/api/admin/viewersOverTime", middleware.RequireAdminAuth(admin.GetViewersOverTime))

	// Get hardware stats
	http.HandleFunc("/api/admin/hardwarestats", middleware.RequireAdminAuth(admin.GetHardwareStats))

	// Get a a detailed list of currently connected clients
	http.HandleFunc("/api/admin/clients", middleware.RequireAdminAuth(controllers.GetConnectedClients))

	// Get all logs
	http.HandleFunc("/api/admin/logs", middleware.RequireAdminAuth(admin.GetLogs))

	// Get warning/error logs
	http.HandleFunc("/api/admin/logs/warnings", middleware.RequireAdminAuth(admin.GetWarnings))

	// Get all chat messages for the admin, unfiltered.
	http.HandleFunc("/api/admin/chat/messages", middleware.RequireAdminAuth(admin.GetChatMessages))

	// Update chat message visibilty
	http.HandleFunc("/api/admin/chat/updatemessagevisibility", middleware.RequireAdminAuth(admin.UpdateMessageVisibility))
	// Update config values

	// Change the current streaming key in memory
	http.HandleFunc("/api/admin/config/key", middleware.RequireAdminAuth(admin.ChangeStreamKey))

	// Change the extra page content in memory
	http.HandleFunc("/api/admin/config/pagecontent", middleware.RequireAdminAuth(admin.ChangeExtraPageContent))

	// Stream title
	http.HandleFunc("/api/admin/config/streamtitle", middleware.RequireAdminAuth(admin.ChangeStreamTitle))

	// Server title
	http.HandleFunc("/api/admin/config/servertitle", middleware.RequireAdminAuth(admin.ChangeServerTitle))

	// Server name
	http.HandleFunc("/api/admin/config/name", middleware.RequireAdminAuth(admin.ChangeServerName))

	// Server summary
	http.HandleFunc("/api/admin/config/serversummary", middleware.RequireAdminAuth(admin.ChangeServerSummary))

	// Return all webhooks
	http.HandleFunc("/api/admin/webhooks", middleware.RequireAdminAuth(admin.GetWebhooks))

	// Delete a single webhook
	http.HandleFunc("/api/admin/webhooks/delete", middleware.RequireAdminAuth(admin.DeleteWebhook))

	// Create a single webhook
	http.HandleFunc("/api/admin/webhooks/create", middleware.RequireAdminAuth(admin.CreateWebhook))

	// Get all access tokens
	http.HandleFunc("/api/admin/accesstokens", middleware.RequireAdminAuth(admin.GetAccessTokens))

	// Delete a single access token
	http.HandleFunc("/api/admin/accesstokens/delete", middleware.RequireAdminAuth(admin.DeleteAccessToken))

	// Create a single access token
	http.HandleFunc("/api/admin/accesstokens/create", middleware.RequireAdminAuth(admin.CreateAccessToken))

	// Send a system message to chat
	http.HandleFunc("/api/integrations/chat/system", middleware.RequireAccessToken(models.ScopeCanSendSystemMessages, admin.SendSystemMessage))

	// Send a user message to chat
	http.HandleFunc("/api/integrations/chat/user", middleware.RequireAccessToken(models.ScopeCanSendUserMessages, admin.SendUserMessage))

	// Send a user action to chat
	http.HandleFunc("/api/integrations/chat/action", middleware.RequireAccessToken(models.ScopeCanSendSystemMessages, admin.SendChatAction))

	// Hide chat message
	http.HandleFunc("/api/integrations/chat/messagevisibility", middleware.RequireAccessToken(models.ScopeHasAdminAccess, admin.UpdateMessageVisibility))

	// Get chat history
	http.HandleFunc("/api/integrations/chat", middleware.RequireAccessToken(models.ScopeHasAdminAccess, controllers.GetChatMessages))

	// Connected clients
	http.HandleFunc("/api/integrations/clients", middleware.RequireAccessToken(models.ScopeHasAdminAccess, controllers.GetConnectedClients))
	// Logo path
	http.HandleFunc("/api/admin/config/logo", middleware.RequireAdminAuth(admin.ChangeLogoPath))

	// Server tags
	http.HandleFunc("/api/admin/config/tags", middleware.RequireAdminAuth(admin.ChangeTags))

	// ffmpeg
	http.HandleFunc("/api/admin/config/ffmpegpath", middleware.RequireAdminAuth(admin.ChangeFfmpegPath))

	// Server http port
	http.HandleFunc("/api/admin/config/webserverport", middleware.RequireAdminAuth(admin.ChangeWebServerPort))

	// Server rtmp port
	http.HandleFunc("/api/admin/config/rtmpserverport", middleware.RequireAdminAuth(admin.ChangeRTMPServerPort))

	// Is server marked as NSFW
	http.HandleFunc("/api/admin/config/nsfw", middleware.RequireAdminAuth(admin.ChangeNSFW))

	// directory enabled
	http.HandleFunc("/api/admin/config/directoryenabled", middleware.RequireAdminAuth(admin.ChangeDirectoryEnabled))

	// social handles
	http.HandleFunc("/api/admin/config/socialhandles", middleware.RequireAdminAuth(admin.SetSocialHandles))

	// disable server upgrade checks
	http.HandleFunc("/api/admin/config/disableupgradechecks", middleware.RequireAdminAuth(admin.ChangeDisableUpgradeChecks))

	// set the number of video segments and duration per segment in a playlist
	http.HandleFunc("/api/admin/config/video/streamlatancylevel", middleware.RequireAdminAuth(admin.SetStreamLatancyLevel))

	// set an array of video output configurations
	http.HandleFunc("/api/admin/config/video/streamoutputvariants", middleware.RequireAdminAuth(admin.SetStreamOutputVariants))

	// set s3 configuration
	http.HandleFunc("/api/admin/config/s3", middleware.RequireAdminAuth(admin.SetS3Configuration))

	port := data.GetHTTPPortNumber()

	log.Tracef("Web server running on port: %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
