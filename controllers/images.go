package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

const (
	contentTypeJPEG = "image/jpeg"
	contentTypeGIF  = "image/gif"
)

// GetThumbnail will return the thumbnail image as a response.
func (s *Service) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	imageFilename := "thumbnail.jpg"
	imagePath := filepath.Join(config.TempDir, imageFilename)

	var imageBytes []byte
	var err error

	if utils.DoesFileExists(imagePath) {
		imageBytes, err = s.getImage(imagePath)
	} else {
		s.GetLogo(w, r)
		return
	}

	if err != nil {
		s.GetLogo(w, r)
		return
	}

	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	s.writeBytesAsImage(imageBytes, contentTypeJPEG, w, cacheTime)
}

// GetPreview will return the preview gif as a response.
func (s *Service) GetPreview(w http.ResponseWriter, r *http.Request) {
	imageFilename := "preview.gif"
	imagePath := filepath.Join(config.TempDir, imageFilename)

	var imageBytes []byte
	var err error

	if utils.DoesFileExists(imagePath) {
		imageBytes, err = s.getImage(imagePath)
	} else {
		s.GetLogo(w, r)
		return
	}

	if err != nil {
		s.GetLogo(w, r)
		return
	}

	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	s.writeBytesAsImage(imageBytes, contentTypeGIF, w, cacheTime)
}
