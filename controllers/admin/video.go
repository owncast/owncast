package admin

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/metrics"
)

// GetVideoPlaybackMetrics returns video playback router.Metrics.
func (c *Controller) GetVideoPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
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
	if c.Core.GetCurrentBroadcast() != nil {
		segmentLength = c.Core.GetCurrentBroadcast().LatencyLevel.SecondsPerSegment
		for _, variants := range c.Core.GetCurrentBroadcast().OutputSettings {
			availableBitrates = append(availableBitrates, variants.VideoBitrate)
		}
	} else {
		segmentLength = c.Data.GetStreamLatencyLevel().SecondsPerSegment
		for _, variants := range c.Data.GetStreamOutputVariants() {
			availableBitrates = append(availableBitrates, variants.VideoBitrate)
		}
	}

	errors := c.Metrics.GetPlaybackErrorCountOverTime()
	medianLatency := c.Metrics.GetMedianLatencyOverTime()
	minimumLatency := c.Metrics.GetMinimumLatencyOverTime()
	maximumLatency := c.Metrics.GetMaximumLatencyOverTime()

	medianDurations := c.Metrics.GetMedianDownloadDurationsOverTime()
	maximumDurations := c.Metrics.GetMaximumDownloadDurationsOverTime()
	minimumDurations := c.Metrics.GetMinimumDownloadDurationsOverTime()

	minPlayerBitrate := c.Metrics.GetSlowestDownloadRateOverTime()
	medianPlayerBitrate := c.Metrics.GetMedianDownloadRateOverTime()
	maxPlayerBitrate := c.Metrics.GetMaxDownloadRateOverTime()
	qualityVariantChanges := c.Metrics.GetQualityVariantChangesOverTime()

	representation := c.Metrics.GetPlaybackMetricsRepresentation()

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
