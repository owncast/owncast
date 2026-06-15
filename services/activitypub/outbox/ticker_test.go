package outbox

import (
	"testing"
	"time"
)

func TestStartStopPingTicker(t *testing.T) {
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
