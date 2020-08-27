package ffmpeg

import (
	"fmt"
	"path/filepath"

	"github.com/gabek/owncast/models"
)

type HLSHandler struct {
	Storage models.StorageProvider
}

func (h *HLSHandler) SegmentWritten(localFilePath string) {
	go func() {
		url, error := h.Storage.Save(localFilePath, 0)
		fmt.Println(url, error)

		playlist := filepath.Join(filepath.Dir(localFilePath), "stream.m3u8")
		_, error = h.Storage.Save(playlist, 0)
	}()
}

func (h *HLSHandler) VariantPlaylistWritten(localFilePath string) {
	go func() {
		url, error := h.Storage.Save(localFilePath, 0)
		fmt.Println(url, error)
	}()
}

func (h *HLSHandler) MasterPlaylistWritten(localFilePath string) {
	go func() {
		h.Storage.GenerateRemotePlaylist(localFilePath)
		url, error := h.Storage.Save(localFilePath, 0)
		fmt.Println(url, error)
	}()
}
