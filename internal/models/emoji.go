package models

// CustomEmoji represents an image that can be used in chat as a custom emoji.
type CustomEmoji struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
