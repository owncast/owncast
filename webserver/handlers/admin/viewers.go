package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// GetViewersOverTime will return the number of viewers at points in time.
func (a *Admin) GetViewersOverTime(w http.ResponseWriter, r *http.Request) {
	windowStartAtStr := r.URL.Query().Get("windowStart")
	windowStartAtUnix, err := strconv.Atoi(windowStartAtStr)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	windowStartAt := time.Unix(int64(windowStartAtUnix), 0)
	windowEnd := time.Now()

	viewersOverTime := a.metrics.GetViewersOverTime(windowStartAt, windowEnd)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(viewersOverTime)
	if err != nil {
		log.Errorln(err)
	}
}

// GetActiveViewers returns currently connected clients.
func (a *Admin) GetActiveViewers(w http.ResponseWriter, r *http.Request) {
	c := a.stream.GetActiveViewers()
	viewers := make([]models.Viewer, 0, len(c))
	for _, v := range c {
		viewers = append(viewers, *v)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(viewers); err != nil {
		webutils.InternalErrorHandler(w, err)
	}
}

// ExternalGetActiveViewers returns currently connected clients.
func (a *Admin) ExternalGetActiveViewers(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	a.GetConnectedChatClients(w, r)
}
