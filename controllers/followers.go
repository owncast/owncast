package controllers

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/persistence"
)

// GetFollowers will handle an API request to fetch the list of followers (non-activitypub response).
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	followers, err := persistence.GetFederationFollowers(-1, 0)
	if err != nil {
		WriteSimpleResponse(w, false, "unable to fetch followers")
		return
	}

	WriteResponse(w, followers)
}
