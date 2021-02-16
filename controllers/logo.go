package controllers

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

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
	return ioutil.ReadFile(path)
}
