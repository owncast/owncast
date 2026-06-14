package admin

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/metrics"
)

// GetVideoPlaybackMetrics returns video playback metrics.
func (a *Admin) GetVideoPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
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
		Representation       int                        `json:"representation"`
	}

	availableBitrates := []int{}
	var segmentLength int
	if broadcast := a.stream.GetCurrentBroadcast(); broadcast != nil {
		segmentLength = broadcast.LatencyLevel.SecondsPerSegment
		for _, variants := range broadcast.OutputSettings {
			availableBitrates = append(availableBitrates, variants.VideoBitrate)
		}
	} else {
		segmentLength = a.configRepository.GetStreamLatencyLevel().SecondsPerSegment
		for _, variants := range a.configRepository.GetStreamOutputVariants() {
			availableBitrates = append(availableBitrates, variants.VideoBitrate)
		}
	}

	errors := a.metrics.GetPlaybackErrorCountOverTime()
	medianLatency := a.metrics.GetMedianLatencyOverTime()
	minimumLatency := a.metrics.GetMinimumLatencyOverTime()
	maximumLatency := a.metrics.GetMaximumLatencyOverTime()

	medianDurations := a.metrics.GetMedianDownloadDurationsOverTime()
	maximumDurations := a.metrics.GetMaximumDownloadDurationsOverTime()
	minimumDurations := a.metrics.GetMinimumDownloadDurationsOverTime()

	minPlayerBitrate := a.metrics.GetSlowestDownloadRateOverTime()
	medianPlayerBitrate := a.metrics.GetMedianDownloadRateOverTime()
	maxPlayerBitrate := a.metrics.GetMaxDownloadRateOverTime()
	qualityVariantChanges := a.metrics.GetQualityVariantChangesOverTime()

	representation := a.metrics.GetPlaybackMetricsRepresentation()

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
		Representation:          representation,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Errorln(err)
	}
}
