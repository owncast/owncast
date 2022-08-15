package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
)

// IndexHandler handles the default index route.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	isIndexRequest := r.URL.Path == "/" || filepath.Base(r.URL.Path) == "index.html" || filepath.Base(r.URL.Path) == ""

	if utils.IsUserAgentAPlayer(r.UserAgent()) && isIndexRequest {
		http.Redirect(w, r, "/hls/stream.m3u8", http.StatusTemporaryRedirect)
		return
	}

	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)

	// Set our global HTTP headers
	middleware.SetHeaders(w)

	serveWeb(w, r)
}
