package static

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	log "github.com/sirupsen/logrus"
)

var (
	emojiCacheMu      sync.Mutex
	emojiCacheData    = make([]models.CustomEmoji, 0)
	emojiCacheModTime time.Time
)

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
			if emojiFS == nil {
				return modTime, fmt.Errorf("unable to open custom emoji directory")
			}

			emojiCacheData = make([]models.CustomEmoji, 0)

			walkFunction := func(path string, d os.DirEntry, err error) error {
				if d == nil || d.IsDir() {
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
