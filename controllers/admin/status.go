package admin

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/app/middleware"
	"github.com/owncast/owncast/models"
)

// Status gets the details of the inbound broadcaster.
func (c *Controller) Status(w http.ResponseWriter, r *http.Request) {
	broadcaster := c.Core.GetBroadcaster()
	status := c.Core.GetStatus()
	currentBroadcast := c.Core.GetCurrentBroadcast()
	health := c.Metrics.GetStreamHealthOverview()
	response := adminStatusResponse{
		Broadcaster:            broadcaster,
		CurrentBroadcast:       currentBroadcast,
		Online:                 status.Online,
		Health:                 health,
		ViewerCount:            status.ViewerCount,
		OverallPeakViewerCount: status.OverallMaxViewerCount,
		SessionPeakViewerCount: status.SessionMaxViewerCount,
		VersionNumber:          status.VersionNumber,
		StreamTitle:            c.Data.GetStreamTitle(),
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
	Online                 bool                         `json:"online"`
	ViewerCount            int                          `json:"viewerCount"`
	OverallPeakViewerCount int                          `json:"overallPeakViewerCount"`
	SessionPeakViewerCount int                          `json:"sessionPeakViewerCount"`
	StreamTitle            string                       `json:"streamTitle"`
	Health                 *models.StreamHealthOverview `json:"health"`
	VersionNumber          string                       `json:"versionNumber"`
}
