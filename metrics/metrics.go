package metrics

import (
	"sync"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/services/stream"
)

// How often we poll for updates.
const (
	hardwareMetricsPollingInterval = 2 * time.Minute
	playbackMetricsPollingInterval = 2 * time.Minute
)

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

// streamSvc and chatSvc are the injected services used by all metrics
// collection routines. Set once by Start; do not access before then.
var (
	streamSvc *stream.Service
	chatSvc   *chat.Service
)

// Start will begin the metrics collection and alerting.
func Start(s *stream.Service, c *chat.Service) {
	streamSvc = s
	chatSvc = c
	host := configRepository.GetServerURL()
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
