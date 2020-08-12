package ffmpeg

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/gabek/owncast/models"
	"github.com/gabek/owncast/utils"
)

var (
	_storage models.StorageProvider
)

func setupVariantObjectStorage() {

}

func segmentWritten(localFilePath string) {
	_storage.SegmentWritten(localFilePath)
}

func variantPlaylistWritten(localFilePath string) {
	_storage.VariantPlaylistWritten(localFilePath)
}

func masterPlaylistWritten(localFilePath string) {
	_storage.MasterPlaylistWritten(localFilePath)
}

func StartTranscoderMonitor(transcoderOutput io.ReadCloser, storage models.StorageProvider) {
	fmt.Println("StartTranscoderMonitor")
	setupVariantObjectStorage()
	_storage = storage
	_storage.Setup()

	scanner := bufio.NewScanner(transcoderOutput)
	scanner.Split(bufio.ScanLines)

	var inProgressTmpFiles = map[string]string{}

	for scanner.Scan() {
		m := scanner.Text()
		// fmt.Println(m)
		if !strings.Contains(m, `Opening '`) {
			continue
		}

		file := getFilenameFromOutputLine(m)
		if file == "" {
			continue
		}

		if strings.HasSuffix(file, ".ts.tmp") {
			variantIndex := utils.GetIndexFromFilePath(file)

			finishedFile := inProgressTmpFiles[variantIndex]
			if finishedFile != "" {
				finalFilename := getFinalNameFromPath(finishedFile)
				// fmt.Println(finalFilename, "finished")
				segmentWritten(finalFilename)
			}

			inProgressTmpFiles[variantIndex] = file
			//fmt.Printf("%+v\n", inProgressTmpFiles)

		} else if strings.HasSuffix(file, ".m3u8.tmp") {
			finalFilename := getFinalNameFromPath(file)
			if finalFilename == "hls/stream.m3u8" {
				// Master Playlist has been written
				masterPlaylistWritten(finalFilename)
			} else {
				// Video variant playlist has been written
				if variantPlaylistWritten != nil {
					variantPlaylistWritten(finalFilename)
				}
				// fmt.Println(finalFilename)
			}

		} else if strings.Contains(m, "progress=end") {
			// End of transcoding.  Report back all the files.
			fmt.Printf("%+v\n", inProgressTmpFiles)
			if segmentWritten != nil {
				for _, inProgressFile := range inProgressTmpFiles {
					segmentWritten(inProgressFile)
				}
			}
		}
	}
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

// //StartVideoContentMonitor starts the video content monitor
// func StartVideoContentMonitor(storage models.ChunkStorageProvider) error {
// 	_storage = storage

// 	pathToMonitor := config.Config.GetPrivateHLSSavePath()

// 	// Create at least one structure to store the segments for the different stream variants
// 	variants = make([]models.Variant, len(config.Config.VideoSettings.StreamQualities))
// 	if len(config.Config.VideoSettings.StreamQualities) > 0 {
// 		for index := range variants {
// 			variants[index] = models.Variant{
// 				VariantIndex: index,
// 				Segments:     make(map[string]*models.Segment),
// 			}
// 		}
// 	} else {
// 		variants[0] = models.Variant{
// 			VariantIndex: 0,
// 			Segments:     make(map[string]*models.Segment),
// 		}
// 	}

// 	// log.Printf("Using directory %s for storing files with %d variants...\n", pathToMonitor, len(variants))

// 	w := watcher.New()

// 	go func() {
// 		for {
// 			select {
// 			case event := <-w.Event:

// 				relativePath := utils.GetRelativePathFromAbsolutePath(event.Path)
// 				if path.Ext(relativePath) == ".tmp" {
// 					continue
// 				}

// 				// Ignore removals
// 				if event.Op == watcher.Remove {
// 					continue
// 				}

// 				// fmt.Println(event.Op, relativePath)

// 				// Handle updates to the master playlist by copying it to webroot
// 				if relativePath == path.Join(config.Config.GetPrivateHLSSavePath(), "stream.m3u8") {
// 					utils.Copy(event.Path, path.Join(config.Config.GetPublicHLSSavePath(), "stream.m3u8"))

// 				} else if filepath.Ext(event.Path) == ".m3u8" {
// 					// Handle updates to playlists, but not the master playlist
// 					updateVariantPlaylist(event.Path)

// 				} else if filepath.Ext(event.Path) == ".ts" {
// 					segment, err := getSegmentFromPath(event.Path)
// 					if err != nil {
// 						log.Error("failed to get the segment from path")
// 						panic(err)
// 					}

// 					newObjectPathChannel := make(chan string, 1)
// 					go func() {
// 						newObjectPath, err := storage.Save(path.Join(config.Config.GetPrivateHLSSavePath(), segment.RelativeUploadPath), 0)
// 						if err != nil {
// 							log.Errorln("failed to save the file to the chunk storage.", err)
// 						}

// 						newObjectPathChannel <- newObjectPath
// 					}()

// 					newObjectPath := <-newObjectPathChannel
// 					segment.RemoteID = newObjectPath
// 					// fmt.Println("Uploaded", segment.RelativeUploadPath, "as", newObjectPath)

// 					variants[segment.VariantIndex].Segments[filepath.Base(segment.RelativeUploadPath)] = &segment

// 					// Force a variant's playlist to be updated after a file is uploaded.
// 					associatedVariantPlaylist := strings.ReplaceAll(event.Path, path.Base(event.Path), "stream.m3u8")
// 					updateVariantPlaylist(associatedVariantPlaylist)
// 				}
// 			case err := <-w.Error:
// 				panic(err)
// 			case <-w.Closed:
// 				return
// 			}
// 		}
// 	}()

// 	// Watch the hls segment storage folder recursively for changes.
// 	w.FilterOps(watcher.Write, watcher.Rename, watcher.Create)

// 	if err := w.AddRecursive(pathToMonitor); err != nil {
// 		return err
// 	}

// 	return w.Start(time.Millisecond * 200)
// }

// func getSegmentFromPath(fullDiskPath string) (models.Segment, error) {
// 	segment := models.Segment{
// 		FullDiskPath:       fullDiskPath,
// 		RelativeUploadPath: utils.GetRelativePathFromAbsolutePath(fullDiskPath),
// 	}

// 	index, err := strconv.Atoi(segment.RelativeUploadPath[0:1])
// 	if err != nil {
// 		return segment, err
// 	}

// 	segment.VariantIndex = index

// 	return segment, nil
// }

// func getVariantIndexFromPath(fullDiskPath string) (int, error) {
// 	return strconv.Atoi(fullDiskPath[0:1])
// }
