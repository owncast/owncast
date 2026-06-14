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

var errorResetDuration = time.Minute * 5

const alertingError = "The %s utilization of %f%% could cause problems with video generation and delivery. Visit the documentation at http://owncast.online/docs/troubleshooting/ if you are experiencing issues."

func (s *Service) handleAlerting() {
	s.handleCPUAlerting()
	s.handleRAMAlerting()
	s.handleDiskAlerting()
}

func (s *Service) handleCPUAlerting() {
	if len(s.metrics.CPUUtilizations) < 2 {
		return
	}

	avg := recentAverage(s.metrics.CPUUtilizations)
	if avg > maxCPUAlertingThresholdPCT && !s.inCPUAlertingState {
		log.Warnf(alertingError, "CPU", avg)
		s.inCPUAlertingState = true

		resetTimer := time.NewTimer(errorResetDuration)
		go func() {
			<-resetTimer.C
			s.inCPUAlertingState = false
		}()
	}
}

func (s *Service) handleRAMAlerting() {
	if len(s.metrics.RAMUtilizations) < 2 {
		return
	}

	avg := recentAverage(s.metrics.RAMUtilizations)
	if avg > maxRAMAlertingThresholdPCT && !s.inRAMAlertingState {
		log.Warnf(alertingError, "memory", avg)
		s.inRAMAlertingState = true

		resetTimer := time.NewTimer(errorResetDuration)
		go func() {
			<-resetTimer.C
			s.inRAMAlertingState = false
		}()
	}
}

func (s *Service) handleDiskAlerting() {
	if len(s.metrics.DiskUtilizations) < 2 {
		return
	}

	avg := recentAverage(s.metrics.DiskUtilizations)

	if avg > maxDiskAlertingThresholdPCT && !s.inDiskAlertingState {
		log.Warnf(alertingError, "disk", avg)
		s.inDiskAlertingState = true

		resetTimer := time.NewTimer(errorResetDuration)
		go func() {
			<-resetTimer.C
			s.inDiskAlertingState = false
		}()
	}
}

func recentAverage(values []TimestampedValue) float64 {
	return (values[len(values)-1].Value + values[len(values)-2].Value) / 2
}
