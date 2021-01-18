package ffmpeg

import (
	log "github.com/sirupsen/logrus"

	"os"
	"path/filepath"
	"sort"

	"github.com/owncast/owncast/core/data"
)

// CleanupOldContent will delete old files from the private dir that are no longer being referenced
// in the stream.
func CleanupOldContent(baseDirectory string) {
	// Determine how many files we should keep on disk
	maxNumber := int(data.GetStreamLatancyLevel().SegmentCount)
	buffer := 10

	files, err := getAllFilesRecursive(baseDirectory)
	if err != nil {
		log.Errorln("Unable to cleanup old video files", err)
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
				log.Errorln(err)
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
		log.Fatalf(err.Error())
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
