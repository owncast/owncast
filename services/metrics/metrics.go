package metrics

import (
	"sync"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	metrics                       *CollectedMetrics
	getStatus                     func() models.Status
	windowedErrorCounts           map[string]float64
	windowedQualityVariantChanges map[string]float64
	windowedBandwidths            map[string]float64
	windowedLatencies             map[string]float64
	windowedDownloadDurations     map[string]float64

	// Prometheus
	labels                  map[string]string
	activeViewerCount       prometheus.Gauge
	activeChatClientCount   prometheus.Gauge
	cpuUsage                prometheus.Gauge
	chatUserCount           prometheus.Gauge
	currentChatMessageCount prometheus.Gauge
	playbackErrorCount      prometheus.Gauge
}

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

	medianSegmentDownloadSeconds  []models.TimestampedValue `json:"-"`
	maximumSegmentDownloadSeconds []models.TimestampedValue `json:"-"`
	DiskUtilizations              []models.TimestampedValue `json:"disk"`

	errorCount      []models.TimestampedValue `json:"-"`
	lowestBitrate   []models.TimestampedValue `json:"-"`
	medianBitrate   []models.TimestampedValue `json:"-"`
	RAMUtilizations []models.TimestampedValue `json:"memory"`

	CPUUtilizations []models.TimestampedValue `json:"cpu"`
	highestBitrate  []models.TimestampedValue `json:"-"`

	minimumSegmentDownloadSeconds []models.TimestampedValue `json:"-"`

	minimumLatency []models.TimestampedValue `json:"-"`
	maximumLatency []models.TimestampedValue `json:"-"`
	medianLatency  []models.TimestampedValue `json:"-"`

	qualityVariantChanges []models.TimestampedValue `json:"-"`

	m sync.Mutex `json:"-"`
}

// New will return a new Metrics instance.
func New() *Metrics {
	return &Metrics{
		windowedErrorCounts:           map[string]float64{},
		windowedQualityVariantChanges: map[string]float64{},
		windowedBandwidths:            map[string]float64{},
		windowedLatencies:             map[string]float64{},
		windowedDownloadDurations:     map[string]float64{},
	}
}

// Metrics is the shared Metrics instance.

// Start will begin the metrics collection and alerting.
func (m *Metrics) Start(getStatus func() models.Status) {
	m.getStatus = getStatus
	host := configRepository.GetServerURL()
	if host == "" {
		host = "unknown"
	}

	c := config.GetConfig()

	m.labels = map[string]string{
		"version": c.VersionNumber,
		"host":    host,
	}

	m.setupPrometheusCollectors()

	m.metrics = new(CollectedMetrics)
	go m.startViewerCollectionMetrics()

	go func() {
		for range time.Tick(hardwareMetricsPollingInterval) {
			m.handlePolling()
		}
	}()

	go func() {
		for range time.Tick(playbackMetricsPollingInterval) {
			m.handlePlaybackPolling()
		}
	}()
}

func (m *Metrics) handlePolling() {
	m.metrics.m.Lock()
	defer m.metrics.m.Unlock()

	// Collect hardware stats
	m.collectCPUUtilization()
	m.collectRAMUtilization()
	m.collectDiskUtilization()

	// Alerting
	m.handleAlerting()
}

// GetMetrics will return the collected metrics.
func (m *Metrics) GetMetrics() *CollectedMetrics {
	return m.metrics
}

var temporaryGlobalInstance *Metrics

func Get() *Metrics {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = new(Metrics)
	}
	return temporaryGlobalInstance
}
