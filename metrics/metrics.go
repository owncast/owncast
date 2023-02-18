package metrics

import (
	"sync"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/models"
)

const (
	// How often we poll for updates.
	viewerMetricsPollingInterval = 2 * time.Minute
	activeChatClientCountKey     = "chat_client_count"
	activeViewerCountKey         = "viewer_count"

	hardwareMetricsPollingInterval = 2 * time.Minute
	playbackMetricsPollingInterval = 2 * time.Minute
)

func New(c *core.Service) (*Service, error) {
	s := &Service{
		Core: c,
	}

	return s, nil
}

// Service stores different collected + timestamped values.
type Service struct {
	Core *core.Service

	Metrics struct {
		mux sync.Mutex `json:"-"`

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

		streamHealthOverview *models.StreamHealthOverview
	}
}

var _getStatus func() models.Status

// Start will begin the metrics collection and alerting.
func (s *Service) Start(getStatus func() models.Status) {
	_getStatus = getStatus
	host := s.Core.Data.GetServerURL()
	if host == "" {
		host = "unknown"
	}
	labels = map[string]string{
		"version": config.VersionNumber,
		"host":    host,
	}

	setupPrometheusCollectors()

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
	s.Metrics.mux.Lock()
	defer s.Metrics.mux.Unlock()

	// Collect hardware stats
	s.collectCPUUtilization()
	s.collectRAMUtilization()
	s.collectDiskUtilization()

	// Alerting
	s.handleAlerting()
}

// GetMetrics will return the collected s.
func (s *Service) GetMetrics() *Service {
	return s
}
