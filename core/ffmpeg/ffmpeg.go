package ffmpeg

import (
	"github.com/gabek/owncast/config"
)

//ShowStreamOfflineState generates and shows the stream's offline state
func ShowStreamOfflineState() error {
	transcoder := CreateTranscoderFromConfig()
	transcoder.SetInput(config.Config.VideoSettings.OfflineContent)
	return transcoder.Start()
}
