package handlers

import (
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/owncast/owncast/config"
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
	relativePath := strings.Replace(requestedPath, "/hls/", "", 1)
	fullPath := filepath.Join(config.HLSStoragePath, relativePath)

	// Defense-in-depth path-traversal guard. http.ServeFile already rejects
	// URL.Path entries containing "..", but we resolve both paths to
	// absolute form and confirm the request stays inside HLSStoragePath
	// rather than relying on that stdlib detail. HLS asset names never
	// contain ".." in practice.
	absBase, err := filepath.Abs(config.HLSStoragePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	absFull, err := filepath.Abs(fullPath)
	if err != nil || (absFull != absBase && !strings.HasPrefix(absFull, absBase+string(filepath.Separator))) {
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
	metricsID := requestClientID
	cmcdKeys := parseCMCDRequest(r)
	if cmcdKeys != nil {
		metricsID = cmcdClientID(r, cmcdKeys)
		h.registerCMCDKeys(metricsID, cmcdKeys)
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

	// Time segment transfers as a server-side observation: for players
	// that report nothing themselves (Safari/iOS native HLS, VLC, mpv, ...)
	// it provides both speed and download duration; for request-mode-only
	// CMCD players it provides the download duration while their reported
	// throughput owns the speed metric. Clients that self-report through
	// the metrics API or the CMCD collector measure their own downloads,
	// so no server observation is needed for them.
	observe := !isPlaylist &&
		h.stream.GetStatus().Online &&
		!h.metrics.IsClientSelfReporting(requestClientID)

	cw := &countingResponseWriter{ResponseWriter: w}
	start := time.Now()
	http.ServeFile(cw, r, absFull) //nolint:gosec // G703: absFull was resolved and verified to be inside HLSStoragePath above
	if observe {
		h.registerServedSegmentMetrics(r, metricsID, cw.bytes, time.Since(start), cmcdKeys != nil)
	}
}
