package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
)

// MetadataPage represents a server-rendered web page for bots and web scrapers.
type MetadataPage struct {
	RequestedURL  string
	Image         string
	Thumbnail     string
	TagsString    string
	Summary       string
	Name          string
	Tags          []string
	SocialHandles []models.SocialHandle
}

// IndexHandler handles the default index route.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)
	isIndexRequest := r.URL.Path == "/" || filepath.Base(r.URL.Path) == "index.html" || filepath.Base(r.URL.Path) == ""

	// For search engine bots and social scrapers return a special
	// server-rendered page.
	if utils.IsUserAgentABot(r.UserAgent()) && isIndexRequest {
		handleScraperMetadataPage(w, r)
		return
	}

	// If the ETags match then return a StatusNotModified
	if responseCode := middleware.ProcessEtags(w, r); responseCode != 0 {
		w.WriteHeader(responseCode)
		return
	}

	// If this is a directory listing request then return a 404
	info, err := os.Stat(path.Join(config.WebRoot, r.URL.Path))
	if err != nil || (info.IsDir() && !isIndexRequest) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if path.Ext(r.URL.Path) == ".m3u8" {
		middleware.DisableCache(w)
	}

	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)

	// Opt-out of Google FLoC
	middleware.DisableFloc(w)

	http.ServeFile(w, r, path.Join(config.WebRoot, r.URL.Path))
}

// Return a basic HTML page with server-rendered metadata from the config file
// to give to Opengraph clients and web scrapers (bots, web crawlers, etc).
func handleScraperMetadataPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join("static", "metadata.html")))
	protocol := "http"

	if r.TLS != nil {
		protocol = "https"
	}

	fullURL, err := url.Parse(fmt.Sprintf("%s://%s%s", protocol, r.Host, r.URL.Path))
	if err != nil {
		log.Panicln(err)
	}
	imageURL, err := url.Parse(fmt.Sprintf("%s://%s%s", protocol, r.Host, "/logo"))
	if err != nil {
		log.Panicln(err)
	}

	status := core.GetStatus()

	// If the thumbnail does not exist or we're offline then just use the logo image
	var thumbnailURL string
	if status.Online && utils.DoesFileExists(filepath.Join(config.WebRoot, "thumbnail.jpg")) {
		thumbnail, err := url.Parse(fmt.Sprintf("%s://%s%s", protocol, r.Host, "/thumbnail.jpg"))
		if err != nil {
			log.Errorln(err)
			thumbnailURL = imageURL.String()
		} else {
			thumbnailURL = thumbnail.String()
		}
	} else {
		thumbnailURL = imageURL.String()
	}

	tagsString := strings.Join(data.GetServerMetadataTags(), ",")
	metadata := MetadataPage{
		Name:          data.GetServerName(),
		RequestedURL:  fullURL.String(),
		Image:         imageURL.String(),
		Summary:       data.GetServerSummary(),
		Thumbnail:     thumbnailURL,
		TagsString:    tagsString,
		Tags:          data.GetServerMetadataTags(),
		SocialHandles: data.GetSocialHandles(),
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, metadata); err != nil {
		log.Panicln(err)
	}
}
