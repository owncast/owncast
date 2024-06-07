package controllers

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/utils"
)

const (
	contentTypeJPEG = "image/jpeg"
	contentTypeGIF  = "image/gif"
)

var previewThumbCache = ttlcache.New(
	ttlcache.WithTTL[string, []byte](15),
	ttlcache.WithCapacity[string, []byte](1),
	ttlcache.WithDisableTouchOnHit[string, []byte](),
)

// GetThumbnail will return the thumbnail image as a response.
func GetThumbnail(w http.ResponseWriter, r *http.Request) {
	imageFilename := "thumbnail.jpg"
	imagePath := filepath.Join(config.TempDir, imageFilename)
	httpCacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	inMemoryCacheTime := time.Duration(15) * time.Second

	var imageBytes []byte
	var err error

	if previewThumbCache.Get(imagePath) != nil {
		ci := previewThumbCache.Get(imagePath)
		imageBytes = ci.Value()
	} else if utils.DoesFileExists(imagePath) {
		imageBytes, err = getImage(imagePath)
		previewThumbCache.Set(imagePath, imageBytes, inMemoryCacheTime)
	} else {
		GetLogo(w, r)
		return
	}

	if err != nil {
		GetLogo(w, r)
		return
	}

	writeBytesAsImage(imageBytes, contentTypeJPEG, w, httpCacheTime)
}

// GetPreview will return the preview gif as a response.
func GetPreview(w http.ResponseWriter, r *http.Request) {
	imageFilename := "preview.gif"
	imagePath := filepath.Join(config.TempDir, imageFilename)
	httpCacheTime := utils.GetCacheDurationSecondsForPath(imagePath)
	inMemoryCacheTime := time.Duration(15) * time.Second

	var imageBytes []byte
	var err error

	if previewThumbCache.Get(imagePath) != nil {
		ci := previewThumbCache.Get(imagePath)
		imageBytes = ci.Value()
	} else if utils.DoesFileExists(imagePath) {
		imageBytes, err = getImage(imagePath)
		previewThumbCache.Set(imagePath, imageBytes, inMemoryCacheTime)
	} else {
		GetLogo(w, r)
		return
	}

	if err != nil {
		GetLogo(w, r)
		return
	}

	writeBytesAsImage(imageBytes, contentTypeGIF, w, httpCacheTime)
}
