package ffmpeg

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/gabek/owncast/config"
)

//StartThumbnailGenerator starts generating thumbnails
func StartThumbnailGenerator(chunkPath string) {
	// Every 20 seconds create a thumbnail from the most
	// recent video segment.
	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := fireThumbnailGenerator(chunkPath); err != nil {
					panic(err)
				}
			case <-quit:
				//TODO: evaluate if this is ever stopped
				log.Println("thumbnail generator has stopped")
				ticker.Stop()
				return
			}
		}
	}()
}

func fireThumbnailGenerator(chunkPath string) error {
	// JPG takes less time to encode than PNG
	outputFile := path.Join("webroot", "thumbnail.jpg")

	framePath := path.Join(chunkPath, "0")
	files, err := ioutil.ReadDir(framePath)
	if err != nil {
		return err
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
		return nil
	}

	mostRecentFile := path.Join(framePath, names[0])

	thumbnailCmdFlags := []string{
		config.Config.FFMpegPath,
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

	if _, err := exec.Command("sh", "-c", ffmpegCmd).Output(); err != nil {
		return err
	}

	return nil
}
