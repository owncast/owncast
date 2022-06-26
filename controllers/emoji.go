package controllers

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

var useCustomEmojiDirectory = utils.DoesFileExists(config.CustomEmojiPath)

// getCustomEmojiList returns a list of custom emoji either from the cache or from the emoji directory.
func getCustomEmojiList() []models.CustomEmoji {
	var emojiFS fs.FS
	if useCustomEmojiDirectory {
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

// GetCustomEmojiList returns a list of custom emoji via the API.
func GetCustomEmojiList(w http.ResponseWriter, r *http.Request) {
	emojiList := getCustomEmojiList()

	if err := json.NewEncoder(w).Encode(emojiList); err != nil {
		InternalErrorHandler(w, err)
	}
}

// GetCustomEmojiImage returns a single emoji image.
func GetCustomEmojiImage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/img/emoji/")
	r.URL.Path = path

	var emojiStaticServer http.Handler
	if useCustomEmojiDirectory {
		emojiFS := os.DirFS(config.CustomEmojiPath)
		emojiStaticServer = http.FileServer(http.FS(emojiFS))
	} else {
		emojiStaticServer = http.FileServer(http.FS(static.GetEmoji()))
	}

	emojiStaticServer.ServeHTTP(w, r)
}
