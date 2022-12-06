package data

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/static"
	log "github.com/sirupsen/logrus"
)

// GetEmojiList returns a list of custom emoji either from the cache or from the emoji directory.
func GetEmojiList(custom bool) []models.CustomEmoji {
	var emojiFS fs.FS

	if custom {
		emojiFS = os.DirFS(config.CustomEmojiPath)
	} else {
		emojiFS = static.GetEmoji()
	}

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
