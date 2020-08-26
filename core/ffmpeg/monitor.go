package ffmpeg

import (
	"fmt"
	"path/filepath"
	"strings"

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
		// playlist := filepath.Join(filepath.Dir(localFilePath), "stream.m3u8")
		// fmt.Println(playlist, nil)
		url, error := h.Storage.Save(localFilePath, 0)
		fmt.Println(url, error)
	}()
}

func getFilenameFromOutputLine(line string) string {
	fileComponents := strings.Split(line, "'")
	file := fileComponents[1]
	return file
}

func getRelativePathFromAbsolutePath(path string) string {
	pathComponents := strings.Split(path, "/")
	variant := pathComponents[len(pathComponents)-2]
	file := pathComponents[len(pathComponents)-1]

	return filepath.Join(variant, file)
}

func getFinalNameFromPath(path string) string {
	ext := filepath.Ext(path)

	if ext != ".tmp" {
		panic("Incorrect file")
	}

	return strings.Replace(path, ext, "", -1)
}
