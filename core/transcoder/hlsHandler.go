package transcoder

import (
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/replays"
	log "github.com/sirupsen/logrus"
)

// HLSHandler gets told about available HLS playlists and segments.
type HLSHandler struct {
	Storage  models.StorageProvider
	Recorder *replays.HLSRecorder
}

// StreamEnded is called when a stream is ended so the end time can be noted
// in the stream's metadata.
func (h *HLSHandler) StreamEnded() {
	if config.EnableReplayFeatures {
		h.Recorder.StreamEnded()
	}
}

func (h *HLSHandler) SetStreamId(streamId string) {
	h.Storage.SetStreamId(streamId)
	if config.EnableReplayFeatures {
		h.Recorder = replays.NewRecording(streamId)
	}
}

// SegmentWritten is fired when a HLS segment is written to disk.
func (h *HLSHandler) SegmentWritten(localFilePath string) {
	remotePath, _, err := h.Storage.SegmentWritten(localFilePath)
	if err != nil {
		log.Errorln(err)
		return
	}

	if h.Recorder != nil {
		h.Recorder.SegmentWritten(remotePath)
	} else {
		log.Debugln("No HLS recorder available to notify of segment written.")
	}
}

// VariantPlaylistWritten is fired when a HLS variant playlist is written to disk.
func (h *HLSHandler) VariantPlaylistWritten(localFilePath string) {
	h.Storage.VariantPlaylistWritten(localFilePath)
}

// MasterPlaylistWritten is fired when a HLS master playlist is written to disk.
func (h *HLSHandler) MasterPlaylistWritten(localFilePath string) {
	h.Storage.MasterPlaylistWritten(localFilePath)
}
