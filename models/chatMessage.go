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

//String converts the chat message to string
//TODO: is this required? or can we remove it
func (s ChatMessage) String() string {
	return s.Author + " says " + s.Body
}
