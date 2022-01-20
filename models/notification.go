package models

type DiscordConfiguration struct {
	Enabled       bool   `json:"enabled"`
	Webhook       string `json:"webhook,omitempty"`
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
}

type BrowserNotificationConfiguration struct {
	Enabled       bool   `json:"enabled"`
	GoLiveMessage string `json:"goLiveMessage,omitempty"`
}

// MailjetConfiguration represents the configuration for using Mailjet.
type MailjetConfiguration struct {
	Enabled       bool   `json:"enabled"`
	ListAddress   string `json:"listAddress,omitempty"`
	ListID        int64  `json:"listID,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	SMTPServer    string `json:"smtpServer,omitempty"`
	FromAddress   string `json:"fromAddress,omitempty"`
	GoLiveSubject string `json:"goLiveSubject,omitempty"`
}
