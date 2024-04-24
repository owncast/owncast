package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CAFxX/httpcompression"
	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/owncast/owncast/activitypub"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/handler"
	"github.com/owncast/owncast/router/middleware"
)

// Start starts the router for the http, ws, and rtmp.
func Start(enableVerboseLogging bool) error {
	// @behlers New Router
	r := chi.NewRouter()

	// Middlewares
	if enableVerboseLogging {
		r.Use(chiMW.Logger)
	}
	r.Use(chiMW.Recoverer)

	// Images
	r.HandleFunc("/thumbnail.jpg", controllers.GetThumbnail)
	r.HandleFunc("/preview.gif", controllers.GetPreview)
	r.HandleFunc("/logo", controllers.GetLogo)
	// return a logo that's compatible with external social networks
	r.HandleFunc("/logo/external", controllers.GetCompatibleLogo)

	// Custom Javascript
	r.HandleFunc("/customjavascript", controllers.ServeCustomJavascript)

	// robots.txt
	r.HandleFunc("/robots.txt", controllers.GetRobotsDotTxt)

	// Return a single emoji image.
	r.HandleFunc(config.EmojiDir, controllers.GetCustomEmojiImage)

	// websocket
	r.HandleFunc("/ws", chat.HandleClientConnection)

	// serve files
	fs := http.FileServer(http.Dir(config.PublicFilesPath))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	// Return HLS video
	r.HandleFunc("/hls/*", controllers.HandleHLSRequest)

	// The admin web app.
	r.HandleFunc("/admin/*", middleware.RequireAdminAuth(controllers.IndexHandler))

	// The primary web app.
	r.HandleFunc("/*", controllers.IndexHandler)

	// mount the api
	r.Mount("/api/", handler.New().Handler())

	// ActivityPub has its own router
	activitypub.Start(data.GetDatastore())

	// Create a custom mux handler to intercept the /debug/vars endpoint.
	// This is a hack because Prometheus enables this endpoint by default
	// due to its use of expvar and we do not want this exposed.
	h2s := &http2.Server{}
	http2Handler := h2c.NewHandler(r, h2s)
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/debug/vars" {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if r.URL.Path == "/embed/chat/" || r.URL.Path == "/embed/chat" {
			// Redirect /embed/chat
			http.Redirect(w, r, "/embed/chat/readonly", http.StatusTemporaryRedirect)
		} else {
			http2Handler.ServeHTTP(w, r)
		}
	})

	port := config.WebServerPort
	ip := config.WebServerIP

	compress, _ := httpcompression.DefaultAdapter() // Use the default configuration
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", ip, port),
		ReadHeaderTimeout: 4 * time.Second,
		Handler:           compress(m),
	}

	if ip != "0.0.0.0" {
		log.Infof("Web server is listening at %s:%d.", ip, port)
	} else {
		log.Infof("Web server is listening on port %d.", port)
	}
	log.Infoln("Configure this server by visiting /admin.")

	return server.ListenAndServe()
}
