package admin

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetHardwareStats will return hardware utilization over time.
func (a *Admin) GetHardwareStats(w http.ResponseWriter, r *http.Request) {
	m := a.metrics.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Errorln(err)
	}
}
