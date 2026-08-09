package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/router/middleware"
)

// clipPageDocument is the exported standalone clip page. Every /clips route is
// served from it; the page reads the clip id out of the path.
const clipPageDocument = "clips/index.html"

// clipMetadataPage is the server-rendered metadata for a shared clip link.
type clipMetadataPage struct {
	Name            string
	Title           string
	Summary         string
	RequestedURL    string
	Thumbnail       string
	PlaylistURL     string
	DurationSeconds int
}

// GetClipPage serves the viewer-facing page for a single clip at
// /clips/{clipId}.
//
// Browsers get the standalone clip page, which plays the clip named in the
// URL. Social scrapers get a server-rendered page carrying that clip's own
// OpenGraph metadata, since the web app is a static export and cannot render
// per-clip meta tags itself.
func (h *Handlers) GetClipPage(w http.ResponseWriter, r *http.Request) {
	if !h.replayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	middleware.EnableCors(w)

	clipID := strings.Trim(strings.TrimPrefix(strings.Trim(r.URL.Path, "/"), "clips"), "/")

	if clipID == "" || !utils.IsUserAgentABot(r.UserAgent()) {
		h.serveClipPageDocument(w, r)
		return
	}

	clip, err := h.replays.NewPlaylistGenerator().GetClip(clipID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tmpl, err := static.GetClipMetadataTemplate()
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if siteURL := h.configRepository.GetServerURL(); siteURL != "" {
		if parsed, err := url.Parse(siteURL); err == nil && parsed.Scheme != "" {
			scheme = parsed.Scheme
		}
	}
	base := fmt.Sprintf("%s://%s", scheme, r.Host)

	title := clip.ClipTitle
	if title == "" {
		title = "Clip"
	}

	// Fall back to the server logo when the clip has no poster, so the link
	// still unfurls with an image.
	thumbnail := base + "/logo/external"
	if clip.Thumbnail != "" {
		thumbnail = base + clip.Thumbnail
	}

	summary := clip.StreamTitle
	if summary == "" {
		summary = h.configRepository.GetServerSummary()
	}

	page := clipMetadataPage{
		Name:            h.configRepository.GetServerName(),
		Title:           title,
		Summary:         summary,
		RequestedURL:    base + "/clips/" + clipID,
		Thumbnail:       thumbnail,
		PlaylistURL:     base + "/clip/" + clipID,
		DurationSeconds: clip.DurationSeconds,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Errorln(err)
	}
}

// serveClipPageDocument returns the exported clip page for any /clips route so
// the client-side app can take over. The static export ships one document per
// route, so every clip is served from the same file.
func (h *Handlers) serveClipPageDocument(w http.ResponseWriter, r *http.Request) {
	middleware.SetCachingHeaders(w, r)

	nonceRandom, _ := utils.GenerateRandomString(5)
	middleware.SetHeaders(w, fmt.Sprintf("nonce-%s", nonceRandom))

	document, err := static.GetWebFile(clipPageDocument)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(document); err != nil {
		log.Errorln(err)
	}
}

// clipThumbnailServer serves files out of the clip poster directory. http.Dir
// confines every lookup to that directory, so a request cannot escape it.
var clipThumbnailServer = http.StripPrefix(
	"/clip-thumbnail/",
	http.FileServer(http.Dir(config.ClipThumbnailsPath)),
)

// GetClipThumbnail serves a clip's poster image.
func (h *Handlers) GetClipThumbnail(w http.ResponseWriter, r *http.Request) {
	if !h.replayFeaturesEnabled() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Posters are flat .jpg files named after their clip; nothing else in the
	// directory is servable.
	if filepath.Ext(r.URL.Path) != ".jpg" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	middleware.EnableCors(w)
	w.Header().Set("Cache-Control", "public, max-age=604800")
	clipThumbnailServer.ServeHTTP(w, r)
}
