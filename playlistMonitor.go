package main

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/radovskyb/watcher"
)

var filesToUpload = make(map[string]string)

func monitorVideoContent(pathToMonitor string, configuration Config, storage ChunkStorage) {
	log.Printf("Using %s for IPFS files...\n", pathToMonitor)

	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.Op != watcher.Write {
					continue
				}
				if filepath.Base(event.Path) == "temp.m3u8" {

					if configuration.IPFS.Enabled {
						for filePath, objectID := range filesToUpload {
							if objectID != "" {
								continue
							}

							newObjectPath := storage.Save(path.Join(configuration.PrivateHLSPath, filePath))
							filesToUpload[filePath] = newObjectPath
						}
					}

					playlistBytes, err := ioutil.ReadFile(event.Path)
					verifyError(err)
					playlistString := string(playlistBytes)

					if configuration.IPFS.Enabled {
						playlistString = storage.GenerateRemotePlaylist(playlistString, filesToUpload)
					}
					writePlaylist(playlistString, path.Join(configuration.PublicHLSPath, "/stream.m3u8"))

				} else if filepath.Ext(event.Path) == ".ts" {
					if configuration.IPFS.Enabled {
						filesToUpload[filepath.Base(event.Path)] = ""
					}
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(pathToMonitor); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
