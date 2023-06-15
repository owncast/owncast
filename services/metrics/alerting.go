package metrics

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	maxCPUAlertingThresholdPCT  = 85
	maxRAMAlertingThresholdPCT  = 85
	maxDiskAlertingThresholdPCT = 90
)

var (
	inCPUAlertingState  = false
	inRAMAlertingState  = false
	inDiskAlertingState = false
)

var errorResetDuration = time.Minute * 5

const alertingError = "The %s utilization of %f%% could cause problems with video generation and delivery. Visit the documentation at http://owncast.online/docs/troubleshooting/ if you are experiencing issues."

func handleAlerting() {
	handleCPUAlerting()
	handleRAMAlerting()
	handleDiskAlerting()
}

func handleCPUAlerting() {
	if len(metrics.CPUUtilizations) < 2 {
		return
	}

	avg := recentAverage(metrics.CPUUtilizations)
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
	if len(metrics.RAMUtilizations) < 2 {
		return
	}

	avg := recentAverage(metrics.RAMUtilizations)
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
	if len(metrics.DiskUtilizations) < 2 {
		return
	}

	avg := recentAverage(metrics.DiskUtilizations)

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

func recentAverage(values []TimestampedValue) float64 {
	return (values[len(values)-1].Value + values[len(values)-2].Value) / 2
}
