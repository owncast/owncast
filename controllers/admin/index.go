package admin

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/static"
	log "github.com/sirupsen/logrus"
)

// ServeAdmin will return admin web assets.
func ServeAdmin(w http.ResponseWriter, r *http.Request) {
	// If the ETags match then return a StatusNotModified
	if responseCode := middleware.ProcessEtags(w, r); responseCode != 0 {
		w.WriteHeader(responseCode)
		return
	}

	adminFiles := static.GetAdmin()
	path := strings.TrimPrefix(r.URL.Path, "/")

	// Determine if the requested path is a directory.
	// If so, append index.html to the request.
	if info, err := fs.Stat(adminFiles, path); err == nil && info.IsDir() {
		path = filepath.Join(path, "index.html")
	} else if _, err := fs.Stat(adminFiles, path+"index.html"); err == nil {
		path = filepath.Join(path, "index.html")
	}

	f, err := adminFiles.Open(path)
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	info, err := f.Stat()
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)
	d, err := adminFiles.ReadFile(path)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), bytes.NewReader(d))
}
