package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
)

// GetWebConfig gets the status of the server.
func GetWebConfig(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	pageContent := utils.RenderPageContentMarkdown(data.GetExtraPageBodyContent())

	configuration := config.InstanceDetails{
		Name:             data.GetServerName(),
		Title:            data.GetStreamTitle(),
		Summary:          data.GetServerSummary(),
		Logo:             data.GetLogoPath(),
		Tags:             data.GetServerMetadataTags(),
		Version:          config.Config.VersionInfo,
		ExtraPageContent: pageContent,
		StreamTitle:      data.GetStreamTitle(),
	}

	if err := json.NewEncoder(w).Encode(configuration); err != nil {
		badRequestHandler(w, err)
	}
}
