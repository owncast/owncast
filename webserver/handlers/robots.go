package handlers

import (
	"net/http"
	"strings"
)

// GetRobotsDotTxt returns the contents of our robots.txt.
func (h *Handlers) GetRobotsDotTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	contents := []string{
		"User-agent: *",
		"Disallow: /admin",
		"Disallow: /api",
	}

	if configRepository.GetDisableSearchIndexing() {
		contents = append(contents, "Disallow: /")
	}

	txt := []byte(strings.Join(contents, "\n"))

	if _, err := w.Write(txt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
