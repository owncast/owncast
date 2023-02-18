package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nanmu42/gzip"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/app/middleware"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers/admin"
	"github.com/owncast/owncast/controllers/auth"
	"github.com/owncast/owncast/controllers/moderation"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/utils"
)

// Start starts the router for the http, ws, and rtmp.
func (a *App) initRoutes() error {
	port := config.WebServerPort
	ip := config.WebServerIP

	// Create a custom mux handler to intercept the /debug/vars endpoint.
	// This is a hack because Prometheus enables this endpoint by default
	// due to its use of expvar and we do not want this exposed.
	//h2s := &http2.Server{}
	//defaultMux := h2c.NewHandler(a.router, h2s)
	a.router = http.NewServeMux()

	// The primary web app.
	a.router.HandleFunc("/", a.Controller.IndexHandler)

	if err := a.initFrontendRoutes("/", a.router); err != nil {
		return fmt.Errorf("initializing frontend routes: %v", err)
	}

	if err := a.initBackendRoutes("/api", a.router); err != nil {
		return fmt.Errorf("initializing frontend routes: %v", err)
	}

	// Return HLS video
	a.router.HandleFunc("/hls/", a.Controller.HandleHLSRequest)

	// websocket
	a.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		a.Core.Chat.HandleClientConnection(w, r)
	})

	// Optional public static files
	a.router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(config.PublicFilesPath))))

	//a.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.Path == "/debug/vars" {
	//		w.WriteHeader(http.StatusNotFound)
	//		return
	//	} else {
	//		defaultMux.ServeHTTP(w, r)
	//	}
	//})

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", ip, port),
		ReadHeaderTimeout: 4 * time.Second,
		Handler:           gzip.DefaultHandler().WrapHandler(a.router),
	}

	log.Infof("Web server is listening on IP %s port %d.", ip, port)
	log.Infoln("The web admin interface is available at /admin.")

	return server.ListenAndServe()
}

func (a *App) initFrontendRoutes(prefix string, mux *http.ServeMux) error {
	routes := map[string]http.HandlerFunc{
		"/admin":            middleware.RequireAdminAuth(a.Controller.IndexHandler, a.Data),
		"/thumbnail.jpg":    a.Controller.GetThumbnail,
		"/preview.gif":      a.Controller.GetPreview,
		"/customjavascript": a.Controller.ServeCustomJavascript,
		"/logo":             a.Controller.GetLogo,
		"/logo/external":    a.Controller.GetCompatibleLogo,
		config.EmojiDir:     a.Controller.GetCustomEmojiImage,
	}

	for suffix, handler := range routes {
		mux.HandleFunc(fmt.Sprintf("%s%s", prefix, suffix), handler)
	}

	return nil
}

func (a *App) initBackendRoutes(prefix string, mux *http.ServeMux) error {
	routes := map[string]http.HandlerFunc{
		"/status":                 a.Controller.GetStatus,
		"/emoji":                  a.Controller.GetCustomEmojiList,
		"/chat":                   middleware.RequireUserAccessToken(a.Controller.GetChatMessages, a.Data),
		"/chat/register":          a.Controller.RegisterAnonymousChatUser,
		"/config":                 a.Controller.GetWebConfig,
		"/YP":                     a.YP.GetYPResponse,
		"/socialplatforms":        a.Controller.GetAllSocialPlatforms,
		"/video/variants":         a.Controller.GetVideoStreamOutputVariants,
		"/ping":                   a.Controller.Ping,
		"/remotefollow":           a.Controller.RemoteFollow,
		"/followers":              middleware.HandlePagination(a.Controller.GetFollowers),
		"/metrics/playback":       a.Controller.ReportPlaybackMetrics,
		"/notifications/register": middleware.RequireUserAccessToken(a.Controller.RegisterForLiveNotifications, a.Data),
	}

	for suffix, handler := range routes {
		mux.HandleFunc(fmt.Sprintf("%s%s", prefix, suffix), handler)
	}

	adminController, err := admin.New(a.Controller)
	if err != nil {
		return fmt.Errorf("initializing admin controller: %v", err)
	}

	prefixAdminApi := fmt.Sprintf("%s/%s", prefix, "admin")
	if err := a.initAdminApiRoutes(prefixAdminApi, mux, adminController); err != nil {
		return fmt.Errorf("initializing backend admin routes: %v", err)
	}

	prefixIntegrations := fmt.Sprintf("%s/%s", prefix, "integrations")
	if err := a.initIntegrationRoutes(prefixIntegrations, mux, adminController); err != nil {
		return fmt.Errorf("initializing backend admin routes: %v", err)
	}

	prefixAuth := fmt.Sprintf("%s/%s", prefix, "auth")
	if err := a.initAuthRoutes(prefixAuth, mux); err != nil {
		return fmt.Errorf("initializing backend admin routes: %v", err)
	}

	// Prometheus metrics
	mux.Handle("/api/admin/prometheus", middleware.RequireAdminAuth(func(rw http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(rw, r)
	}, a.Data))

	return nil
}

func (a *App) initIntegrationRoutes(prefix string, mux *http.ServeMux, admin *admin.Controller) error {
	type externalApiHandler = func(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request)
	routes := make(map[string]map[string]externalApiHandler)

	routes[user.ScopeCanSendSystemMessages] = map[string]externalApiHandler{
		"/chat/system": admin.SendSystemMessage,
		"/chat/action": admin.SendChatAction,
	}

	// don't know why this one is so nasty
	mux.Handle(utils.RestEndpoint(
		fmt.Sprintf("%s%s", prefix, "/chat/system/client/{clientId}"),
		middleware.RequireExternalAPIAccessToken(user.ScopeCanSendSystemMessages, admin.SendSystemMessageToConnectedClient, a.Data)),
	)

	routes[user.ScopeCanSendChatMessages] = map[string]externalApiHandler{
		"/chat/user": admin.SendUserMessage,
		"/chat/send": admin.SendIntegrationChatMessage,
	}

	routes[user.ScopeHasAdminAccess] = map[string]externalApiHandler{
		"/chat/messagevisibility": admin.ExternalUpdateMessageVisibility,
		"/streamtitle":            admin.ExternalSetStreamTitle,
		"/chat":                   a.Controller.ExternalGetChatMessages,
		"/clients":                admin.ExternalGetConnectedChatClients,
	}

	for scope, routeHandler := range routes {
		for suffix, fn := range routeHandler {
			handler := middleware.RequireExternalAPIAccessToken(scope, fn, a.Data)
			mux.HandleFunc(fmt.Sprintf("%s%s", prefix, suffix), handler)
		}
	}

	return nil
}

func (a *App) initAuthRoutes(prefix string, mux *http.ServeMux) error {
	authController, err := auth.New(a.Controller)
	if err != nil {
		return fmt.Errorf("initializing auth controller: %v", err)
	}

	http.HandleFunc("/indieauth", middleware.RequireUserAccessToken(authController.IndieAuth.StartAuthFlow, a.Data))
	http.HandleFunc("/indieauth/callback", authController.IndieAuth.HandleRedirect)
	http.HandleFunc("/provider/indieauth", authController.IndieAuth.HandleAuthEndpoint)

	http.HandleFunc("/fediverse", middleware.RequireUserAccessToken(authController.Fediverse.RegisterFediverseOTPRequest, a.Data))
	http.HandleFunc("/fediverse/verify", authController.Fediverse.VerifyFediverseOTPRequest)

	return nil
}

func (a *App) initAdminApiRoutes(prefix string, mux *http.ServeMux, adminController *admin.Controller) error {
	routes := map[string]http.HandlerFunc{
		"/status":          adminController.Status,
		"/disconnect":      adminController.DisconnectInboundConnection,
		"/serverconfig":    adminController.GetServerConfig,
		"/viewersOverTime": adminController.GetViewersOverTime,
		"/viewers":         adminController.GetActiveViewers,
		"/hardwarestats":   adminController.GetHardwareStats,
		"/chat/clients":    adminController.GetConnectedChatClients,

		"/logs":          adminController.GetLogs,
		"/logs/warnings": adminController.GetWarnings,

		"/chat/messages":            adminController.GetChatMessages,
		"/chat/messagevisibility":   adminController.UpdateMessageVisibility,
		"/chat/users/setenabled":    adminController.UpdateUserEnabled,
		"/chat/users/ipbans/create": adminController.BanIPAddress,
		"/chat/users/ipbans/remove": adminController.UnBanIPAddress,
		"/chat/users/ipbans":        adminController.GetIPAddressBans,
		"/chat/users/disabled":      adminController.GetDisabledUsers,
		"/chat/users/setmoderator":  adminController.UpdateUserModerator,
		"/chat/users/moderators":    adminController.GetModerators,

		"/followers":         middleware.HandlePagination(a.Controller.GetFollowers),
		"/followers/pending": adminController.GetPendingFollowRequests,
		"/followers/blocked": adminController.GetBlockedAndRejectedFollowers,
		"/followers/approve": adminController.ApproveFollower,

		"/emoji/upload": adminController.UploadCustomEmoji,
		"/emoji/delete": adminController.DeleteCustomEmoji,

		"/config/adminpass":                  adminController.SetAdminPassword,
		"/config/streamkeys":                 adminController.SetStreamKeys,
		"/config/pagecontent":                adminController.SetExtraPageContent,
		"/config/streamtitle":                adminController.SetStreamTitle,
		"/config/name":                       adminController.SetServerName,
		"/config/serversummary":              adminController.SetServerSummary,
		"/config/offlinemessage":             adminController.SetCustomOfflineMessage,
		"/config/welcomemessage":             adminController.SetServerWelcomeMessage,
		"/config/chat/disable":               adminController.SetChatDisabled,
		"/config/chat/joinmessagesenabled":   adminController.SetChatJoinMessagesEnabled,
		"/config/chat/establishedusermode":   adminController.SetEnableEstablishedChatUserMode,
		"/config/chat/forbiddenusernames":    adminController.SetForbiddenUsernameList,
		"/config/chat/suggestedusernames":    adminController.SetSuggestedUsernameList,
		"/config/video/codec":                adminController.SetVideoCodec,
		"/config/appearance":                 adminController.SetCustomColorVariableValues,
		"/config/logo":                       adminController.SetLogo,
		"/config/tags":                       adminController.SetTags,
		"/config/ffmpegpath":                 adminController.SetFfmpegPath,
		"/config/webserverport":              adminController.SetWebServerPort,
		"/config/webserverip":                adminController.SetWebServerIP,
		"/config/rtmpserverport":             adminController.SetRTMPServerPort,
		"/config/sockethostoverride":         adminController.SetSocketHostOverride,
		"/config/nsfw":                       adminController.SetNSFW,
		"/config/directoryenabled":           adminController.SetDirectoryEnabled,
		"/config/socialhandles":              adminController.SetSocialHandles,
		"/config/video/streamlatencylevel":   adminController.SetStreamLatencyLevel,
		"/config/video/streamoutputvariants": adminController.SetStreamOutputVariants,
		"/config/s3":                         adminController.SetS3Configuration,
		"/config/serverurl":                  adminController.SetServerURL,
		"/config/externalactions":            adminController.SetExternalActions,
		"/config/customstyles":               adminController.SetCustomStyles,
		"/config/customjavascript":           adminController.SetCustomJavascript,
		"/config/hideviewercount":            adminController.SetHideViewerCount,
		"/config/federation/enable":          adminController.SetFederationEnabled,
		"/config/federation/private":         adminController.SetFederationActivityPrivate,
		"/config/federation/showengagement":  adminController.SetFederationShowEngagement,
		"/config/federation/username":        adminController.SetFederationUsername,
		"/config/federation/livemessage":     adminController.SetFederationGoLiveMessage,
		"/config/federation/blockdomains":    adminController.SetFederationBlockDomains,
		"/config/notifications/discord":      adminController.SetDiscordNotificationConfiguration,
		"/config/notifications/browser":      adminController.SetBrowserNotificationConfiguration,

		"/webhooks":        adminController.GetWebhooks,
		"/webhooks/delete": adminController.DeleteWebhook,
		"/webhooks/create": adminController.CreateWebhook,

		"/accesstokens":        adminController.GetExternalAPIUsers,
		"/accesstokens/delete": adminController.DeleteExternalAPIUser,
		"/accesstokens/create": adminController.CreateExternalAPIUser,

		"/update/options":   adminController.AutoUpdateOptions,
		"/update/start":     adminController.AutoUpdateStart,
		"/update/forcequit": adminController.AutoUpdateForceQuit,

		"/YP/reset": adminController.ResetYPRegistration,

		"/metrics/video": adminController.GetVideoPlaybackMetrics,

		"/federation/send":    adminController.SendFederatedMessage,
		"/federation/actions": middleware.HandlePagination(adminController.GetFederatedActions),
	}

	for suffix, handler := range routes {
		mux.HandleFunc(fmt.Sprintf("%s%s", prefix, suffix), middleware.RequireAdminAuth(handler, a.Data))
	}

	moderationController, err := moderation.New(a.Controller)
	if err != nil {
		return fmt.Errorf("initializing admin controller: %v", err)
	}

	// special case, this does not use the prefix that was passed in
	// (the routes already exist!), and this is coupled to the admin router
	for absoluteRoute, handler := range map[string]http.HandlerFunc{
		"/api/chat/messagevisibility": middleware.RequireUserModerationScopeAccesstoken(adminController.UpdateMessageVisibility, a.Data),
		"/api/chat/users/setenabled":  middleware.RequireUserModerationScopeAccesstoken(adminController.UpdateUserEnabled, a.Data),
		"/api/moderation/chat/user/":  middleware.RequireUserModerationScopeAccesstoken(moderationController.GetUserDetails, a.Data),
	} {
		mux.HandleFunc(absoluteRoute, middleware.RequireAdminAuth(handler, a.Data))
	}

	return nil
}
