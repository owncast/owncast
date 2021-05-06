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
		0: {Level: 0, SecondsPerSegment: 1, SegmentCount: 3}, // Approx 5 seconds
		1: {Level: 1, SecondsPerSegment: 2, SegmentCount: 2}, // Approx 7-8 seconds
		2: {Level: 2, SecondsPerSegment: 3, SegmentCount: 3}, // Default Approx 11 seconds
		3: {Level: 3, SecondsPerSegment: 4, SegmentCount: 3}, // Approx 15 seconds
		4: {Level: 4, SecondsPerSegment: 5, SegmentCount: 4}, // Approx 18 seconds
	}
}

// GetLatencyLevel will return the latency level at index.
func GetLatencyLevel(index int) LatencyLevel {
	return GetLatencyConfigs()[index]
}
