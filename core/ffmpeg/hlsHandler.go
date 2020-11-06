package ffmpeg

import (
	"context"
	"path/filepath"
	"time"

	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// HLSHandler gets told about available HLS playlists and segments
type HLSHandler struct {
	Storage models.StorageProvider
	Monitor *hlsVariantWriteMonitor
}

// NewHLSHandler returns an initialized HLSHandler with monitor running.
func NewHLSHandler(storage models.StorageProvider, segmentLength int) *HLSHandler {
	monitor := newHlsVariantWriteMonitor(
		context.TODO(),
		log.StandardLogger(),
		time.Duration(segmentLength*3)*time.Second,
	)

	return &HLSHandler{
		Storage: storage,
		Monitor: monitor,
	}
}

// SegmentWritten is fired when a HLS segment is written to disk
func (h *HLSHandler) SegmentWritten(localFilePath string) {
	h.Storage.SegmentWritten(localFilePath)
	h.Monitor.SegmentWritten(filepath.Dir(localFilePath), time.Now())
}

// VariantPlaylistWritten is fired when a HLS variant playlist is written to disk
func (h *HLSHandler) VariantPlaylistWritten(localFilePath string) {
	h.Storage.VariantPlaylistWritten(localFilePath)
	h.Monitor.VariantPlaylistWritten(filepath.Dir(localFilePath), time.Now())
}

// MasterPlaylistWritten is fired when a HLS master playlist is written to disk
func (h *HLSHandler) MasterPlaylistWritten(localFilePath string) {
	h.Storage.MasterPlaylistWritten(localFilePath)
}
