package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/models"
)

// Status gets the details of the inbound broadcaster
func Status(w http.ResponseWriter, r *http.Request) {
	broadcaster := core.GetBroadcaster()
	status := core.GetStatus()

	response := adminStatusResponse{
		Broadcaster:            broadcaster,
		Online:                 status.Online,
		ViewerCount:            status.ViewerCount,
		OverallPeakViewerCount: status.OverallMaxViewerCount,
		SessionPeakViewerCount: status.SessionMaxViewerCount,
		VersionNumber:          status.VersionNumber,
		DisableUpgradeChecks:   config.Config.DisableUpgradeChecks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type adminStatusResponse struct {
	Broadcaster            *models.Broadcaster `json:"broadcaster"`
	Online                 bool                `json:"online"`
	ViewerCount            int                 `json:"viewerCount"`
	OverallPeakViewerCount int                 `json:"overallPeakViewerCount"`
	SessionPeakViewerCount int                 `json:"sessionPeakViewerCount"`

	VersionNumber        string `json:"versionNumber"`
	DisableUpgradeChecks bool   `json:"disableUpgradeChecks"`
}
