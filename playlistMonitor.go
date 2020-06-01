package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/radovskyb/watcher"
)

var filesToUpload = make(map[string]string)

func monitorVideoContent(path string, ipfs *icore.CoreAPI) {
	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:
				if event.Op != watcher.Write {
					continue
				}
				if filepath.Base(event.Path) == "temp.m3u8" {

					for filePath, objectID := range filesToUpload {
						if objectID != "" {
							continue
						}

						newObjectPath := save("hls/"+filePath, ipfs)
						fmt.Println(filePath, newObjectPath)

						filesToUpload[filePath] = newObjectPath
					}
					playlistBytes, err := ioutil.ReadFile(event.Path)
					verifyError(err)
					playlistString := string(playlistBytes)

					if false {
						playlistString = generateRemotePlaylist(playlistString, filesToUpload)
					}
					writePlaylist(playlistString, "webroot/stream.m3u8")

				} else if filepath.Ext(event.Path) == ".ts" {
					filesToUpload[filepath.Base(event.Path)] = ""
					// copy(event.Path, "webroot/"+filepath.Base(event.Path))
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(path); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
