package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	if emoji.Data == nil {
		webutils.WriteSimpleResponse(w, false, "data field is required")
		return
	}

	if emoji.Name == nil {
		webutils.WriteSimpleResponse(w, false, "name field is required")
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

	// The admin sends the emoji path relative to the emoji directory with a
	// leading slash, and os.Root only accepts relative names.
	name := strings.TrimPrefix(emoji.Name, "/")

	// Check the name that was sent rather than the result of joining it onto
	// the emoji directory. Join cleans the ../ away first, which is what let
	// traversal through the equivalent check here before.
	if !filepath.IsLocal(name) {
		webutils.WriteSimpleResponse(w, false, "Emoji path is not valid")
		return
	}

	// Remove through a root handle so the emoji directory boundary holds even
	// when a name that looks local resolves out of it through a symlink.
	root, err := os.OpenRoot(config.CustomEmojiPath)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}
	defer root.Close()

	if err := root.Remove(name); err != nil {
		if os.IsNotExist(err) {
			webutils.WriteSimpleResponse(w, false, fmt.Sprintf("Emoji %q doesn't exist", emoji.Name))
		} else {
			webutils.WriteSimpleResponse(w, false, err.Error())
		}
		return
	}

	webutils.WriteSimpleResponse(w, true, fmt.Sprintf("Emoji %q has been deleted", emoji.Name))
}
