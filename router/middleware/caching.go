package middleware

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/amalfra/etag"
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
	setCacheSeconds(getCacheDurationSecondsForPath(r.URL.Path), w)
}

func getCacheDurationSecondsForPath(filePath string) int {
	if path.Base(filePath) == "thumbnail.jpg" {
		// Thumbnails re-generate during live
		return 20
	} else if path.Ext(filePath) == ".js" || path.Ext(filePath) == ".css" {
		// Cache javascript & CSS
		return 60
	} else if path.Ext(filePath) == ".ts" {
		// Cache video segments as long as you want. They can't change.
		// This matters most for local hosting of segments for recordings
		// and not for live or 3rd party storage.
		return 31557600
	}

	// Default cache length in seconds
	return 30 * 60
}
