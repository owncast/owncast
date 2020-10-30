package admin

import (
	"io/ioutil"
	"net/http"

	"github.com/markbates/pkger"
	"github.com/owncast/owncast/router/middleware"
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
	if err == nil {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return
		}

		w.Write(b)
	}
}
