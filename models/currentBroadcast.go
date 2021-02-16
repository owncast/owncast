package models

// CurrentBroadcast represents the configuration associated with the currently active stream.
type CurrentBroadcast struct {
	OutputSettings []StreamOutputVariant `json:"outputSettings"`
	LatencyLevel   LatencyLevel          `json:"latencyLevel"`
}
