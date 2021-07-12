package logging

// Custom logging hooks for powering our logs API.
// Modeled after https://github.com/sirupsen/logrus/blob/master/hooks/test/test.go

import (
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/owncast/owncast/utils"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
)

const maxLogEntries = 500

type OCLogger struct {
	Entries  []logrus.Entry
	Warnings []logrus.Entry
	mu       sync.RWMutex
}

var Logger *OCLogger

// Setup configures our custom logging destinations.
func Setup(enableDebugOptions bool, enableVerboseLogging bool) {
	// Create the logging directory if needed
	loggingDirectory := filepath.Dir(getLogFilePath())
	if !utils.DoesFileExists(loggingDirectory) {
		if err := os.Mkdir(loggingDirectory, 0700); err != nil {
			logger.Errorln("unable to create logs directory", loggingDirectory, err)
		}
	}

	// Write logs to a file
	path := getLogFilePath()
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)

	logMapping := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.DebugLevel: writer,
		logrus.TraceLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
	}

	logger.AddHook(lfshook.NewHook(
		logMapping,
		&logger.TextFormatter{},
	))

	if enableVerboseLogging {
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Write to stdout console
	logger.SetOutput(os.Stdout)

	// Write to our custom logging hook for the log API
	_logger := new(OCLogger)
	logger.AddHook(_logger)

	if enableDebugOptions {
		logrus.SetReportCaller(true)
	}

	Logger = _logger
}

// Fire runs for every logging request.
func (l *OCLogger) Fire(e *logger.Entry) error {
	// Store all log messages to return back in the logging API
	l.mu.Lock()
	defer l.mu.Unlock()

	// Append to log entries
	if len(l.Entries) > maxLogEntries {
		l.Entries = l.Entries[1:]
	}
	l.Entries = append(l.Entries, *e)

	if e.Level <= logger.WarnLevel {
		if len(l.Warnings) > maxLogEntries {
			l.Warnings = l.Warnings[1:]
		}
		l.Warnings = append(l.Warnings, *e)
	}

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
	logCount := int(math.Min(float64(len(l.Entries)), maxLogEntries))
	entries := make([]*logrus.Entry, logCount)
	for i := 0; i < len(entries); i++ {
		// Make a copy, for safety
		entries[len(entries)-logCount:][i] = &l.Entries[i]
	}

	return entries
}

// WarningEntries returns all warning or greater that were logged.
func (l *OCLogger) WarningEntries() []*logrus.Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Make a copy so the returned value won't race with future log requests
	logCount := int(math.Min(float64(len(l.Warnings)), maxLogEntries))
	entries := make([]*logrus.Entry, logCount)
	for i := 0; i < len(entries); i++ {
		// Make a copy, for safety
		entries[len(entries)-logCount:][i] = &l.Warnings[i]
	}

	return entries
}
