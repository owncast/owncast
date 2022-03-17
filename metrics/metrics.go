package metrics

import (
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
)

// How often we poll for updates.
const hardwareMetricsPollingInterval = 1 * time.Minute

const (
	// How often we poll for updates.
	viewerMetricsPollingInterval = 2 * time.Minute
	activeChatClientCountKey     = "chat_client_count"
	activeViewerCountKey         = "viewer_count"
)

// CollectedMetrics stores different collected + timestamped values.
type CollectedMetrics struct {
	CPUUtilizations  []TimestampedValue `json:"cpu"`
	RAMUtilizations  []TimestampedValue `json:"memory"`
	DiskUtilizations []TimestampedValue `json:"disk"`

	errorCount     []TimestampedValue `json:"-"`
	lowestBitrate  []TimestampedValue `json:"-"`
	medianBitrate  []TimestampedValue `json:"-"`
	highestBitrate []TimestampedValue `json:"-"`

	medianSegmentDownloadSeconds  []TimestampedValue `json:"-"`
	maximumSegmentDownloadSeconds []TimestampedValue `json:"-"`
	minimumSegmentDownloadSeconds []TimestampedValue `json:"-"`

	minimumLatency []TimestampedValue `json:"-"`
	maximumLatency []TimestampedValue `json:"-"`
	medianLatency  []TimestampedValue `json:"-"`

	qualityVariantChanges []TimestampedValue `json:"-"`
}

// Metrics is the shared Metrics instance.
var metrics *CollectedMetrics

// Start will begin the metrics collection and alerting.
func Start(getStatus func() models.Status) {
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

	for range time.Tick(hardwareMetricsPollingInterval) {
		handlePolling()
	}
}

func handlePolling() {
	// Collect hardware stats
	collectCPUUtilization()
	collectRAMUtilization()
	collectDiskUtilization()

	collectPlaybackErrorCount()
	collectLatencyValues()
	collectSegmentDownloadDuration()
	collectLowestBandwidth()
	collectQualityVariantChanges()

	// Alerting
	handleAlerting()
}

// GetMetrics will return the collected metrics.
func GetMetrics() *CollectedMetrics {
	return metrics
}
