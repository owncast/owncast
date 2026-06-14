package handlers

import (
	"net/http"

	webutils "github.com/owncast/owncast/webserver/utils"
)

// GetFollowers will handle an API request to fetch the list of followers (non-activitypub response).
func (h *Handlers) GetFollowers(offset int, limit int, w http.ResponseWriter, r *http.Request) {
	followers, total, err := h.followersRepository.GetFollowers(limit, offset)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, "unable to fetch followers")
		return
	}

	response := webutils.PaginatedResponse{
		Total:   total,
		Results: followers,
	}
	webutils.WriteResponse(w, response)
}
