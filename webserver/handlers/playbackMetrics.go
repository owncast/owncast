package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// ReportPlaybackMetrics will accept playback metrics from a client and save
// them for future video health reporting.
func (h *Handlers) ReportPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request generated.PlaybackMetrics
	if err := decoder.Decode(&request); err != nil {
		log.Errorln("error decoding playback metrics payload:", err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if request.Errors == nil {
		webutils.WriteSimpleResponse(w, false, "errors field is required")
		return
	}

	if request.QualityVariantChanges == nil {
		webutils.WriteSimpleResponse(w, false, "qualityVariantChanges field is required")
		return
	}

	// Only a valid report earns registration and server-side observation
	// suppression.
	clientID := utils.GenerateClientIDFromRequest(r)
	h.metrics.RegisterSelfReportingClient(clientID)

	// This endpoint predates CMCD, so the client has no session identity
	// of its own: it is both the playback client and the viewer.
	report := metrics.PlaybackReport{
		ClientID:                 clientID,
		ViewerID:                 clientID,
		ErrorCount:               *request.Errors,
		HasErrorCount:            true,
		QualityVariantChanges:    *request.QualityVariantChanges,
		HasQualityVariantChanges: true,
		ClientReported:           true,
	}

	if request.Bandwidth != nil {
		report.BandwidthKbps = *request.Bandwidth
	}
	if request.Latency != nil {
		report.LatencySeconds = *request.Latency
	}
	if request.DownloadDuration != nil {
		report.DownloadSeconds = *request.DownloadDuration
	}

	h.metrics.RegisterPlaybackReport(report)
}
