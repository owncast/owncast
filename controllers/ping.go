package controllers

import (
	"net/http"

	"github.com/owncast/owncast/models"
)

// Ping is fired by a client to show they are still an active viewer.
func (s *Service) Ping(w http.ResponseWriter, r *http.Request) {
	viewer := models.GenerateViewerFromRequest(r)
	s.Core.SetViewerActive(&viewer)
	w.WriteHeader(http.StatusOK)
}
