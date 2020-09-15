package config

import (
	"log"
	"os/exec"
	"strings"
)

func getDefaults() config {
	defaults := config{}
	defaults.WebServerPort = 8080
	defaults.FFMpegPath = getDefaultFFMpegPath()
	defaults.VideoSettings.ChunkLengthInSeconds = 4
	defaults.Files.MaxNumberInPlaylist = 5
	defaults.PublicHLSPath = "webroot/hls"
	defaults.PrivateHLSPath = "hls"
	defaults.VideoSettings.OfflineContent = "static/offline.m4v"
	defaults.InstanceDetails.ExtraInfoFile = "/static/content.md"
	defaults.YP.Enabled = false
	defaults.YP.YPServiceURL = "https://owncast-yp-test.gabek.vercel.app/home" // TODO: CHANGE

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
	cmd := exec.Command("which", "ffmpeg")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicln("Unable to determine path to ffmpeg. Please specify it in the config file.")
	}

	path := strings.TrimSpace(string(out))

	return path
}
