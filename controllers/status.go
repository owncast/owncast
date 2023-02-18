package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/owncast/owncast/app/middleware"
	"github.com/owncast/owncast/utils"
)

// GetStatus gets the status of the server.
func (s *Service) GetStatus(w http.ResponseWriter, r *http.Request) {
	response := s.getStatusResponse()

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.InternalErrorHandler(w, err)
	}
}

func (s *Service) getStatusResponse() webStatusResponse {
	status := s.Core.GetStatus()
	response := webStatusResponse{
		Online:             status.Online,
		ServerTime:         time.Now(),
		LastConnectTime:    status.LastConnectTime,
		LastDisconnectTime: status.LastDisconnectTime,
		VersionNumber:      status.VersionNumber,
		StreamTitle:        status.StreamTitle,
	}
	if !s.Data.GetHideViewerCount() {
		response.ViewerCount = status.ViewerCount
	}
	return response
}

type webStatusResponse struct {
	Online             bool            `json:"online"`
	ViewerCount        int             `json:"viewerCount,omitempty"`
	ServerTime         time.Time       `json:"serverTime"`
	LastConnectTime    *utils.NullTime `json:"lastConnectTime"`
	LastDisconnectTime *utils.NullTime `json:"lastDisconnectTime"`

	VersionNumber string `json:"versionNumber"`
	StreamTitle   string `json:"streamTitle"`
}
