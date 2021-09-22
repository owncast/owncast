package metrics

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const maxCPUAlertingThresholdPCT = 85
const maxRAMAlertingThresholdPCT = 85
const maxDiskAlertingThresholdPCT = 90

var inCPUAlertingState = false
var inRAMAlertingState = false
var inDiskAlertingState = false

var errorResetDuration = time.Minute * 5

const alertingError = "The %s utilization of %d%% could cause problems with video generation and delivery. Visit the documentation at http://owncast.online/docs/troubleshooting/ if you are experiencing issues."

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
	if avg > maxCPUAlertingThresholdPCT && !inCPUAlertingState {
		log.Warnf(alertingError, "CPU", avg)
		inCPUAlertingState = true

		resetTimer := time.NewTimer(errorResetDuration)
		go func() {
			<-resetTimer.C
			inCPUAlertingState = false
		}()
	}
}

func handleRAMAlerting() {
	if len(Metrics.RAMUtilizations) < 2 {
		return
	}

	avg := recentAverage(Metrics.RAMUtilizations)
	if avg > maxRAMAlertingThresholdPCT && !inRAMAlertingState {
		log.Warnf(alertingError, "memory", avg)
		inRAMAlertingState = true

		resetTimer := time.NewTimer(errorResetDuration)
		go func() {
			<-resetTimer.C
			inRAMAlertingState = false
		}()
	}
}

func handleDiskAlerting() {
	if len(Metrics.DiskUtilizations) < 2 {
		return
	}

	avg := recentAverage(Metrics.DiskUtilizations)

	if avg > maxDiskAlertingThresholdPCT && !inDiskAlertingState {
		log.Warnf(alertingError, "disk", avg)
		inDiskAlertingState = true

		resetTimer := time.NewTimer(errorResetDuration)
		go func() {
			<-resetTimer.C
			inDiskAlertingState = false
		}()
	}
}

func recentAverage(values []timestampedValue) int {
	return (values[len(values)-1].Value + values[len(values)-2].Value) / 2
}
