package controllers

import (
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/utils"
)

// Ping is fired by a client to show they are still an active viewer.
func Ping(w http.ResponseWriter, r *http.Request) {
	id := utils.GenerateClientIDFromRequest(r)
	core.SetViewerIDActive(id)
	WriteSimpleResponse(w, true, "OK")
}
