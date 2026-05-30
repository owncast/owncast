package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
)

// GetFavicon will return the favicon image as a response.
func (h *Handlers) GetFavicon(w http.ResponseWriter, r *http.Request) {
	faviconFilename := h.configRepository.GetFaviconPath()
	if faviconFilename == "" {
		returnDefaultFavicon(w)
		return
	}

	faviconPath := filepath.Join(config.DataDirectory, faviconFilename)
	faviconBytes, err := os.ReadFile(faviconPath) //nolint:gosec
	if err != nil {
		returnDefaultFavicon(w)
		return
	}

	contentType := "image/x-icon"
	if filepath.Ext(faviconFilename) == ".png" {
		contentType = "image/png" //nolint:goconst
	}

	cacheTime := utils.GetCacheDurationSecondsForPath(faviconPath)
	writeFaviconResponse(faviconBytes, contentType, w, cacheTime)
}

func returnDefaultFavicon(w http.ResponseWriter) {
	faviconBytes := static.GetFavicon()
	cacheTime := utils.GetCacheDurationSecondsForPath("favicon.png")
	writeFaviconResponse(faviconBytes, "image/png", w, cacheTime)
}

func writeFaviconResponse(data []byte, contentType string, w http.ResponseWriter, cacheSeconds int) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(cacheSeconds))

	if _, err := w.Write(data); err != nil {
		log.Println("unable to write favicon.")
	}
}
