package models

type SocialHandle struct {
	Platform string `yaml:"platform" json:"platform"`
	URL      string `yaml:"url" json:"url"`
}
