package router

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/controllers"
	"github.com/gabek/owncast/core/chat"
	"github.com/gabek/owncast/core/rtmp"
)

//Start starts the router for the http, ws, and rtmp
func Start() error {
	// start the rtmp server
	go rtmp.Start()

	// static files
	http.HandleFunc("/", controllers.IndexHandler)

	// status of the system
	http.HandleFunc("/status", controllers.GetStatus)

	// custom emoji supported in the chat
	http.HandleFunc("/emoji", controllers.GetCustomEmoji)

	if !config.Config.DisableWebFeatures {
		// websocket chat server
		go chat.Start()

		// chat rest api
		http.HandleFunc("/chat", controllers.GetChatMessages)

		// web config api
		http.HandleFunc("/config", controllers.GetWebConfig)

		// chat embed
		http.HandleFunc("/embed/chat", controllers.GetChatEmbed)

		// video embed
		http.HandleFunc("/embed/video", controllers.GetVideoEmbed)
	}

	port := config.Config.GetPublicWebServerPort()

	log.Infof("Web server running on port: %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
