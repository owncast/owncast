package handlers

import (
	"net/http"
)

// ServeCustomJavascript will serve optional custom Javascript.
func (h *Handlers) ServeCustomJavascript(w http.ResponseWriter, r *http.Request) {
	js := configRepository.GetCustomJavascript()
	_, _ = w.Write([]byte(js))
}
