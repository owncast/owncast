package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/webserver/responses"
)

// GetCustomEmojiList returns a list of emoji via the API.
func (h *Handlers) GetCustomEmojiList(w http.ResponseWriter, r *http.Request) {
	emojiList := data.GetEmojiList()
	middleware.SetCachingHeaders(w, r)

	if err := json.NewEncoder(w).Encode(emojiList); err != nil {
		responses.InternalErrorHandler(w, err)
	}
}

// GetCustomEmojiImage returns a single emoji image.
func (h *Handlers) GetCustomEmojiImage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/img/emoji/")
	r.URL.Path = path

	c := config.GetConfig()
	emojiFS := os.DirFS(c.CustomEmojiPath)
	middleware.SetCachingHeaders(w, r)

	http.FileServer(http.FS(emojiFS)).ServeHTTP(w, r)
}
