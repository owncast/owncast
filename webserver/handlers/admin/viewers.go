package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/metrics"
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
		viewers = append(viewers, v)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(viewers); err != nil {
		webutils.InternalErrorHandler(w, err)
	}
}

// GetPlaybackClients returns one entry per client currently playing back
// video, joined with the most recent playback measurements for that
// client. A viewer whose player reports nothing and whose segment
// transfers weren't measurable is still listed, with no playback health,
// because "we can't see this viewer" is itself something a streamer needs
// to know.
func (a *Admin) GetPlaybackClients(w http.ResponseWriter, r *http.Request) {
	viewers := a.stream.GetActiveViewers()
	states := a.metrics.GetActivePlaybackClients()

	clients := make([]models.PlaybackClient, 0, len(viewers))
	measured := make(map[string]bool, len(viewers))

	for _, state := range states {
		viewer, ok := viewers[state.ViewerID]
		if !ok {
			// The client stopped requesting video recently enough that
			// its measurements are still fresh but its viewer record has
			// already been pruned. It isn't watching now.
			continue
		}
		measured[state.ViewerID] = true
		clients = append(clients, newPlaybackClient(state.ClientID, viewer, playbackHealth(state)))
	}

	for id, viewer := range viewers {
		if !measured[id] {
			clients = append(clients, newPlaybackClient(id, viewer, nil))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(clients); err != nil {
		webutils.InternalErrorHandler(w, err)
	}
}

func newPlaybackClient(clientID string, viewer models.Viewer, health *models.PlaybackClientHealth) models.PlaybackClient {
	return models.PlaybackClient{
		ClientID:  clientID,
		ViewerID:  viewer.ClientID,
		FirstSeen: viewer.FirstSeen,
		Geo:       viewer.Geo,
		UserAgent: viewer.UserAgent,
		IPAddress: viewer.IPAddress,
		Playback:  health,
	}
}

func playbackHealth(state metrics.ClientPlaybackState) *models.PlaybackClientHealth {
	source := models.PlaybackSourceServer
	if state.IsSelfReporting() {
		source = models.PlaybackSourceClient
	}

	health := &models.PlaybackClientHealth{
		Source:            source,
		LastUpdate:        state.LastUpdate,
		PlayerState:       state.PlayerState,
		MeasurementStatus: state.ServerMeasurementStatus,
		// A value the client stopped reporting goes back to unknown
		// instead of standing in for the current one: our own player
		// omits live latency once it exceeds 100 seconds, which is
		// exactly when its last good number would mislead the most.
		BandwidthKbps:   state.Measurement(state.BandwidthKbps, state.BandwidthAt),
		LatencySeconds:  state.Measurement(state.LatencySeconds, state.LatencyAt),
		DownloadSeconds: state.Measurement(state.DownloadSeconds, state.DownloadAt),
		BitrateKbps:     state.Measurement(state.BitrateKbps, state.BitrateAt),
	}

	// The counters are session totals rather than point measurements, so
	// they only go unknown when the client never reported them at all.
	if state.HasErrorCount {
		health.ErrorCount = &state.ErrorCount
	}
	if state.HasQualityChanges {
		health.QualityVariantChanges = &state.QualityVariantChanges
	}

	return health
}

// ExternalGetActiveViewers returns currently connected clients.
func (a *Admin) ExternalGetActiveViewers(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	a.GetConnectedChatClients(w, r)
}
