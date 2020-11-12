package yp

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

type ypDetailsResponse struct {
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	Logo                  string   `json:"logo"`
	NSFW                  bool     `json:"nsfw"`
	Tags                  []string `json:"tags"`
	Online                bool     `json:"online"`
	ViewerCount           int      `json:"viewerCount"`
	OverallMaxViewerCount int      `json:"overallMaxViewerCount"`
	SessionMaxViewerCount int      `json:"sessionMaxViewerCount"`

	LastConnectTime utils.NullTime `json:"lastConnectTime"`
}

// GetYPResponse gets the status of the server for YP purposes.
func GetYPResponse(w http.ResponseWriter, r *http.Request) {
	status := getStatus()

	response := ypDetailsResponse{
		Name:                  config.Config.InstanceDetails.Name,
		Description:           config.Config.InstanceDetails.Summary,
		Logo:                  config.Config.InstanceDetails.Logo.Large,
		NSFW:                  config.Config.InstanceDetails.NSFW,
		Tags:                  config.Config.InstanceDetails.Tags,
		Online:                status.Online,
		ViewerCount:           status.ViewerCount,
		OverallMaxViewerCount: status.OverallMaxViewerCount,
		SessionMaxViewerCount: status.SessionMaxViewerCount,
		LastConnectTime:       status.LastConnectTime,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}
