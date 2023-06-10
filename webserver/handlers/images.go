package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/responses"
)

const (
	contentTypeJPEG = "image/jpeg"
	contentTypeGIF  = "image/gif"
)

// GetThumbnail will return the thumbnail image as a response.
func (h *Handlers) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	imageFilename := "thumbnail.jpg"
	imagePath := filepath.Join(config.TempDir, imageFilename)

	var imageBytes []byte
	var err error

	if utils.DoesFileExists(imagePath) {
		imageBytes, err = getImage(imagePath)
	} else {
		h.GetLogo(w, r)
		return
	}

	if err != nil {
		h.GetLogo(w, r)
		return
	}

	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	responses.WriteBytesAsImage(imageBytes, contentTypeJPEG, w, cacheTime)
}

// GetPreview will return the preview gif as a response.
func (h *Handlers) GetPreview(w http.ResponseWriter, r *http.Request) {
	imageFilename := "preview.gif"
	imagePath := filepath.Join(config.TempDir, imageFilename)

	var imageBytes []byte
	var err error

	if utils.DoesFileExists(imagePath) {
		imageBytes, err = getImage(imagePath)
	} else {
		h.GetLogo(w, r)
		return
	}

	if err != nil {
		h.GetLogo(w, r)
		return
	}

	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	responses.WriteBytesAsImage(imageBytes, contentTypeGIF, w, cacheTime)
}
