package config

import "path/filepath"

const (
	WebRoot                = "webroot"
	PrivateHLSStoragePath  = "hls"
	GeoIPDatabasePath      = "data/GeoLite2-City.mmdb"
	ExtraInfoFile          = "data/content.md"
	StatsFile              = "data/stats.json"
	FfmpegSuggestedVersion = "v4.1.5" // Requires the v
)

var (
	PublicHLSStoragePath = filepath.Join(WebRoot, "hls")
)
