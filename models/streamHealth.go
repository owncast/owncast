package models

// StreamHealthOverview represents an overview of the current stream health.
type StreamHealthOverview struct {
	Healthy           bool   `json:"healthy"`
	HealthyPercentage int    `json:"healthPercentage"`
	Message           string `json:"message"`
	Representation    int    `json:"representation"`
}
