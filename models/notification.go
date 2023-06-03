package models

// DiscordConfiguration represents the configuration for the discord
// notification service.
type DiscordConfiguration struct {
	Webhook       string `json:"webhook,omitempty"`
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
	Enabled       bool   `json:"enabled"`
}

// BrowserNotificationConfiguration represents the configuration for
// browser notifications.
type BrowserNotificationConfiguration struct {
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
	Enabled       bool   `json:"enabled"`
}
