package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/metrics"
	"github.com/owncast/owncast/services/status"
	"github.com/owncast/owncast/webserver/middleware"
	log "github.com/sirupsen/logrus"
)

// Status gets the details of the inbound broadcaster.
func (h *Handlers) GetAdminStatus(w http.ResponseWriter, r *http.Request) {
	s := status.Get()
	m := metrics.Get()

	broadcaster := s.GetBroadcaster()
	currentBroadcast := s.GetCurrentBroadcast()
	health := m.GetStreamHealthOverview()
	response := adminStatusResponse{
		Broadcaster:            broadcaster,
		CurrentBroadcast:       currentBroadcast,
		Online:                 s.Online,
		Health:                 health,
		ViewerCount:            s.ViewerCount,
		OverallPeakViewerCount: s.OverallMaxViewerCount,
		SessionPeakViewerCount: s.SessionMaxViewerCount,
		VersionNumber:          s.VersionNumber,
		StreamTitle:            configRepository.GetStreamTitle(),
	}

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}

type adminStatusResponse struct {
	Broadcaster            *models.Broadcaster          `json:"broadcaster"`
	CurrentBroadcast       *models.CurrentBroadcast     `json:"currentBroadcast"`
	Health                 *models.StreamHealthOverview `json:"health"`
	StreamTitle            string                       `json:"streamTitle"`
	VersionNumber          string                       `json:"versionNumber"`
	ViewerCount            int                          `json:"viewerCount"`
	OverallPeakViewerCount int                          `json:"overallPeakViewerCount"`
	SessionPeakViewerCount int                          `json:"sessionPeakViewerCount"`
	Online                 bool                         `json:"online"`
}
