package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CAFxX/httpcompression"
	"github.com/owncast/owncast/webserver/handlers"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type webServer struct {
	router   *http.ServeMux
	handlers *handlers.Handlers
	server   *http.Server
}

func New() *webServer {
	s := &webServer{
		router: http.NewServeMux(),
	}

	s.setupRoutes()

	return s
}

func (s *webServer) Start(listenIP string, listenPort int) error {
	compress, _ := httpcompression.DefaultAdapter() // Use the default configuration
	h2s := &http2.Server{}
	http2Router := h2c.NewHandler(s.router, h2s)

	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", listenIP, listenPort),
		ReadHeaderTimeout: 4 * time.Second,
		Handler:           compress(http2Router),
	}

	if listenIP != "0.0.0.0" {
		log.Infof("Web server is listening at %s:%d.", listenIP, listenPort)
	} else {
		log.Infof("Web server is listening on port %d.", listenPort)
	}

	return s.server.ListenAndServe()
}

// ServeHTTP is the entry point for all web requests.
func (s *webServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
