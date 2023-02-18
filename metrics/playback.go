package metrics

import (
	"math"
	"time"

	"github.com/owncast/owncast/utils"
)

// Playback error counts reported since the last time we collected s.Metrics.
var (
	windowedErrorCounts           = map[string]float64{}
	windowedQualityVariantChanges = map[string]float64{}
	windowedBandwidths            = map[string]float64{}
	windowedLatencies             = map[string]float64{}
	windowedDownloadDurations     = map[string]float64{}
)

func (s *Service) handlePlaybackPolling() {
	s.Metrics.mux.Lock()
	defer s.Metrics.mux.Unlock()

	// Make sure this is fired first before all the values get cleared below.
	if _getStatus().Online {
		s.generateStreamHealthOverview()
	}

	s.collectPlaybackErrorCount()
	s.collectLatencyValues()
	s.collectSegmentDownloadDuration()
	s.collectLowestBandwidth()
	s.collectQualityVariantChanges()
}

// RegisterPlaybackErrorCount will add to the windowed playback error count.
func (s *Service) RegisterPlaybackErrorCount(clientID string, count float64) {
	s.Metrics.mux.Lock()
	defer s.Metrics.mux.Unlock()

	windowedErrorCounts[clientID] = count
	// windowedErrorCounts = append(windowedErrorCounts, count)
}

// RegisterQualityVariantChangesCount will add to the windowed quality variant
// change count.
func (s *Service) RegisterQualityVariantChangesCount(clientID string, count float64) {
	s.Metrics.mux.Lock()
	defer s.Metrics.mux.Unlock()

	windowedQualityVariantChanges[clientID] = count
}

// RegisterPlayerBandwidth will add to the windowed playback bandwidth.
func (s *Service) RegisterPlayerBandwidth(clientID string, kbps float64) {
	s.Metrics.mux.Lock()
	defer s.Metrics.mux.Unlock()

	windowedBandwidths[clientID] = kbps
}

// RegisterPlayerLatency will add to the windowed player latency values.
func (s *Service) RegisterPlayerLatency(clientID string, seconds float64) {
	s.Metrics.mux.Lock()
	defer s.Metrics.mux.Unlock()

	windowedLatencies[clientID] = seconds
}

// RegisterPlayerSegmentDownloadDuration will add to the windowed player segment
// download duration values.
func (s *Service) RegisterPlayerSegmentDownloadDuration(clientID string, seconds float64) {
	s.Metrics.mux.Lock()
	defer s.Metrics.mux.Unlock()

	windowedDownloadDurations[clientID] = seconds
}

// collectPlaybackErrorCount will take all of the error counts each individual
// player reported and average them into a single metric. This is done so
// one person with bad connectivity doesn't make it look like everything is
// horrible for everyone.
func (s *Service) collectPlaybackErrorCount() {
	valueSlice := utils.Float64MapToSlice(windowedErrorCounts)
	count := utils.Sum(valueSlice)
	windowedErrorCounts = map[string]float64{}

	s.Metrics.errorCount = append(s.Metrics.errorCount, TimestampedValue{
		Time:  time.Now(),
		Value: count,
	})

	if len(s.Metrics.errorCount) > maxCollectionValues {
		s.Metrics.errorCount = s.Metrics.errorCount[1:]
	}

	// Save to Prometheus collector.
	playbackErrorCount.Set(count)
}

func (s *Service) collectSegmentDownloadDuration() {
	median := 0.0
	max := 0.0
	min := 0.0

	valueSlice := utils.Float64MapToSlice(windowedDownloadDurations)

	if len(valueSlice) > 0 {
		median = utils.Median(valueSlice)
		min, max = utils.MinMax(valueSlice)
		windowedDownloadDurations = map[string]float64{}
	}

	s.Metrics.medianSegmentDownloadSeconds = append(s.Metrics.medianSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: median,
	})

	if len(s.Metrics.medianSegmentDownloadSeconds) > maxCollectionValues {
		s.Metrics.medianSegmentDownloadSeconds = s.Metrics.medianSegmentDownloadSeconds[1:]
	}

	s.Metrics.minimumSegmentDownloadSeconds = append(s.Metrics.minimumSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: min,
	})

	if len(s.Metrics.minimumSegmentDownloadSeconds) > maxCollectionValues {
		s.Metrics.minimumSegmentDownloadSeconds = s.Metrics.minimumSegmentDownloadSeconds[1:]
	}

	s.Metrics.maximumSegmentDownloadSeconds = append(s.Metrics.maximumSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: max,
	})

	if len(s.Metrics.maximumSegmentDownloadSeconds) > maxCollectionValues {
		s.Metrics.maximumSegmentDownloadSeconds = s.Metrics.maximumSegmentDownloadSeconds[1:]
	}
}

// GetMedianDownloadDurationsOverTime will return a window of durations errors over time.
func (s *Service) GetMedianDownloadDurationsOverTime() []TimestampedValue {
	return s.Metrics.medianSegmentDownloadSeconds
}

// GetMaximumDownloadDurationsOverTime will return a maximum durations errors over time.
func (s *Service) GetMaximumDownloadDurationsOverTime() []TimestampedValue {
	return s.Metrics.maximumSegmentDownloadSeconds
}

// GetMinimumDownloadDurationsOverTime will return a maximum durations errors over time.
func (s *Service) GetMinimumDownloadDurationsOverTime() []TimestampedValue {
	return s.Metrics.minimumSegmentDownloadSeconds
}

// GetPlaybackErrorCountOverTime will return a window of playback errors over time.
func (s *Service) GetPlaybackErrorCountOverTime() []TimestampedValue {
	return s.Metrics.errorCount
}

func (s *Service) collectLatencyValues() {
	median := 0.0
	min := 0.0
	max := 0.0

	valueSlice := utils.Float64MapToSlice(windowedLatencies)
	windowedLatencies = map[string]float64{}

	if len(valueSlice) > 0 {
		median = utils.Median(valueSlice)
		min, max = utils.MinMax(valueSlice)
		windowedLatencies = map[string]float64{}
	}

	s.Metrics.medianLatency = append(s.Metrics.medianLatency, TimestampedValue{
		Time:  time.Now(),
		Value: median,
	})

	if len(s.Metrics.medianLatency) > maxCollectionValues {
		s.Metrics.medianLatency = s.Metrics.medianLatency[1:]
	}

	s.Metrics.minimumLatency = append(s.Metrics.minimumLatency, TimestampedValue{
		Time:  time.Now(),
		Value: min,
	})

	if len(s.Metrics.minimumLatency) > maxCollectionValues {
		s.Metrics.minimumLatency = s.Metrics.minimumLatency[1:]
	}

	s.Metrics.maximumLatency = append(s.Metrics.maximumLatency, TimestampedValue{
		Time:  time.Now(),
		Value: max,
	})

	if len(s.Metrics.maximumLatency) > maxCollectionValues {
		s.Metrics.maximumLatency = s.Metrics.maximumLatency[1:]
	}
}

// GetMedianLatencyOverTime will return the median latency values over time.
func (s *Service) GetMedianLatencyOverTime() []TimestampedValue {
	if len(s.Metrics.medianLatency) == 0 {
		return []TimestampedValue{}
	}

	return s.Metrics.medianLatency
}

// GetMinimumLatencyOverTime will return the min latency values over time.
func (s *Service) GetMinimumLatencyOverTime() []TimestampedValue {
	if len(s.Metrics.minimumLatency) == 0 {
		return []TimestampedValue{}
	}

	return s.Metrics.minimumLatency
}

// GetMaximumLatencyOverTime will return the max latency values over time.
func (s *Service) GetMaximumLatencyOverTime() []TimestampedValue {
	if len(s.Metrics.maximumLatency) == 0 {
		return []TimestampedValue{}
	}

	return s.Metrics.maximumLatency
}

// collectLowestBandwidth will collect the bandwidth currently collected
// so we can report to the streamer the worst possible streaming condition
// being experienced.
func (s *Service) collectLowestBandwidth() {
	min := 0.0
	median := 0.0
	max := 0.0

	valueSlice := utils.Float64MapToSlice(windowedBandwidths)

	if len(windowedBandwidths) > 0 {
		min, max = utils.MinMax(valueSlice)
		min = math.Round(min)
		max = math.Round(max)
		median = utils.Median(valueSlice)
		windowedBandwidths = map[string]float64{}
	}

	s.Metrics.lowestBitrate = append(s.Metrics.lowestBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(min),
	})

	if len(s.Metrics.lowestBitrate) > maxCollectionValues {
		s.Metrics.lowestBitrate = s.Metrics.lowestBitrate[1:]
	}

	s.Metrics.medianBitrate = append(s.Metrics.medianBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(median),
	})

	if len(s.Metrics.medianBitrate) > maxCollectionValues {
		s.Metrics.medianBitrate = s.Metrics.medianBitrate[1:]
	}

	s.Metrics.highestBitrate = append(s.Metrics.highestBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(max),
	})

	if len(s.Metrics.highestBitrate) > maxCollectionValues {
		s.Metrics.highestBitrate = s.Metrics.highestBitrate[1:]
	}
}

// GetSlowestDownloadRateOverTime will return the collected lowest bandwidth values
// over time.
func (s *Service) GetSlowestDownloadRateOverTime() []TimestampedValue {
	if len(s.Metrics.lowestBitrate) == 0 {
		return []TimestampedValue{}
	}

	return s.Metrics.lowestBitrate
}

// GetMedianDownloadRateOverTime will return the collected median bandwidth values.
func (s *Service) GetMedianDownloadRateOverTime() []TimestampedValue {
	if len(s.Metrics.medianBitrate) == 0 {
		return []TimestampedValue{}
	}
	return s.Metrics.medianBitrate
}

// GetMaximumDownloadRateOverTime will return the collected maximum bandwidth values.
func (s *Service) GetMaximumDownloadRateOverTime() []TimestampedValue {
	if len(s.Metrics.maximumLatency) == 0 {
		return []TimestampedValue{}
	}
	return s.Metrics.maximumLatency
}

// GetMinimumDownloadRateOverTime will return the collected minimum bandwidth values.
func (s *Service) GetMinimumDownloadRateOverTime() []TimestampedValue {
	if len(s.Metrics.minimumLatency) == 0 {
		return []TimestampedValue{}
	}
	return s.Metrics.minimumLatency
}

// GetMaxDownloadRateOverTime will return the collected highest bandwidth values.
func (s *Service) GetMaxDownloadRateOverTime() []TimestampedValue {
	if len(s.Metrics.highestBitrate) == 0 {
		return []TimestampedValue{}
	}
	return s.Metrics.highestBitrate
}

func (s *Service) collectQualityVariantChanges() {
	valueSlice := utils.Float64MapToSlice(windowedQualityVariantChanges)
	count := utils.Sum(valueSlice)
	windowedQualityVariantChanges = map[string]float64{}

	s.Metrics.qualityVariantChanges = append(s.Metrics.qualityVariantChanges, TimestampedValue{
		Time:  time.Now(),
		Value: count,
	})
}

// GetQualityVariantChangesOverTime will return the collected quality variant
// changes.
func (s *Service) GetQualityVariantChangesOverTime() []TimestampedValue {
	return s.Metrics.qualityVariantChanges
}

// GetPlaybackMetricsRepresentation returns what percentage of all known players
// the metrics represent.
func (s *Service) GetPlaybackMetricsRepresentation() int {
	totalPlayerCount := len(s.Core.GetActiveViewers())
	representation := utils.IntPercentage(len(windowedBandwidths), totalPlayerCount)
	return representation
}
