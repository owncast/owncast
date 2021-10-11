package static

import "embed"

//go:embed admin/*
//go:embed admin/_next/static
//go:embed admin/_next/static/chunks/pages/*.js
//go:embed admin/_next/static/*/*.js
var adminFiles embed.FS

// GetAdmin will return an embedded filesystem reference to the admin web app.
func GetAdmin() embed.FS {
	return adminFiles
}
