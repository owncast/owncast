package config

import (
	"time"

	"github.com/owncast/owncast/models"
)

// Defaults will hold default configuration values.
type Defaults struct {
	PageBodyContent string

	FederationGoLiveMessage string

	Summary              string
	ServerWelcomeMessage string
	Logo                 string
	YPServer             string

	Title string

	DatabaseFilePath string

	FederationUsername string
	WebServerIP        string
	Name               string
	AdminPassword      string
	StreamKeys         []models.StreamKey

	StreamVariants []models.StreamOutputVariant

	Tags               []string
	RTMPServerPort     int
	SegmentsInPlaylist int

	SegmentLengthSeconds int
	WebServerPort        int

	ChatEstablishedUserModeTimeDuration time.Duration

	YPEnabled bool
}

// GetDefaults will return default configuration values.
func GetDefaults() Defaults {
	return Defaults{
		Name:                 "New Owncast Server",
		Summary:              "This is a new live video streaming server powered by Owncast.",
		ServerWelcomeMessage: "",
		Logo:                 "logo.svg",
		AdminPassword:        "abc123",
		StreamKeys: []models.StreamKey{
			{Key: "abc123", Comment: "Default stream key"},
		},
		Tags: []string{
			"owncast",
			"streaming",
		},

		PageBodyContent: `
# Welcome to Owncast!

- This is a live stream powered by [Owncast](https://owncast.online), a free and open source live streaming server.

- To discover more examples of streams, visit [Owncast's directory](https://directory.owncast.online).

- If you're the owner of this server you should visit the admin and customize the content on this page.

<hr/>

<video id="video" controls preload="metadata" style="width: 60vw; max-width: 600px; min-width: 200px;" poster="https://videos.owncast.online/t/xaJ3xNn9Y6pWTdB25m9ai3">
  <source src="https://assets.owncast.tv/video/owncast-embed.mp4" type="video/mp4" />
</video>
	`,

		DatabaseFilePath: "data/owncast.db",

		YPEnabled: false,
		YPServer:  "https://directory.owncast.online",

		WebServerPort:  8080,
		WebServerIP:    "0.0.0.0",
		RTMPServerPort: 1935,

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
