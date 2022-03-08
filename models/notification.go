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

// MailjetConfiguration represents the configuration for using Mailjet.
type MailjetConfiguration struct {
	Enabled bool  `json:"enabled"`
	ListID  int64 `json:"listID,omitempty"`
}

// SMTPConfiguration represents the configuration for using SMTP notifications.
type SMTPConfiguration struct {
	Enabled       bool   `json:"enabled"`
	ListAddress   string `json:"listAddress,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Server        string `json:"server,omitempty"`
	FromAddress   string `json:"fromAddress,omitempty"`
	GoLiveSubject string `json:"goLiveSubject,omitempty"`
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
