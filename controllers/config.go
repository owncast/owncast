package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/router/middleware"
)

// GetWebConfig gets the status of the server.
func GetWebConfig(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	configuration := config.Config.InstanceDetails
	configuration.Version = config.Config.VersionInfo
	if err := json.NewEncoder(w).Encode(configuration); err != nil {
		badRequestHandler(w, err)
	}
}
