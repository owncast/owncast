package transcoder

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
)

// Recorder is notified of recorded HLS activity so a stream can be replayed
// later. The replay subsystem implements it; it is optional and only set when
// replay features are enabled. Defined here (rather than imported from the
// replays package) to keep the transcoder free of a dependency on it.
type Recorder interface {
	// SegmentWritten is called with the public path of each stored segment
	// and its size in bytes.
	SegmentWritten(path string, size int64)
	// VariantPlaylistWritten is called with the local path of each written
	// variant playlist, which carries the real EXTINF segment durations.
	VariantPlaylistWritten(localFilePath string)
	// StreamEnded is called when the live stream ends.
	StreamEnded()
}

// HLSHandler gets told about available HLS playlists and segments.
type HLSHandler struct {
	Storage models.StorageProvider

	// Recorder, when non-nil, records each written segment for later replay.
	Recorder Recorder
}

// StreamEnded notes the end of the stream in the recorder, if one is set.
func (h *HLSHandler) StreamEnded() {
	if h.Recorder != nil {
		h.Recorder.StreamEnded()
	}
}

// SegmentWritten is fired when a HLS segment is written to disk.
func (h *HLSHandler) SegmentWritten(localFilePath string) {
	// Capture the size before storage providers get a chance to move or
	// clean up the local file.
	var size int64
	if info, err := os.Stat(localFilePath); err == nil {
		size = info.Size()
	}

	remotePath, err := h.Storage.SegmentWritten(localFilePath)
	if err != nil {
		log.Debugln(err, localFilePath)
		return
	}

	if h.Recorder != nil {
		h.Recorder.SegmentWritten(remotePath, size)
	}
}

// VariantPlaylistWritten is fired when a HLS variant playlist is written to disk.
func (h *HLSHandler) VariantPlaylistWritten(localFilePath string) {
	h.Storage.VariantPlaylistWritten(localFilePath)

	if h.Recorder != nil {
		h.Recorder.VariantPlaylistWritten(localFilePath)
	}
}

// MasterPlaylistWritten is fired when a HLS master playlist is written to disk.
func (h *HLSHandler) MasterPlaylistWritten(localFilePath string) {
	h.Storage.MasterPlaylistWritten(localFilePath)
}
