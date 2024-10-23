package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
	log "github.com/sirupsen/logrus"
)

// ReportPlaybackMetrics will accept playback metrics from a client and save
// them for future video health reporting.
func ReportPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
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

	metrics.RegisterPlaybackErrorCount(clientID, *request.Errors)
	if *request.Bandwidth != 0.0 {
		metrics.RegisterPlayerBandwidth(clientID, *request.Bandwidth)
	}

	if *request.Latency != 0.0 {
		metrics.RegisterPlayerLatency(clientID, *request.Latency)
	}

	if *request.DownloadDuration != 0.0 {
		metrics.RegisterPlayerSegmentDownloadDuration(clientID, *request.DownloadDuration)
	}

	metrics.RegisterQualityVariantChangesCount(clientID, *request.QualityVariantChanges)
}
