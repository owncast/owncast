package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/router/middleware"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// GetStatus gets the status of the server.
func (h *Handlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	response := h.getStatusResponse()

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		webutils.InternalErrorHandler(w, err)
	}
}

func (h *Handlers) getStatusResponse() webStatusResponse {
	status := h.stream.GetStatus()
	response := webStatusResponse{
		Online:             status.Online,
		ServerTime:         time.Now(),
		LastConnectTime:    status.LastConnectTime,
		LastDisconnectTime: status.LastDisconnectTime,
		VersionNumber:      status.VersionNumber,
		StreamTitle:        status.StreamTitle,
	}
	if !h.configRepository.GetHideViewerCount() {
		response.ViewerCount = status.ViewerCount
	}

	// The live broadcast's id is only useful to a viewer who can clip it, so
	// it is published only while clipping is actually available.
	if h.configRepository.GetReplayFeaturesEnabled() && h.configRepository.GetClipsEnabled() {
		response.StreamID = status.StreamID
	}

	return response
}

type webStatusResponse struct {
	ServerTime         time.Time       `json:"serverTime"`
	LastConnectTime    *utils.NullTime `json:"lastConnectTime"`
	LastDisconnectTime *utils.NullTime `json:"lastDisconnectTime"`

	VersionNumber string `json:"versionNumber"`
	StreamTitle   string `json:"streamTitle"`
	// StreamID identifies the live broadcast so a viewer can clip it. Only
	// present while a stream is live and clipping is enabled.
	StreamID    string `json:"streamId,omitempty"`
	ViewerCount int    `json:"viewerCount,omitempty"`
	Online      bool   `json:"online"`
}
