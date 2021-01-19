package models

type CurrentBroadcast struct {
	OutputSettings []StreamOutputVariant `json:"outputSettings"`
	LatencyLevel   LatencyLevel          `json:"latencyLevel"`
}
