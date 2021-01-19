package models

type LatencyLevel struct {
	Level             int `json:"level"`
	SecondsPerSegment int `json:"-"`
	SegmentCount      int `json:"-"`
}

func GetLatencyConfigs() map[int]LatencyLevel {
	return map[int]LatencyLevel{
		1: {Level: 1, SecondsPerSegment: 1, SegmentCount: 2},
		2: {Level: 2, SecondsPerSegment: 2, SegmentCount: 2},
		3: {Level: 3, SecondsPerSegment: 4, SegmentCount: 3},
		4: {Level: 4, SecondsPerSegment: 4, SegmentCount: 4},
		5: {Level: 5, SecondsPerSegment: 8, SegmentCount: 4},
	}
}

func GetLatencyLevel(index int) LatencyLevel {
	return GetLatencyConfigs()[index]
}
