package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/models"
)

// GetViewersOverTime will return the number of viewers at points in time.
func (c *Controller) GetViewersOverTime(w http.ResponseWriter, r *http.Request) {
	windowStartAtStr := r.URL.Query().Get("windowStart")
	windowStartAtUnix, err := strconv.Atoi(windowStartAtStr)
	if err != nil {
		c.Service.WriteSimpleResponse(w, false, err.Error())
		return
	}

	windowStartAt := time.Unix(int64(windowStartAtUnix), 0)
	windowEnd := time.Now()

	viewersOverTime := c.Metrics.GetViewersOverTime(windowStartAt, windowEnd)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(viewersOverTime)
	if err != nil {
		log.Errorln(err)
	}
}

// GetActiveViewers returns currently connected clients.
func (c *Controller) GetActiveViewers(w http.ResponseWriter, r *http.Request) {
	viewersMap := c.Core.GetActiveViewers()
	viewers := make([]models.Viewer, 0, len(viewersMap))
	for _, v := range viewersMap {
		viewers = append(viewers, *v)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(viewers); err != nil {
		c.Service.InternalErrorHandler(w, err)
	}
}

// ExternalGetActiveViewers returns currently connected clients.
func (c *Controller) ExternalGetActiveViewers(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	c.GetConnectedChatClients(w, r)
}
