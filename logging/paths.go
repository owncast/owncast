package logging

import (
	"path/filepath"

	"github.com/owncast/owncast/config"
)

// GetTranscoderLogFilePath returns the logging path for the transcoder log output.
func GetTranscoderLogFilePath() string {
	return filepath.Join(config.LogDirectory, "transcoder.log")
}

func getLogFilePath() string {
	return filepath.Join(config.LogDirectory, "owncast.log")
}
