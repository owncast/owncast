package data

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// GetEmojiList returns a list of custom emoji from the emoji directory.
func GetEmojiList() []models.CustomEmoji {
	var emojiFS = os.DirFS(config.CustomEmojiPath)

	emojiResponse := make([]models.CustomEmoji, 0)

	files, err := fs.Glob(emojiFS, "*")
	if err != nil {
		log.Errorln(err)
		return emojiResponse
	}

	for _, name := range files {
		emojiPath := filepath.Join(config.EmojiDir, name)
		singleEmoji := models.CustomEmoji{Name: name, URL: emojiPath}
		emojiResponse = append(emojiResponse, singleEmoji)
	}

	return emojiResponse
}

// SetupEmojiDirectory sets up the custom emoji directory by copying all built-in
// emojis if the directory does not yet exist.
func SetupEmojiDirectory() (err error) {
	if utils.DoesFileExists(config.CustomEmojiPath) {
		return nil
	}

	if err = os.MkdirAll(config.CustomEmojiPath, 0o750); err != nil {
		return fmt.Errorf("unable to create custom emoji directory: %w", err)
	}

	staticFS := static.GetEmoji()
	files, err := fs.Glob(staticFS, "*")
	if err != nil {
		return fmt.Errorf("unable to read built-in emoji files: %w", err)
	}

	// Now copy all built-in emojis to the custom emoji directory
	for _, name := range files {
		emojiPath := filepath.Join(config.CustomEmojiPath, filepath.Base(name))

		// nolint:gosec
		diskFile, err := os.Create(emojiPath)
		if err != nil {
			return fmt.Errorf("unable to create custom emoji file on disk: %w", err)
		}

		memFile, err := staticFS.Open(name)
		if err != nil {
			_ = diskFile.Close()
			return fmt.Errorf("unable to open built-in emoji file: %w", err)
		}

		if _, err = io.Copy(diskFile, memFile); err != nil {
			_ = diskFile.Close()
			_ = os.Remove(emojiPath)
			return fmt.Errorf("unable to copy built-in emoji file to disk: %w", err)
		}

		if err = diskFile.Close(); err != nil {
			_ = os.Remove(emojiPath)
			return fmt.Errorf("unable to close custom emoji file on disk: %w", err)
		}
	}

	return nil
}
