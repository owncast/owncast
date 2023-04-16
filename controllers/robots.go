package controllers

import (
	"net/http"
	"strings"

	"github.com/owncast/owncast/core/data"
)

// GetRobotsDotTxt returns the contents of our robots.txt.
func GetRobotsDotTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	contents := []string{
		"User-agent: *",
		"Disallow: /admin",
		"Disallow: /api",
	}

	if data.GetDisableSearchIndexing() {
		contents = append(contents, "Disallow: /")
	}

	txt := []byte(strings.Join(contents, "\n"))

	if _, err := w.Write(txt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
