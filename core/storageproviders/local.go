package storageproviders

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/ffmpeg"
	"github.com/owncast/owncast/utils"
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
	_, error := s.Save(localFilePath, 0)
	if error != nil {
		log.Errorln(error)
		return
	}
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
		newPath = filepath.Join(config.PublicHLSStoragePath, filepath.Base(filePath))
	} else {
		newPath = filepath.Join(config.WebRoot, filePath)
	}

	// Move video segments to the destination directory.
	// Copy playlists to the destination directory so they can still be referenced in
	// the private hls working directory.
	if filepath.Ext(filePath) == ".m3u8" {
		utils.Copy(filePath, newPath)
	} else {
		utils.Move(filePath, newPath)
		ffmpeg.Cleanup(filepath.Dir(newPath))
	}

	return newPath, nil
}
