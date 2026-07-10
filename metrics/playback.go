package metrics

import (
	"math"
	"time"

	"github.com/owncast/owncast/utils"
)

// How long after a direct client playback report we keep suppressing
// lower-fidelity server-derived measurements for that client. Players
// report every ~10s, so a client silent for this long has stopped
// self-reporting.
const selfReportSuppressionWindow = 30 * time.Second

func (s *Service) handlePlaybackPolling() {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()

	// Make sure this is fired first before all the values get cleared below.
	if s.stream.GetStatus().Online {
		s.generateStreamHealthOverview()
	}

	s.collectPlaybackErrorCount()
	s.collectLatencyValues()
	s.collectSegmentDownloadDuration()
	s.collectLowestBandwidth()
	s.collectQualityVariantChanges()

	s.pruneSelfReportingClients()
}

// RegisterSelfReportingClient marks a client as directly reporting its own
// playback metrics, which suppresses server-derived observation for it —
// the client's own measurements are richer.
func (s *Service) RegisterSelfReportingClient(clientID string) {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()
	s.selfReportingClients[clientID] = time.Now()
}

// IsClientSelfReporting returns true if the client recently reported its
// own playback metrics.
func (s *Service) IsClientSelfReporting(clientID string) bool {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()
	t, ok := s.selfReportingClients[clientID]
	return ok && time.Since(t) < selfReportSuppressionWindow
}

// pruneSelfReportingClients drops clients that have stopped reporting.
// Callers must hold s.metrics.m.
func (s *Service) pruneSelfReportingClients() {
	for id, t := range s.selfReportingClients {
		if time.Since(t) >= selfReportSuppressionWindow {
			delete(s.selfReportingClients, id)
		}
	}
}

// PlaybackReport is a batch of one client's playback measurements applied
// in a single locked update. Zero values are skipped; ErrorCount is applied
// whenever HasErrorCount is set so a zero count still marks the client as a
// healthy participant in the stream health overview.
type PlaybackReport struct {
	ClientID        string
	BandwidthKbps   float64
	LatencySeconds  float64
	DownloadSeconds float64
	BitrateKbps     float64
	ErrorCount      float64
	HasErrorCount   bool
}

// RegisterPlaybackReport applies a client's playback measurements under one
// lock acquisition. This is the hot-path entry point used per media request
// (CMCD request mode and server-side observation); avoid adding per-field
// lock round-trips here.
func (s *Service) RegisterPlaybackReport(report PlaybackReport) {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()

	if report.BandwidthKbps > 0 {
		s.windowedBandwidths[report.ClientID] = report.BandwidthKbps
	}
	if report.LatencySeconds > 0 {
		s.windowedLatencies[report.ClientID] = report.LatencySeconds
	}
	if report.DownloadSeconds > 0 {
		s.windowedDownloadDurations[report.ClientID] = report.DownloadSeconds
	}
	// A change in the encoded bitrate being played is a quality variant
	// change, giving variant switch tracking for every CMCD player.
	if report.BitrateKbps > 0 {
		if prev, ok := s.lastClientBitrates[report.ClientID]; ok && prev != report.BitrateKbps {
			s.windowedQualityVariantChanges[report.ClientID]++
		}
		s.lastClientBitrates[report.ClientID] = report.BitrateKbps
	}
	if report.HasErrorCount {
		s.windowedErrorCounts[report.ClientID] += report.ErrorCount
	}
}

// RegisterPlaybackErrorCount will add to the windowed playback error count.
// Players report deltas every ~10s but the window is only harvested every
// 2 minutes, so counts must accumulate — assignment would drop every report
// that is followed by a zero before the collector runs.
func (s *Service) RegisterPlaybackErrorCount(clientID string, count float64) {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()
	s.windowedErrorCounts[clientID] += count
}

// RegisterQualityVariantChangesCount will add to the windowed quality variant
// change count.
func (s *Service) RegisterQualityVariantChangesCount(clientID string, count float64) {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()
	s.windowedQualityVariantChanges[clientID] += count
}

// RegisterPlayerBandwidth will add to the windowed playback bandwidth.
func (s *Service) RegisterPlayerBandwidth(clientID string, kbps float64) {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()
	s.windowedBandwidths[clientID] = kbps
}

// RegisterPlayerLatency will add to the windowed player latency values.
func (s *Service) RegisterPlayerLatency(clientID string, seconds float64) {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()
	s.windowedLatencies[clientID] = seconds
}

// RegisterPlayerSegmentDownloadDuration will add to the windowed player segment
// download duration values.
func (s *Service) RegisterPlayerSegmentDownloadDuration(clientID string, seconds float64) {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()
	s.windowedDownloadDurations[clientID] = seconds
}

// collectPlaybackErrorCount will take all of the error counts each individual
// player reported and average them into a single metric. This is done so
// one person with bad connectivity doesn't make it look like everything is
// horrible for everyone.
func (s *Service) collectPlaybackErrorCount() {
	valueSlice := utils.Float64MapToSlice(s.windowedErrorCounts)
	count := utils.Sum(valueSlice)
	s.windowedErrorCounts = map[string]float64{}

	s.metrics.errorCount = append(s.metrics.errorCount, TimestampedValue{
		Time:  time.Now(),
		Value: count,
	})

	if len(s.metrics.errorCount) > maxCollectionValues {
		s.metrics.errorCount = s.metrics.errorCount[1:]
	}

	// Save to Prometheus collector.
	s.playbackErrorCount.Set(count)
}

func (s *Service) collectSegmentDownloadDuration() {
	median := 0.0
	max := 0.0
	min := 0.0

	valueSlice := utils.Float64MapToSlice(s.windowedDownloadDurations)

	if len(valueSlice) > 0 {
		median = utils.Median(valueSlice)
		min, max = utils.MinMax(valueSlice)
		s.windowedDownloadDurations = map[string]float64{}
	}

	s.metrics.medianSegmentDownloadSeconds = append(s.metrics.medianSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: median,
	})

	if len(s.metrics.medianSegmentDownloadSeconds) > maxCollectionValues {
		s.metrics.medianSegmentDownloadSeconds = s.metrics.medianSegmentDownloadSeconds[1:]
	}

	s.metrics.minimumSegmentDownloadSeconds = append(s.metrics.minimumSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: min,
	})

	if len(s.metrics.minimumSegmentDownloadSeconds) > maxCollectionValues {
		s.metrics.minimumSegmentDownloadSeconds = s.metrics.minimumSegmentDownloadSeconds[1:]
	}

	s.metrics.maximumSegmentDownloadSeconds = append(s.metrics.maximumSegmentDownloadSeconds, TimestampedValue{
		Time:  time.Now(),
		Value: max,
	})

	if len(s.metrics.maximumSegmentDownloadSeconds) > maxCollectionValues {
		s.metrics.maximumSegmentDownloadSeconds = s.metrics.maximumSegmentDownloadSeconds[1:]
	}
}

// GetMedianDownloadDurationsOverTime will return a window of durations errors over time.
func (s *Service) GetMedianDownloadDurationsOverTime() []TimestampedValue {
	return s.metrics.medianSegmentDownloadSeconds
}

// GetMaximumDownloadDurationsOverTime will return a maximum durations errors over time.
func (s *Service) GetMaximumDownloadDurationsOverTime() []TimestampedValue {
	return s.metrics.maximumSegmentDownloadSeconds
}

// GetMinimumDownloadDurationsOverTime will return a maximum durations errors over time.
func (s *Service) GetMinimumDownloadDurationsOverTime() []TimestampedValue {
	return s.metrics.minimumSegmentDownloadSeconds
}

// GetPlaybackErrorCountOverTime will return a window of playback errors over time.
func (s *Service) GetPlaybackErrorCountOverTime() []TimestampedValue {
	return s.metrics.errorCount
}

func (s *Service) collectLatencyValues() {
	median := 0.0
	min := 0.0
	max := 0.0

	valueSlice := utils.Float64MapToSlice(s.windowedLatencies)
	s.windowedLatencies = map[string]float64{}

	if len(valueSlice) > 0 {
		median = utils.Median(valueSlice)
		min, max = utils.MinMax(valueSlice)
		s.windowedLatencies = map[string]float64{}
	}

	s.metrics.medianLatency = append(s.metrics.medianLatency, TimestampedValue{
		Time:  time.Now(),
		Value: median,
	})

	if len(s.metrics.medianLatency) > maxCollectionValues {
		s.metrics.medianLatency = s.metrics.medianLatency[1:]
	}

	s.metrics.minimumLatency = append(s.metrics.minimumLatency, TimestampedValue{
		Time:  time.Now(),
		Value: min,
	})

	if len(s.metrics.minimumLatency) > maxCollectionValues {
		s.metrics.minimumLatency = s.metrics.minimumLatency[1:]
	}

	s.metrics.maximumLatency = append(s.metrics.maximumLatency, TimestampedValue{
		Time:  time.Now(),
		Value: max,
	})

	if len(s.metrics.maximumLatency) > maxCollectionValues {
		s.metrics.maximumLatency = s.metrics.maximumLatency[1:]
	}
}

// GetMedianLatencyOverTime will return the median latency values over time.
func (s *Service) GetMedianLatencyOverTime() []TimestampedValue {
	if len(s.metrics.medianLatency) == 0 {
		return []TimestampedValue{}
	}

	return s.metrics.medianLatency
}

// GetMinimumLatencyOverTime will return the min latency values over time.
func (s *Service) GetMinimumLatencyOverTime() []TimestampedValue {
	if len(s.metrics.minimumLatency) == 0 {
		return []TimestampedValue{}
	}

	return s.metrics.minimumLatency
}

// GetMaximumLatencyOverTime will return the max latency values over time.
func (s *Service) GetMaximumLatencyOverTime() []TimestampedValue {
	if len(s.metrics.maximumLatency) == 0 {
		return []TimestampedValue{}
	}

	return s.metrics.maximumLatency
}

// collectLowestBandwidth will collect the bandwidth currently collected
// so we can report to the streamer the worst possible streaming condition
// being experienced.
func (s *Service) collectLowestBandwidth() {
	min := 0.0
	median := 0.0
	max := 0.0

	valueSlice := utils.Float64MapToSlice(s.windowedBandwidths)

	if len(s.windowedBandwidths) > 0 {
		min, max = utils.MinMax(valueSlice)
		min = math.Round(min)
		max = math.Round(max)
		median = utils.Median(valueSlice)
		s.windowedBandwidths = map[string]float64{}
	}

	s.metrics.lowestBitrate = append(s.metrics.lowestBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(min),
	})

	if len(s.metrics.lowestBitrate) > maxCollectionValues {
		s.metrics.lowestBitrate = s.metrics.lowestBitrate[1:]
	}

	s.metrics.medianBitrate = append(s.metrics.medianBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(median),
	})

	if len(s.metrics.medianBitrate) > maxCollectionValues {
		s.metrics.medianBitrate = s.metrics.medianBitrate[1:]
	}

	s.metrics.highestBitrate = append(s.metrics.highestBitrate, TimestampedValue{
		Time:  time.Now(),
		Value: math.Round(max),
	})

	if len(s.metrics.highestBitrate) > maxCollectionValues {
		s.metrics.highestBitrate = s.metrics.highestBitrate[1:]
	}
}

// GetSlowestDownloadRateOverTime will return the collected lowest bandwidth values
// over time.
func (s *Service) GetSlowestDownloadRateOverTime() []TimestampedValue {
	if len(s.metrics.lowestBitrate) == 0 {
		return []TimestampedValue{}
	}

	return s.metrics.lowestBitrate
}

// GetMedianDownloadRateOverTime will return the collected median bandwidth values.
func (s *Service) GetMedianDownloadRateOverTime() []TimestampedValue {
	if len(s.metrics.medianBitrate) == 0 {
		return []TimestampedValue{}
	}
	return s.metrics.medianBitrate
}

// GetMaximumDownloadRateOverTime will return the collected maximum bandwidth values.
func (s *Service) GetMaximumDownloadRateOverTime() []TimestampedValue {
	if len(s.metrics.maximumLatency) == 0 {
		return []TimestampedValue{}
	}
	return s.metrics.maximumLatency
}

// GetMinimumDownloadRateOverTime will return the collected minimum bandwidth values.
func (s *Service) GetMinimumDownloadRateOverTime() []TimestampedValue {
	if len(s.metrics.minimumLatency) == 0 {
		return []TimestampedValue{}
	}
	return s.metrics.minimumLatency
}

// GetMaxDownloadRateOverTime will return the collected highest bandwidth values.
func (s *Service) GetMaxDownloadRateOverTime() []TimestampedValue {
	if len(s.metrics.highestBitrate) == 0 {
		return []TimestampedValue{}
	}
	return s.metrics.highestBitrate
}

func (s *Service) collectQualityVariantChanges() {
	valueSlice := utils.Float64MapToSlice(s.windowedQualityVariantChanges)
	count := utils.Sum(valueSlice)
	s.windowedQualityVariantChanges = map[string]float64{}
	// Resetting the baselines each window bounds the map at the
	// cost of missing a variant switch that spans a harvest boundary.
	s.lastClientBitrates = map[string]float64{}

	s.metrics.qualityVariantChanges = append(s.metrics.qualityVariantChanges, TimestampedValue{
		Time:  time.Now(),
		Value: count,
	})
}

// GetQualityVariantChangesOverTime will return the collected quality variant
// changes.
func (s *Service) GetQualityVariantChangesOverTime() []TimestampedValue {
	return s.metrics.qualityVariantChanges
}

// GetPlaybackMetricsRepresentation returns what percentage of all known players
// the metrics represent.
func (s *Service) GetPlaybackMetricsRepresentation() int {
	totalPlayerCount := len(s.stream.GetActiveViewers())
	representation := utils.IntPercentage(len(s.windowedBandwidths), totalPlayerCount)
	return representation
}
