package admin

import (
	"encoding/json"
	"net/http"

	"github.com/gabek/owncast/controllers"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/models"
	"github.com/gabek/owncast/router/middleware"
)

// GetInboundBroadasterDetails gets the details of the inbound broadcaster
func GetInboundBroadasterDetails(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	broadcaster := core.GetBroadcaster()
	if broadcaster == nil {
		controllers.WriteSimpleResponse(w, false, "no broadcaster connected")
		return
	}

	response := inboundBroadasterDetailsResponse{
		models.BaseAPIResponse{
			true,
			"",
		},
		broadcaster,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type inboundBroadasterDetailsResponse struct {
	models.BaseAPIResponse
	Broadcaster *models.Broadcaster `json:"broadcaster"`
}
