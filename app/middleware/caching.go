package middleware

import (
	"net/http"
	"strconv"

	"github.com/owncast/owncast/utils"
)

// DisableCache writes the disable cache header on the responses.
func DisableCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Expires", "Thu, 1 Jan 1970 00:00:00 GMT")
}

func setCacheSeconds(seconds int, w http.ResponseWriter) {
	secondsStr := strconv.Itoa(seconds)
	w.Header().Set("Cache-Control", "public, max-age="+secondsStr)
}

// SetCachingHeaders will set the cache control header of a response.
func SetCachingHeaders(w http.ResponseWriter, r *http.Request) {
	setCacheSeconds(utils.GetCacheDurationSecondsForPath(r.URL.Path), w)
}
