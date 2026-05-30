package config

import "path/filepath"

const (
	// StaticVersionNumber is the version of Owncast that is used when it's not overwritten via build-time settings.
	StaticVersionNumber = "0.3.0" // Shown when you build from develop
	// FfmpegSuggestedVersion is the version of ffmpeg we suggest.
	FfmpegSuggestedVersion = "v4.1.5" // Requires the v
	// DataDirectory is the directory we save data to.
	DataDirectory = "data"
	// EmojiDir defines the URL route prefix for emoji requests.
	EmojiDir = "/img/emoji/"
	// MaxUserColor is the largest color value available to assign to users.
	// They start at 0 and can be treated as IDs more than colors themselves.
	MaxUserColor = 7
	// MaxChatDisplayNameLength is the maximum length of a chat display name.
	MaxChatDisplayNameLength = 30

	// ActivityPub namespace properties for Owncast metadata.
	APOwncastNamespaceStreamStatus      = "https://owncast.online/ns#streamStatus"
	APOwncastNamespaceStreamTitle       = "https://owncast.online/ns#streamTitle"
	APOwncastNamespaceServerName        = "https://owncast.online/ns#serverName"
	APOwncastNamespaceStreamDescription = "https://owncast.online/ns#streamDescription"
	APOwncastNamespaceLogoURL           = "https://owncast.online/ns#logoUrl"
	APOwncastNamespaceThumbnailURL      = "https://owncast.online/ns#thumbnailUrl"
	APOwncastNamespaceStreamTags        = "https://owncast.online/ns#streamTags"

	// APStreamStatusLive is the stream status value for a live stream.
	APStreamStatusLive = "live"
	// APStreamStatusOffline is the stream status value for an offline stream.
	APStreamStatusOffline = "offline"
)

var (
	// HLSStoragePath is the directory HLS video is written to.
	HLSStoragePath = filepath.Join(DataDirectory, "hls")

	// CustomEmojiPath is the emoji directory.
	CustomEmojiPath = filepath.Join(DataDirectory, "emoji")

	// PublicFilesPath is the optional directory for hosting public files.
	PublicFilesPath = filepath.Join(DataDirectory, "public")
)
