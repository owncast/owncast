package metrics

import (
	"time"

	"github.com/owncast/owncast/core"
)

// How often we poll for updates.
const viewerMetricsPollingInterval = 5 * time.Minute

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

	count := core.GetStatus().ViewerCount
	value := timestampedValue{
		Value: count,
		Time:  time.Now(),
	}
	Metrics.Viewers = append(Metrics.Viewers, value)
}
