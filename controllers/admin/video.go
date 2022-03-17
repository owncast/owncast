package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/metrics"
	log "github.com/sirupsen/logrus"
)

// GetVideoPlaybackMetrics returns video playback metrics.
func GetVideoPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Errors                  []metrics.TimestampedValue `json:"errors"`
		QualityVariantChanges   []metrics.TimestampedValue `json:"qualityVariantChanges"`
		Latency                 []metrics.TimestampedValue `json:"latency"`
		SegmentDownloadDuration []metrics.TimestampedValue `json:"segmentDownloadDuration"`
		SlowestDownloadRate     []metrics.TimestampedValue `json:"minPlayerBitrate"`
		AvailableBitrates       []int                      `json:"availableBitrates"`
		SegmentLength           int                        `json:"segmentLength"`
	}

	availableBitrates := []int{}
	var segmentLength int
	if core.GetCurrentBroadcast() != nil {
		segmentLength = core.GetCurrentBroadcast().LatencyLevel.SecondsPerSegment
		for _, variants := range core.GetCurrentBroadcast().OutputSettings {
			availableBitrates = append(availableBitrates, variants.VideoBitrate)
		}
	} else {
		segmentLength = data.GetStreamLatencyLevel().SecondsPerSegment
		for _, variants := range data.GetStreamOutputVariants() {
			availableBitrates = append(availableBitrates, variants.VideoBitrate)
		}
	}

	errors := metrics.GetPlaybackErrorCountOverTime()
	latency := metrics.GetLatencyOverTime()
	durations := metrics.GetDownloadDurationsOverTime()
	minPlayerBitrate := metrics.GetSlowestDownloadRateOverTime()
	qualityVariantChanges := metrics.GetQualityVariantChangesOverTime()

	resp := response{
		AvailableBitrates:       availableBitrates,
		Errors:                  errors,
		Latency:                 latency,
		SegmentLength:           segmentLength,
		SegmentDownloadDuration: durations,
		SlowestDownloadRate:     minPlayerBitrate,
		QualityVariantChanges:   qualityVariantChanges,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Errorln(err)
	}
}
