package config

const (
	// StaticVersionNumber is the version of Owncast that is used when it's not overwritten via build-time settings.
	StaticVersionNumber = "0.1.3" // Shown when you build from develop
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
)
