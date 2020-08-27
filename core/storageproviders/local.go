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
	// Locally generated playlist doesn't need to be rewritten.  Noop.
	return nil
}
