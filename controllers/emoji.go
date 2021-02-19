package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// Make this path configurable if somebody has a valid reason
// to need it to be.  The config is getting a bit bloated.
const emojiDir = "/img/emoji" // Relative to webroot

// GetCustomEmoji returns a list of custom emoji via the API.
func GetCustomEmoji(w http.ResponseWriter, r *http.Request) {
	emojiList := make([]models.CustomEmoji, 0)

	fullPath := filepath.Join(config.WebRoot, emojiDir)
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		log.Errorln(err)
		// Throw HTTP 500
		return
	}

	// Memoize this result somewhere?  Right now it iterates through the
	// filesystem every time the API is hit.  Should you need to restart
	// the server to add emoji?
	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
		emojiPath := filepath.Join(emojiDir, f.Name())
		singleEmoji := models.CustomEmoji{Name: name, Emoji: emojiPath}
		emojiList = append(emojiList, singleEmoji)
	}

	if err := json.NewEncoder(w).Encode(emojiList); err != nil {
		InternalErrorHandler(w, err)
	}
}
