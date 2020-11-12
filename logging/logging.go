package logging

// Custom logging hooks for powering our logs API.
// Modeled after https://github.com/sirupsen/logrus/blob/master/hooks/test/test.go

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
)

const maxLogEntries = 500

type OCLogger struct {
	Entries []logrus.Entry
	mu      sync.RWMutex
}

var Logger *OCLogger

// Setup configures our custom logging destinations.
func Setup() {
	logger.SetOutput(os.Stdout) // Send all logs to console

	_logger := new(OCLogger)
	logger.AddHook(_logger)

	Logger = _logger
}

// Fire runs for every logging request.
func (l *OCLogger) Fire(e *logger.Entry) error {
	// Store all log messages to return back in the logging API
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.Entries) > maxLogEntries {
		l.Entries = l.Entries[1:]
	}
	l.Entries = append(l.Entries, *e)

	return nil
}

// Levels specifies what log levels we care about.
func (l *OCLogger) Levels() []logrus.Level {
	return logrus.AllLevels
}

// AllEntries returns all entries that were logged.
func (l *OCLogger) AllEntries() []*logrus.Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	// Make a copy so the returned value won't race with future log requests
	entries := make([]*logrus.Entry, len(l.Entries))
	for i := 0; i < len(l.Entries); i++ {
		// Make a copy, for safety
		entries[i] = &l.Entries[i]
	}

	return entries
}

// WarningEntries returns all warning or greater that were logged.
func (l *OCLogger) WarningEntries() []*logrus.Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	// Make a copy so the returned value won't race with future log requests
	entries := make([]*logrus.Entry, 0)
	for i := 0; i < len(l.Entries); i++ {
		if l.Entries[i].Level <= logrus.WarnLevel {
			// Make a copy, for safety
			entries = append(entries, &l.Entries[i])
		}
	}

	return entries
}
