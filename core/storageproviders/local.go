package storageproviders

import (
	"path/filepath"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/models"
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
	// Move the file to the public directory
	newPath := filepath.Join(config.Config.GetPublicHLSSavePath(), utils.GetRelativePathFromAbsolutePath(filePath))
	go utils.Move(filePath, newPath)

	// Move the associated playlist as well
	playlist := filepath.Join(filepath.Dir(filePath), "stream.m3u8")
	newPlaylistPath := filepath.Join(config.Config.GetPublicHLSSavePath(), utils.GetRelativePathFromAbsolutePath(playlist))

	go utils.Move(playlist, newPlaylistPath)

	return newPath, nil
}

// GenerateRemotePlaylist will rewrite the playlist using storage provider details
func (s *LocalStorage) GenerateRemotePlaylist(playlist string, variant models.Variant) string {
	// Locally generated playlist is valid
	return playlist
}

// SegmentWritten is a callback when a single video segment is written by the transcoder
func (s *LocalStorage) SegmentWritten(localFilePath string) {
	s.Save(localFilePath, 0)
}

// VariantPlaylistWritten is a callback when a single video variant's playlist is written
// by the transcoder
func (s *LocalStorage) VariantPlaylistWritten(localFilePath string) {
}

// MasterPlaylistWritten is a callback when the master playlist is written by the transcoder
func (s *LocalStorage) MasterPlaylistWritten(localFilePath string) {
	newPath := filepath.Join(config.Config.GetPublicHLSSavePath(), filepath.Base(localFilePath))
	go utils.Move(localFilePath, newPath)
}
