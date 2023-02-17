package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/internal/controllers"
	"github.com/owncast/owncast/internal/core"
	"github.com/owncast/owncast/internal/core/user"
	"github.com/owncast/owncast/internal/metrics"
	"github.com/owncast/owncast/internal/models"
)

// GetViewersOverTime will return the number of viewers at points in time.
func GetViewersOverTime(w http.ResponseWriter, r *http.Request) {
	windowStartAtStr := r.URL.Query().Get("windowStart")
	windowStartAtUnix, err := strconv.Atoi(windowStartAtStr)
	if err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
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
func GetActiveViewers(w http.ResponseWriter, r *http.Request) {
	c := core.GetActiveViewers()
	viewers := make([]models.Viewer, 0, len(c))
	for _, v := range c {
		viewers = append(viewers, *v)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(viewers); err != nil {
		controllers.InternalErrorHandler(w, err)
	}
}

// ExternalGetActiveViewers returns currently connected clients.
func ExternalGetActiveViewers(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	GetConnectedChatClients(w, r)
}
