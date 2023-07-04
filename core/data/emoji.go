package data

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var emojiCacheMu sync.Mutex
var emojiCacheData = make([]models.CustomEmoji, 0)
var emojiCacheModTime time.Time

// UpdateEmojiList will update the cache (if required) and
// return the modifiation time.
func UpdateEmojiList(force bool) (time.Time, error) {
	var modTime time.Time

	emojiPathInfo, err := os.Stat(config.CustomEmojiPath)
	if err != nil {
		return modTime, err
	}

	modTime = emojiPathInfo.ModTime()

	if modTime.After(emojiCacheModTime) || force {
		emojiCacheMu.Lock()
		defer emojiCacheMu.Unlock()

		// double-check that another thread didn't update this while waiting.
		if modTime.After(emojiCacheModTime) || force {
			emojiCacheModTime = modTime
			if force {
				emojiCacheModTime = time.Now()
			}
			emojiFS := os.DirFS(config.CustomEmojiPath)

			emojiCacheData = make([]models.CustomEmoji, 0)

			walkFunction := func(path string, d os.DirEntry, err error) error {
				if d.IsDir() {
					return nil
				}

				emojiPath := filepath.Join(config.EmojiDir, path)
				fileName := d.Name()
				fileBase := fileName[:len(fileName)-len(filepath.Ext(fileName))]
				singleEmoji := models.CustomEmoji{Name: fileBase, URL: emojiPath}
				emojiCacheData = append(emojiCacheData, singleEmoji)
				return nil
			}

			if err := fs.WalkDir(emojiFS, ".", walkFunction); err != nil {
				log.Errorln("unable to fetch emojis: " + err.Error())
			}
		}
	}

	return modTime, nil
}

// GetEmojiList returns a list of custom emoji from the emoji directory.
func GetEmojiList() []models.CustomEmoji {
	_, err := UpdateEmojiList(false)
	if err != nil {
		return nil
	}

	// Lock to make sure this doesn't get updated in the middle of reading
	emojiCacheMu.Lock()
	defer emojiCacheMu.Unlock()

	// return a copy of cache data, ensures underlying slice isn't affected
	// by future update
	emojiData := make([]models.CustomEmoji, len(emojiCacheData))
	copy(emojiData, emojiCacheData)

	return emojiData
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
