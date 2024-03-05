package controllers

import (
	"net/http"

	"github.com/owncast/owncast/core/data"
)

// ServeCustomJavascript will serve optional custom Javascript.
func ServeCustomJavascript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")

	js := data.GetCustomJavascript()
	_, _ = w.Write([]byte(js))
}
