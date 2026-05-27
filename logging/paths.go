package logging

import (
	"path/filepath"
)

// GetTranscoderLogFilePath returns the logging path for the transcoder log
// output relative to the provided log directory.
func GetTranscoderLogFilePath(logDirectory string) string {
	return filepath.Join(logDirectory, "transcoder.log")
}

func getLogFilePath(logDirectory string) string {
	return filepath.Join(logDirectory, "owncast.log")
}
