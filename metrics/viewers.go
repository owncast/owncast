package metrics

import (
	"time"

	"github.com/nakabonne/tstorage"
	log "github.com/sirupsen/logrus"
)

var storage tstorage.Storage

func (s *Service) startViewerCollectionMetrics() {
	storage, _ = tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithDataPath("./data/metrics"),
	)
	defer storage.Close()

	s.collectViewerCount()

	for range time.Tick(viewerMetricsPollingInterval) {
		s.collectViewerCount()
		s.collectChatClientCount()
	}
}

func (s *Service) collectViewerCount() {
	// Don't collect metrics for viewers if there's no stream active.
	if !s.Core.GetStatus().Online {
		return
	}

	count := s.Core.GetStatus().ViewerCount

	// Save active viewer count to our Prometheus collector.
	activeViewerCount.Set(float64(count))

	// Insert active viewer count into our on-disk time series storage.
	if err := storage.InsertRows([]tstorage.Row{
		{
			Metric:    activeViewerCountKey,
			DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: float64(count)},
		},
	}); err != nil {
		log.Errorln(err)
	}
}

func (s *Service) collectChatClientCount() {
	count := len(s.Core.Chat.GetClients())
	activeChatClientCount.Set(float64(count))

	// Total message count
	cmc := s.Core.Data.GetMessagesCount()
	// Insert message count into Prometheus collector.
	currentChatMessageCount.Set(float64(cmc))

	// Total user count
	uc := s.Core.Data.GetUsersCount()
	// Insert user count into Prometheus collector.
	chatUserCount.Set(float64(uc))

	// Insert active chat user count into our on-disk time series storage.
	if err := storage.InsertRows([]tstorage.Row{
		{
			Metric:    activeChatClientCountKey,
			DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: float64(count)},
		},
	}); err != nil {
		log.Errorln(err)
	}
}

// GetViewersOverTime will return a window of viewer counts over time.
func (s *Service) GetViewersOverTime(start, end time.Time) []TimestampedValue {
	p, err := storage.Select(activeViewerCountKey, nil, start.Unix(), end.Unix())
	if err != nil && err != tstorage.ErrNoDataPoints {
		log.Errorln(err)
	}
	datapoints := makeTimestampedValuesFromDatapoints(p)

	return datapoints
}

// GetChatClientCountOverTime will return a window of connected chat clients over time.
func (s *Service) GetChatClientCountOverTime(start, end time.Time) []TimestampedValue {
	p, err := storage.Select(activeChatClientCountKey, nil, start.Unix(), end.Unix())
	if err != nil && err != tstorage.ErrNoDataPoints {
		log.Errorln(err)
	}
	datapoints := makeTimestampedValuesFromDatapoints(p)

	return datapoints
}
