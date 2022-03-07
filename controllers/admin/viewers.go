package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
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
	viewers := []models.Viewer{}
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
