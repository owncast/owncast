package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/metrics"
)

// GetViewersOverTime will return the number of viewers at points in time.
func GetViewersOverTime(w http.ResponseWriter, r *http.Request) {
	viewersOverTime := metrics.Metrics.Viewers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(viewersOverTime)
}
