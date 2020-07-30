package log

import (
	"io"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/termui"
	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
)

// Logger is the termui logger
type Logger struct {
	writer *io.PipeWriter
}

// Setup will setup the termui logger
func (l *Logger) Setup() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if config.Config.EnableTerminalUI {
		l.writer = logger.New().Writer()
		logger.SetOutput(l)
	}
}

// Write will Write to the termui logger instead of to standard console.
func (l *Logger) Write(p []byte) (n int, err error) {
	termui.Log(string(p))
	return 0, nil
}
