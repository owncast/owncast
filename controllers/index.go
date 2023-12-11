package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

var (
	indexCacheAdapter   *cache.Adapter
	indexBotSearchCache *cache.Client
)

// IndexHandler handles the default index route.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	if indexCacheAdapter == nil {
		ca, err := memory.NewAdapter(
			memory.AdapterWithAlgorithm(memory.LFU),
			memory.AdapterWithCapacity(50),
			memory.AdapterWithStorageCapacity(104_857_600),
		)
		indexCacheAdapter = &ca
		if err != nil {
			log.Warn("unable to create web cache", err)
		}
	}

	if indexBotSearchCache == nil {
		rc, err := cache.NewClient(
			cache.ClientWithAdapter(*indexCacheAdapter),
			cache.ClientWithTTL(30*time.Second),
			cache.ClientWithExpiresHeader(),
			cache.ClientWithRefreshKey("bot-search-page"),
			cache.ClientWithExpiresHeader(),
		)
		indexBotSearchCache = rc
		if err != nil {
			log.Warn("unable to create web cache client", err)
		}
	}

	isIndexRequest := r.URL.Path == "/" || filepath.Base(r.URL.Path) == "index.html" || filepath.Base(r.URL.Path) == ""

	if utils.IsUserAgentAPlayer(r.UserAgent()) && isIndexRequest {
		http.Redirect(w, r, "/hls/stream.m3u8", http.StatusTemporaryRedirect)
		return
	}

	// For search engine bots and social scrapers return a special
	// server-rendered page.
	if utils.IsUserAgentABot(r.UserAgent()) && isIndexRequest {
		indexBotSearchCache.Middleware(http.HandlerFunc(handleScraperMetadataPage))
		return
	}

	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)

	nonceRandom, _ := utils.GenerateRandomString(5)

	// Set our global HTTP headers
	middleware.SetHeaders(w, fmt.Sprintf("nonce-%s", nonceRandom))

	if isIndexRequest {
		renderIndexHtml(w, nonceRandom)
		return
	}

	serveWeb(w, r)
}

func renderIndexHtml(w http.ResponseWriter, nonce string) {
	type serverSideContent struct {
		Name             string
		Summary          string
		RequestedURL     string
		TagsString       string
		ThumbnailURL     string
		Thumbnail        string
		Image            string
		StatusJSON       string
		ServerConfigJSON string
		EmbedVideo       string
		Nonce            string
	}

	status := getStatusResponse()
	sb, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config := getConfigResponse()
	cb, err := json.Marshal(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := serverSideContent{
		Name:             data.GetServerName(),
		Summary:          data.GetServerSummary(),
		RequestedURL:     fmt.Sprintf("%s%s", data.GetServerURL(), "/"),
		TagsString:       strings.Join(data.GetServerMetadataTags(), ","),
		ThumbnailURL:     "thumbnail.jpg",
		Thumbnail:        "thumbnail.jpg",
		Image:            "logo/external",
		StatusJSON:       string(sb),
		ServerConfigJSON: string(cb),
		EmbedVideo:       "embed/video",
		Nonce:            nonce,
	}

	index, err := static.GetWebIndexTemplate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := index.Execute(w, content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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

// Return a basic HTML page with server-rendered metadata from the config
// to give to Opengraph clients and web scrapers (bots, web crawlers, etc).
func handleScraperMetadataPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := static.GetBotMetadataTemplate()
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	scheme := "http"

	if siteURL := data.GetServerURL(); siteURL != "" {
		if parsed, err := url.Parse(siteURL); err == nil && parsed.Scheme != "" {
			scheme = parsed.Scheme
		}
	}

	fullURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path))
	if err != nil {
		log.Errorln(err)
	}
	imageURL, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, r.Host, "/logo/external"))
	if err != nil {
		log.Errorln(err)
	}

	status := core.GetStatus()

	// If the thumbnail does not exist or we're offline then just use the logo image
	var thumbnailURL string
	if status.Online && utils.DoesFileExists(filepath.Join(config.DataDirectory, "tmp", "thumbnail.jpg")) {
		thumbnail, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, r.Host, "/thumbnail.jpg"))
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

	// Set a cache header
	middleware.SetCachingHeaders(w, r)

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, metadata); err != nil {
		log.Errorln(err)
	}
}
