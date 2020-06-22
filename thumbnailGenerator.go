package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func startThumbnailGenerator(chunkPath string) {
	// Every 20 seconds create a thumbnail from the most
	// recent video segment.
	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fireThumbnailGenerator(chunkPath)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func fireThumbnailGenerator(chunkPath string) {
	framePath := path.Join(chunkPath, "0")
	files, err := ioutil.ReadDir(framePath)

	// JPG takes less time to encode than PNG
	outputFile := path.Join("webroot", "thumbnail.jpg")

	// fmt.Println("Generating thumbnail from", framePath, "to", outputFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var modTime time.Time
	var names []string
	for _, fi := range files {
		if path.Ext(fi.Name()) != ".ts" {
			continue
		}

		if fi.Mode().IsRegular() {
			if !fi.ModTime().Before(modTime) {
				if fi.ModTime().After(modTime) {
					modTime = fi.ModTime()
					names = names[:0]
				}
				names = append(names, fi.Name())
			}
		}
	}

	if len(names) == 0 {
		return
	}

	mostRecentFile := path.Join(framePath, names[0])

	thumbnailCmdFlags := []string{
		configuration.FFMpegPath,
		"-y",                 // Overwrite file
		"-threads 1",         // Low priority processing
		"-t 1",               // Pull from frame 1
		"-i", mostRecentFile, // Input
		"-f image2",  // format
		"-vframes 1", // Single frame
		outputFile,
	}

	ffmpegCmd := strings.Join(thumbnailCmdFlags, " ")

	// fmt.Println(ffmpegCmd)

	_, err = exec.Command("sh", "-c", ffmpegCmd).Output()
	if err != nil {
		log.Errorln("Unable to generate thumbnail: ", err)
	}
}
