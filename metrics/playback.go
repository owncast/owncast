package metrics

import (
	"time"

	"github.com/nakabonne/tstorage"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// Playback error counts reported since the last time we collected metrics.
var windowedErrorCounts = []float64{}

// IncrementPlaybackCount will add to the windowed playback error count.
func IncrementPlaybackCount(count float64) {
	windowedErrorCounts = append(windowedErrorCounts, count)
}

// collectPlaybackErrorCount will take all of the error counts each individual
// player reported and average them into a single metric. This is done so
// one person with bad connectivity doesn't make it look like everything is
// horrible for everyone.
func collectPlaybackErrorCount() {
	count := utils.Avg(windowedErrorCounts)
	windowedErrorCounts = []float64{}

	// Save to Prometheus collector.
	playbackErrorCount.Set(count)

	// Save to on disk time series storage.
	if err := storage.InsertRows([]tstorage.Row{
		{
			Metric:    playbackErrorCountKey,
			DataPoint: tstorage.DataPoint{Timestamp: time.Now().Unix(), Value: count},
		},
	}); err != nil {
		log.Errorln(err)
	}
}

// GetPlaybackErrorCountOverTime will return a window of playback errors over time.
func GetPlaybackErrorCountOverTime(start, end time.Time) []timestampedValue {
	p, err := storage.Select(playbackErrorCountKey, nil, start.Unix(), end.Unix())
	if err != nil && err != tstorage.ErrNoDataPoints {
		log.Errorln(err)
	}
	datapoints := makeTimestampedValuesFromDatapoints(p)

	return datapoints
}
