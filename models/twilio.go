package models

type TwilioConfiguration struct {
	Enabled       bool   `json:"enabled"`
	AccountSid    string `json:"accountSid,omitempty"`
	AuthToken     string `json:"authToken,omitempty"`
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
}

type DiscordConfiguration struct {
	Enabled       bool   `json:"enabled"`
	Webhook       string `json:"webhook,omitempty"`
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
}

type BrowserNotificationConfiguration struct {
	Enabled       bool   `json:"enabled"`
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
}
