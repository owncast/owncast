package controllers

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/persistence"
)

// GetFollowers will handle an API request to fetch the list of followers (non-activitypub response).
func GetFollowers(offset int, limit int, w http.ResponseWriter, r *http.Request) {
	followers, total, err := persistence.GetFederationFollowers(limit, offset)
	if err != nil {
		WriteSimpleResponse(w, false, "unable to fetch followers")
		return
	}

	response := PaginatedResponse{
		Total:   total,
		Results: followers,
	}
	WriteResponse(w, response)
}
