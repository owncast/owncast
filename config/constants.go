package config

import "path/filepath"

const (
	// StaticVersionNumber is the version of Owncast that is used when it's not overwritten via build-time settings.
	StaticVersionNumber = "0.0.8" // Shown when you build from develop
	// WebRoot is the web server root directory.
	WebRoot = "webroot"
	// PrivateHLSStoragePath is the HLS write directory.
	PrivateHLSStoragePath = "hls"
	// FfmpegSuggestedVersion is the version of ffmpeg we suggest.
	FfmpegSuggestedVersion = "v4.1.5" // Requires the v
	// DataDirectory is the directory we save data to.
	DataDirectory = "data"
)

var (
	// PublicHLSStoragePath is the directory we write public HLS files to for distribution.
	PublicHLSStoragePath = filepath.Join(WebRoot, "hls")
	// BackupDirectory is the directory we write backup files to.
	BackupDirectory = filepath.Join(DataDirectory, "backup")
)
