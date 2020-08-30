package middleware

import (
	"net/http"
	"strconv"
)

//DisableCache writes the disable cache header on the responses
func DisableCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Expires", "0")
}

//SetCache will set the cache control header of a response
func SetCache(days int, w http.ResponseWriter) {
	seconds := strconv.Itoa(days * 86400)
	w.Header().Set("Cache-Control", "max-age="+seconds)
	w.Header().Set("Expires", seconds)
}
