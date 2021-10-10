package static

import "embed"

//go:embed admin/*
//go:embed admin/_next/static
//go:embed admin/_next/static/chunks/pages/*.js
//go:embed admin/_next/static/*/*.js
var adminFiles embed.FS

func GetAdmin() embed.FS {
	return adminFiles
}
