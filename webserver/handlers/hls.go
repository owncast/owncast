package handlers

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/router/middleware"
)

// HandleHLSRequest will manage all requests to HLS content.
func (h *Handlers) HandleHLSRequest(w http.ResponseWriter, r *http.Request) {
	// Sanity check to limit requests to HLS file types.
	if filepath.Ext(r.URL.Path) != ".m3u8" && filepath.Ext(r.URL.Path) != ".ts" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	requestedPath := r.URL.Path
	relativePath := strings.TrimPrefix(requestedPath, "/hls/")
	if !filepath.IsLocal(relativePath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If using external storage then only allow requests for the
	// master playlist at stream.m3u8, no variants or segments.
	if h.configRepository.GetS3Config().Enabled && relativePath != "stream.m3u8" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Players that support CMCD attach their own playback measurements to
	// media requests; harvest them from every request that carries them.
	// The CMCD identity, when present, is the client's metrics identity.
	requestClientID := utils.GenerateClientIDFromRequest(r)
	cmcdKeys := parseCMCDRequest(r)
	if cmcdKeys != nil {
		h.registerCMCDKeys(cmcdClientID(r, cmcdKeys), requestClientID, cmcdKeys)
	}

	isPlaylist := path.Ext(r.URL.Path) == ".m3u8"
	if isPlaylist {
		// Playlists should never be cached.
		middleware.DisableCache(w)

		// Force the correct content type
		w.Header().Set("Content-Type", "application/x-mpegURL")

		// Use this as an opportunity to mark this viewer as active.
		viewer := models.GenerateViewerFromRequest(r)
		h.stream.SetViewerActive(&viewer)
	} else {
		cacheTime := utils.GetCacheDurationSecondsForPath(relativePath)
		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(cacheTime))
	}
	middleware.EnableCors(w)

	// Time online segment transfers as a server-side observation. Playlist
	// requests and offline streams cannot be observed, so avoid their metrics
	// lock acquisition on this hot path.
	observe := !isPlaylist && h.stream.GetStatus().Online
	var observation metrics.ServerObservation
	if observe {
		observation = h.metrics.ServerObservationFor(requestClientID)
		observe = !observation.Suppressed
	}

	file, err := os.OpenInRoot(config.HLSStoragePath, relativePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cw := &countingResponseWriter{ResponseWriter: w}
	start := time.Now()
	http.ServeContent(cw, r, relativePath, info.ModTime(), file)
	if observe {
		h.registerServedSegmentMetrics(r, observation.ClientID, requestClientID, cw.bytes, time.Since(start), observation.SpeedKnown)
	}
}
