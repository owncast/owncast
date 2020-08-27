package metrics

import (
	log "github.com/sirupsen/logrus"
)

const maxCPUAlertingThresholdPCT = 95
const maxRAMAlertingThresholdPCT = 95

const alertingError = "The %s utilization of %d%% is higher than the alerting threshold of %d%%.  This can cause issues with video generation and delivery. Please visit the documentation at http://owncast.online/docs/troubleshooting/ to help troubleshoot this issue."

func handleAlerting() {
	handleCPUAlerting()
	handleRAMAlerting()
}

func handleCPUAlerting() {
	if len(Metrics.CPUUtilizations) < 2 {
		return
	}

	avg := recentAverage(Metrics.CPUUtilizations)
	if avg > maxCPUAlertingThresholdPCT {
		log.Errorf(alertingError, "CPU", avg, maxCPUAlertingThresholdPCT)
	}
}

func handleRAMAlerting() {
	if len(Metrics.RAMUtilizations) < 2 {
		return
	}

	avg := recentAverage(Metrics.RAMUtilizations)
	if avg > maxRAMAlertingThresholdPCT {
		log.Errorf(alertingError, "memory", avg, maxRAMAlertingThresholdPCT)
	}
}

func recentAverage(values []value) int {
	return int((values[len(values)-1].Value + values[len(values)-2].Value) / 2)
}
