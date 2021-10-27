package storageproviders

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"

	"os"
	"path/filepath"
	"sort"
)

// LocalStorage represents an instance of the local storage provider for HLS video.
type LocalStorage struct {
	// Cleanup old public HLS content every N min from the webroot.
	onlineCleanupTicker *time.Ticker
}

// NewLocalStorage returns a new LocalStorage instance.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

// Setup configures this storage provider.
func (s *LocalStorage) Setup() error {
	// NOTE: This cleanup timer will have to be disabled to support recordings in the future
	// as all HLS segments have to be publicly available on disk to keep a recording of them.
	s.onlineCleanupTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for range s.onlineCleanupTicker.C {
			s.CleanupOldContent(config.HLSStoragePath)
		}
	}()
	return nil
}

// SegmentWritten is called when a single segment of video is written.
func (s *LocalStorage) SegmentWritten(localFilePath string) {
	if _, err := s.Save(localFilePath, 0); err != nil {
		log.Warnln(err)
	}
}

// VariantPlaylistWritten is called when a variant hls playlist is written.
func (s *LocalStorage) VariantPlaylistWritten(localFilePath string) {
	if _, err := s.Save(localFilePath, 0); err != nil {
		log.Errorln(err)
		return
	}
}

// MasterPlaylistWritten is called when the master hls playlist is written.
func (s *LocalStorage) MasterPlaylistWritten(localFilePath string) {
	if _, err := s.Save(localFilePath, 0); err != nil {
		log.Warnln(err)
	}
}

// Save will save a local filepath using the storage provider.
func (s *LocalStorage) Save(filePath string, retryCount int) (string, error) {
	return filePath, nil
}

// CleanupOldContent will delete old files from the private dir that are no longer being referenced
// in the stream.
func (s *LocalStorage) CleanupOldContent(baseDirectory string) {
	// Determine how many files we should keep on disk
	maxNumber := data.GetStreamLatencyLevel().SegmentCount
	buffer := 10

	files, err := getAllFilesRecursive(baseDirectory)
	if err != nil {
		log.Debugln("Unable to cleanup old video files", err)
		return
	}

	// Delete old private HLS files on disk
	for directory := range files {
		files := files[directory]
		if len(files) < maxNumber+buffer {
			continue
		}

		filesToDelete := files[maxNumber+buffer:]
		log.Traceln("Deleting", len(filesToDelete), "old files from", baseDirectory, "for video variant", directory)

		for _, file := range filesToDelete {
			fileToDelete := filepath.Join(baseDirectory, directory, file.Name())
			err := os.Remove(fileToDelete)
			if err != nil {
				log.Debugln(err)
			}
		}
	}
}

func getAllFilesRecursive(baseDirectory string) (map[string][]os.FileInfo, error) {
	var files = make(map[string][]os.FileInfo)

	var directory string
	err := filepath.Walk(baseDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			directory = info.Name()
		}

		if filepath.Ext(info.Name()) == ".ts" {
			files[directory] = append(files[directory], info)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by date so we can delete old files
	for directory := range files {
		sort.Slice(files[directory], func(i, j int) bool {
			return files[directory][i].ModTime().UnixNano() > files[directory][j].ModTime().UnixNano()
		})
	}

	return files, nil
}
