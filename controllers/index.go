package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/router/middleware"
	"github.com/gabek/owncast/utils"
)

type MetadataPage struct {
	Config       config.InstanceDetails
	RequestedURL string
	Image        string
	Thumbnail    string
	TagsString   string
}

//IndexHandler handles the default index route
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	isIndexRequest := r.URL.Path == "/" || r.URL.Path == "/index.html"

	// Reject requests for the web UI if it's disabled.
	if isIndexRequest && config.Config.DisableWebFeatures {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - y u ask 4 this?  If this is an error let us know: https://github.com/gabek/owncast/issues"))
		return
	}

	if utils.IsUserAgentABot(r.UserAgent()) && isIndexRequest {
		handleScraperMetadataPage(w, r)
		return
	}

	http.ServeFile(w, r, path.Join("webroot", r.URL.Path))

	if path.Ext(r.URL.Path) == ".m3u8" {
		middleware.DisableCache(&w)

		clientID := utils.GenerateClientIDFromRequest(r)
		core.SetClientActive(clientID)
	}
}

// Return a basic HTML page with server-rendered metadata from the config file
// to give to Opengraph clients and web scrapers (bots, web crawlers, etc).
func handleScraperMetadataPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join("static", "metadata.html")))

	fullURL, err := url.Parse(fmt.Sprintf("http://%s%s", r.Host, r.URL.Path))
	imageURL, err := url.Parse(fmt.Sprintf("http://%s%s", r.Host, config.Config.InstanceDetails.Logo["large"]))

	status := core.GetStatus()

	// If the thumbnail does not exist or we're offline then just use the logo image
	var thumbnailURL string
	if status.Online && utils.DoesFileExists("webroot/thumbnail.jpg") {
		thumbnail, err := url.Parse(fmt.Sprintf("http://%s%s", r.Host, "/thumbnail.jpg"))
		if err != nil {
			log.Errorln(err)
		}
		thumbnailURL = thumbnail.String()
	} else {
		thumbnailURL = imageURL.String()
	}

	tagsString := strings.Join(config.Config.InstanceDetails.Tags, ",")
	metadata := MetadataPage{config.Config.InstanceDetails, fullURL.String(), imageURL.String(), thumbnailURL, tagsString}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, metadata)

	if err != nil {
		log.Panicln(err)
	}
}
