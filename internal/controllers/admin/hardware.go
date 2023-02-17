package admin

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/metrics"
)

// GetHardwareStats will return hardware utilization over time.
func GetHardwareStats(w http.ResponseWriter, r *http.Request) {
	m := metrics.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Errorln(err)
	}
}
