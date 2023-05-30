package metrics

import (
	"sync"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

// How often we poll for updates.
const hardwareMetricsPollingInterval = 2 * time.Minute
const playbackMetricsPollingInterval = 2 * time.Minute

const (
	// How often we poll for updates.
	viewerMetricsPollingInterval = 2 * time.Minute
	activeChatClientCountKey     = "chat_client_count"
	activeViewerCountKey         = "viewer_count"
)

// CollectedMetrics stores different collected + timestamped values.
type CollectedMetrics struct {
	streamHealthOverview *models.StreamHealthOverview

	medianSegmentDownloadSeconds  []TimestampedValue `json:"-"`
	maximumSegmentDownloadSeconds []TimestampedValue `json:"-"`
	DiskUtilizations              []TimestampedValue `json:"disk"`

	errorCount      []TimestampedValue `json:"-"`
	lowestBitrate   []TimestampedValue `json:"-"`
	medianBitrate   []TimestampedValue `json:"-"`
	RAMUtilizations []TimestampedValue `json:"memory"`

	CPUUtilizations []TimestampedValue `json:"cpu"`
	highestBitrate  []TimestampedValue `json:"-"`

	minimumSegmentDownloadSeconds []TimestampedValue `json:"-"`

	minimumLatency []TimestampedValue `json:"-"`
	maximumLatency []TimestampedValue `json:"-"`
	medianLatency  []TimestampedValue `json:"-"`

	qualityVariantChanges []TimestampedValue `json:"-"`

	m sync.Mutex `json:"-"`
}

// Metrics is the shared Metrics instance.
var metrics *CollectedMetrics

var _getStatus func() models.Status

// Start will begin the metrics collection and alerting.
func Start(getStatus func() models.Status) {
	_getStatus = getStatus
	host := data.GetServerURL()
	if host == "" {
		host = "unknown"
	}
	labels = map[string]string{
		"version": config.VersionNumber,
		"host":    host,
	}

	setupPrometheusCollectors()

	metrics = new(CollectedMetrics)
	go startViewerCollectionMetrics()

	go func() {
		for range time.Tick(hardwareMetricsPollingInterval) {
			handlePolling()
		}
	}()

	go func() {
		for range time.Tick(playbackMetricsPollingInterval) {
			handlePlaybackPolling()
		}
	}()
}

func handlePolling() {
	metrics.m.Lock()
	defer metrics.m.Unlock()

	// Collect hardware stats
	collectCPUUtilization()
	collectRAMUtilization()
	collectDiskUtilization()

	// Alerting
	handleAlerting()
}

// GetMetrics will return the collected metrics.
func GetMetrics() *CollectedMetrics {
	return metrics
}
