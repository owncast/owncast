package utils

import (
	"path"

	log "github.com/sirupsen/logrus"
)

// MigrateCustomEmojiLocations migrates custom emoji from the old location to the new location.
func MigrateCustomEmojiLocations() {
	oldLocation := path.Join("webroot", "img", "emoji")
	newLocation := path.Join("data", "emoji")

	if !DoesFileExists(oldLocation) {
		return
	}

	log.Println("Moving custom emoji directory from", oldLocation, "to", newLocation)

	if err := Move(oldLocation, newLocation); err != nil {
		log.Errorln("error moving custom emoji directory", err)
	}
}
