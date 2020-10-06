package metrics

import "time"

type timestampedValue struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
}
