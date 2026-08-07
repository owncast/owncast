package replays

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

// ClipThumbnailFilename returns the poster image filename for a clip.
func ClipThumbnailFilename(clipID string) string {
	return clipID + ".jpg"
}

// ClipThumbnailPath returns the on-disk path of a clip's poster image.
func ClipThumbnailPath(clipID string) string {
	return filepath.Join(config.ClipThumbnailsPath, ClipThumbnailFilename(clipID))
}

// GenerateClipThumbnail extracts a single frame from the first segment of a
// clip and stores it as the clip's poster image, used for clip listings and
// link previews. Best effort: a clip without a poster still plays, so
// failures are logged rather than returned.
func (s *Service) GenerateClipThumbnail(clipID string) {
	generator := s.NewPlaylistGenerator()

	clip, err := generator.GetClip(clipID)
	if err != nil {
		log.Debugln("unable to generate clip thumbnail:", err)
		return
	}

	configs, err := generator.GetConfigurationsForStream(clip.StreamID)
	if err != nil || len(configs) == 0 {
		log.Debugln("unable to generate clip thumbnail: no output configurations")
		return
	}

	// Use the highest-bitrate variant available so the poster is the best
	// quality version of the frame.
	best := configs[0]
	for _, c := range configs {
		if c.VideoBitrate > best.VideoBitrate {
			best = c
		}
	}

	segments, err := generator.GetAllSegmentsForOutputConfigurationAndWindow(best.ID, clip.RelativeStartTime, clip.RelativeEndTime)
	if err != nil || len(segments) == 0 {
		log.Debugln("unable to generate clip thumbnail: no segments in clip window")
		return
	}

	first := segments[0]

	// Only local segments can be read directly by ffmpeg here; remote
	// (object storage) segments are served by URL and are skipped.
	sourcePath := filepath.Join(config.DataDirectory, first.Path)
	if !utils.DoesFileExists(sourcePath) {
		log.Debugln("unable to generate clip thumbnail: segment not available locally", sourcePath)
		return
	}

	if err := os.MkdirAll(config.ClipThumbnailsPath, 0o750); err != nil {
		log.Errorln("unable to create clip thumbnail directory:", err)
		return
	}

	// Seek to where the clip actually starts inside its first segment so the
	// poster shows the clipped moment, not the segment boundary.
	seekSeconds := float64(clip.RelativeStartTime) - first.MediaOffset
	if seekSeconds < 0 {
		seekSeconds = 0
	}

	ffmpegPath := utils.ValidatedFfmpegPath(s.configRepository.GetFfMpegPath())
	outputPath := ClipThumbnailPath(clipID)

	flags := []string{
		"-y",            // Overwrite existing
		"-threads", "1", // Low priority processing
		"-ss", formatSeconds(seekSeconds), // Seek within the segment
		"-i", sourcePath,
		"-f", "image2",
		"-vframes", "1",
		"-vf", "scale=480:-1",
		outputPath,
	}

	if _, err := exec.Command(ffmpegPath, flags...).Output(); err != nil { // #nosec G204 -- ffmpeg path is validated and arguments are server-generated
		log.Debugln("unable to generate clip thumbnail:", err)
	}
}

// RemoveClipThumbnail deletes a clip's poster image, if one exists.
func RemoveClipThumbnail(clipID string) {
	path := ClipThumbnailPath(clipID)
	if !utils.DoesFileExists(path) {
		return
	}

	if err := os.Remove(path); err != nil {
		log.Debugln("unable to remove clip thumbnail:", err)
	}
}

// formatSeconds renders a seek offset for ffmpeg with millisecond precision.
func formatSeconds(seconds float64) string {
	return strconv.FormatFloat(seconds, 'f', 3, 64)
}
