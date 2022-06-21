package controllers

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/static"
	log "github.com/sirupsen/logrus"
)

var emojiStaticServer = http.FileServer(http.FS(static.GetEmoji()))

// getCustomEmojiList returns a list of custom emoji either from the cache or from the emoji directory.
func getCustomEmojiList() []models.CustomEmoji {
	bundledEmoji := static.GetEmoji()
	emojiResponse := make([]models.CustomEmoji, 0)

	files, err := fs.Glob(bundledEmoji, "*")
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
	emojiStaticServer.ServeHTTP(w, r)
}
