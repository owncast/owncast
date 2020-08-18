package metrics

import (
	"time"
)

// How often we poll for updates
const metricsPollingInterval = 15 * time.Second

type value struct {
	Time  time.Time
	Value int
}

type metrics struct {
	CPUUtilizations []value
	RAMUtilizations []value
}

// Metrics is the shared Metrics instance
var Metrics *metrics

// Start will begin the metrics collection and alerting
func Start() {
	Metrics = new(metrics)

	for range time.Tick(metricsPollingInterval) {
		handlePolling()
	}
}

func handlePolling() {
	// Collect hardware stats
	collectCPUUtilization()
	collectRAMUtilization()

	// Alerting
	handleAlerting()
}
