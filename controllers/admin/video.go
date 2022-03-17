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
		Errors                []metrics.TimestampedValue `json:"errors"`
		QualityVariantChanges []metrics.TimestampedValue `json:"qualityVariantChanges"`

		HighestLatency []metrics.TimestampedValue `json:"highestLatency"`
		MedianLatency  []metrics.TimestampedValue `json:"medianLatency"`
		LowestLatency  []metrics.TimestampedValue `json:"lowestLatency"`

		MedianDownloadDuration  []metrics.TimestampedValue `json:"medianSegmentDownloadDuration"`
		MaximumDownloadDuration []metrics.TimestampedValue `json:"maximumSegmentDownloadDuration"`
		MinimumDownloadDuration []metrics.TimestampedValue `json:"minimumSegmentDownloadDuration"`

		SlowestDownloadRate  []metrics.TimestampedValue `json:"minPlayerBitrate"`
		MedianDownloadRate   []metrics.TimestampedValue `json:"medianPlayerBitrate"`
		HighestDownloadRater []metrics.TimestampedValue `json:"maxPlayerBitrate"`
		AvailableBitrates    []int                      `json:"availableBitrates"`
		SegmentLength        int                        `json:"segmentLength"`
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
	medianLatency := metrics.GetMedianLatencyOverTime()
	minimumLatency := metrics.GetMinimumLatencyOverTime()
	maximumLatency := metrics.GetMaximumLatencyOverTime()

	medianDurations := metrics.GetMedianDownloadDurationsOverTime()
	maximumDurations := metrics.GetMaximumDownloadDurationsOverTime()
	minimumDurations := metrics.GetMinimumDownloadDurationsOverTime()

	minPlayerBitrate := metrics.GetSlowestDownloadRateOverTime()
	medianPlayerBitrate := metrics.GetMedianDownloadRateOverTime()
	maxPlayerBitrate := metrics.GetMaxDownloadRateOverTime()
	qualityVariantChanges := metrics.GetQualityVariantChangesOverTime()

	resp := response{
		AvailableBitrates:       availableBitrates,
		Errors:                  errors,
		MedianLatency:           medianLatency,
		HighestLatency:          maximumLatency,
		LowestLatency:           minimumLatency,
		SegmentLength:           segmentLength,
		MedianDownloadDuration:  medianDurations,
		MaximumDownloadDuration: maximumDurations,
		MinimumDownloadDuration: minimumDurations,
		SlowestDownloadRate:     minPlayerBitrate,
		MedianDownloadRate:      medianPlayerBitrate,
		HighestDownloadRater:    maxPlayerBitrate,
		QualityVariantChanges:   qualityVariantChanges,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Errorln(err)
	}
}
