package handlers

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/persistence/followersrepository"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// GetFollowers will handle an API request to fetch the list of followers (non-activitypub response).
func GetFollowers(offset int, limit int, w http.ResponseWriter, r *http.Request) {
	followersRepo := followersrepository.Get()
	followers, total, err := followersRepo.GetFollowers(limit, offset)
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
