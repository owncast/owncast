package middleware

import (
	"net/http"

	"github.com/owncast/owncast/utils"
)

// RequireUserAccessToken will validate a provided user's access token and make sure the associated user is enabled.
// Not to be used for validating 3rd party access.
func RequireActivityPubOrRedirect(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptedContentTypes := []string{"application/json", "application/json+ld", "application/activity+json", "application/activity+json, application/ld+json"}
		accept := r.Header.Get("Accept")
		contentType := r.Header.Get("Content-Type")

		_, acceptMatches := utils.FindInSlice(acceptedContentTypes, accept)
		_, contentTypeMatches := utils.FindInSlice(acceptedContentTypes, contentType)
		if !acceptMatches && !contentTypeMatches {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		handler(w, r)
	})
}
