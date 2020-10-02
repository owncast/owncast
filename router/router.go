package router

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/controllers"
	"github.com/gabek/owncast/controllers/admin"

	"github.com/gabek/owncast/core/chat"
	"github.com/gabek/owncast/core/rtmp"
	"github.com/gabek/owncast/router/middleware"
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

	if !config.Config.DisableWebFeatures {
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
	}

	// Authenticated admin requests

	// Disconnect inbound stream
	http.HandleFunc("/api/admin/disconnect", middleware.RequireAdminAuth(admin.DisconnectInboundConnection))

	// Change the current streaming key in memory
	http.HandleFunc("/api/admin/changekey", middleware.RequireAdminAuth(admin.ChangeStreamKey))

	port := config.Config.GetPublicWebServerPort()

	log.Infof("Web server running on port: %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
