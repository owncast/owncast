package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/logging"
)

// GetLogs will return all logs.
func (c *Controller) GetLogs(w http.ResponseWriter, r *http.Request) {
	logs := logging.Logger.AllEntries()
	response := make([]logsResponse, 0)

	for i := 0; i < len(logs); i++ {
		response = append(response, c.fromEntry(logs[i]))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorln(err)
	}
}

// GetWarnings will return only warning and error logs.
func (c *Controller) GetWarnings(w http.ResponseWriter, r *http.Request) {
	logs := logging.Logger.WarningEntries()
	response := make([]logsResponse, 0)

	for i := 0; i < len(logs); i++ {
		logEntry := logs[i]
		if logEntry != nil {
			response = append(response, c.fromEntry(logEntry))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorln(err)
	}
}

type logsResponse struct {
	Message string    `json:"message"`
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
}

func (c *Controller) fromEntry(e *logrus.Entry) logsResponse {
	return logsResponse{
		Message: e.Message,
		Level:   e.Level.String(),
		Time:    e.Time,
	}
}
