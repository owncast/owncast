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

// ServeAdmin will return admin web assets.
func ServeAdmin(w http.ResponseWriter, r *http.Request) {
	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)

	// Determine if the requested path is a directory.
	// If so, append index.html to the request.
	path := r.URL.Path
	dirCheck, err := pkger.Stat(path)
	if dirCheck != nil && err == nil && dirCheck.IsDir() {
		path = filepath.Join(path, "index.html")
	}

	f, err := pkger.Open(path)
	if err != nil {
		log.Debugln(err, path)
		errorHandler(w, http.StatusNotFound)
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Warnln(err)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", mimeType)
	if _, err = w.Write(b); err != nil {
		log.Errorln(err)
	}
}

func errorHandler(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
