package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/static"
)

// GetCustomEmojiList returns a list of emoji via the API.
func GetCustomEmojiList(w http.ResponseWriter, r *http.Request) {
	emojiList := data.GetEmojiList(data.GetUseCustomEmojis())

	if err := json.NewEncoder(w).Encode(emojiList); err != nil {
		InternalErrorHandler(w, err)
	}
}

// GetCustomEmojiImage returns a single emoji image.
func GetCustomEmojiImage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/img/emoji/")
	r.URL.Path = path

	var emojiStaticServer http.Handler
	if data.GetUseCustomEmojis() {
		emojiFS := os.DirFS(config.CustomEmojiPath)
		emojiStaticServer = http.FileServer(http.FS(emojiFS))
	} else {
		emojiStaticServer = http.FileServer(http.FS(static.GetEmoji()))
	}

	emojiStaticServer.ServeHTTP(w, r)
}
