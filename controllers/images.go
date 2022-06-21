package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

// GetThumbnail will return the thumbnail image as a response.
func GetThumbnail(w http.ResponseWriter, r *http.Request) {
	imageFilename := "thumbnail.jpg"
	imagePath := filepath.Join(config.TempDir, imageFilename)

	var imageBytes []byte
	var err error

	if utils.DoesFileExists(imagePath) {
		imageBytes, err = getImage(imagePath)
	} else {
		GetLogo(w, r)
		return
	}

	if err != nil {
		GetLogo(w, r)
		return
	}

	contentType := "image/jpeg"
	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	writeBytesAsImage(imageBytes, contentType, w, cacheTime)
}

// GetPreview will return the preview gif as a response.
func GetPreview(w http.ResponseWriter, r *http.Request) {
	imageFilename := "preview.gif"
	imagePath := filepath.Join(config.TempDir, imageFilename)

	var imageBytes []byte
	var err error

	if utils.DoesFileExists(imagePath) {
		imageBytes, err = getImage(imagePath)
	} else {
		GetLogo(w, r)
		return
	}

	if err != nil {
		GetLogo(w, r)
		return
	}

	contentType := "image/jpeg"
	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	writeBytesAsImage(imageBytes, contentType, w, cacheTime)
}
