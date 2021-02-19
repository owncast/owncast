package models

// LatencyLevel is a representation of HLS configuration values.
type LatencyLevel struct {
	Level             int `json:"level"`
	SecondsPerSegment int `json:"-"`
	SegmentCount      int `json:"-"`
}

// GetLatencyConfigs will return the available latency level options.
func GetLatencyConfigs() map[int]LatencyLevel {
	return map[int]LatencyLevel{
		1: {Level: 1, SecondsPerSegment: 1, SegmentCount: 2},
		2: {Level: 2, SecondsPerSegment: 2, SegmentCount: 2},
		3: {Level: 3, SecondsPerSegment: 3, SegmentCount: 3},
		4: {Level: 4, SecondsPerSegment: 3, SegmentCount: 4}, // Default
		5: {Level: 5, SecondsPerSegment: 4, SegmentCount: 5},
		6: {Level: 6, SecondsPerSegment: 6, SegmentCount: 10},
	}
}

// GetLatencyLevel will return the latency level at index.
func GetLatencyLevel(index int) LatencyLevel {
	return GetLatencyConfigs()[index]
}
