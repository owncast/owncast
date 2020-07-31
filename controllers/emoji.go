package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabek/owncast/models"
	log "github.com/sirupsen/logrus"
)

// Make this path configurable if somebody has a valid reason
// to need it to be.  The config is getting a bit bloated.
const emojiPath = "/img/emoji" // Relative to webroot

//GetCustomEmoji returns a list of custom emoji via the API
func GetCustomEmoji(w http.ResponseWriter, r *http.Request) {
	emojiList := make([]models.CustomEmoji, 0)

	fullPath := filepath.Join("webroot", emojiPath)
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
		path := filepath.Join(emojiPath, f.Name())
		singleEmoji := models.CustomEmoji{name, path}
		emojiList = append(emojiList, singleEmoji)
	}

	json.NewEncoder(w).Encode(emojiList)
}
