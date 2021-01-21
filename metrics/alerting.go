package metrics

import (
	log "github.com/sirupsen/logrus"
)

const maxCPUAlertingThresholdPCT = 80
const maxRAMAlertingThresholdPCT = 80
const maxDiskAlertingThresholdPCT = 90

const alertingError = "The %s utilization of %d%% can cause issues with video generation and delivery. Please visit the documentation at http://owncast.online/docs/troubleshooting/ to help troubleshoot this issue."

func handleAlerting() {
	handleCPUAlerting()
	handleRAMAlerting()
	handleDiskAlerting()
}

func handleCPUAlerting() {
	if len(Metrics.CPUUtilizations) < 2 {
		return
	}

	avg := recentAverage(Metrics.CPUUtilizations)
	if avg > maxCPUAlertingThresholdPCT {
		log.Warnf(alertingError, "CPU", maxCPUAlertingThresholdPCT)
	}
}

func handleRAMAlerting() {
	if len(Metrics.RAMUtilizations) < 2 {
		return
	}

	avg := recentAverage(Metrics.RAMUtilizations)
	if avg > maxRAMAlertingThresholdPCT {
		log.Warnf(alertingError, "memory", maxRAMAlertingThresholdPCT)
	}
}

func handleDiskAlerting() {
	if len(Metrics.DiskUtilizations) < 2 {
		return
	}

	avg := recentAverage(Metrics.DiskUtilizations)

	if avg > maxDiskAlertingThresholdPCT {
		log.Warnf(alertingError, "disk", maxRAMAlertingThresholdPCT)
	}
}

func recentAverage(values []timestampedValue) int {
	return (values[len(values)-1].Value + values[len(values)-2].Value) / 2
}
