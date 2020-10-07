package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/amalfra/etag"
	"github.com/owncast/owncast/utils"
)

//DisableCache writes the disable cache header on the responses
func DisableCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Expires", "Thu, 1 Jan 1970 00:00:00 GMT")
}

func setCacheSeconds(seconds int, w http.ResponseWriter) {
	secondsStr := strconv.Itoa(seconds)
	w.Header().Set("Cache-Control", "public, max-age="+secondsStr)
}

// ProcessEtags gets and sets ETags for caching purposes
func ProcessEtags(w http.ResponseWriter, r *http.Request) int {
	info, err := os.Stat(filepath.Join("webroot", r.URL.Path))
	if err != nil {
		return 0
	}

	localContentEtag := etag.Generate(info.ModTime().String(), true)
	if remoteEtagHeader := r.Header.Get("If-None-Match"); remoteEtagHeader != "" {
		if remoteEtagHeader == localContentEtag {
			return http.StatusNotModified
		}
	}

	w.Header().Set("Etag", localContentEtag)

	return 0
}

// SetCachingHeaders will set the cache control header of a response
func SetCachingHeaders(w http.ResponseWriter, r *http.Request) {
	setCacheSeconds(utils.GetCacheDurationSecondsForPath(r.URL.Path), w)
}
