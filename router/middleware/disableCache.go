package middleware

import (
	"net/http"
)

//DisableCache writes the disable cache header on the responses
func DisableCache(w *http.ResponseWriter) {
	(*w).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*w).Header().Set("Expires", "0")
}
