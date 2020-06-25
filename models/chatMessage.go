package models

//ChatMessage represents a single chat message
type ChatMessage struct {
	ClientID string `json:"-"`

	Author      string `json:"author"`
	Body        string `json:"body"`
	Image       string `json:"image"`
	ID          string `json:"id"`
	MessageType string `json:"type"`
}

//Valid checks to ensure the message is valid
func (s ChatMessage) Valid() bool {
	return s.Author != "" && s.Body != "" && s.ID != ""
}
