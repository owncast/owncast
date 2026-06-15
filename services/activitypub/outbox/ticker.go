package outbox

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const streamPingInterval = 5 * time.Minute

// StartStreamPingTicker starts a recurring stream-ping Offer activity.
// Safe to call repeatedly; subsequent calls are no-ops while the ticker
// is already running.
func (s *Service) StartStreamPingTicker() {
	s.pingTickerMu.Lock()
	defer s.pingTickerMu.Unlock()

	if s.pingTicker != nil {
		log.Debugln("Stream ping ticker already running")
		return
	}

	t := time.NewTicker(streamPingInterval)
	done := make(chan bool)
	s.pingTicker = t
	s.pingTickerDone = done

	go func() {
		for {
			select {
			case <-t.C:
				if err := s.SendStreamPing(); err != nil {
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
func (s *Service) StopStreamPingTicker() {
	s.pingTickerMu.Lock()
	defer s.pingTickerMu.Unlock()

	if s.pingTicker != nil {
		s.pingTicker.Stop()
		close(s.pingTickerDone)
		s.pingTicker = nil
		s.pingTickerDone = nil
		log.Infoln("Stopped stream ping ticker")
	}
}
