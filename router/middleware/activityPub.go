package middleware

import (
	"net/http"
	"strings"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
)

// RequireActivityPubOrRedirect will validate a provided user's access token and make sure the associated user is enabled.
// Not to be used for validating 3rd party access.
func RequireActivityPubOrRedirect(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !data.GetFederationEnabled() {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handleAccepted := func() {
			handler(w, r)
		}

		acceptedContentTypes := []string{"application/json", "application/json+ld", "application/activity+json"}
		acceptString := r.Header.Get("Accept")
		accept := strings.Split(acceptString, ",")

		for _, singleType := range accept {
			if _, accepted := utils.FindInSlice(acceptedContentTypes, strings.TrimSpace(singleType)); accepted {
				handleAccepted()
				return
			}
		}

		contentTypeString := r.Header.Get("Content-Type")
		contentTypes := strings.Split(contentTypeString, ",")
		for _, singleType := range contentTypes {
			if _, accepted := utils.FindInSlice(acceptedContentTypes, strings.TrimSpace(singleType)); accepted {
				handleAccepted()
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
}
