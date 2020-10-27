package config

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getDefaults() config {
	defaults := config{}
	defaults.WebServerPort = 8080
	defaults.FFMpegPath = getDefaultFFMpegPath()
	defaults.VideoSettings.ChunkLengthInSeconds = 4
	defaults.Files.MaxNumberInPlaylist = 5
	defaults.YP.Enabled = false
	defaults.YP.YPServiceURL = "https://yp.owncast.online"
	defaults.DatabaseFilePath = "data/owncast.db"

	defaultQuality := StreamQuality{
		IsAudioPassthrough: true,
		VideoBitrate:       1200,
		EncoderPreset:      "veryfast",
		Framerate:          24,
	}
	defaults.VideoSettings.StreamQualities = []StreamQuality{defaultQuality}

	return defaults
}

func getDefaultFFMpegPath() string {
	// First look to see if ffmpeg is in the current working directory
	localCopy := "./ffmpeg"
	hasLocalCopyError := verifyFFMpegPath(localCopy)
	if hasLocalCopyError == nil {
		// No error, so all is good.  Use the local copy.
		return localCopy
	}

	cmd := exec.Command("which", "ffmpeg")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Debugln("Unable to determine path to ffmpeg. Please specify it in the config file.")
	}

	path := strings.TrimSpace(string(out))

	return path
}
