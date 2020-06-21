package main

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/radovskyb/watcher"
)

type Segment struct {
	VariantIndex       int    // The bitrate variant
	FullDiskPath       string // Where it lives on disk
	RelativeUploadPath string // Path it should have remotely
	RemoteID           string // Used for IPFS
}

type Variant struct {
	VariantIndex int
	Segments     map[string]*Segment
}

func (v *Variant) getSegmentForFilename(filename string) *Segment {
	return v.Segments[filename]
	// for _, segment := range v.Segments {
	// 	if path.Base(segment.FullDiskPath) == filename {
	// 		return &segment
	// 	}
	// }
	// return nil
}

func getSegmentFromPath(fullDiskPath string) Segment {
	segment := Segment{}
	segment.FullDiskPath = fullDiskPath
	segment.RelativeUploadPath = getRelativePathFromAbsolutePath(fullDiskPath)
	index, error := strconv.Atoi(segment.RelativeUploadPath[0:1])
	verifyError(error)
	segment.VariantIndex = index

	return segment
}

func getVariantIndexFromPath(fullDiskPath string) int {
	index, error := strconv.Atoi(fullDiskPath[0:1])
	verifyError(error)
	return index
}

var variants []Variant

func updateVariantPlaylist(fullPath string) {
	relativePath := getRelativePathFromAbsolutePath(fullPath)
	variantIndex := getVariantIndexFromPath(relativePath)
	variant := variants[variantIndex]

	playlistBytes, err := ioutil.ReadFile(fullPath)
	verifyError(err)
	playlistString := string(playlistBytes)
	// fmt.Println("Rewriting playlist", relativePath, "to", path.Join(configuration.PublicHLSPath, relativePath))

	playlistString = storage.GenerateRemotePlaylist(playlistString, variant)

	writePlaylist(playlistString, path.Join(configuration.PublicHLSPath, relativePath))
}

func monitorVideoContent(pathToMonitor string, configuration Config, storage ChunkStorage) {
	// Create at least one structure to store the segments for the different stream variants
	variants = make([]Variant, len(configuration.VideoSettings.StreamQualities))
	if len(configuration.VideoSettings.StreamQualities) > 0 && !configuration.VideoSettings.EnablePassthrough {
		for index := range variants {
			variants[index] = Variant{index, make(map[string]*Segment)}
		}
	} else {
		variants[0] = Variant{0, make(map[string]*Segment)}
	}
	// log.Printf("Using directory %s for storing files with %d variants...\n", pathToMonitor, len(variants))

	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:

				relativePath := getRelativePathFromAbsolutePath(event.Path)
				if path.Ext(relativePath) == ".tmp" {
					continue
				}

				// Ignore removals
				if event.Op == watcher.Remove {
					continue
				}

				// fmt.Println(event.Op, relativePath)

				// Handle updates to the master playlist by copying it to webroot
				if relativePath == path.Join(configuration.PrivateHLSPath, "stream.m3u8") {
					copy(event.Path, path.Join(configuration.PublicHLSPath, "stream.m3u8"))

				} else if filepath.Ext(event.Path) == ".m3u8" {
					// Handle updates to playlists, but not the master playlist
					updateVariantPlaylist(event.Path)
				} else if filepath.Ext(event.Path) == ".ts" {
					segment := getSegmentFromPath(event.Path)

					newObjectPathChannel := make(chan string, 1)
					go func() {
						newObjectPath := storage.Save(path.Join(configuration.PrivateHLSPath, segment.RelativeUploadPath), 0)
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
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch the hls segment storage folder recursively for changes.
	w.FilterOps(watcher.Write, watcher.Rename, watcher.Create)

	if err := w.AddRecursive(pathToMonitor); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 200); err != nil {
		log.Fatalln(err)
	}
}
