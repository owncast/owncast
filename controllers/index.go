package controllers

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/router/middleware"
	"github.com/gabek/owncast/utils"

	"xojoc.pw/useragent"
)

type MetadataPage struct {
	Config       config.InstanceDetails
	RequestedURL string
}

//IndexHandler handles the default index route
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	ua := useragent.Parse(r.UserAgent())

	if ua.Type == useragent.Crawler && r.URL.Path == "/" {
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

	fullURL := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)
	metadata := MetadataPage{config.Config.InstanceDetails, fullURL}

	w.Header().Set("Content-Type", "text/html")
	err := tmpl.Execute(w, metadata)

	if err != nil {
		log.Panicln(err)
	}
}
