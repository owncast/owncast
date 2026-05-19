package handlers

import (
	"net/http"
)

// ServeCustomJavascript will serve optional custom Javascript.
func (h *Handlers) ServeCustomJavascript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	_, _ = w.Write([]byte(h.configRepository.GetCustomJavascript()))
}
