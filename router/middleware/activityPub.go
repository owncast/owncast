package middleware

import "net/http"

// RequireUserAccessToken will validate a provided user's access token and make sure the associated user is enabled.
// Not to be used for validating 3rd party access.
func RequireActivityPubOrRedirect(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		if accept != "application/json" && accept != "application/json+ld" && accept != "application/activity+json" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		handler(w, r)
	})
}
