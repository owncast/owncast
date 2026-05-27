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

	"github.com/owncast/owncast/config"
	apcontrollers "github.com/owncast/owncast/services/activitypub/controllers"
	"github.com/owncast/owncast/webserver/handlers"
	"github.com/owncast/owncast/webserver/router/middleware"
)

// Start starts the router for the http, ws, and rtmp.
//
// cfg supplies the bound port + IP for the public web server. h carries
// dependency-injected handler methods. Free-function handlers (the
// majority during the migration) are referenced directly by package
// path; methods on *handlers.Handlers are registered via the h receiver.
// mw carries the methodified HTTP middleware (admin basic-auth, federation
// content-type gating). apc carries the methodified ActivityPub HTTP
// handler set.
func Start(cfg *config.Config, enableVerboseLogging bool, h *handlers.Handlers, mw *middleware.Middleware, apc *apcontrollers.Controllers) error {
	// @behlers New Router
	r := chi.NewRouter()

	// Middlewares
	if enableVerboseLogging {
		r.Use(chiMW.RequestLogger(&chiMW.DefaultLogFormatter{Logger: log.StandardLogger(), NoColor: true}))
	}
	r.Use(chiMW.Recoverer)

	addStaticFileEndpoints(r, h, apc)

	// websocket
	r.HandleFunc("/ws", h.HandleWebsocketConnection)

	// serve files
	fs := http.FileServer(http.Dir(config.PublicFilesPath))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	// Return HLS video
	r.HandleFunc("/hls/*", h.HandleHLSRequest)

	// The admin web app.
	r.HandleFunc("/admin/*", mw.RequireAdminAuth(h.IndexHandler))

	// Single ActivityPub Actor
	r.HandleFunc("/federation/user/*", mw.RequireActivityPubOrRedirect(apc.ActorHandler))

	// Single AP object
	r.HandleFunc("/federation/*", mw.RequireActivityPubOrRedirect(apc.ObjectHandler))

	// The primary web app.
	r.HandleFunc("/*", h.IndexHandler)

	// mount the api
	r.Mount("/api/", handlers.New(h).Handler())

	// Create a custom mux handler to intercept the /debug/vars endpoint.
	// This is a hack because Prometheus enables this endpoint by default
	// due to its use of expvar and we do not want this exposed.
	rootHandler := r
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
			rootHandler.ServeHTTP(w, r)
		}
	})

	port := cfg.WebServerPort
	ip := cfg.WebServerIP

	// Allow cleartext (unencrypted) HTTP/2 in addition to HTTP/1, replacing the
	// previously used and now-deprecated golang.org/x/net/http2/h2c wrapper.
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	compress, _ := httpcompression.DefaultAdapter() // Use the default configuration
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", ip, port),
		ReadHeaderTimeout: 4 * time.Second,
		Handler:           compress(m),
		Protocols:         protocols,
	}

	if ip != "0.0.0.0" {
		log.Infof("Web server is listening at %s:%d.", ip, port)
	} else {
		log.Infof("Web server is listening on port %d.", port)
	}
	log.Infoln("Configure this server by visiting /admin.")

	return server.ListenAndServe()
}

func addStaticFileEndpoints(r chi.Router, h *handlers.Handlers, apc *apcontrollers.Controllers) {
	// Images
	r.HandleFunc("/thumbnail.jpg", h.GetThumbnail)
	r.HandleFunc("/preview.gif", h.GetPreview)
	r.HandleFunc("/logo", h.GetLogo)
	r.HandleFunc("/favicon.ico", h.GetFavicon)
	// return a logo that's compatible with external social networks
	r.HandleFunc("/logo/external", h.GetCompatibleLogo)

	// Custom Javascript
	r.HandleFunc("/customjavascript", h.ServeCustomJavascript)

	// robots.txt
	r.HandleFunc("/robots.txt", h.GetRobotsDotTxt)

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
