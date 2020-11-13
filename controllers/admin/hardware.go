package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/metrics"
	log "github.com/sirupsen/logrus"
)

// GetHardwareStats will return hardware utilization over time.
func GetHardwareStats(w http.ResponseWriter, r *http.Request) {
	m := metrics.Metrics

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Errorln(err)
	}
}
