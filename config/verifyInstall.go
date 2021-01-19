package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

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
	out, _ := cmd.CombinedOutput()

	response := string(out)
	if response == "" {
		return fmt.Errorf("unable to determine the version of your ffmpeg installation at %s. you may experience issues with video.", path)
	}

	responseComponents := strings.Split(response, " ")
	if len(responseComponents) < 3 {
		log.Debugf("unable to determine the version of your ffmpeg installation at %s. you may experience issues with video.", path)
		return nil
	}

	fullVersionString := responseComponents[2]

	versionString := "v" + strings.Split(fullVersionString, "-")[0]

	// Some builds of ffmpeg have wierd build numbers that we can't parse
	if !semver.IsValid(versionString) {
		log.Debugf("unable to determine if ffmpeg version %s is recent enough. if you experience issues with video you may want to look into updating", fullVersionString)
		return nil
	}

	if semver.Compare(versionString, FfmpegSuggestedVersion) == -1 {
		return fmt.Errorf("your %s version of ffmpeg at %s may be older than the suggested version of %s. you may experience issues with video.", versionString, path, FfmpegSuggestedVersion)
	}

	return nil
}
