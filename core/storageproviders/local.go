package storageproviders

import (
	"path/filepath"

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

// SegmentWritten is called when a single segment of video is written
func (s *LocalStorage) SegmentWritten(localFilePath string) {
	s.Save(localFilePath, 0)
}

// VariantPlaylistWritten is called when a variant hls playlist is written
func (s *LocalStorage) VariantPlaylistWritten(localFilePath string) {
	s.Save(localFilePath, 0)
}

// MasterPlaylistWritten is called when the master hls playlist is written
func (s *LocalStorage) MasterPlaylistWritten(localFilePath string) {
	s.Save(localFilePath, 0)
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
