package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
)

// GetCustomEmojiList returns a list of emoji via the API.
func (s *Service) GetCustomEmojiList(w http.ResponseWriter, r *http.Request) {
	emojiList := data.GetEmojiList()

	if err := json.NewEncoder(w).Encode(emojiList); err != nil {
		s.InternalErrorHandler(w, err)
	}
}

// GetCustomEmojiImage returns a single emoji image.
func (s *Service) GetCustomEmojiImage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/img/emoji/")
	r.URL.Path = path

	emojiFS := os.DirFS(config.CustomEmojiPath)
	http.FileServer(http.FS(emojiFS)).ServeHTTP(w, r)
}
