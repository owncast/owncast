package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

// ExternalAccessTokenHandlerFunc is a function that is called after validing access.
type ExternalAccessTokenHandlerFunc func(models.ExternalAPIUser, http.ResponseWriter, *http.Request)

// UserAccessTokenHandlerFunc is a function that is called after validing user access.
type UserAccessTokenHandlerFunc func(models.User, http.ResponseWriter, *http.Request)

// adminAuthRealm is the WWW-Authenticate realm string used by every admin
// auth challenge in Owncast. Anything that gates on admin Basic Auth (the
// main admin API, the plugin management API, plugin admin pages) must
// challenge with this exact realm so the browser shares one credential
// cache across all of them.
const adminAuthRealm = "Owncast Authenticated Request"

// IsAdminRequest reports whether r carries valid admin credentials. The
// Owncast admin UI sends Basic Auth on every API call; embedded contexts
// that cannot inject a custom Authorization header (notably plugin admin
// iframes) authenticate via the admin session cookie instead. Shared by
// RequireAdminAuth and any caller that needs the same check without
// rejecting the request (e.g. a handler filling an "authenticated"
// boolean on a downstream payload).
func (m *Middleware) IsAdminRequest(r *http.Request) bool {
	if user, pass, ok := r.BasicAuth(); ok {
		if subtle.ConstantTimeCompare([]byte(user), []byte("admin")) == 1 &&
			utils.CompareHash(m.configRepository.GetAdminPassword(), pass) == nil {
			return true
		}
	}
	return m.hasValidAdminSessionCookie(r)
}

// RequireAdminAuth wraps a handler requiring HTTP basic auth for it using
// the admin password as the password and a hardcoded "admin" for username.
//
// As a side effect, a valid Basic Auth request also primes an admin
// session cookie. Embedded contexts that can't inject the Authorization
// header (notably plugin admin iframes) authenticate via that cookie on
// their next same-origin request, so the user isn't prompted by the
// browser's native Basic Auth dialog.
func (m *Middleware) RequireAdminAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow CORS only for localhost:3000 to support Owncast development.
		validAdminHost := "http://localhost:3000"
		w.Header().Set("Access-Control-Allow-Origin", validAdminHost)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// For request needing CORS, send a 204.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if !m.IsAdminRequest(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+adminAuthRealm+`"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Debugln("Failed admin authentication")
			return
		}

		// Prime the admin session cookie if the request authenticated via
		// Basic Auth and doesn't already carry a valid cookie. No-op when
		// the cookie is already fresh, so most admin API calls don't pay
		// for it (and don't rotate the token unnecessarily).
		m.ensureAdminSessionCookie(w, r)

		handler(w, r)
	}
}

func accessDenied(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized) //nolint
	w.Write([]byte("unauthorized"))        //nolint
}

// RequireExternalAPIAccessToken will validate a 3rd party access token.
func (m *Middleware) RequireExternalAPIAccessToken(scope string, handler ExternalAccessTokenHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We should accept 3rd party preflight OPTIONS requests.
		if r.Method == "OPTIONS" {
			// All OPTIONS requests should have a wildcard CORS header.
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		authHeader := r.Header.Get("Authorization")
		token := ""
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = authHeader[len("bearer "):]
		}

		if token == "" {
			log.Warnln("invalid access token")
			accessDenied(w)
			return
		}

		integration, err := m.userRepository.GetExternalAPIUserForAccessTokenAndScope(token, scope)
		if integration == nil || err != nil {
			accessDenied(w)
			return
		}

		// All auth'ed 3rd party requests should have a wildcard CORS header.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		handler(*integration, w, r)

		if err := m.userRepository.SetExternalAPIUserAccessTokenAsUsed(token); err != nil {
			log.Debugln("token not found when updating last_used timestamp")
		}
	})
}

// RequireUserAccessToken will validate a provided user's access token and make sure the associated user is enabled.
// Not to be used for validating 3rd party access.
func (m *Middleware) RequireUserAccessToken(handler UserAccessTokenHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.URL.Query().Get("accessToken")
		if accessToken == "" {
			accessDenied(w)
			return
		}

		ipAddress := utils.GetIPAddressFromRequest(r)
		// Check if this client's IP address is banned.
		if blocked, err := m.authRepository.IsIPAddressBanned(ipAddress); blocked {
			log.Debugln("Client ip address has been blocked. Rejecting.")
			accessDenied(w)
			return
		} else if err != nil {
			log.Errorln("error determining if IP address is blocked: ", err)
		}

		// A user is required to use the websocket
		user := m.userRepository.GetUserByToken(accessToken)
		if user == nil || !user.IsEnabled() {
			accessDenied(w)
			return
		}

		handler(*user, w, r)
	})
}

// RequireUserModerationScopeAccesstoken will validate a provided user's access token and make sure the associated user is enabled
// and has "MODERATOR" scope assigned to the user.
func (m *Middleware) RequireUserModerationScopeAccesstoken(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.URL.Query().Get("accessToken")
		if accessToken == "" {
			accessDenied(w)
			return
		}

		// A user is required to use the websocket
		user := m.userRepository.GetUserByToken(accessToken)
		if user == nil || !user.IsEnabled() || !user.IsModerator() {
			accessDenied(w)
			return
		}

		handler(w, r)
	})
}
