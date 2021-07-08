package transcoder

import (
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

var _lastTranscoderLogMessage = ""
var l = &sync.RWMutex{}

var errorMap = map[string]string{
	"Unrecognized option 'vaapi_device'":        "you are likely trying to utilize a vaapi codec, but your version of ffmpeg or your hardware doesn't support it. change your codec to libx264 and restart your stream",
	"unable to open display":                    "your copy of ffmpeg is likely installed via snap packages. please uninstall and re-install via a non-snap method.  https://owncast.online/docs/troubleshooting/#misc-video-issues",
	"Failed to open file 'http://127.0.0.1":     "error transcoding. make sure your version of ffmpeg is compatible with your selected codec or is recent enough https://owncast.online/docs/troubleshooting/#codecs",
	"can't configure encoder":                   "error with codec. if your copy of ffmpeg or your hardware does not support your selected codec you may need to select another",
	"Unable to parse option value":              "you are likely trying to utilize a specific codec, but your version of ffmpeg or your hardware doesn't support it. either fix your ffmpeg install or try changing your codec to libx264 and restart your stream",
	"OpenEncodeSessionEx failed: out of memory": "your NVIDIA gpu is limiting the number of concurrent stream qualities you can support. remove a stream output variant and try again.",
	"Cannot use rename on non file protocol, this may lead to races and temporary partial files": "",
	"No VA display found for device": "vaapi not enabled. either your copy of ffmpeg does not support it, your hardware does not support it, or you need to install additional drivers for your hardware.",
	"Could not find a valid device":  "your codec is either not supported or not configured properly",
	"H.264 bitstream error":          "transcoding content error playback issues may arise. you may want to use the default codec if you are not already.",
	"intel_enc_hw_context_init: Assertion 'encoder_context->mfc_context' failed": "if you are using Intel graphics you may be missing the i965-va-driver-shader drivers",
	`Unknown encoder 'h264_qsv'`:       "your copy of ffmpeg does not have support for Intel QuickSync encoding (h264_qsv). change the selected codec in your video settings",
	`Unknown encoder 'h264_vaapi'`:     "your copy of ffmpeg does not have support for VA-API encoding (h264_vaapi). change the selected codec in your video settings",
	`Unknown encoder 'h264_nvenc'`:     "your copy of ffmpeg does not have support for NVIDIA hardware encoding (h264_nvenc). change the selected codec in your video settings",
	`Unknown encoder 'h264_x264'`:      "your copy of ffmpeg does not have support for the default x264 codec (h264_x264). download a version of ffmpeg that supports this.",
	`Unrecognized option 'x264-params`: "your copy of ffmpeg does not have support for the default libx264 codec (h264_x264). download a version of ffmpeg that supports this.",
	`Failed to set value '/dev/dri/renderD128' for option 'vaapi_device': Invalid argument`: "failed to set va-api device to /dev/dri/renderD128. your system is likely not properly configured for va-api",

	// Generic error for a codec
	"Unrecognized option": "error with codec. if your copy of ffmpeg or your hardware does not support your selected codec you may need to select another",
}

var ignoredErrors = []string{
	"Duplicated segment filename detected",
	"Error while opening encoder for output stream",
	"Unable to parse option value",
	"Last message repeated",
	"Option not found",
	"use of closed network connection",
	"URL read error: End of file",
	"upload playlist failed, will retry with a new http session",
	"VBV underflow",
	"Cannot use rename on non file protocol",
	"Device creation failed",
	"Error parsing global options",
	"maybe the hls segment duration will not precise",
}

func handleTranscoderMessage(message string) {
	log.Debugln(message)

	l.Lock()
	defer l.Unlock()

	// Ignore certain messages that we don't care about.
	for _, error := range ignoredErrors {
		if strings.Contains(message, error) {
			return
		}
	}

	// Convert specific transcoding messages to human-readable messages.
	for error, displayMessage := range errorMap {
		if strings.Contains(message, error) {
			message = displayMessage
			break
		}
	}

	if message == "" {
		return
	}

	// No good comes from a flood of repeated messages.
	if message == _lastTranscoderLogMessage {
		return
	}

	log.Error(message)

	_lastTranscoderLogMessage = message
}

func createVariantDirectories() {
	// Create private hls data dirs
	utils.CleanupDirectory(config.PublicHLSStoragePath)
	utils.CleanupDirectory(config.PrivateHLSStoragePath)
	
	if len(data.GetStreamOutputVariants()) != 0 {
		for index := range data.GetStreamOutputVariants() {
			err := os.MkdirAll(path.Join(config.PrivateHLSStoragePath, strconv.Itoa(index)), 0777)
			if err != nil {
				log.Fatalln(err)
			}
			dir := path.Join(config.PublicHLSStoragePath, strconv.Itoa(index))
			log.Traceln("Creating", dir)
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		dir := path.Join(config.PrivateHLSStoragePath, strconv.Itoa(0))
		log.Traceln("Creating", dir)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatalln(err)
		}
		dir = path.Join(config.PublicHLSStoragePath, strconv.Itoa(0))
		log.Traceln("Creating", dir)
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
