package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// Status gets the details of the inbound broadcaster.
func Status(w http.ResponseWriter, r *http.Request) {
	broadcaster := core.GetBroadcaster()
	status := core.GetStatus()
	currentBroadcast := core.GetCurrentBroadcast()

	response := adminStatusResponse{
		Broadcaster:            broadcaster,
		CurrentBroadcast:       currentBroadcast,
		Online:                 status.Online,
		ViewerCount:            status.ViewerCount,
		OverallPeakViewerCount: status.OverallMaxViewerCount,
		SessionPeakViewerCount: status.SessionMaxViewerCount,
		VersionNumber:          status.VersionNumber,
		DisableUpgradeChecks:   data.GetDisableUpgradeChecks(),
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}

type adminStatusResponse struct {
	Broadcaster            *models.Broadcaster      `json:"broadcaster"`
	CurrentBroadcast       *models.CurrentBroadcast `json:"currentBroadcast"`
	Online                 bool                     `json:"online"`
	ViewerCount            int                      `json:"viewerCount"`
	OverallPeakViewerCount int                      `json:"overallPeakViewerCount"`
	SessionPeakViewerCount int                      `json:"sessionPeakViewerCount"`

	VersionNumber        string `json:"versionNumber"`
	DisableUpgradeChecks bool   `json:"disableUpgradeChecks"`
}
