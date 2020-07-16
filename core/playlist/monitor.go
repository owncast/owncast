package playlist

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/radovskyb/watcher"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/models"
	"github.com/gabek/owncast/utils"
)

var (
	_storage models.ChunkStorageProvider
	variants []models.Variant
)

//StartVideoContentMonitor starts the video content monitor
func StartVideoContentMonitor(storage models.ChunkStorageProvider) error {
	_storage = storage

	pathToMonitor := config.Config.GetPrivateHLSSavePath()

	// Create at least one structure to store the segments for the different stream variants
	variants = make([]models.Variant, len(config.Config.VideoSettings.StreamQualities))
	if len(config.Config.VideoSettings.StreamQualities) > 0 {
		for index := range variants {
			variants[index] = models.Variant{
				VariantIndex: index,
				Segments:     make(map[string]*models.Segment),
			}
		}
	} else {
		variants[0] = models.Variant{
			VariantIndex: 0,
			Segments:     make(map[string]*models.Segment),
		}
	}

	// log.Printf("Using directory %s for storing files with %d variants...\n", pathToMonitor, len(variants))

	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:

				relativePath := utils.GetRelativePathFromAbsolutePath(event.Path)
				if path.Ext(relativePath) == ".tmp" {
					continue
				}

				// Ignore removals
				if event.Op == watcher.Remove {
					continue
				}

				// fmt.Println(event.Op, relativePath)

				// Handle updates to the master playlist by copying it to webroot
				if relativePath == path.Join(config.Config.GetPrivateHLSSavePath(), "stream.m3u8") {
					utils.Copy(event.Path, path.Join(config.Config.GetPublicHLSSavePath(), "stream.m3u8"))

				} else if filepath.Ext(event.Path) == ".m3u8" {
					// Handle updates to playlists, but not the master playlist
					updateVariantPlaylist(event.Path)

				} else if filepath.Ext(event.Path) == ".ts" {
					segment, err := getSegmentFromPath(event.Path)
					if err != nil {
						log.Error("failed to get the segment from path")
						panic(err)
					}

					newObjectPathChannel := make(chan string, 1)
					go func() {
						newObjectPath, err := storage.Save(path.Join(config.Config.GetPrivateHLSSavePath(), segment.RelativeUploadPath), 0)
						if err != nil {
							log.Errorln("failed to save the file to the chunk storage.", err)
						}

						newObjectPathChannel <- newObjectPath
					}()

					newObjectPath := <-newObjectPathChannel
					segment.RemoteID = newObjectPath
					// fmt.Println("Uploaded", segment.RelativeUploadPath, "as", newObjectPath)

					variants[segment.VariantIndex].Segments[filepath.Base(segment.RelativeUploadPath)] = &segment

					// Force a variant's playlist to be updated after a file is uploaded.
					associatedVariantPlaylist := strings.ReplaceAll(event.Path, path.Base(event.Path), "stream.m3u8")
					updateVariantPlaylist(associatedVariantPlaylist)
				}
			case err := <-w.Error:
				panic(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch the hls segment storage folder recursively for changes.
	w.FilterOps(watcher.Write, watcher.Rename, watcher.Create)

	if err := w.AddRecursive(pathToMonitor); err != nil {
		return err
	}

	return w.Start(time.Millisecond * 200)
}

func getSegmentFromPath(fullDiskPath string) (models.Segment, error) {
	segment := models.Segment{
		FullDiskPath:       fullDiskPath,
		RelativeUploadPath: utils.GetRelativePathFromAbsolutePath(fullDiskPath),
	}

	index, err := strconv.Atoi(segment.RelativeUploadPath[0:1])
	if err != nil {
		return segment, err
	}

	segment.VariantIndex = index

	return segment, nil
}

func getVariantIndexFromPath(fullDiskPath string) (int, error) {
	return strconv.Atoi(fullDiskPath[0:1])
}

func updateVariantPlaylist(fullPath string) error {
	relativePath := utils.GetRelativePathFromAbsolutePath(fullPath)
	variantIndex, err := getVariantIndexFromPath(relativePath)
	if err != nil {
		return err
	}

	variant := variants[variantIndex]

	playlistBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}

	playlistString := string(playlistBytes)
	playlistString = _storage.GenerateRemotePlaylist(playlistString, variant)

	return WritePlaylist(playlistString, path.Join(config.Config.GetPublicHLSSavePath(), relativePath))
}
