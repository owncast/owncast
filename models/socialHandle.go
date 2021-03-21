package models

// SocialHandle represents an external link.
type SocialHandle struct {
	Platform string `yaml:"platform" json:"platform,omitempty"`
	URL      string `yaml:"url" json:"url,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

// GetSocialHandle will return the details for a supported platform.
func GetSocialHandle(platform string) *SocialHandle {
	allPlatforms := GetAllSocialHandles()
	if platform, ok := allPlatforms[platform]; ok {
		return &platform
	}

	return nil
}

// GetAllSocialHandles will return a list of all the social platforms we support.
func GetAllSocialHandles() map[string]SocialHandle {
	socialHandlePlatforms := map[string]SocialHandle{
		"bandcamp": {
			Platform: "Bandcamp",
			Icon:     "/img/platformlogos/bandcamp.svg",
		},
		"discord": {
			Platform: "Discord",
			Icon:     "/img/platformlogos/discord.svg",
		},
		"facebook": {
			Platform: "Facebook",
			Icon:     "/img/platformlogos/facebook.svg",
		},
		"github": {
			Platform: "GitHub",
			Icon:     "/img/platformlogos/github.svg",
		},
		"gitlab": {
			Platform: "GitLab",
			Icon:     "/img/platformlogos/gitlab.svg",
		},
		"instagram": {
			Platform: "Instagram",
			Icon:     "/img/platformlogos/instagram.svg",
		},
		"keyoxide": {
			Platform: "Keyoxide",
			Icon:     "/img/platformlogos/keyoxide.png",
		},
		"kofi": {
			Platform: "Ko-Fi",
			Icon:     "/img/platformlogos/ko-fi.svg",
		},
		"linkedin": {
			Platform: "LinkedIn",
			Icon:     "/img/platformlogos/linkedin.svg",
		},
		"mastodon": {
			Platform: "Mastodon",
			Icon:     "/img/platformlogos/mastodon.svg",
		},
		"patreon": {
			Platform: "Patreon",
			Icon:     "/img/platformlogos/patreon.svg",
		},
		"paypal": {
			Platform: "Paypal",
			Icon:     "/img/platformlogos/paypal.svg",
		},
		"snapchat": {
			Platform: "Snapchat",
			Icon:     "/img/platformlogos/snapchat.svg",
		},
		"soundcloud": {
			Platform: "Soundcloud",
			Icon:     "/img/platformlogos/soundcloud.svg",
		},
		"spotify": {
			Platform: "Spotify",
			Icon:     "/img/platformlogos/spotify.svg",
		},
		"tiktok": {
			Platform: "TikTok",
			Icon:     "/img/platformlogos/tiktok.svg",
		},
		"twitch": {
			Platform: "Twitch",
			Icon:     "/img/platformlogos/twitch.svg",
		},
		"twitter": {
			Platform: "Twitter",
			Icon:     "/img/platformlogos/twitter.svg",
		},
		"youtube": {
			Platform: "YouTube",
			Icon:     "/img/platformlogos/youtube.svg",
		},
		"odysee": {
			Platform: "Odysee",
			Icon:     "/img/platformlogos/odysee.svg",
		},
		"lbry": {
			Platform: "LBRY",
			Icon:     "/img/platformlogos/lbry.svg",
		},
		"liberapay": {
			Platform: "LiberaPay",
			Icon:     "/img/platformlogos/liberapay.svg",
		},
		"donate": {
			Platform: "Donations",
			Icon:     "/img/platformlogos/donate.svg",
		},
		"follow": {
			Platform: "Follow",
			Icon:     "/img/platformlogos/follow.svg",
		},
	}

	return socialHandlePlatforms
}
