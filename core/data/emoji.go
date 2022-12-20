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
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// GetEmojiList returns a list of custom emoji from the emoji directory.
func GetEmojiList() []models.CustomEmoji {
	emojiFS := os.DirFS(config.CustomEmojiPath)

	emojiResponse := make([]models.CustomEmoji, 0)

	walkFunction := func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		emojiPath := filepath.Join(config.EmojiDir, path)
		singleEmoji := models.CustomEmoji{Name: d.Name(), URL: emojiPath}
		emojiResponse = append(emojiResponse, singleEmoji)
		return nil
	}

	if err := fs.WalkDir(emojiFS, ".", walkFunction); err != nil {
		log.Errorln("unable to fetch emojis: " + err.Error())
		return emojiResponse
	}

	return emojiResponse
}

// SetupEmojiDirectory sets up the custom emoji directory by copying all built-in
// emojis if the directory does not yet exist.
func SetupEmojiDirectory() (err error) {
	type emojiDirectory struct {
		path  string
		isDir bool
	}

	if utils.DoesFileExists(config.CustomEmojiPath) {
		return nil
	}

	if err = os.MkdirAll(config.CustomEmojiPath, 0o750); err != nil {
		return fmt.Errorf("unable to create custom emoji directory: %w", err)
	}

	staticFS := static.GetEmoji()
	files := []emojiDirectory{}

	walkFunction := func(path string, d os.DirEntry, err error) error {
		if path == "." {
			return nil
		}

		if d.Name() == "LICENSE.md" {
			return nil
		}

		files = append(files, emojiDirectory{path: path, isDir: d.IsDir()})
		return nil
	}

	if err := fs.WalkDir(staticFS, ".", walkFunction); err != nil {
		log.Errorln("unable to fetch emojis: " + err.Error())
		return errors.Wrap(err, "unable to fetch embedded emoji files")
	}

	if err != nil {
		return fmt.Errorf("unable to read built-in emoji files: %w", err)
	}

	// Now copy all built-in emojis to the custom emoji directory
	for _, path := range files {
		emojiPath := filepath.Join(config.CustomEmojiPath, path.path)

		if path.isDir {
			if err := os.Mkdir(emojiPath, 0o700); err != nil {
				return errors.Wrap(err, "unable to create emoji directory, check permissions?: "+path.path)
			}
			continue
		}

		memFile, staticOpenErr := staticFS.Open(path.path)
		if staticOpenErr != nil {
			return errors.Wrap(staticOpenErr, "unable to open emoji file from embedded filesystem")
		}

		// nolint:gosec
		diskFile, err := os.Create(emojiPath)
		if err != nil {
			return fmt.Errorf("unable to create custom emoji file on disk: %w", err)
		}

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
