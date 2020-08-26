package storageproviders

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/utils"
)

type LocalStorage struct {
}

// Setup configures this storage provider
func (s *LocalStorage) Setup() error {
	// no-op
	return nil
}

// Save will save a local filepath using the storage provider
func (s *LocalStorage) Save(filePath string, retryCount int) (string, error) {
	newPath := ""

	// This is a hack
	if filePath == "hls/stream.m3u8" {
		newPath = filepath.Join(config.Config.GetPublicHLSSavePath(), filepath.Base(filePath))
	} else {
		newPath = filepath.Join("webroot", filePath)
	}
	go utils.Move(filePath, newPath)

	return newPath, nil
}

// GenerateRemotePlaylist will rewrite the playlist using storage provider details
func (s *LocalStorage) GenerateRemotePlaylist(filePath string) error {
	// Locally generated playlist is valid
	return nil
}

// SegmentWritten is a callback when a single video segment is written by the transcoder
func (s *LocalStorage) SegmentWritten(localFilePath string) {
	s.Save(localFilePath, 0)
}

// VariantPlaylistWritten is a callback when a single video variant's playlist is written
// by the transcoder
func (s *LocalStorage) VariantPlaylistWritten(localFilePath string) {
	s.Save(localFilePath, 0)
}

// MasterPlaylistWritten is a callback when the master playlist is written by the transcoder
func (s *LocalStorage) MasterPlaylistWritten(localFilePath string) {
}

func getVariantIndexFromPath(fullDiskPath string) (int, error) {
	return strconv.Atoi(fullDiskPath[0:1])
}

func getRelativePathFromAbsolutePath(path string) string {
	pathComponents := strings.Split(path, "/")
	variant := pathComponents[len(pathComponents)-2]
	file := pathComponents[len(pathComponents)-1]

	return filepath.Join(variant, file)
}
