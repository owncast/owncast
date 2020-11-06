package admin

import (
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/markbates/pkger"
	"github.com/owncast/owncast/router/middleware"
	log "github.com/sirupsen/logrus"
)

// ServeAdmin will return admin web assets
func ServeAdmin(w http.ResponseWriter, r *http.Request) {
	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)

	path := r.URL.Path
	if path == "/admin" || path == "/admin/" {
		path = "/admin/index.html"
	}

	f, err := pkger.Open(path)
	if err != nil {
		log.Warnln(err, path)
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Warnln(err)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", mimeType)
	w.Write(b)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
}
