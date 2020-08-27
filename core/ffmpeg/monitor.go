package ffmpeg

import (
	"github.com/gabek/owncast/models"
)

type HLSHandler struct {
	Storage models.StorageProvider
}

func (h *HLSHandler) SegmentWritten(localFilePath string) {
	h.Storage.SegmentWritten(localFilePath)
}

func (h *HLSHandler) VariantPlaylistWritten(localFilePath string) {
	h.Storage.VariantPlaylistWritten(localFilePath)
}

func (h *HLSHandler) MasterPlaylistWritten(localFilePath string) {
	h.Storage.MasterPlaylistWritten(localFilePath)
}
