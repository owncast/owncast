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

	return Avg(_durationStorage[key])
}

func removeHighValue(values []float64) []float64 {
	sort.Float64s(values)
	return values[:len(values)-1]
}

// Avg will return the average value from a slice of float64s.
func Avg(values []float64) float64 {
	total := 0.0
	for _, number := range values {
		total += number
	}
	average := total / float64(len(values))
	return average
}

// Sum returns the sum of a slice of values.
func Sum(values []float64) float64 {
	total := 0.0
	for _, number := range values {
		total += number
	}
	return total
}

// MinMax will return the min and max values from a slice of float64s.
func MinMax(array []float64) (float64, float64) {
	max := array[0]
	min := array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func mean(input []float64) float64 {
	sum := Sum(input)

	return sum / float64(len(input))
}

// Median gets the median number in a slice of numbers.
func Median(input []float64) float64 {
	if len(input) == 1 {
		return input[0]
	}

	c := make([]float64, len(input))
	copy(c, input)

	var median float64
	l := len(c)
	if l%2 == 0 {
		median = mean(c[l/2-1 : l/2+1])
	} else {
		median = c[l/2]
	}

	return median
}
