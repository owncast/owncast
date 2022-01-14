package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

var _hasWarnedSVGLogo = false

// GetLogo will return the logo image as a response.
func GetLogo(w http.ResponseWriter, r *http.Request) {
	imageFilename := data.GetLogoPath()
	if imageFilename == "" {
		returnDefault(w)
		return
	}
	imagePath := filepath.Join("data", imageFilename)
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
	writeBytesAsImage(imageBytes, contentType, w, cacheTime)
}

// GetCompatibleLogo will return the logo unless it's a SVG
// and in that case will return a default placeholder.
// Used for sharing to external social networks that generally
// don't support SVG.
func GetCompatibleLogo(w http.ResponseWriter, r *http.Request) {
	imageFilename := data.GetLogoPath()

	// If the logo image is not a SVG then we can return it
	// without any problems.
	if imageFilename != "" && filepath.Ext(imageFilename) != ".svg" {
		GetLogo(w, r)
		return
	}

	// Otherwise use a fallback logo.png.
	imagePath := filepath.Join(config.WebRoot, "img", "logo.png")
	contentType := "image/png"
	imageBytes, err := getImage(imagePath)
	if err != nil {
		returnDefault(w)
		return
	}

	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	writeBytesAsImage(imageBytes, contentType, w, cacheTime)

	referrer := r.Referer()
	if referrer == "" {
		referrer = "an external site"
	}
	if !_hasWarnedSVGLogo {
		log.Warnf("%s requested your logo. because many social networks do not support SVGs we returned a placeholder instead. change your current logo to a png or jpeg to be most compatible with external social networking sites.", referrer)
		_hasWarnedSVGLogo = true
	}
}

func returnDefault(w http.ResponseWriter) {
	imagePath := filepath.Join(config.WebRoot, "img", "logo.svg")
	imageBytes, err := getImage(imagePath)
	if err != nil {
		log.Errorln(err)
		return
	}
	cacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	writeBytesAsImage(imageBytes, "image/svg+xml", w, cacheTime)
}

func writeBytesAsImage(data []byte, contentType string, w http.ResponseWriter, cacheSeconds int) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(cacheSeconds))

	if _, err := w.Write(data); err != nil {
		log.Println("unable to write image.")
	}
}

func getImage(path string) ([]byte, error) {
	return os.ReadFile(path) // nolint
}
