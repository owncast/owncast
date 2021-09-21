package middleware

import (
	"net/http"
	"strings"
)

// SetHeaders will set our global headers for web resources.
func SetHeaders(w http.ResponseWriter) {
	// Tell Google to not use this response in their FLoC tracking.
	w.Header().Set("Permissions-Policy", "interest-cohort=()")

	// Content security policy
	csp := []string{
		"script-src 'self' 'sha256-rnxPrBaD0OuYxsCdrll4QJwtDLcBJqFh0u27CoX5jZ8='",
		"worker-src 'self' blob:", // No single quotes around blob:
	}
	w.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
}
