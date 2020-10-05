package ffmpeg

import (
	"github.com/owncast/owncast/config"
)

//ShowStreamOfflineState generates and shows the stream's offline state
func ShowStreamOfflineState() {
	transcoder := NewTranscoder()
	transcoder.SetSegmentLength(10)
	transcoder.SetAppendToStream(true)
	transcoder.SetInput(config.Config.GetOfflineContentPath())
	transcoder.Start()
}
