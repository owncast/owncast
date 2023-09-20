package core

import (
	"io"

	"github.com/owncast/owncast/core/transcoder"
)

func setupVideoComponentsForId(streamId string) {
}

func setupLiveTranscoderForId(streamId string, rtmpOut *io.PipeReader) {
	_storage.SetStreamId(streamId)
	handler.SetStreamId(streamId)

	go func() {
		_transcoder = transcoder.NewTranscoder(streamId)
		_transcoder.TranscoderCompleted = func(error) {
			SetStreamAsDisconnected()
			_transcoder = nil
			_currentBroadcast = nil
		}
		_transcoder.SetStdin(rtmpOut)
		_transcoder.Start(true)
	}()
}
