package config

import (
	"time"

	"github.com/owncast/owncast/models"
	"golang.org/x/crypto/bcrypt"
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
func GetDefaults() *Defaults {
	return &Defaults{
		Name:                 "New Owncast Server",
		Summary:              "This is a new live video streaming server powered by Owncast.",
		ServerWelcomeMessage: "",
		Logo:                 "logo.svg",
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

<video id="video" controls preload="metadata" width="40%" poster="https://videos.owncast.online/t/xaJ3xNn9Y6pWTdB25m9ai3">
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

func (d *Defaults) GetAdminPasswordHash() ([]byte, error) {
	return d.getDefaultPasswordHash()
}

func (d *Defaults) GetStreamKeysHashed() ([]models.StreamKeyHashed, error) {
	hash, err := d.getDefaultPasswordHash()
	if err != nil {
		return nil, err
	}
	return []models.StreamKeyHashed{
		{Key: hash, Comment: "default stream key"},
	}, nil
}

func (d *Defaults) getDefaultPasswordHash() ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("abc123"), bcrypt.DefaultCost)
	return hash, err
}
