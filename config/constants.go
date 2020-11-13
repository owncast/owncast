package config

import "path/filepath"

const (
	WebRoot               = "webroot"
	PrivateHLSStoragePath = "hls"
	GeoIPDatabasePath     = "data/GeoLite2-City.mmdb"
	ExtraInfoFile         = "data/content.md"
	StatsFile             = "data/stats.json"
)

var (
	PublicHLSStoragePath = filepath.Join(WebRoot, "hls")
)
