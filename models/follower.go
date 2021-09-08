package models

import "net/url"

type Follower struct {
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Link  string  `json:"link"`
	Inbox url.URL `json:"-"`
}
