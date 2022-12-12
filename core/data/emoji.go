package data

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
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
