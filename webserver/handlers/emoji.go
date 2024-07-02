package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/webserver/router/middleware"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// GetCustomEmojiList returns a list of emoji via the API.
func GetCustomEmojiList(w http.ResponseWriter, r *http.Request) {
	emojiList := data.GetEmojiList()
	middleware.SetCachingHeaders(w, r)

	if err := json.NewEncoder(w).Encode(emojiList); err != nil {
		webutils.InternalErrorHandler(w, err)
	}
}

// GetCustomEmojiImage returns a single emoji image.
func GetCustomEmojiImage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/img/emoji/")
	r.URL.Path = path

	emojiFS := os.DirFS(config.CustomEmojiPath)
	middleware.SetCachingHeaders(w, r)
	http.FileServer(http.FS(emojiFS)).ServeHTTP(w, r)
}
