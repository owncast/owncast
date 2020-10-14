package utils

import (
	"sort"
	"time"
)

// The "start" timestamp of a timing event
var _pointsInTime = make(map[string]time.Time)

// A collection of timestamp durations for returning the average of
var _durationStorage = make(map[string][]float64)

// StartPerformanceMonitor will keep track of the start time of this event
func StartPerformanceMonitor(key string) {
	if len(_durationStorage[key]) > 30 {
		_durationStorage[key] = removeHighAndLow(_durationStorage[key])
	}
	_pointsInTime[key] = time.Now()
}

// GetAveragePerformance will return the average durations for the event
func GetAveragePerformance(key string) float64 {
	timestamp := _pointsInTime[key]
	if timestamp.IsZero() {
		return 0
	}

	delta := time.Since(timestamp).Seconds()
	_durationStorage[key] = append(_durationStorage[key], delta)
	if len(_durationStorage[key]) < 10 {
		return 0
	}
	_durationStorage[key] = removeHighAndLow(_durationStorage[key])
	return avg(_durationStorage[key])
}

func removeHighAndLow(values []float64) []float64 {
	sort.Float64s(values)
	return values[1 : len(values)-1]
}

func avg(values []float64) float64 {
	total := 0.0
	for _, number := range values {
		total = total + number
	}
	average := total / float64(len(values))
	return average
}
