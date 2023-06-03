package yp

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

type ypDetailsResponse struct {
	LastConnectTime *utils.NullTime       `json:"lastConnectTime"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	StreamTitle     string                `json:"streamTitle,omitempty"`
	Logo            string                `json:"logo"`
	Tags            []string              `json:"tags"`
	Social          []models.SocialHandle `json:"social"`

	ViewerCount           int  `json:"viewerCount"`
	OverallMaxViewerCount int  `json:"overallMaxViewerCount"`
	SessionMaxViewerCount int  `json:"sessionMaxViewerCount"`
	NSFW                  bool `json:"nsfw"`
	Online                bool `json:"online"`
}

// GetYPResponse gets the status of the server for YP purposes.
func GetYPResponse(w http.ResponseWriter, r *http.Request) {
	if !data.GetDirectoryEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	status := getStatus()

	streamTitle := data.GetStreamTitle()

	response := ypDetailsResponse{
		Name:                  data.GetServerName(),
		Description:           data.GetServerSummary(),
		StreamTitle:           streamTitle,
		Logo:                  "/logo",
		NSFW:                  data.GetNSFW(),
		Tags:                  data.GetServerMetadataTags(),
		Online:                status.Online,
		ViewerCount:           status.ViewerCount,
		OverallMaxViewerCount: status.OverallMaxViewerCount,
		SessionMaxViewerCount: status.SessionMaxViewerCount,
		LastConnectTime:       status.LastConnectTime,
		Social:                data.GetSocialHandles(),
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}
