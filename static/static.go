package static

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
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
func GetWeb() fs.FS {
	wf, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	return wf
}

//go:embed img/emoji/*
var emojiFiles embed.FS

// GetEmoji will return the emoji files.
func GetEmoji() fs.FS {
	ef, err := fs.Sub(emojiFiles, "img/emoji")
	if err != nil {
		log.Fatal(err)
	}
	return ef
}

// GetWebIndexTemplate will return the bot/scraper metadata template.
func GetWebIndexTemplate() (*template.Template, error) {
	webFiles := GetWeb()
	name := "index.html"
	t, err := template.ParseFS(webFiles, name)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open index.html template")
	}

	tmpl := template.Must(t, err)
	return tmpl, err
}

//go:embed offline-v2.ts
var offlineVideoSegment []byte

// GetOfflineSegment will return the offline video segment data.
func GetOfflineSegment() []byte {
	return getFileSystemStaticFileOrDefault("offline-v2.ts", offlineVideoSegment)
}

//go:embed img/logo.png
var logo []byte

// GetLogo will return the logo data.
func GetLogo() []byte {
	return getFileSystemStaticFileOrDefault("img/logo.png", logo)
}

func getFileSystemStaticFileOrDefault(path string, defaultData []byte) []byte {
	fullPath := filepath.Join("static", path)
	data, err := os.ReadFile(fullPath) //nolint: gosec
	if err != nil {
		return defaultData
	}

	return data
}

//go:embed metadata.html.tmpl
var botMetadataTemplate embed.FS

// GetBotMetadataTemplate will return the bot/scraper metadata template.
func GetBotMetadataTemplate() (*template.Template, error) {
	name := "metadata.html.tmpl"
	t, err := template.ParseFS(botMetadataTemplate, name)
	tmpl := template.Must(t, err)
	return tmpl, err
}
