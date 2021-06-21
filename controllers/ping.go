package controllers

import (
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/utils"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	id := utils.GenerateClientIDFromRequest(r)
	core.SetViewerIdActive(id)
}
