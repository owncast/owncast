// Package metrics owns the Owncast metrics service: hardware, playback,
// chat and viewer metrics collection plus Prometheus exposition.
// Construct via New(Deps) and call Start to launch collection goroutines.
package metrics

import (
	"sync"
	"time"

	"github.com/nakabonne/tstorage"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/userrepository"
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

// Deps lists the explicit dependencies of the metrics Service.
type Deps struct {
	Stream                *stream.Service
	Chat                  *chat.Service
	ConfigRepository      configrepository.ConfigRepository
	ChatMessageRepository chatmessagerepository.ChatMessageRepository
	UserRepository        userrepository.UserRepository
}

// Service owns the per-instance metrics collection state: the collected
// time-series window, Prometheus gauges, the on-disk tstorage handle, and
// the windowed per-client maps. Construct via New(Deps); call Start to
// launch the collection goroutines.
type Service struct {
	stream                *stream.Service
	chat                  *chat.Service
	configRepository      configrepository.ConfigRepository
	chatMessageRepository chatmessagerepository.ChatMessageRepository
	userRepository        userrepository.UserRepository

	// Metrics is the collected window of hardware/playback timestamped
	// values used by the admin health endpoints.
	metrics *CollectedMetrics

	// Prometheus collectors. Built by setupPrometheusCollectors() during
	// Start; do not access before Start.
	labels                  map[string]string
	activeViewerCount       prometheus.Gauge
	activeChatClientCount   prometheus.Gauge
	cpuUsage                prometheus.Gauge
	chatUserCount           prometheus.Gauge
	currentChatMessageCount prometheus.Gauge
	playbackErrorCount      prometheus.Gauge

	// On-disk time-series store used by GetViewersOverTime and
	// GetChatClientCountOverTime. Opened by startViewerCollectionMetrics.
	storage tstorage.Storage

	// Windowed per-client playback values consumed by the periodic
	// playback aggregation pass.
	windowedErrorCounts           map[string]float64
	windowedQualityVariantChanges map[string]float64
	windowedBandwidths            map[string]float64
	windowedLatencies             map[string]float64
	windowedDownloadDurations     map[string]float64

	// Alerting state: whether we've already logged a threshold-exceeded
	// warning for the given resource (debounced by errorResetDuration).
	inCPUAlertingState  bool
	inRAMAlertingState  bool
	inDiskAlertingState bool
}

// New constructs an idle metrics Service. Call Start to launch the
// collection goroutines and register the Prometheus collectors.
func New(deps Deps) *Service {
	return &Service{
		stream:                        deps.Stream,
		chat:                          deps.Chat,
		configRepository:              deps.ConfigRepository,
		chatMessageRepository:         deps.ChatMessageRepository,
		userRepository:                deps.UserRepository,
		windowedErrorCounts:           map[string]float64{},
		windowedQualityVariantChanges: map[string]float64{},
		windowedBandwidths:            map[string]float64{},
		windowedLatencies:             map[string]float64{},
		windowedDownloadDurations:     map[string]float64{},
	}
}

// Start will begin the metrics collection and alerting.
func (s *Service) Start() {
	host := s.configRepository.GetServerURL()
	if host == "" {
		host = "unknown"
	}
	s.labels = map[string]string{
		"version": config.VersionNumber,
		"host":    host,
	}

	s.setupPrometheusCollectors()

	s.metrics = new(CollectedMetrics)
	go s.startViewerCollectionMetrics()

	go func() {
		for range time.Tick(hardwareMetricsPollingInterval) {
			s.handlePolling()
		}
	}()

	go func() {
		for range time.Tick(playbackMetricsPollingInterval) {
			s.handlePlaybackPolling()
		}
	}()
}

func (s *Service) handlePolling() {
	s.metrics.m.Lock()
	defer s.metrics.m.Unlock()

	// Collect hardware stats
	s.collectCPUUtilization()
	s.collectRAMUtilization()
	s.collectDiskUtilization()

	// Alerting
	s.handleAlerting()
}

// GetMetrics will return the collected metrics.
func (s *Service) GetMetrics() *CollectedMetrics {
	return s.metrics
}
