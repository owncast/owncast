package metrics

import (
	"time"

	"github.com/nakabonne/tstorage"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

// How often we poll for updates.
const viewerMetricsPollingInterval = 2 * time.Minute

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

	// Save to our Prometheus collector.
	activeViewerCount.Set(float64(core.GetStatus().ViewerCount))
}

func collectChatClientCount() {
	count := len(chat.GetClients())
	activeChatClientCount.Set(float64(count))

	// Total message count
	cmc := data.GetMessagesCount()
	currentChatMessageCount.Set(float64(cmc))

	// Total user count
	uc := data.GetUsersCount()
	chatUserCount.Set(float64(uc))

	if err := storage.InsertRows([]tstorage.Row{
		{
			Metric:    "viewercount",
			DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: float64(count)},
		},
	}); err != nil {
		log.Errorln(err)
	}
}

// GetViewersOverTime will return a window of viewer counts over time.
func GetViewersOverTime(start, end time.Time) []timestampedValue {
	p, err := storage.Select("viewercount", nil, start.Unix(), end.Unix())
	if err != nil && err != tstorage.ErrNoDataPoints {
		log.Errorln(err)
	}
	datapoints := makeTimestampedValuesFromDatapoints(p)

	return datapoints
}
