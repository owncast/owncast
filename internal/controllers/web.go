package controllers

import (
	"net/http"

	"github.com/owncast/owncast/internal/static"
)

var staticServer = http.FileServer(http.FS(static.GetWeb()))

// serveWeb will serve web assets.
func serveWeb(w http.ResponseWriter, r *http.Request) {
	staticServer.ServeHTTP(w, r)
}
