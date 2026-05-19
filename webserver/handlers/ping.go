package handlers

import (
	"net/http"

	"github.com/owncast/owncast/models"
)

// Ping is fired by a client to show they are still an active viewer.
func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	viewer := models.GenerateViewerFromRequest(r)
	h.stream.SetViewerActive(&viewer)
	w.WriteHeader(http.StatusOK)
}
