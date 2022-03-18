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

// TwitterConfiguration represents the configuration for Twitter access.
type TwitterConfiguration struct {
	Enabled           bool   `json:"enabled"`
	APIKey            string `json:"apiKey"`    // aka consumer key
	APISecret         string `json:"apiSecret"` // aka consumer secret
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	BearerToken       string `json:"bearerToken"`
	GoLiveMessage     string `json:"goLiveMessage,omitempty"`
}
