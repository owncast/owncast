package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
)

type webConfigResponse struct {
	Name                 string                  `json:"name"`
	Summary              string                  `json:"summary"`
	Logo                 string                  `json:"logo"`
	Tags                 []string                `json:"tags"`
	Version              string                  `json:"version"`
	NSFW                 bool                    `json:"nsfw"`
	ExtraPageContent     string                  `json:"extraPageContent"`
	StreamTitle          string                  `json:"streamTitle,omitempty"` // What's going on with the current stream
	SocialHandles        []models.SocialHandle   `json:"socialHandles"`
	ChatDisabled         bool                    `json:"chatDisabled"`
	ExternalActions      []models.ExternalAction `json:"externalActions"`
	CustomStyles         string                  `json:"customStyles"`
	MaxSocketPayloadSize int                     `json:"maxSocketPayloadSize"`
}

// GetWebConfig gets the status of the server.
func GetWebConfig(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	middleware.DisableCache(w)
	w.Header().Set("Content-Type", "application/json")

	pageContent := utils.RenderPageContentMarkdown(data.GetExtraPageBodyContent())
	socialHandles := data.GetSocialHandles()
	for i, handle := range socialHandles {
		platform := models.GetSocialHandle(handle.Platform)
		if platform != nil {
			handle.Icon = platform.Icon
			socialHandles[i] = handle
		}
	}

	serverSummary := data.GetServerSummary()
	serverSummary = utils.RenderPageContentMarkdown(serverSummary)

	configuration := webConfigResponse{
		Name:                 data.GetServerName(),
		Summary:              serverSummary,
		Logo:                 "/logo",
		Tags:                 data.GetServerMetadataTags(),
		Version:              config.GetReleaseString(),
		NSFW:                 data.GetNSFW(),
		ExtraPageContent:     pageContent,
		StreamTitle:          data.GetStreamTitle(),
		SocialHandles:        socialHandles,
		ChatDisabled:         data.GetChatDisabled(),
		ExternalActions:      data.GetExternalActions(),
		CustomStyles:         data.GetCustomStyles(),
		MaxSocketPayloadSize: config.MaxSocketPayloadSize,
	}

	if err := json.NewEncoder(w).Encode(configuration); err != nil {
		BadRequestHandler(w, err)
	}
}

// GetAllSocialPlatforms will return a list of all social platform types.
func GetAllSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	w.Header().Set("Content-Type", "application/json")

	platforms := models.GetAllSocialHandles()
	if err := json.NewEncoder(w).Encode(platforms); err != nil {
		InternalErrorHandler(w, err)
	}
}
