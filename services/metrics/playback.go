package metrics

import (
	"math"
	"time"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/utils"
)

// Playback error counts reported since the last time we collected metrics.
var (
	windowedErrorCounts           = map[string]float64{}
	windowedQualityVariantChanges = map[string]float64{}
	windowedBandwidths            = map[string]float64{}
	windowedLatencies             = map[string]float64{}
	windowedDownloadDurations     = map[string]float64{}
)

func handlePlaybackPolling() {
	metrics.m.Lock()
	defer metrics.m.Unlock()

	// Make sure this is fired first before all the values get cleared below.
	if _getStatus().Online {
		generateStreamHealthOverview()
	}

	collectPlaybackErrorCount()
	collectLatencyValues()
	collectSegmentDownloadDuration()
	collectLowestBandwidth()
	collectQualityVariantChanges()
}

// RegisterPlaybackErrorCount will add to the windowed playback error count.
func RegisterPlaybackErrorCount(clientID string, count float64) {
	metrics.m.Lock()
	defer metrics.m.Unlock()
	windowedErrorCounts[clientID] = count
	// windowedErrorCounts = append(windowedErrorCounts, count)
}

// RegisterQualityVariantChangesCount will add to the windowed quality variant
// change count.
func RegisterQualityVariantChangesCount(clientID string, count float64) {
	metrics.m.Lock()
	defer metrics.m.Unlock()
	windowedQualityVariantChanges[clientID] = count
}

// RegisterPlayerBandwidth will add to the windowed playback bandwidth.
func RegisterPlayerBandwidth(clientID string, kbps float64) {
	metrics.m.Lock()
	defer metrics.m.Unlock()
	windowedBandwidths[clientID] = kbps
}

// RegisterPlayerLatency will add to the windowed player latency values.
func RegisterPlayerLatency(clientID string, seconds float64) {
	metrics.m.Lock()
	defer metrics.m.Unlock()
	windowedLatencies[clientID] = seconds
}

// RegisterPlayerSegmentDownloadDuration will add to the windowed player segment
// download duration values.
func RegisterPlayerSegmentDownloadDuration(clientID string, seconds float64) {
	metrics.m.Lock()
	defer metrics.m.Unlock()
	windowedDownloadDurations[clientID] = seconds
}

// collectPlaybackErrorCount will take all of the error counts each individual
// player reported and average them into a single metric. This is done so
// one person with bad connectivity doesn't make it look like everything is
// horrible for everyone.
func collectPlaybackErrorCount() {
	valueSlice := utils.Float64MapToSlice(windowedErrorCounts)
	count := utils.Sum(valueSlice)
	windowedErrorCounts = map[string]float64{}

	metrics.errorCount = append(metrics.errorCount, TimestampedValue{
		Time:  time.Now(),
		Value: count,
	})

	if len(metrics.errorCount) > maxCollectionValues {
		metrics.errorCount = metrics.errorCount[1:]
	}

	// Save to Prometheus collector.
	playbackErrorCount.Set(count)
}

func collectSegmentDownloadDuration() {
	median := 0.0
	max := 0.0
	min := 0.0

	valueSlice := utils.Float64MapToSlice(windowedDownloadDurations)

	if len(valueSlice) > 0 {
		median = utils.Median(valueSlice)
		min, max = utils.MinMax(valueSlice)
		windowedDownloadDurations = map[string]float64{}
	}

	metrics.medianSegmentDownloadSeconds = append(metrics.medianSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: median,
	})

	if len(metrics.medianSegmentDownloadSeconds) > maxCollectionValues {
		metrics.medianSegmentDownloadSeconds = metrics.medianSegmentDownloadSeconds[1:]
	}

	metrics.minimumSegmentDownloadSeconds = append(metrics.minimumSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: min,
	})

	if len(metrics.minimumSegmentDownloadSeconds) > maxCollectionValues {
		metrics.minimumSegmentDownloadSeconds = metrics.minimumSegmentDownloadSeconds[1:]
	}

	metrics.maximumSegmentDownloadSeconds = append(metrics.maximumSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: max,
	})

	if len(metrics.maximumSegmentDownloadSeconds) > maxCollectionValues {
		metrics.maximumSegmentDownloadSeconds = metrics.maximumSegmentDownloadSeconds[1:]
	}
}

// GetMedianDownloadDurationsOverTime will return a window of durations errors over time.
func GetMedianDownloadDurationsOverTime() []TimestampedValue {
	return metrics.medianSegmentDownloadSeconds
}

// GetMaximumDownloadDurationsOverTime will return a maximum durations errors over time.
func GetMaximumDownloadDurationsOverTime() []TimestampedValue {
	return metrics.maximumSegmentDownloadSeconds
}

// GetMinimumDownloadDurationsOverTime will return a maximum durations errors over time.
func GetMinimumDownloadDurationsOverTime() []TimestampedValue {
	return metrics.minimumSegmentDownloadSeconds
}

// GetPlaybackErrorCountOverTime will return a window of playback errors over time.
func GetPlaybackErrorCountOverTime() []TimestampedValue {
	return metrics.errorCount
}

func collectLatencyValues() {
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

	metrics.medianLatency = append(metrics.medianLatency, TimestampedValue{
		Time:  time.Now(),
		Value: median,
	})

	if len(metrics.medianLatency) > maxCollectionValues {
		metrics.medianLatency = metrics.medianLatency[1:]
	}

	metrics.minimumLatency = append(metrics.minimumLatency, TimestampedValue{
		Time:  time.Now(),
		Value: min,
	})

	if len(metrics.minimumLatency) > maxCollectionValues {
		metrics.minimumLatency = metrics.minimumLatency[1:]
	}

	metrics.maximumLatency = append(metrics.maximumLatency, TimestampedValue{
		Time:  time.Now(),
		Value: max,
	})

	if len(metrics.maximumLatency) > maxCollectionValues {
		metrics.maximumLatency = metrics.maximumLatency[1:]
	}
}

// GetMedianLatencyOverTime will return the median latency values over time.
func GetMedianLatencyOverTime() []TimestampedValue {
	if len(metrics.medianLatency) == 0 {
		return []TimestampedValue{}
	}

	return metrics.medianLatency
}

// GetMinimumLatencyOverTime will return the min latency values over time.
func GetMinimumLatencyOverTime() []TimestampedValue {
	if len(metrics.minimumLatency) == 0 {
		return []TimestampedValue{}
	}

	return metrics.minimumLatency
}

// GetMaximumLatencyOverTime will return the max latency values over time.
func GetMaximumLatencyOverTime() []TimestampedValue {
	if len(metrics.maximumLatency) == 0 {
		return []TimestampedValue{}
	}

	return metrics.maximumLatency
}

// collectLowestBandwidth will collect the bandwidth currently collected
// so we can report to the streamer the worst possible streaming condition
// being experienced.
func collectLowestBandwidth() {
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

	metrics.lowestBitrate = append(metrics.lowestBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(min),
	})

	if len(metrics.lowestBitrate) > maxCollectionValues {
		metrics.lowestBitrate = metrics.lowestBitrate[1:]
	}

	metrics.medianBitrate = append(metrics.medianBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(median),
	})

	if len(metrics.medianBitrate) > maxCollectionValues {
		metrics.medianBitrate = metrics.medianBitrate[1:]
	}

	metrics.highestBitrate = append(metrics.highestBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(max),
	})

	if len(metrics.highestBitrate) > maxCollectionValues {
		metrics.highestBitrate = metrics.highestBitrate[1:]
	}
}

// GetSlowestDownloadRateOverTime will return the collected lowest bandwidth values
// over time.
func GetSlowestDownloadRateOverTime() []TimestampedValue {
	if len(metrics.lowestBitrate) == 0 {
		return []TimestampedValue{}
	}

	return metrics.lowestBitrate
}

// GetMedianDownloadRateOverTime will return the collected median bandwidth values.
func GetMedianDownloadRateOverTime() []TimestampedValue {
	if len(metrics.medianBitrate) == 0 {
		return []TimestampedValue{}
	}
	return metrics.medianBitrate
}

// GetMaximumDownloadRateOverTime will return the collected maximum bandwidth values.
func GetMaximumDownloadRateOverTime() []TimestampedValue {
	if len(metrics.maximumLatency) == 0 {
		return []TimestampedValue{}
	}
	return metrics.maximumLatency
}

// GetMinimumDownloadRateOverTime will return the collected minimum bandwidth values.
func GetMinimumDownloadRateOverTime() []TimestampedValue {
	if len(metrics.minimumLatency) == 0 {
		return []TimestampedValue{}
	}
	return metrics.minimumLatency
}

// GetMaxDownloadRateOverTime will return the collected highest bandwidth values.
func GetMaxDownloadRateOverTime() []TimestampedValue {
	if len(metrics.highestBitrate) == 0 {
		return []TimestampedValue{}
	}
	return metrics.highestBitrate
}

func collectQualityVariantChanges() {
	valueSlice := utils.Float64MapToSlice(windowedQualityVariantChanges)
	count := utils.Sum(valueSlice)
	windowedQualityVariantChanges = map[string]float64{}

	metrics.qualityVariantChanges = append(metrics.qualityVariantChanges, TimestampedValue{
		Time:  time.Now(),
		Value: count,
	})
}

// GetQualityVariantChangesOverTime will return the collected quality variant
// changes.
func GetQualityVariantChangesOverTime() []TimestampedValue {
	return metrics.qualityVariantChanges
}

// GetPlaybackMetricsRepresentation returns what percentage of all known players
// the metrics represent.
func GetPlaybackMetricsRepresentation() int {
	totalPlayerCount := len(core.GetActiveViewers())
	representation := utils.IntPercentage(len(windowedBandwidths), totalPlayerCount)
	return representation
}
