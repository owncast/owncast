package admin

import (
	"encoding/json"
	"net/http"

	"github.com/gabek/owncast/metrics"
)

// GetHardwareStats will return hardware utilization over time
func GetHardwareStats(w http.ResponseWriter, r *http.Request) {
	metrics := metrics.Metrics

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
