package ffmpeg

import (
	"github.com/gabek/owncast/config"
)

//ShowStreamOfflineState generates and shows the stream's offline state
func ShowStreamOfflineState() {
	transcoder := NewTranscoder()
	transcoder.SetInput(config.Config.VideoSettings.OfflineContent)
	transcoder.Start()
}
