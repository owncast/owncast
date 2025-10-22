package outbox

import (
	"testing"
	"time"
)

func TestStartStopPingTicker(t *testing.T) {
	// Test that ticker starts
	StartStreamPingTicker()

	// Verify ticker is running
	pingTickerMutex.Lock()
	if pingTicker == nil {
		t.Error("Expected ticker to be running after start")
	}
	pingTickerMutex.Unlock()

	// Test that multiple starts don't create multiple tickers
	StartStreamPingTicker()

	// Stop the ticker
	StopStreamPingTicker()

	// Verify ticker is stopped
	pingTickerMutex.Lock()
	if pingTicker != nil {
		t.Error("Expected ticker to be stopped")
	}
	pingTickerMutex.Unlock()

	// Test that stop on already stopped ticker doesn't panic
	StopStreamPingTicker()
}

func TestPingTickerThreadSafety(t *testing.T) {
	// Test concurrent start/stop calls
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			StartStreamPingTicker()
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			StopStreamPingTicker()
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	<-done
	<-done

	// Cleanup
	StopStreamPingTicker()
}
