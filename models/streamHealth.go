package models

// StreamHealthOverview represents an overview of the current stream health.
type StreamHealthOverview struct {
	Message           string `json:"message"`
	HealthyPercentage int    `json:"healthPercentage"`
	Representation    int    `json:"representation"`
	Healthy           bool   `json:"healthy"`
}
