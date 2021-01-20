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
	Name             string                `json:"name"`
	Title            string                `json:"title"`
	Summary          string                `json:"summary"`
	Logo             string                `json:"logo"`
	Tags             []string              `json:"tags"`
	Version          string                `json:"version"`
	NSFW             bool                  `json:"nsfw"`
	ExtraPageContent string                `json:"extraPageContent"`
	StreamTitle      string                `json:"streamTitle,omitempty"` // What's going on with the current stream
	SocialHandles    []models.SocialHandle `json:"socialHandles"`
}

// GetWebConfig gets the status of the server.
func GetWebConfig(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	pageContent := utils.RenderPageContentMarkdown(data.GetExtraPageBodyContent())

	configuration := webConfigResponse{
		Name:             data.GetServerName(),
		Title:            data.GetServerTitle(),
		Summary:          data.GetServerSummary(),
		Logo:             data.GetLogoPath(),
		Tags:             data.GetServerMetadataTags(),
		Version:          config.VersionInfo,
		ExtraPageContent: pageContent,
		StreamTitle:      data.GetStreamTitle(),
		SocialHandles:    data.GetSocialHandles(),
	}

	if err := json.NewEncoder(w).Encode(configuration); err != nil {
		BadRequestHandler(w, err)
	}
}

func GetAllSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	platforms := models.GetAllSocialHandles()
	if err := json.NewEncoder(w).Encode(platforms); err != nil {
		internalErrorHandler(w, err)
	}
}
