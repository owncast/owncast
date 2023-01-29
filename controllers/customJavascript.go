package controllers

import (
	"net/http"

	"github.com/owncast/owncast/core/data"
)

// ServeCustomJavascript will serve optional custom Javascript.
func ServeCustomJavascript(w http.ResponseWriter, r *http.Request) {
	js := data.GetCustomJavascript()
	_, _ = w.Write([]byte(js))
}
