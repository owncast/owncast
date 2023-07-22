package metrics

import (
	"time"

	"github.com/nakabonne/tstorage"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/status"
	"github.com/owncast/owncast/storage/chatrepository"
	"github.com/owncast/owncast/storage/userrepository"
	log "github.com/sirupsen/logrus"
)

var storage tstorage.Storage

func (m *Metrics) startViewerCollectionMetrics() {
	storage, _ = tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithDataPath("./data/metrics"),
	)
	defer storage.Close()

	m.collectViewerCount()

	for range time.Tick(viewerMetricsPollingInterval) {
		m.collectViewerCount()
		m.collectChatClientCount()
	}
}

func (m *Metrics) collectViewerCount() {
	s := status.Get()

	// Don't collect metrics for viewers if there's no stream active.
	if !s.Online {
		return
	}

	count := s.ViewerCount

	// Save active viewer count to our Prometheus collector.
	m.activeViewerCount.Set(float64(count))

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

func (m *Metrics) collectChatClientCount() {
	count := len(m.chatService.GetClients())
	m.activeChatClientCount.Set(float64(count))
	chatRepository := chatrepository.Get()
	usersRepository := userrepository.Get()

	// Total message count
	cmc := chatRepository.GetMessagesCount()
	// Insert message count into Prometheus collector.
	m.currentChatMessageCount.Set(float64(cmc))

	// Total user count
	uc := usersRepository.GetUsersCount()
	// Insert user count into Prometheus collector.
	m.chatUserCount.Set(float64(uc))

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
func (m *Metrics) GetViewersOverTime(start, end time.Time) []models.TimestampedValue {
	p, err := storage.Select(activeViewerCountKey, nil, start.Unix(), end.Unix())
	if err != nil && err != tstorage.ErrNoDataPoints {
		log.Errorln(err)
	}
	datapoints := models.MakeTimestampedValuesFromDatapoints(p)

	return datapoints
}

// GetChatClientCountOverTime will return a window of connected chat clients over time.
func (m *Metrics) GetChatClientCountOverTime(start, end time.Time) []models.TimestampedValue {
	p, err := storage.Select(activeChatClientCountKey, nil, start.Unix(), end.Unix())
	if err != nil && err != tstorage.ErrNoDataPoints {
		log.Errorln(err)
	}
	datapoints := models.MakeTimestampedValuesFromDatapoints(p)

	return datapoints
}
