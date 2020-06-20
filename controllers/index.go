package controllers

import (
	"net/http"
	"path"

	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/router/middleware"
	"github.com/gabek/owncast/utils"
)

//IndexHandler handles the default index route
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	http.ServeFile(w, r, path.Join("webroot", r.URL.Path))

	if path.Ext(r.URL.Path) == ".m3u8" {
		middleware.DisableCache(&w)

		clientID := utils.GenerateClientIDFromRequest(r)
		core.SetClientActive(clientID)
	}
}
