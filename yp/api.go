package yp

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

type ypDetailsResponse struct {
	Name                  string                `json:"name"`
	Description           string                `json:"description"`
	StreamTitle           string                `json:"streamTitle,omitempty"`
	Logo                  string                `json:"logo"`
	NSFW                  bool                  `json:"nsfw"`
	Tags                  []string              `json:"tags"`
	Online                bool                  `json:"online"`
	ViewerCount           int                   `json:"viewerCount"`
	OverallMaxViewerCount int                   `json:"overallMaxViewerCount"`
	SessionMaxViewerCount int                   `json:"sessionMaxViewerCount"`
	Social                []models.SocialHandle `json:"social"`

	LastConnectTime *utils.NullTime `json:"lastConnectTime"`
}

// GetYPResponse gets the status of the server for YP purposes.
func (yp *YP) GetYPResponse(w http.ResponseWriter, r *http.Request) {
	if !yp.data.GetDirectoryEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	status := getStatus()

	streamTitle := yp.data.GetStreamTitle()

	response := ypDetailsResponse{
		Name:                  yp.data.GetServerName(),
		Description:           yp.data.GetServerSummary(),
		StreamTitle:           streamTitle,
		Logo:                  "/logo",
		NSFW:                  yp.data.GetNSFW(),
		Tags:                  yp.data.GetServerMetadataTags(),
		Online:                status.Online,
		ViewerCount:           status.ViewerCount,
		OverallMaxViewerCount: status.OverallMaxViewerCount,
		SessionMaxViewerCount: status.SessionMaxViewerCount,
		LastConnectTime:       status.LastConnectTime,
		Social:                yp.data.GetSocialHandles(),
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}
