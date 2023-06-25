package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/metrics"
	"github.com/owncast/owncast/webserver/responses"
	log "github.com/sirupsen/logrus"
)

// GetViewersOverTime will return the number of viewers at points in time.
func (h *Handlers) GetViewersOverTime(w http.ResponseWriter, r *http.Request) {
	windowStartAtStr := r.URL.Query().Get("windowStart")
	windowStartAtUnix, err := strconv.Atoi(windowStartAtStr)
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	windowStartAt := time.Unix(int64(windowStartAtUnix), 0)
	windowEnd := time.Now()

	viewersOverTime := metrics.GetViewersOverTime(windowStartAt, windowEnd)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(viewersOverTime)
	if err != nil {
		log.Errorln(err)
	}
}

// GetActiveViewers returns currently connected clients.
func (h *Handlers) GetActiveViewers(w http.ResponseWriter, r *http.Request) {
	c := core.GetActiveViewers()
	viewers := make([]models.Viewer, 0, len(c))
	for _, v := range c {
		viewers = append(viewers, *v)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(viewers); err != nil {
		responses.InternalErrorHandler(w, err)
	}
}

// ExternalGetActiveViewers returns currently connected clients.
func (h *Handlers) ExternalGetActiveViewers(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	h.GetConnectedChatClients(w, r)
}
