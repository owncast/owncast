package metrics

import (
	"time"

	"github.com/nakabonne/tstorage"
)

// TimestampedValue is a value with a timestamp.
type TimestampedValue struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

func makeTimestampedValuesFromDatapoints(dp []*tstorage.DataPoint) []TimestampedValue {
	tv := []TimestampedValue{}
	for _, d := range dp {
		tv = append(tv, TimestampedValue{Time: time.Unix(d.Timestamp, 0), Value: d.Value})
	}

	return tv
}
