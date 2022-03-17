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
	val := 0.0

	if len(windowedDownloadDurations) > 0 {
		val = utils.Avg(windowedDownloadDurations)
		windowedDownloadDurations = []float64{}
	}
	metrics.segmentDownloadSeconds = append(metrics.segmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: val,
	})

	if len(metrics.segmentDownloadSeconds) > maxCollectionValues {
		metrics.segmentDownloadSeconds = metrics.segmentDownloadSeconds[1:]
	}
}

// GetDownloadDurationsOverTime will return a window of durations errors over time.
func GetDownloadDurationsOverTime() []TimestampedValue {
	return metrics.segmentDownloadSeconds
}

// GetPlaybackErrorCountOverTime will return a window of playback errors over time.
func GetPlaybackErrorCountOverTime() []TimestampedValue {
	return metrics.errorCount
}

func collectLatencyValues() {
	val := 0.0

	if len(windowedLatencies) > 0 {
		val = utils.Avg(windowedLatencies)
		val = math.Round(val)
		windowedLatencies = []float64{}
	}

	metrics.averageLatency = append(metrics.averageLatency, TimestampedValue{
		Time:  time.Now(),
		Value: val,
	})

	if len(metrics.averageLatency) > maxCollectionValues {
		metrics.averageLatency = metrics.averageLatency[1:]
	}
}

// GetLatencyOverTime will return the min, max and avg latency values over time.
func GetLatencyOverTime() []TimestampedValue {
	if len(metrics.averageLatency) == 0 {
		return []TimestampedValue{}
	}

	return metrics.averageLatency
}

// collectLowestBandwidth will collect the lowest bandwidth currently collected
// so we can report to the streamer the worst possible streaming condition
// being experienced.
func collectLowestBandwidth() {
	val := 0.0

	if len(windowedBandwidths) > 0 {
		val, _ = utils.MinMax(windowedBandwidths)
		val = math.Round(val)
		windowedBandwidths = []float64{}
	}

	metrics.lowestBitrate = append(metrics.lowestBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(val),
	})

	if len(metrics.lowestBitrate) > maxCollectionValues {
		metrics.lowestBitrate = metrics.lowestBitrate[1:]
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
