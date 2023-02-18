package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/owncast/owncast/app/middleware"
	"github.com/owncast/owncast/static"
	"github.com/owncast/owncast/utils"
)

// IndexHandler handles the default index route.
func (s *Service) IndexHandler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	isIndexRequest := r.URL.Path == "/" || filepath.Base(r.URL.Path) == "index.html" || filepath.Base(r.URL.Path) == ""

	if utils.IsUserAgentAPlayer(r.UserAgent()) && isIndexRequest {
		http.Redirect(w, r, "/hls/stream.m3u8", http.StatusTemporaryRedirect)
		return
	}

	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)

	nonceRandom, _ := utils.GenerateRandomString(5)

	// Set our global HTTP headers
	middleware.SetHeaders(w, fmt.Sprintf("nonce-%s", nonceRandom))

	if isIndexRequest {
		s.renderIndexHtml(w, nonceRandom)
		return
	}

	s.serveWeb(w, r)
}

func (s *Service) renderIndexHtml(w http.ResponseWriter, nonce string) {
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
		Nonce            string
	}

	status := s.getStatusResponse()
	sb, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config := s.getConfigResponse()
	cb, err := json.Marshal(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := serverSideContent{
		Name:             s.Data.GetServerName(),
		Summary:          s.Data.GetServerSummary(),
		RequestedURL:     s.Data.GetServerURL(),
		TagsString:       strings.Join(s.Data.GetServerMetadataTags(), ","),
		ThumbnailURL:     "/thumbnail",
		Thumbnail:        "/thumbnail",
		Image:            "/logo/external",
		StatusJSON:       string(sb),
		ServerConfigJSON: string(cb),
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
