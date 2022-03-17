package metrics

import (
	"math"
	"time"

	"github.com/owncast/owncast/utils"
)

// Playback error counts reported since the last time we collected metrics.
var (
	windowedErrorCounts           = []float64{}
	windowedQualityVariantChanges = []float64{}
	windowedBandwidths            = []float64{}
	windowedLatencies             = []float64{}
	windowedDownloadDurations     = []float64{}
)

// RegisterPlaybackErrorCount will add to the windowed playback error count.
func RegisterPlaybackErrorCount(count float64) {
	windowedErrorCounts = append(windowedErrorCounts, count)
}

// RegisterQualityVariantChangesCount will add to the windowed quality variant
// change count.
func RegisterQualityVariantChangesCount(count float64) {
	windowedQualityVariantChanges = append(windowedQualityVariantChanges, count)
}

// RegisterPlayerBandwidth will add to the windowed playback bandwidth.
func RegisterPlayerBandwidth(kbps float64) {
	windowedBandwidths = append(windowedBandwidths, kbps)
}

// RegisterPlayerLatency will add to the windowed player latency values.
func RegisterPlayerLatency(seconds float64) {
	windowedLatencies = append(windowedLatencies, seconds)
}

// RegisterPlayerSegmentDownloadDuration will add to the windowed player segment
// download duration values.
func RegisterPlayerSegmentDownloadDuration(seconds float64) {
	windowedDownloadDurations = append(windowedDownloadDurations, seconds)
}

// collectPlaybackErrorCount will take all of the error counts each individual
// player reported and average them into a single metric. This is done so
// one person with bad connectivity doesn't make it look like everything is
// horrible for everyone.
func collectPlaybackErrorCount() {
	count := utils.Sum(windowedErrorCounts)
	windowedErrorCounts = []float64{}

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

	if len(windowedDownloadDurations) > 0 {
		median = utils.Median(windowedDownloadDurations)
		min, max = utils.MinMax(windowedDownloadDurations)
		windowedDownloadDurations = []float64{}
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

	if len(windowedLatencies) > 0 {
		median = utils.Median(windowedLatencies)
		min, max = utils.MinMax(windowedLatencies)
		windowedLatencies = []float64{}
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

	if len(windowedBandwidths) > 0 {
		min, max = utils.MinMax(windowedBandwidths)
		min = math.Round(min)
		max = math.Round(max)
		median = utils.Median(windowedBandwidths)
		windowedBandwidths = []float64{}
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
	count := utils.Sum(windowedQualityVariantChanges)
	windowedQualityVariantChanges = []float64{}

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
