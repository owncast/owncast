package logging

import (
	"path/filepath"

	"github.com/owncast/owncast/services/config"
)

// GetTranscoderLogFilePath returns the logging path for the transcoder log output.
func GetTranscoderLogFilePath() string {
	c := config.GetConfig()
	return filepath.Join(c.LogDirectory, "transcoder.log")
}

func getLogFilePath() string {
	c := config.GetConfig()
	return filepath.Join(c.LogDirectory, "owncast.log")
}
