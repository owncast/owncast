package models

// DiscordConfiguration represents the configuration for the discord
// notification service.
type DiscordConfiguration struct {
	Enabled       bool   `json:"enabled"`
	Webhook       string `json:"webhook,omitempty"`
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
}

// BrowserNotificationConfiguration represents the configuration for
// browser notifications.
type BrowserNotificationConfiguration struct {
	Enabled       bool   `json:"enabled"`
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
}
