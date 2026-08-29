package outbox

import (
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
)

func TestStartStopPingTicker(t *testing.T) {
	originalLevel := log.GetLevel()
	log.SetLevel(log.DebugLevel)
	defer log.SetLevel(originalLevel)

	hook := logtest.NewGlobal()
	defer hook.Reset()

	s := &Service{}

	// Test that ticker starts
	s.StartStreamPingTicker()

	// Verify ticker is running
	s.pingTickerMu.Lock()
	if s.pingTicker == nil {
		t.Error("Expected ticker to be running after start")
	}
	s.pingTickerMu.Unlock()

	// Test that multiple starts don't create multiple tickers
	s.StartStreamPingTicker()

	// Stop the ticker
	s.StopStreamPingTicker()

	// Verify ticker is stopped
	s.pingTickerMu.Lock()
	if s.pingTicker != nil {
		t.Error("Expected ticker to be stopped")
	}
	s.pingTickerMu.Unlock()

	// Test that stop on already stopped ticker doesn't panic
	s.StopStreamPingTicker()

	var tickerLogs int
	for _, entry := range hook.AllEntries() {
		if strings.Contains(strings.ToLower(entry.Message), "stream ping ticker") {
			tickerLogs++
			if entry.Level != log.DebugLevel {
				t.Errorf("ticker log %q level = %s, want debug", entry.Message, entry.Level)
			}
		}
	}
	if tickerLogs != 3 {
		t.Errorf("ticker log count = %d, want 3", tickerLogs)
	}
}

func TestPingTickerThreadSafety(t *testing.T) {
	s := &Service{}

	// Test concurrent start/stop calls
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			s.StartStreamPingTicker()
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			s.StopStreamPingTicker()
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	<-done
	<-done

	// Cleanup
	s.StopStreamPingTicker()
}
