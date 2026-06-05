package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

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
	var request generated.ReportPlaybackMetricsJSONRequestBody
	if err := decoder.Decode(&request); err != nil {
		log.Errorln("error decoding playback metrics payload:", err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	clientID := utils.GenerateClientIDFromRequest(r)

	if request.Errors == nil {
		webutils.WriteSimpleResponse(w, false, "errors field is required")
		return
	}
	h.metrics.RegisterPlaybackErrorCount(clientID, *request.Errors)

	if request.Bandwidth != nil && *request.Bandwidth != 0.0 {
		h.metrics.RegisterPlayerBandwidth(clientID, *request.Bandwidth)
	}

	if request.Latency != nil && *request.Latency != 0.0 {
		h.metrics.RegisterPlayerLatency(clientID, *request.Latency)
	}

	if request.DownloadDuration != nil && *request.DownloadDuration != 0.0 {
		h.metrics.RegisterPlayerSegmentDownloadDuration(clientID, *request.DownloadDuration)
	}

	if request.QualityVariantChanges == nil {
		webutils.WriteSimpleResponse(w, false, "qualityVariantChanges field is required")
		return
	}
	h.metrics.RegisterQualityVariantChangesCount(clientID, *request.QualityVariantChanges)
}
