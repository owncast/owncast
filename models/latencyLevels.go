package models

type LatencyLevel struct {
	Level             int `json:"level"`
	SecondsPerSegment int `json:"-"`
	SegmentCount      int `json:"-"`
}

func GetLatencyConfigs() map[int]LatencyLevel {
	return map[int]LatencyLevel{
		1: {Level: 1, SecondsPerSegment: 2, SegmentCount: 2},
		2: {Level: 2, SecondsPerSegment: 2, SegmentCount: 3},
		3: {Level: 3, SecondsPerSegment: 3, SegmentCount: 3},
		4: {Level: 4, SecondsPerSegment: 4, SegmentCount: 4},
		5: {Level: 5, SecondsPerSegment: 5, SegmentCount: 5},
		6: {Level: 6, SecondsPerSegment: 7, SegmentCount: 10},
	}
}

func GetLatencyLevel(index int) LatencyLevel {
	return GetLatencyConfigs()[index]
}
