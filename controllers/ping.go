package controllers

import (
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/models"
)

// Ping is fired by a client to show they are still an active viewer.
func Ping(w http.ResponseWriter, r *http.Request) {
	viewer := models.GenerateViewerFromRequest(r)
	core.SetViewerActive(&viewer)
	w.WriteHeader(http.StatusOK)
}
