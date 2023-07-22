package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/metrics"
	"github.com/owncast/owncast/services/status"
	log "github.com/sirupsen/logrus"
)

// GetVideoPlaybackMetrics returns video playback metrics.
func (h *Handlers) GetVideoPlaybackMetrics(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Errors                []models.TimestampedValue `json:"errors"`
		QualityVariantChanges []models.TimestampedValue `json:"qualityVariantChanges"`

		HighestLatency []models.TimestampedValue `json:"highestLatency"`
		MedianLatency  []models.TimestampedValue `json:"medianLatency"`
		LowestLatency  []models.TimestampedValue `json:"lowestLatency"`

		MedianDownloadDuration  []models.TimestampedValue `json:"medianSegmentDownloadDuration"`
		MaximumDownloadDuration []models.TimestampedValue `json:"maximumSegmentDownloadDuration"`
		MinimumDownloadDuration []models.TimestampedValue `json:"minimumSegmentDownloadDuration"`

		SlowestDownloadRate  []models.TimestampedValue `json:"minPlayerBitrate"`
		MedianDownloadRate   []models.TimestampedValue `json:"medianPlayerBitrate"`
		HighestDownloadRater []models.TimestampedValue `json:"maxPlayerBitrate"`
		AvailableBitrates    []int                     `json:"availableBitrates"`
		SegmentLength        int                       `json:"segmentLength"`
		Representation       int                       `json:"representation"`
	}

	s := status.Get()

	availableBitrates := []int{}
	var segmentLength int
	if s.GetCurrentBroadcast() != nil {
		segmentLength = s.GetCurrentBroadcast().LatencyLevel.SecondsPerSegment
		for _, variants := range s.GetCurrentBroadcast().OutputSettings {
			availableBitrates = append(availableBitrates, variants.VideoBitrate)
		}
	} else {
		segmentLength = configRepository.GetStreamLatencyLevel().SecondsPerSegment
		for _, variants := range configRepository.GetStreamOutputVariants() {
			availableBitrates = append(availableBitrates, variants.VideoBitrate)
		}
	}

	m := metrics.Get()

	errors := m.GetPlaybackErrorCountOverTime()
	medianLatency := m.GetMedianLatencyOverTime()
	minimumLatency := m.GetMinimumLatencyOverTime()
	maximumLatency := m.GetMaximumLatencyOverTime()

	medianDurations := m.GetMedianDownloadDurationsOverTime()
	maximumDurations := m.GetMaximumDownloadDurationsOverTime()
	minimumDurations := m.GetMinimumDownloadDurationsOverTime()

	minPlayerBitrate := m.GetSlowestDownloadRateOverTime()
	medianPlayerBitrate := m.GetMedianDownloadRateOverTime()
	maxPlayerBitrate := m.GetMaxDownloadRateOverTime()
	qualityVariantChanges := m.GetQualityVariantChangesOverTime()

	representation := m.GetPlaybackMetricsRepresentation()

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
