package transcoder

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func handleTranscoderMessage(message string) {
	errors := map[string]string{
		"Unrecognized option 'vaapi_device'":    "you are likely trying to utilize a vaapi codec, but your version of ffmpeg or your hardware doesn't support it. change your codec to libx264 and restart your stream",
		"unable to open display":                "your copy of ffmpeg is likely installed via snap packages. please uninstall and re-install via a non-snap method.  https://owncast.online/docs/troubleshooting/#misc-video-issues",
		"Failed to open file 'http://127.0.0.1": "error transcoding. make sure your version of ffmpeg is compatible with your selected codec.",
		"can't configure encoder":               "error with selected codec. if your copy of ffmpeg or your hardware does not support your selected codec you may need to select another",
	}

	for error, displayMessage := range errors {
		if strings.Contains(message, error) {
			log.Error(displayMessage)
			return
		}
	}

	log.Error(message)
}
