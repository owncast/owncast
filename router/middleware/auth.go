package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gabek/owncast/config"
	log "github.com/sirupsen/logrus"
)

// RequireAdminAuth wraps a handler requiring HTTP basic auth for it using the given
// the stream key as the password and and a hardcoded "admin" for username.
func RequireAdminAuth(handler http.HandlerFunc) http.HandlerFunc {
	username := "admin"
	password := config.Config.VideoSettings.StreamingKey

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()
		realm := "Owncast Authenticated Request"

		// Failed
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Warnln("Failed authentication for", r.URL.Path, "from", r.RemoteAddr, r.UserAgent())
			return
		}

		// Success
		log.Traceln("Authenticated request OK for", r.URL.Path, "from", r.RemoteAddr, r.UserAgent())
		handler(w, r)
	}
}
