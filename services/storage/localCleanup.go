package storage

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

// protectedSegments returns the segment filenames a cleanup pass must keep.
// A nil protector, or one that fails to answer, protects nothing: cleanup
// then behaves as it does on a server with no recorded video.
func protectedSegments(protector models.SegmentProtector) map[string]bool {
	if protector == nil {
		return nil
	}

	protected, err := protector.ProtectedSegmentFilenames()
	if err != nil {
		log.Warnln("unable to determine which video segments are in use, keeping none:", err)
		return nil
	}

	return protected
}

// localCleanup deletes all but the newest maxNumber segments of each variant,
// keeping any segment a clip or replay still references.
func localCleanup(maxNumber int, protector models.SegmentProtector) error {
	baseDirectory := config.HLSStoragePath

	files, err := getAllFilesRecursive(baseDirectory)
	if err != nil {
		return errors.Wrap(err, "unable find old video files for cleanup")
	}

	protected := protectedSegments(protector)

	// Delete old private HLS files on disk
	for directory := range files {
		files := files[directory]
		if len(files) < maxNumber {
			continue
		}

		filesToDelete := files[maxNumber:]
		log.Traceln("Deleting", len(filesToDelete), "old files from", baseDirectory, "for video variant", directory)

		for _, file := range filesToDelete {
			if protected[file.Name()] {
				continue
			}

			fileToDelete := filepath.Join(baseDirectory, directory, file.Name())
			err := os.Remove(fileToDelete)
			if err != nil {
				return errors.Wrap(err, "unable to delete old video files")
			}
		}
	}
	return nil
}
