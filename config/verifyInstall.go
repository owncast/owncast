package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/mod/semver"
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

	cmd := exec.Command(path)
	out, err := cmd.CombinedOutput()

	response := string(out)
	if response == "" {
		fmt.Println(err)
		return fmt.Errorf("unable to determine the version of your ffmpeg installation at %s. you may experience issues with video.", path)
	}

	responseComponents := strings.Split(response, " ")
	fullVersionString := responseComponents[2]
	versionString := "v" + strings.Split(fullVersionString, "-")[0]
	if !semver.IsValid(versionString) || semver.Compare(versionString, FfmpegSuggestedVersion) == -1 {
		return fmt.Errorf("your %s version of ffmpeg at %s may be older than the suggested version of %s. you may experience issues with video.", versionString, path, FfmpegSuggestedVersion)
	}

	return nil
}
