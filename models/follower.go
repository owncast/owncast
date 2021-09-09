package models

type Follower struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Image    string `json:"image"`
	Link     string `json:"link"`
	Inbox    string `json:"-"`
}
