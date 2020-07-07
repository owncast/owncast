package ffmpeg

import (
	"io/ioutil"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
)

//StartThumbnailGenerator starts generating thumbnails
func StartThumbnailGenerator(chunkPath string, variantIndex int) {
	// Every 20 seconds create a thumbnail from the most
	// recent video segment.
	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := fireThumbnailGenerator(chunkPath, variantIndex); err != nil {
					log.Errorln("Unable to generate thumbnail:", err)
				}
			case <-quit:
				//TODO: evaluate if this is ever stopped
				log.Debug("thumbnail generator has stopped")
				ticker.Stop()
				return
			}
		}
	}()
}

func fireThumbnailGenerator(chunkPath string, variantIndex int) error {
	// JPG takes less time to encode than PNG
	outputFile := path.Join("webroot", "thumbnail.jpg")

	framePath := path.Join(chunkPath, strconv.Itoa(variantIndex))
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
