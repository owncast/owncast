package config

import (
	"errors"
	"fmt"
	"os"
)

// verifyFFMpegPath verifies that the path exists, is a file, and is executable.
func verifyFFMpegPath(path string) error {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return errors.New("ffmpeg path does not exist")
	}

	if err != nil {
		return fmt.Errorf("error while verifying the ffmpeg path: %s", err.Error())
	}

	if stat.IsDir() {
		return errors.New("ffmpeg path can not be a folder")
	}

	mode := stat.Mode()
	//source: https://stackoverflow.com/a/60128480
	if mode&0111 == 0 {
		return errors.New("ffmpeg path is not executable")
	}

	return nil
}
