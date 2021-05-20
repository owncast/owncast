package config

import "path/filepath"

const (
	// StaticVersionNumber is the version of Owncast that is used when it's not overwritten via build-time settings.
	StaticVersionNumber = "0.0.8" // Shown when you build from develop
	// WebRoot is the web server root directory.
	WebRoot = "webroot"
	// PrivateHLSStoragePath is the HLS write directory.
	PrivateHLSStoragePath = "hls"
	// ExtraInfoFile is the markdown file for page content.  Remove this after the migrator is removed.
	ExtraInfoFile = "data/content.md"
	// StatsFile is the json file we used to save stats in.  Remove this after the migrator is removed.
	StatsFile = "data/stats.json"
	// FfmpegSuggestedVersion is the version of ffmpeg we suggest.
	FfmpegSuggestedVersion = "v4.1.5" // Requires the v
	// BackupDirectory is the directory we write backup files to.
	BackupDirectory = "backup"
	// DataDirectory is the directory we save data to.
	DataDirectory = "data"
)

var (
	// PublicHLSStoragePath is the directory we write public HLS files to for distribution.
	PublicHLSStoragePath = filepath.Join(WebRoot, "hls")
)
