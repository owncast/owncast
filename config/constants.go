package config

import "path/filepath"

const (
	WebRoot               = "webroot"
	PrivateHLSStoragePath = "hls"
)

var (
	PublicHLSStoragePath = filepath.Join(WebRoot, "hls")
)
