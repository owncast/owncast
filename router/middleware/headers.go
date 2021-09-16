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
		"script-src 'self' 'sha256-2HPCfJIJHnY0NrRDPTOdC7AOSJIcQyNxzUuut3TsYRY=' 'sha256-qYEKg5UMg/KbbMBkyPIGsxtkfn/safeLBT08DK3592g=' 'sha256-2erOadwY1DsoNdxVjGlxldMJrFEUzr5sLDdB8lmm9m8=' 'sha256-DgrU+KwEGMFcB8B2ZdQyuxWWvTm7LeGpc+8SkxbSxGA='",
		"worker-src 'self' blob:", // No single quotes around blob:
	}
	w.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
}
