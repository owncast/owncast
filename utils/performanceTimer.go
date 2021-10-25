package utils

import (
	"sort"
	"sync"
	"time"
)

var l = sync.Mutex{}

// The "start" timestamp of a timing event.
var _pointsInTime = make(map[string]time.Time)

// A collection of timestamp durations for returning the average of.
var _durationStorage = make(map[string][]float64)

// StartPerformanceMonitor will keep track of the start time of this event.
func StartPerformanceMonitor(key string) {
	l.Lock()
	if len(_durationStorage[key]) > 20 {
		_durationStorage[key] = removeHighValue(_durationStorage[key])
	}
	_pointsInTime[key] = time.Now()
	l.Unlock()
}

// GetAveragePerformance will return the average durations for the event.
func GetAveragePerformance(key string) float64 {
	timestamp := _pointsInTime[key]
	if timestamp.IsZero() {
		return 0
	}

	l.Lock()
	defer l.Unlock()

	delta := time.Since(timestamp).Seconds()
	_durationStorage[key] = append(_durationStorage[key], delta)
	if len(_durationStorage[key]) < 8 {
		return 0
	}
	_durationStorage[key] = removeHighValue(_durationStorage[key])

	return avg(_durationStorage[key])
}

func removeHighValue(values []float64) []float64 {
	sort.Float64s(values)
	return values[:len(values)-1]
}

func avg(values []float64) float64 {
	total := 0.0
	for _, number := range values {
		total += number
	}
	average := total / float64(len(values))
	return average
}
