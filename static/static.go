package static

import (
	"embed"
	"html/template"
)

//go:embed admin/*
//go:embed admin/_next/static
//go:embed admin/_next/static/chunks/pages/*.js
//go:embed admin/_next/static/*/*.js
var adminFiles embed.FS

// GetAdmin will return an embedded filesystem reference to the admin web app.
func GetAdmin() embed.FS {
	return adminFiles
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

//go:embed offline.ts
var offlineVideoSegment []byte

// GetOfflineSegment will return the offline video segment data.
func GetOfflineSegment() []byte {
	return offlineVideoSegment
}
