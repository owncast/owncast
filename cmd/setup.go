package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/owncast/owncast/logging"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (app *Application) createDirectories() {
	// Create the data directory if needed
	if !utils.DoesFileExists("data") {
		if err := os.Mkdir("./data", 0o700); err != nil {
			log.Fatalln("Cannot create data directory", err)
		}
	}

	// Recreate the temp dir
	if utils.DoesFileExists(app.configservice.TempDir) {
		err := os.RemoveAll(app.configservice.TempDir)
		if err != nil {
			log.Fatalln("Unable to remove temp dir! Check permissions.", app.configservice.TempDir, err)
		}
	}
	if err := os.Mkdir(app.configservice.TempDir, 0o700); err != nil {
		log.Fatalln("Unable to create temp dir!", err)
	}
}

func (app *Application) configureLogging(enableDebugFeatures bool, enableVerboseLogging bool, logDirectory string) {
	logging.Setup(enableDebugFeatures, enableVerboseLogging, logDirectory)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

// setupEmojiDirectory sets up the custom emoji directory by copying all built-in
// emojis if the directory does not yet exist.
func (app *Application) setupEmojiDirectory() (err error) {
	type emojiDirectory struct {
		path  string
		isDir bool
	}

	// Migrate old (pre 0.1.0) emoji to new location if they exist.
	app.migrateCustomEmojiLocations()

	if utils.DoesFileExists(app.configservice.CustomEmojiPath) {
		return nil
	}

	if err = os.MkdirAll(app.configservice.CustomEmojiPath, 0o750); err != nil {
		return fmt.Errorf("unable to create custom emoji directory: %w", err)
	}

	staticFS := static.GetEmoji()
	files := []emojiDirectory{}

	walkFunction := func(path string, d os.DirEntry, err error) error {
		if path == "." {
			return nil
		}

		if d.Name() == "LICENSE.md" {
			return nil
		}

		files = append(files, emojiDirectory{path: path, isDir: d.IsDir()})
		return nil
	}

	if err := fs.WalkDir(staticFS, ".", walkFunction); err != nil {
		log.Errorln("unable to fetch emojis: " + err.Error())
		return errors.Wrap(err, "unable to fetch embedded emoji files")
	}

	if err != nil {
		return fmt.Errorf("unable to read built-in emoji files: %w", err)
	}

	// Now copy all built-in emojis to the custom emoji directory
	for _, path := range files {
		emojiPath := filepath.Join(app.configservice.CustomEmojiPath, path.path)

		if path.isDir {
			if err := os.Mkdir(emojiPath, 0o700); err != nil {
				return errors.Wrap(err, "unable to create emoji directory, check permissions?: "+path.path)
			}
			continue
		}

		memFile, staticOpenErr := staticFS.Open(path.path)
		if staticOpenErr != nil {
			return errors.Wrap(staticOpenErr, "unable to open emoji file from embedded filesystem")
		}

		// nolint:gosec
		diskFile, err := os.Create(emojiPath)
		if err != nil {
			return fmt.Errorf("unable to create custom emoji file on disk: %w", err)
		}

		if err != nil {
			_ = diskFile.Close()
			return fmt.Errorf("unable to open built-in emoji file: %w", err)
		}

		if _, err = io.Copy(diskFile, memFile); err != nil {
			_ = diskFile.Close()
			_ = os.Remove(emojiPath)
			return fmt.Errorf("unable to copy built-in emoji file to disk: %w", err)
		}

		if err = diskFile.Close(); err != nil {
			_ = os.Remove(emojiPath)
			return fmt.Errorf("unable to close custom emoji file on disk: %w", err)
		}
	}

	return nil
}

// MigrateCustomEmojiLocations migrates custom emoji from the old location to the new location.
func (app *Application) migrateCustomEmojiLocations() {
	oldLocation := path.Join("webroot", "img", "emoji")
	newLocation := path.Join("data", "emoji")

	if !utils.DoesFileExists(oldLocation) {
		return
	}

	log.Println("Moving custom emoji directory from", oldLocation, "to", newLocation)

	if err := utils.Move(oldLocation, newLocation); err != nil {
		log.Errorln("error moving custom emoji directory", err)
	}
}

func (app *Application) resetDirectories() {
	log.Trace("Resetting file directories to a clean slate.")

	// Wipe hls data directory
	utils.CleanupDirectory(app.configservice.HLSStoragePath)

	// Remove the previous thumbnail
	logo := app.configRepository.GetLogoPath()
	if utils.DoesFileExists(logo) {
		err := utils.Copy(path.Join("data", logo), filepath.Join(config.DataDirectory, "thumbnail.jpg"))
		if err != nil {
			log.Warnln(err)
		}
	}
}
