package chat

import (
	"time"
)

const (
	maxBacklogHours = 2 // Keep backlog max hours worth of messages
)

func (s *Service) setupPersistence() {
	chatDataPruner := time.NewTicker(5 * time.Minute)
	go func() {
		s.runPruner()
		for range chatDataPruner.C {
			s.runPruner()
		}
	}()
}
