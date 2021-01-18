package models

type LatancyLevel struct {
	Level             int
	SecondsPerSegment int
	SegmentCount      int
}

func GetLatancyConfigs() map[int]LatancyLevel {
	return map[int]LatancyLevel{
		1: {Level: 1, SecondsPerSegment: 1, SegmentCount: 2},
		2: {Level: 2, SecondsPerSegment: 2, SegmentCount: 2},
		3: {Level: 3, SecondsPerSegment: 4, SegmentCount: 3},
		4: {Level: 4, SecondsPerSegment: 4, SegmentCount: 4},
		5: {Level: 5, SecondsPerSegment: 8, SegmentCount: 4},
	}
}

func GetLatancyLevel(index int) LatancyLevel {
	return GetLatancyConfigs()[index]
}
