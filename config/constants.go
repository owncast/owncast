package config

import "path/filepath"

const (
	CurrentBuildString     = "0.0.6" // Shown when you build from master
	WebRoot                = "webroot"
	PrivateHLSStoragePath  = "hls"
	ExtraInfoFile          = "data/content.md"
	StatsFile              = "data/stats.json"
	FfmpegSuggestedVersion = "v4.1.5" // Requires the v
	BackupDirectory        = "backup"
)

var (
	PublicHLSStoragePath = filepath.Join(WebRoot, "hls")
)
