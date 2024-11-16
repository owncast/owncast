package yp

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
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
	configRepository := configrepository.Get()
	if !configRepository.GetDirectoryEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	status := getStatus()

	streamTitle := configRepository.GetStreamTitle()

	response := ypDetailsResponse{
		Name:                  configRepository.GetServerName(),
		Description:           configRepository.GetServerSummary(),
		StreamTitle:           streamTitle,
		Logo:                  "/logo",
		NSFW:                  configRepository.GetNSFW(),
		Tags:                  configRepository.GetServerMetadataTags(),
		Online:                status.Online,
		ViewerCount:           status.ViewerCount,
		OverallMaxViewerCount: status.OverallMaxViewerCount,
		SessionMaxViewerCount: status.SessionMaxViewerCount,
		LastConnectTime:       status.LastConnectTime,
		Social:                configRepository.GetSocialHandles(),
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}
