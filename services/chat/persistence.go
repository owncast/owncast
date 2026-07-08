package chat

import (
	"context"
	"time"
)

const (
	maxBacklogHours = 2 // Keep backlog max hours worth of messages

	// DataPruneInterval is how often RunDataPruner should run, registered
	// on the central scheduler by main.go.
	DataPruneInterval = 5 * time.Minute
)

// RunDataPruner is the scheduler job body that trims chat history down to
// the backlog window.
func (s *Service) RunDataPruner(_ context.Context) {
	s.runPruner()
}
