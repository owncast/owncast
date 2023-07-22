package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/services/metrics"
	log "github.com/sirupsen/logrus"
)

// GetHardwareStats will return hardware utilization over time.
func (h *Handlers) GetHardwareStats(w http.ResponseWriter, r *http.Request) {
	m := metrics.Get()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		log.Errorln(err)
	}
}
