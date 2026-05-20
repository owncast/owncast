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
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/utils"
)

const maxLogEntries = 500

// OCLogger represents the owncast internal logging.
type OCLogger struct {
	Entries  []log.Entry
	Warnings []log.Entry
	mu       sync.RWMutex
}

// Logger is the shared instance of the internal logger.
var Logger *OCLogger

// Setup configures our custom logging destinations.
func Setup(logDirectory string, enableDebugOptions bool, enableVerboseLogging bool) {
	// Create the logging directory if needed
	loggingDirectory := filepath.Dir(getLogFilePath(logDirectory))
	if !utils.DoesFileExists(loggingDirectory) {
		if err := os.Mkdir(loggingDirectory, 0o700); err != nil {
			log.Errorln("unable to create logs directory", loggingDirectory, err)
		}
	}

	// Write logs to a file
	path := getLogFilePath(logDirectory)
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)

	logMapping := lfshook.WriterMap{
		log.InfoLevel:  writer,
		log.DebugLevel: writer,
		log.TraceLevel: writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
	}

	log.AddHook(lfshook.NewHook(
		logMapping,
		&log.TextFormatter{},
	))

	if enableVerboseLogging {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Write to stdout console
	log.SetOutput(os.Stdout)

	// Write to our custom logging hook for the log API
	_logger := new(OCLogger)
	log.AddHook(_logger)

	if enableDebugOptions {
		log.SetReportCaller(true)
	}

	Logger = _logger
}

// Fire runs for every logging request.
func (l *OCLogger) Fire(e *log.Entry) error {
	// Store all log messages to return back in the logging API
	l.mu.Lock()
	defer l.mu.Unlock()

	// Append to log entries
	if len(l.Entries) > maxLogEntries {
		l.Entries = l.Entries[1:]
	}
	l.Entries = append(l.Entries, *e)

	if e.Level <= log.WarnLevel {
		if len(l.Warnings) > maxLogEntries {
			l.Warnings = l.Warnings[1:]
		}
		l.Warnings = append(l.Warnings, *e)
	}

	return nil
}

// Levels specifies what log levels we care about.
func (l *OCLogger) Levels() []log.Level {
	return log.AllLevels
}

// AllEntries returns all entries that were logged.
func (l *OCLogger) AllEntries() []*log.Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Make a copy so the returned value won't race with future log requests
	logCount := int(math.Min(float64(len(l.Entries)), maxLogEntries))
	entries := make([]*log.Entry, logCount)
	for i := 0; i < len(entries); i++ {
		// Make a copy, for safety
		entries[len(entries)-logCount:][i] = &l.Entries[i]
	}

	return entries
}

// WarningEntries returns all warning or greater that were logged.
func (l *OCLogger) WarningEntries() []*log.Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Make a copy so the returned value won't race with future log requests
	logCount := int(math.Min(float64(len(l.Warnings)), maxLogEntries))
	entries := make([]*log.Entry, logCount)
	for i := 0; i < len(entries); i++ {
		// Make a copy, for safety
		entries[len(entries)-logCount:][i] = &l.Warnings[i]
	}

	return entries
}
