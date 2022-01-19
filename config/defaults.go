package config

import (
	"time"

	"github.com/owncast/owncast/models"
)

// Defaults will hold default configuration values.
type Defaults struct {
	Name                 string
	Title                string
	Summary              string
	ServerWelcomeMessage string
	Logo                 string
	Tags                 []string
	PageBodyContent      string

	DatabaseFilePath string
	WebServerPort    int
	WebServerIP      string
	RTMPServerPort   int
	StreamKey        string

	YPEnabled bool
	YPServer  string

	SegmentLengthSeconds int
	SegmentsInPlaylist   int
	StreamVariants       []models.StreamOutputVariant

	FederationUsername      string
	FederationGoLiveMessage string

	ChatEstablishedUserModeTimeDuration time.Duration
}

// GetDefaults will return default configuration values.
func GetDefaults() Defaults {
	return Defaults{
		Name:                 "Owncast",
		Title:                "My Owncast Server",
		Summary:              "This is brief summary of whom you are or what your stream is. You can edit this description in the admin.",
		ServerWelcomeMessage: "",
		Logo:                 "logo.svg",
		Tags: []string{
			"owncast",
			"streaming",
		},

		PageBodyContent: "# This is your page content that can be edited from the admin.",

		DatabaseFilePath: "data/owncast.db",

		YPEnabled: false,
		YPServer:  "https://directory.owncast.online",

		WebServerPort:  8080,
		WebServerIP:    "0.0.0.0",
		RTMPServerPort: 1935,
		StreamKey:      "abc123",

		ChatEstablishedUserModeTimeDuration: time.Minute * 15,

		StreamVariants: []models.StreamOutputVariant{
			{
				IsAudioPassthrough: true,
				VideoBitrate:       1200,
				Framerate:          24,
				CPUUsageLevel:      2,
			},
		},

		FederationUsername:      "streamer",
		FederationGoLiveMessage: "I've gone live!",
	}
}
