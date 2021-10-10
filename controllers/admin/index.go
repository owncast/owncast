package admin

import (
	"bytes"
	_ "embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/static"
)

// ServeAdmin will return admin web assets.
func ServeAdmin(w http.ResponseWriter, r *http.Request) {
	adminFiles := static.GetAdmin()

	// Determine if the requested path is a directory.
	// If so, append index.html to the request.
	path := strings.TrimPrefix(r.URL.Path, "/")
	if strings.HasSuffix(path, "/") {
		path = filepath.Join(path, "index.html")
	}

	f, err := adminFiles.Open(path)
	info, err := f.Stat()

	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Set a cache control max-age header
	middleware.SetCachingHeaders(w, r)

	d, err := adminFiles.ReadFile(path)
	if err != nil {
		panic(err)
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), bytes.NewReader(d))
}

func errorHandler(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
