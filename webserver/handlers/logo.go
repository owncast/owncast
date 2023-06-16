package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/responses"
	log "github.com/sirupsen/logrus"
)

var _hasWarnedSVGLogo = false

// GetLogo will return the logo image as a response.
func (h *Handlers) GetLogo(w http.ResponseWriter, r *http.Request) {
	imageFilename := data.GetLogoPath()
	if imageFilename == "" {
		returnDefault(w)
		return
	}
	imagePath := filepath.Join(config.DataDirectory, imageFilename)
	imageBytes, err := getImage(imagePath)
	if err != nil {
		returnDefault(w)
		return
	}

	contentType := "image/jpeg"
	if filepath.Ext(imageFilename) == ".svg" {
		contentType = "image/svg+xml"
	} else if filepath.Ext(imageFilename) == ".gif" {
		contentType = "image/gif"
	} else if filepath.Ext(imageFilename) == ".png" {
		contentType = "image/png"
	}

	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	responses.WriteBytesAsImage(imageBytes, contentType, w, cacheTime)
}

// GetCompatibleLogo will return the logo unless it's a SVG
// and in that case will return a default placeholder.
// Used for sharing to external social networks that generally
// don't support SVG.
func (h *Handlers) GetCompatibleLogo(w http.ResponseWriter, r *http.Request) {
	imageFilename := data.GetLogoPath()

	// If the logo image is not a SVG then we can return it
	// without any problems.
	if imageFilename != "" && filepath.Ext(imageFilename) != ".svg" {
		h.GetLogo(w, r)
		return
	}

	// Otherwise use a fallback logo.png.
	imagePath := filepath.Join(config.DataDirectory, "logo.png")
	contentType := "image/png"
	imageBytes, err := getImage(imagePath)
	if err != nil {
		returnDefault(w)
		return
	}

	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	responses.WriteBytesAsImage(imageBytes, contentType, w, cacheTime)

	if !_hasWarnedSVGLogo {
		log.Warnf("an external site requested your logo. because many social networks do not support SVGs we returned a placeholder instead. change your current logo to a png or jpeg to be most compatible with external social networking sites.")
		_hasWarnedSVGLogo = true
	}
}

func returnDefault(w http.ResponseWriter) {
	imageBytes := static.GetLogo()
	cacheTime := utils.GetCacheDurationSecondsForPath("logo.png")
	responses.WriteBytesAsImage(imageBytes, "image/png", w, cacheTime)
}

func getImage(path string) ([]byte, error) {
	return os.ReadFile(path) // nolint
}
