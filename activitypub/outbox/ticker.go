package outbox

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	pingTicker      *time.Ticker
	pingTickerDone  chan bool
	pingTickerMutex sync.Mutex
)

// StartStreamPingTicker starts a ticker that sends stream ping Offer activities every 5 minutes.
func StartStreamPingTicker() {
	pingTickerMutex.Lock()
	defer pingTickerMutex.Unlock()

	// Don't start if already running
	if pingTicker != nil {
		log.Debugln("Stream ping ticker already running")
		return
	}

	pingTicker = time.NewTicker(5 * time.Minute)
	pingTickerDone = make(chan bool)

	// Capture the done channel in a local variable to avoid race conditions
	done := pingTickerDone
	ticker := pingTicker

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := SendStreamPing(); err != nil {
					log.Errorf("Error sending stream ping: %v", err)
				}
			case <-done:
				return
			}
		}
	}()

	log.Infoln("Started stream ping ticker (5 minute interval)")
}

// StopStreamPingTicker stops the stream ping ticker if it is running.
func StopStreamPingTicker() {
	pingTickerMutex.Lock()
	defer pingTickerMutex.Unlock()

	if pingTicker != nil {
		pingTicker.Stop()
		close(pingTickerDone)
		pingTicker = nil
		pingTickerDone = nil
		log.Infoln("Stopped stream ping ticker")
	}
}
