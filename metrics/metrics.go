package metrics

import (
	"time"
)

// How often we poll for updates
const metricsPollingInterval = 1 * time.Minute

// CollectedMetrics stores different collected + timestamped values
type CollectedMetrics struct {
	CPUUtilizations  []timestampedValue `json:"cpu"`
	RAMUtilizations  []timestampedValue `json:"memory"`
	DiskUtilizations []timestampedValue `json:"disk"`

	Viewers []timestampedValue `json:"-"`
}

// Metrics is the shared Metrics instance
var Metrics *CollectedMetrics

// Start will begin the metrics collection and alerting
func Start() {
	Metrics = new(CollectedMetrics)
	go startViewerCollectionMetrics()
	handlePolling()

	for range time.Tick(metricsPollingInterval) {
		handlePolling()
	}
}

func handlePolling() {
	// Collect hardware stats
	collectCPUUtilization()
	collectRAMUtilization()
	collectDiskUtilization()

	// Alerting
	handleAlerting()
}
