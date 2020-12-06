package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/owncast/owncast/logging"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// GetLogs will return all logs.
func GetLogs(w http.ResponseWriter, r *http.Request) {
	logs := logging.Logger.AllEntries()
	response := make([]logsResponse, 0)

	for i := 0; i < len(logs); i++ {
		response = append(response, fromEntry(logs[i]))
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}

// GetWarnings will return only warning and error logs.
func GetWarnings(w http.ResponseWriter, r *http.Request) {
	logs := logging.Logger.WarningEntries()
	response := make([]logsResponse, 0)

	for i := 0; i < len(logs); i++ {
		logEntry := logs[i]
		if logEntry != nil {
			response = append(response, fromEntry(logEntry))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}

type logsResponse struct {
	Message string    `json:"message"`
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
}

func fromEntry(e *logrus.Entry) logsResponse {
	return logsResponse{
		Message: e.Message,
		Level:   e.Level.String(),
		Time:    e.Time,
	}
}
