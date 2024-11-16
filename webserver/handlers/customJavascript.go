package handlers

import (
	"net/http"

	"github.com/owncast/owncast/persistence/configrepository"
)

// ServeCustomJavascript will serve optional custom Javascript.
func ServeCustomJavascript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")

	configRepository := configrepository.Get()
	js := configRepository.GetCustomJavascript()
	_, _ = w.Write([]byte(js))
}
