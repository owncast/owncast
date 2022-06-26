package config

import "path/filepath"

const (
	// StaticVersionNumber is the version of Owncast that is used when it's not overwritten via build-time settings.
	StaticVersionNumber = "0.1.0" // Shown when you build from develop
	// FfmpegSuggestedVersion is the version of ffmpeg we suggest.
	FfmpegSuggestedVersion = "v4.1.5" // Requires the v
	// DataDirectory is the directory we save data to.
	DataDirectory = "data"
	// EmojiDir is relative to the static directory.
	EmojiDir = "/img/emoji"
)

var (
	// BackupDirectory is the directory we write backup files to.
	BackupDirectory = filepath.Join(DataDirectory, "backup")

	// HLSStoragePath is the directory HLS video is written to.
	HLSStoragePath = filepath.Join(DataDirectory, "hls")

	// CustomEmojiPath is the optional emoji override directory.
	CustomEmojiPath = filepath.Join(DataDirectory, "emoji")
)
