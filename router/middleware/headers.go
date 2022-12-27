package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// SetHeaders will set our global headers for web resources.
func SetHeaders(w http.ResponseWriter, nonce string) {
	// Content security policy
	csp := []string{
		fmt.Sprintf("script-src '%s' 'self'", nonce),
		"worker-src 'self' blob:", // No single quotes around blob:
	}
	w.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
}
