package handlers

import (
	"net/http"

	"github.com/owncast/owncast/storage/federationrepository"
	"github.com/owncast/owncast/webserver/responses"
)

// GetFollowers will handle an API request to fetch the list of followers (non-activitypub response).
func (h *Handlers) GetFollowers(offset int, limit int, w http.ResponseWriter, r *http.Request) {
	federationRepository := federationrepository.Get()
	followers, total, err := federationRepository.GetFederationFollowers(limit, offset)
	if err != nil {
		responses.WriteSimpleResponse(w, false, "unable to fetch followers")
		return
	}

	response := responses.PaginatedResponse{
		Total:   total,
		Results: followers,
	}
	responses.WriteResponse(w, response)
}
