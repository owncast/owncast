package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/router/middleware"
)

//GetStatus gets the status of the server
func GetStatus(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	status := core.GetStatus()
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(status)
}
