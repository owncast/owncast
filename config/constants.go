package config

import "path/filepath"

const (
	WebRoot               = "webroot"
	PrivateHLSStoragePath = "hls"
	GeoIPDatabasePath     = "data/GeoLite2-City.mmdb"
)

var (
	PublicHLSStoragePath = filepath.Join(WebRoot, "hls")
)
