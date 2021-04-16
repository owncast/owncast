package middleware

import "net/http"

// DisableFloc will tell Google to not use this response in their FLoC tracking.
func DisableFloc(w http.ResponseWriter) {
	w.Header().Set("Permissions-Policy", "interest-cohort=()")
}
