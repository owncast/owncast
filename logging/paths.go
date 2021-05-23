package logging

import (
	"path/filepath"

	"github.com/owncast/owncast/config"
)

func GetTranscoderLogFilePath() string {
	return filepath.Join(config.LogDirectory, "transcoder.log")
}

func getLogFilePath() string {
	return filepath.Join(config.LogDirectory, "owncast.log")
}
