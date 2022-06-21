package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
)

// MetadataPage represents a server-rendered web page for bots and web scrapers.
type MetadataPage struct {
	RequestedURL  string
	Image         string
	Thumbnail     string
	TagsString    string
	Summary       string
	Name          string
	Tags          []string
	SocialHandles []models.SocialHandle
}

// IndexHandler handles the default index route.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	isIndexRequest := r.URL.Path == "/" || filepath.Base(r.URL.Path) == "index.html" || filepath.Base(r.URL.Path) == ""

	if utils.IsUserAgentAPlayer(r.UserAgent()) && isIndexRequest {
		http.Redirect(w, r, "/hls/stream.m3u8", http.StatusTemporaryRedirect)
		return
	}

	// If the ETags match then return a StatusNotModified
	// if responseCode := middleware.ProcessEtags(w, r); responseCode != 0 {
	// 	w.WriteHeader(responseCode)
	// 	return
	// }

	// If this is a directory listing request then return a 404
	// info, err := os.Stat(path.Join(config.WebRoot, r.URL.Path))
	// if err != nil || (info.IsDir() && !isIndexRequest) {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)

	// Set our global HTTP headers
	middleware.SetHeaders(w)

	serveWeb(w, r)
}
