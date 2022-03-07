package metrics

import (
	"time"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
)

// How often we poll for updates.
const viewerMetricsPollingInterval = 2 * time.Minute

func startViewerCollectionMetrics() {
	collectViewerCount()

	for range time.Tick(viewerMetricsPollingInterval) {
		collectViewerCount()
	}
}

func collectViewerCount() {
	if len(Metrics.Viewers) > maxCollectionValues {
		Metrics.Viewers = Metrics.Viewers[1:]
	}

	count := _getStatus().ViewerCount
	value := timestampedValue{
		Value: count,
		Time:  time.Now(),
	}
	Metrics.Viewers = append(Metrics.Viewers, value)

	// Save to our Prometheus collector.
	activeViewerCount.Set(float64(count))

	// Total message count
	cmc := data.GetMessagesCount()
	currentChatMessageCount.Set(float64(cmc))

	// Total user count
	uc := data.GetUsersCount()
	chatUserCount.Set(float64(uc))
}

func collectChatClientCount() {
	count := len(chat.GetClients())
	activeChatClientCount.Set(float64(count))
}
