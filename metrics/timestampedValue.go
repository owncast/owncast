package metrics

import (
	"time"

	"github.com/nakabonne/tstorage"
)

type timestampedValue struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
}

func makeTimestampedValuesFromDatapoints(dp []*tstorage.DataPoint) []timestampedValue {
	tv := []timestampedValue{}
	for _, d := range dp {
		tv = append(tv, timestampedValue{Time: time.Unix(d.Timestamp, 0), Value: int(d.Value)})
	}

	return tv
}
