package metrics

import (
	"time"

	"github.com/nakabonne/tstorage"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

var storage tstorage.Storage

func startViewerCollectionMetrics() {
	storage, _ = tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithDataPath("./data/metrics"),
	)
	defer storage.Close()

	collectViewerCount()

	for range time.Tick(viewerMetricsPollingInterval) {
		collectViewerCount()
		collectChatClientCount()
	}
}

func collectViewerCount() {
	// Don't collect metrics for viewers if there's no stream active.
	if !core.GetStatus().Online {
		return
	}

	count := core.GetStatus().ViewerCount

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

func collectChatClientCount() {
	count := len(chat.GetClients())
	activeChatClientCount.Set(float64(count))

	// Total message count
	cmc := data.GetMessagesCount()
	// Insert message count into Prometheus collector.
	currentChatMessageCount.Set(float64(cmc))

	// Total user count
	uc := data.GetUsersCount()
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
func GetViewersOverTime(start, end time.Time) []TimestampedValue {
	p, err := storage.Select(activeViewerCountKey, nil, start.Unix(), end.Unix())
	if err != nil && err != tstorage.ErrNoDataPoints {
		log.Errorln(err)
	}
	datapoints := makeTimestampedValuesFromDatapoints(p)

	return datapoints
}

// GetChatClientCountOverTime will return a window of connected chat clients over time.
func GetChatClientCountOverTime(start, end time.Time) []TimestampedValue {
	p, err := storage.Select(activeChatClientCountKey, nil, start.Unix(), end.Unix())
	if err != nil && err != tstorage.ErrNoDataPoints {
		log.Errorln(err)
	}
	datapoints := makeTimestampedValuesFromDatapoints(p)

	return datapoints
}
