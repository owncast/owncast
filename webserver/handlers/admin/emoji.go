package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// UploadCustomEmoji allows POSTing a new custom emoji to the server.
func UploadCustomEmoji(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	emoji := new(generated.UploadCustomEmojiJSONBody)

	if err := json.NewDecoder(r.Body).Decode(emoji); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	bytes, _, err := utils.DecodeBase64Image(*emoji.Data)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Prevent path traversal attacks
	emojiFileName := filepath.Base(*emoji.Name)
	targetPath := filepath.Join(config.CustomEmojiPath, emojiFileName)

	err = os.MkdirAll(config.CustomEmojiPath, 0o700)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if utils.DoesFileExists(targetPath) {
		webutils.WriteSimpleResponse(w, false, fmt.Sprintf("An emoji with the name %q already exists", emojiFileName))
		return
	}

	if err = os.WriteFile(targetPath, bytes, 0o600); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, fmt.Sprintf("Emoji %q has been uploaded", emojiFileName))
}

// DeleteCustomEmoji deletes a custom emoji.
func DeleteCustomEmoji(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	type deleteEmoji struct {
		Name string `json:"name"`
	}

	emoji := new(deleteEmoji)

	if err := json.NewDecoder(r.Body).Decode(emoji); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	targetPath := filepath.Join(config.CustomEmojiPath, emoji.Name)

	if !filepath.IsLocal(targetPath) {
		webutils.WriteSimpleResponse(w, false, "Emoji path is not valid")
		return
	}

	if err := os.Remove(targetPath); err != nil {
		if os.IsNotExist(err) {
			webutils.WriteSimpleResponse(w, false, fmt.Sprintf("Emoji %q doesn't exist", emoji.Name))
		} else {
			webutils.WriteSimpleResponse(w, false, err.Error())
		}
		return
	}

	webutils.WriteSimpleResponse(w, true, fmt.Sprintf("Emoji %q has been deleted", emoji.Name))
}
