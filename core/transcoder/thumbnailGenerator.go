package transcoder

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
)

var _timer *time.Ticker

// StopThumbnailGenerator will stop the periodic generating of a thumbnail from video.
func StopThumbnailGenerator() {
	if _timer != nil {
		_timer.Stop()
	}
}

// StartThumbnailGenerator starts generating thumbnails.
func StartThumbnailGenerator(chunkPath string, variantIndex int) {
	// Every 20 seconds create a thumbnail from the most
	// recent video segment.
	_timer = time.NewTicker(20 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-_timer.C:
				if err := fireThumbnailGenerator(chunkPath, variantIndex); err != nil {
					log.Errorln("Unable to generate thumbnail:", err)
				}
			case <-quit:
				log.Debug("thumbnail generator has stopped")
				_timer.Stop()
				return
			}
		}
	}()
}

func fireThumbnailGenerator(segmentPath string, variantIndex int) error {
	// JPG takes less time to encode than PNG
	outputFile := path.Join(config.WebRoot, "thumbnail.jpg")
	previewGifFile := path.Join(config.WebRoot, "preview.gif")

	framePath := path.Join(segmentPath, strconv.Itoa(variantIndex))
	files, err := os.ReadDir(framePath)
	if err != nil {
		return err
	}

	var modTime time.Time
	var names []string
	for _, f := range files {
		if path.Ext(f.Name()) != ".ts" {
			continue
		}

		fi, err := f.Info()
		if err != nil {
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
	ffmpegPath := utils.ValidatedFfmpegPath(data.GetFfMpegPath())
	outputFileTemp := path.Join(config.WebRoot, "tempthumbnail.jpg")

	thumbnailCmdFlags := []string{
		ffmpegPath,
		"-y",                 // Overwrite file
		"-threads 1",         // Low priority processing
		"-t 1",               // Pull from frame 1
		"-i", mostRecentFile, // Input
		"-f image2",  // format
		"-vframes 1", // Single frame
		outputFileTemp,
	}

	ffmpegCmd := strings.Join(thumbnailCmdFlags, " ")
	if _, err := exec.Command("sh", "-c", ffmpegCmd).Output(); err != nil {
		return err
	}

	// rename temp file
	if err := os.Rename(outputFileTemp, outputFile); err != nil {
		log.Errorln(err)
	}

	makeAnimatedGifPreview(mostRecentFile, previewGifFile)

	return nil
}

func makeAnimatedGifPreview(sourceFile string, outputFile string) {
	ffmpegPath := utils.ValidatedFfmpegPath(data.GetFfMpegPath())
	outputFileTemp := path.Join(config.WebRoot, "temppreview.gif")

	// Filter is pulled from https://engineering.giphy.com/how-to-make-gifs-with-ffmpeg/
	animatedGifFlags := []string{
		ffmpegPath,
		"-y",             // Overwrite file
		"-threads 1",     // Low priority processing
		"-i", sourceFile, // Input
		"-t 1", // Output is one second in length
		"-filter_complex", "\"[0:v] fps=8,scale=w=480:h=-1:flags=lanczos,split [a][b];[a] palettegen=stats_mode=full [p];[b][p] paletteuse=new=1\"",
		outputFileTemp,
	}

	ffmpegCmd := strings.Join(animatedGifFlags, " ")
	if _, err := exec.Command("sh", "-c", ffmpegCmd).Output(); err != nil {
		log.Errorln(err)
		// rename temp file
	} else if err := os.Rename(outputFileTemp, outputFile); err != nil {
		log.Errorln(err)
	}
}
