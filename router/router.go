package router

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/controllers/admin"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/yp"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/router/middleware"

	"github.com/owncast/owncast/yp"
)

//Start starts the router for the http, ws, and rtmp
func Start() error {
	// start the rtmp server
	go rtmp.Start()

	// static files
	http.HandleFunc("/", controllers.IndexHandler)

	// status of the system
	http.HandleFunc("/api/status", controllers.GetStatus)

	// custom emoji supported in the chat
	http.HandleFunc("/api/emoji", controllers.GetCustomEmoji)

	// websocket chat server
	go chat.Start()

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
	http.HandleFunc("/api/admin/broadcaster", middleware.RequireAdminAuth(admin.GetInboundBroadasterDetails))

	// Disconnect inbound stream
	http.HandleFunc("/api/admin/disconnect", middleware.RequireAdminAuth(admin.DisconnectInboundConnection))

	// Change the current streaming key in memory
	http.HandleFunc("/api/admin/changekey", middleware.RequireAdminAuth(admin.ChangeStreamKey))

	// Server config
	http.HandleFunc("/api/admin/serverconfig", middleware.RequireAdminAuth(admin.GetServerConfig))

	// Get viewer count over time
	http.HandleFunc("/api/admin/viewersOverTime", middleware.RequireAdminAuth(admin.GetViewersOverTime))

	// Get hardware stats
	http.HandleFunc("/api/admin/hardwarestats", middleware.RequireAdminAuth(admin.GetHardwareStats))

	port := config.Config.GetPublicWebServerPort()

	log.Infof("Web server running on port: %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
