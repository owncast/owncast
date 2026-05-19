package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CAFxX/httpcompression"
	"github.com/go-chi/chi/v5"
	chiMW "github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	apcontrollers "github.com/owncast/owncast/services/activitypub/controllers"
	"github.com/owncast/owncast/webserver/handlers"
	"github.com/owncast/owncast/webserver/router/middleware"
)

// Start starts the router for the http, ws, and rtmp.
//
// h carries dependency-injected handler methods. Free-function handlers
// (the majority during the migration) are referenced directly by package
// path; methods on *handlers.Handlers are registered via the h receiver.
// apc carries the methodified ActivityPub HTTP handler set.
func Start(enableVerboseLogging bool, h *handlers.Handlers, apc *apcontrollers.Controllers) error {
	// @behlers New Router
	r := chi.NewRouter()

	// Middlewares
	if enableVerboseLogging {
		r.Use(chiMW.RequestLogger(&chiMW.DefaultLogFormatter{Logger: log.StandardLogger(), NoColor: true}))
	}
	r.Use(chiMW.Recoverer)

	addStaticFileEndpoints(r, apc)

	// websocket
	r.HandleFunc("/ws", chat.HandleClientConnection)

	// serve files
	fs := http.FileServer(http.Dir(config.PublicFilesPath))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	// Return HLS video
	r.HandleFunc("/hls/*", h.HandleHLSRequest)

	// The admin web app.
	r.HandleFunc("/admin/*", middleware.RequireAdminAuth(h.IndexHandler))

	// Single ActivityPub Actor
	r.HandleFunc("/federation/user/*", middleware.RequireActivityPubOrRedirect(apc.ActorHandler))

	// Single AP object
	r.HandleFunc("/federation/*", middleware.RequireActivityPubOrRedirect(apc.ObjectHandler))

	// The primary web app.
	r.HandleFunc("/*", h.IndexHandler)

	// mount the api
	r.Mount("/api/", handlers.New(h).Handler())

	// Create a custom mux handler to intercept the /debug/vars endpoint.
	// This is a hack because Prometheus enables this endpoint by default
	// due to its use of expvar and we do not want this exposed.
	h2s := &http2.Server{}
	http2Handler := h2c.NewHandler(r, h2s)
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/debug/vars":
			w.WriteHeader(http.StatusNotFound)
			return
		case "/embed/chat/", "/embed/chat":
			// Redirect /embed/chat
			http.Redirect(w, r, "/embed/chat/readonly", http.StatusTemporaryRedirect)
		default:
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

func addStaticFileEndpoints(r chi.Router, apc *apcontrollers.Controllers) {
	// Images
	r.HandleFunc("/thumbnail.jpg", handlers.GetThumbnail)
	r.HandleFunc("/preview.gif", handlers.GetPreview)
	r.HandleFunc("/logo", handlers.GetLogo)
	r.HandleFunc("/favicon.ico", handlers.GetFavicon)
	// return a logo that's compatible with external social networks
	r.HandleFunc("/logo/external", handlers.GetCompatibleLogo)

	// Custom Javascript
	r.HandleFunc("/customjavascript", handlers.ServeCustomJavascript)

	// robots.txt
	r.HandleFunc("/robots.txt", handlers.GetRobotsDotTxt)

	// Return a single emoji image.
	emojiDir := config.EmojiDir
	if !strings.HasSuffix(emojiDir, "*") {
		emojiDir += "*"
	}
	r.HandleFunc(emojiDir, handlers.GetCustomEmojiImage)

	// WebFinger
	r.HandleFunc("/.well-known/webfinger", apc.WebfingerHandler)

	// Host Metadata
	r.HandleFunc("/.well-known/host-meta", apc.HostMetaController)

	// Nodeinfo v1
	r.HandleFunc("/.well-known/nodeinfo", apc.NodeInfoController)

	// x-nodeinfo v2
	r.HandleFunc("/.well-known/x-nodeinfo2", apc.XNodeInfo2Controller)

	// Nodeinfo v2
	r.HandleFunc("/nodeinfo/2.0", apc.NodeInfoV2Controller)

	// Instance details
	r.HandleFunc("/api/v1/instance", apc.InstanceV1Controller)
}
