package metrics

import (
	"math"
	"time"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/status"
	"github.com/owncast/owncast/utils"
)

func (m *Metrics) handlePlaybackPolling() {
	m.metrics.m.Lock()
	defer m.metrics.m.Unlock()

	s := status.Get()

	// Make sure this is fired first before all the values get cleared below.
	if s.Online {
		m.generateStreamHealthOverview()
	}

	m.collectPlaybackErrorCount()
	m.collectLatencyValues()
	m.collectSegmentDownloadDuration()
	m.collectLowestBandwidth()
	m.collectQualityVariantChanges()
}

// RegisterPlaybackErrorCount will add to the windowed playback error count.
func (m *Metrics) RegisterPlaybackErrorCount(clientID string, count float64) {
	m.metrics.m.Lock()
	defer m.metrics.m.Unlock()
	m.windowedErrorCounts[clientID] = count
	// windowedErrorCounts = append(windowedErrorCounts, count)
}

// RegisterQualityVariantChangesCount will add to the windowed quality variant
// change count.
func (m *Metrics) RegisterQualityVariantChangesCount(clientID string, count float64) {
	m.metrics.m.Lock()
	defer m.metrics.m.Unlock()
	m.windowedQualityVariantChanges[clientID] = count
}

// RegisterPlayerBandwidth will add to the windowed playback bandwidth.
func (m *Metrics) RegisterPlayerBandwidth(clientID string, kbps float64) {
	m.metrics.m.Lock()
	defer m.metrics.m.Unlock()
	m.windowedBandwidths[clientID] = kbps
}

// RegisterPlayerLatency will add to the windowed player latency values.
func (m *Metrics) RegisterPlayerLatency(clientID string, seconds float64) {
	m.metrics.m.Lock()
	defer m.metrics.m.Unlock()
	m.windowedLatencies[clientID] = seconds
}

// RegisterPlayerSegmentDownloadDuration will add to the windowed player segment
// download duration values.
func (m *Metrics) RegisterPlayerSegmentDownloadDuration(clientID string, seconds float64) {
	m.metrics.m.Lock()
	defer m.metrics.m.Unlock()
	m.windowedDownloadDurations[clientID] = seconds
}

// collectPlaybackErrorCount will take all of the error counts each individual
// player reported and average them into a single metric. This is done so
// one person with bad connectivity doesn't make it look like everything is
// horrible for everyone.
func (m *Metrics) collectPlaybackErrorCount() {
	valueSlice := utils.Float64MapToSlice(m.windowedErrorCounts)
	count := utils.Sum(valueSlice)
	m.windowedErrorCounts = map[string]float64{}

	m.metrics.errorCount = append(m.metrics.errorCount, models.TimestampedValue{
		Time:  time.Now(),
		Value: count,
	})

	if len(m.metrics.errorCount) > maxCollectionValues {
		m.metrics.errorCount = m.metrics.errorCount[1:]
	}

	// Save to Prometheus collector.
	m.playbackErrorCount.Set(count)
}

func (m *Metrics) collectSegmentDownloadDuration() {
	median := 0.0
	max := 0.0
	min := 0.0

	valueSlice := utils.Float64MapToSlice(m.windowedDownloadDurations)

	if len(valueSlice) > 0 {
		median = utils.Median(valueSlice)
		min, max = utils.MinMax(valueSlice)
		m.windowedDownloadDurations = map[string]float64{}
	}

	m.metrics.medianSegmentDownloadSeconds = append(m.metrics.medianSegmentDownloadSeconds, models.TimestampedValue{
		Time:  time.Now(),
		Value: median,
	})

	if len(m.metrics.medianSegmentDownloadSeconds) > maxCollectionValues {
		m.metrics.medianSegmentDownloadSeconds = m.metrics.medianSegmentDownloadSeconds[1:]
	}

	m.metrics.minimumSegmentDownloadSeconds = append(m.metrics.minimumSegmentDownloadSeconds, models.TimestampedValue{
		Time:  time.Now(),
		Value: min,
	})

	if len(m.metrics.minimumSegmentDownloadSeconds) > maxCollectionValues {
		m.metrics.minimumSegmentDownloadSeconds = m.metrics.minimumSegmentDownloadSeconds[1:]
	}

	m.metrics.maximumSegmentDownloadSeconds = append(m.metrics.maximumSegmentDownloadSeconds, models.TimestampedValue{
		Time:  time.Now(),
		Value: max,
	})

	if len(m.metrics.maximumSegmentDownloadSeconds) > maxCollectionValues {
		m.metrics.maximumSegmentDownloadSeconds = m.metrics.maximumSegmentDownloadSeconds[1:]
	}
}

// GetMedianDownloadDurationsOverTime will return a window of durations errors over time.
func (m *Metrics) GetMedianDownloadDurationsOverTime() []models.TimestampedValue {
	return m.metrics.medianSegmentDownloadSeconds
}

// GetMaximumDownloadDurationsOverTime will return a maximum durations errors over time.
func (m *Metrics) GetMaximumDownloadDurationsOverTime() []models.TimestampedValue {
	return m.metrics.maximumSegmentDownloadSeconds
}

// GetMinimumDownloadDurationsOverTime will return a maximum durations errors over time.
func (m *Metrics) GetMinimumDownloadDurationsOverTime() []models.TimestampedValue {
	return m.metrics.minimumSegmentDownloadSeconds
}

// GetPlaybackErrorCountOverTime will return a window of playback errors over time.
func (m *Metrics) GetPlaybackErrorCountOverTime() []models.TimestampedValue {
	return m.metrics.errorCount
}

func (m *Metrics) collectLatencyValues() {
	median := 0.0
	min := 0.0
	max := 0.0

	valueSlice := utils.Float64MapToSlice(m.windowedLatencies)
	m.windowedLatencies = map[string]float64{}

	if len(valueSlice) > 0 {
		median = utils.Median(valueSlice)
		min, max = utils.MinMax(valueSlice)
		m.windowedLatencies = map[string]float64{}
	}

	m.metrics.medianLatency = append(m.metrics.medianLatency, models.TimestampedValue{
		Time:  time.Now(),
		Value: median,
	})

	if len(m.metrics.medianLatency) > maxCollectionValues {
		m.metrics.medianLatency = m.metrics.medianLatency[1:]
	}

	m.metrics.minimumLatency = append(m.metrics.minimumLatency, models.TimestampedValue{
		Time:  time.Now(),
		Value: min,
	})

	if len(m.metrics.minimumLatency) > maxCollectionValues {
		m.metrics.minimumLatency = m.metrics.minimumLatency[1:]
	}

	m.metrics.maximumLatency = append(m.metrics.maximumLatency, models.TimestampedValue{
		Time:  time.Now(),
		Value: max,
	})

	if len(m.metrics.maximumLatency) > maxCollectionValues {
		m.metrics.maximumLatency = m.metrics.maximumLatency[1:]
	}
}

// GetMedianLatencyOverTime will return the median latency values over time.
func (m *Metrics) GetMedianLatencyOverTime() []models.TimestampedValue {
	if len(m.metrics.medianLatency) == 0 {
		return []models.TimestampedValue{}
	}

	return m.metrics.medianLatency
}

// GetMinimumLatencyOverTime will return the min latency values over time.
func (m *Metrics) GetMinimumLatencyOverTime() []models.TimestampedValue {
	if len(m.metrics.minimumLatency) == 0 {
		return []models.TimestampedValue{}
	}

	return m.metrics.minimumLatency
}

// GetMaximumLatencyOverTime will return the max latency values over time.
func (m *Metrics) GetMaximumLatencyOverTime() []models.TimestampedValue {
	if len(m.metrics.maximumLatency) == 0 {
		return []models.TimestampedValue{}
	}

	return m.metrics.maximumLatency
}

// collectLowestBandwidth will collect the bandwidth currently collected
// so we can report to the streamer the worst possible streaming condition
// being experienced.
func (m *Metrics) collectLowestBandwidth() {
	min := 0.0
	median := 0.0
	max := 0.0

	valueSlice := utils.Float64MapToSlice(m.windowedBandwidths)

	if len(m.windowedBandwidths) > 0 {
		min, max = utils.MinMax(valueSlice)
		min = math.Round(min)
		max = math.Round(max)
		median = utils.Median(valueSlice)
		m.windowedBandwidths = map[string]float64{}
	}

	m.metrics.lowestBitrate = append(m.metrics.lowestBitrate, models.TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(min),
	})

	if len(m.metrics.lowestBitrate) > maxCollectionValues {
		m.metrics.lowestBitrate = m.metrics.lowestBitrate[1:]
	}

	m.metrics.medianBitrate = append(m.metrics.medianBitrate, models.TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(median),
	})

	if len(m.metrics.medianBitrate) > maxCollectionValues {
		m.metrics.medianBitrate = m.metrics.medianBitrate[1:]
	}

	m.metrics.highestBitrate = append(m.metrics.highestBitrate, models.TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(max),
	})

	if len(m.metrics.highestBitrate) > maxCollectionValues {
		m.metrics.highestBitrate = m.metrics.highestBitrate[1:]
	}
}

// GetSlowestDownloadRateOverTime will return the collected lowest bandwidth values
// over time.
func (m *Metrics) GetSlowestDownloadRateOverTime() []models.TimestampedValue {
	if len(m.metrics.lowestBitrate) == 0 {
		return []models.TimestampedValue{}
	}

	return m.metrics.lowestBitrate
}

// GetMedianDownloadRateOverTime will return the collected median bandwidth values.
func (m *Metrics) GetMedianDownloadRateOverTime() []models.TimestampedValue {
	if len(m.metrics.medianBitrate) == 0 {
		return []models.TimestampedValue{}
	}
	return m.metrics.medianBitrate
}

// GetMaximumDownloadRateOverTime will return the collected maximum bandwidth values.
func (m *Metrics) GetMaximumDownloadRateOverTime() []models.TimestampedValue {
	if len(m.metrics.maximumLatency) == 0 {
		return []models.TimestampedValue{}
	}
	return m.metrics.maximumLatency
}

// GetMinimumDownloadRateOverTime will return the collected minimum bandwidth values.
func (m *Metrics) GetMinimumDownloadRateOverTime() []models.TimestampedValue {
	if len(m.metrics.minimumLatency) == 0 {
		return []models.TimestampedValue{}
	}
	return m.metrics.minimumLatency
}

// GetMaxDownloadRateOverTime will return the collected highest bandwidth values.
func (m *Metrics) GetMaxDownloadRateOverTime() []models.TimestampedValue {
	if len(m.metrics.highestBitrate) == 0 {
		return []models.TimestampedValue{}
	}
	return m.metrics.highestBitrate
}

func (m *Metrics) collectQualityVariantChanges() {
	valueSlice := utils.Float64MapToSlice(m.windowedQualityVariantChanges)
	count := utils.Sum(valueSlice)
	m.windowedQualityVariantChanges = map[string]float64{}

	m.metrics.qualityVariantChanges = append(m.metrics.qualityVariantChanges, models.TimestampedValue{
		Time:  time.Now(),
		Value: count,
	})
}

// GetQualityVariantChangesOverTime will return the collected quality variant
// changes.
func (m *Metrics) GetQualityVariantChangesOverTime() []models.TimestampedValue {
	return m.metrics.qualityVariantChanges
}

// GetPlaybackMetricsRepresentation returns what percentage of all known players
// the metrics represent.
func (m *Metrics) GetPlaybackMetricsRepresentation() int {
	totalPlayerCount := len(core.GetActiveViewers())
	representation := utils.IntPercentage(len(m.windowedBandwidths), totalPlayerCount)
	return representation
}
