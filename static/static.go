package static

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

//go:embed web/*
//go:embed web/_next/static
//go:embed web/_next/static/chunks/pages/*.js
//go:embed web/_next/static/*/*.js
var webFiles embed.FS

// GetWeb will return an embedded filesystem reference to the admin web app.
func GetWeb() embed.FS {
	return webFiles
}

// GetWebIndexTemplate will return the bot/scraper metadata template.
func GetWebIndexTemplate() (*template.Template, error) {
	webFiles := GetWeb()
	name := "web/index.html"
	t, err := template.ParseFS(webFiles, name)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open index.html template")
	}

	tmpl := template.Must(t, err)
	return tmpl, err
}

//go:embed offline.ts
var offlineVideoSegment []byte

// GetOfflineSegment will return the offline video segment data.
func GetOfflineSegment() []byte {
	return getFileSystemStaticFileOrDefault("offline.ts", offlineVideoSegment)
}

func getFileSystemStaticFileOrDefault(path string, defaultData []byte) []byte {
	fullPath := filepath.Join("static", path)
	data, err := os.ReadFile(fullPath) //nolint: gosec
	if err != nil {
		return defaultData
	}

	return data
}
