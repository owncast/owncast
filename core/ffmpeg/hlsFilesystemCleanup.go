package ffmpeg

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/gabek/owncast/config"
)

// Cleanup will delete old files off disk that are no longer being referenced
// in the stream.
func Cleanup(directoryPath string) {
	// Determine how many files we should keep on disk
	maxNumber := config.Config.GetMaxNumberOfReferencedSegmentsInPlaylist()
	buffer := 10

	files, err := getSegmentFiles(directoryPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) < maxNumber+buffer {
		return
	}

	// Delete old files on disk
	filesToDelete := files[maxNumber+buffer:]
	for _, file := range filesToDelete {
		os.Remove(filepath.Join(directoryPath, file.Name()))
	}
}

func getSegmentFiles(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	filteredList := make([]os.FileInfo, 0)

	// Filter out playlists because we don't want to clean them up
	for _, file := range list {
		if filepath.Ext(file.Name()) == ".m3u8" {
			continue
		}
		filteredList = append(filteredList, file)
	}

	// Sort by date so we can delete old files
	sort.Slice(filteredList, func(i, j int) bool {
		return filteredList[i].ModTime().UnixNano() > filteredList[j].ModTime().UnixNano()
	})

	return filteredList, nil
}
