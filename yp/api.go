package yp

import (
	"encoding/json"
	"net/http"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/utils"
)

type ypDetailsResponse struct {
	NSFW                  bool `json:"nsfw"`
	Online                bool `json:"online"`
	ViewerCount           int  `json:"viewerCount"`
	OverallMaxViewerCount int  `json:"overallMaxViewerCount"`
	SessionMaxViewerCount int  `json:"sessionMaxViewerCount"`

	LastConnectTime utils.NullTime `json:"lastConnectTime"`
}

//GetYPResponse gets the status of the server for YP purposes
func GetYPResponse(w http.ResponseWriter, r *http.Request) {
	status := getStatus()

	response := ypDetailsResponse{
		NSFW:                  config.Config.InstanceDetails.NSFW,
		Online:                status.Online,
		ViewerCount:           status.ViewerCount,
		OverallMaxViewerCount: status.OverallMaxViewerCount,
		SessionMaxViewerCount: status.SessionMaxViewerCount,
		LastConnectTime:       status.LastConnectTime,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}
