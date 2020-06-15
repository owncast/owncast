package main

type ChatMessage struct {
	Author      string `json:"author"`
	Body        string `json:"body"`
	Image       string `json:"image"`
	ID          string `json:"id"`
	MessageType string `json:"type"`
}

func (s *ChatMessage) AsJson() string {
	return s.Author + " says " + s.Body
}

func (s *ChatMessage) String() string {
	return s.Author + " says " + s.Body
}

type PingMessage struct {
	MessageType string `json:"type"`
}

func (self *PingMessage) AsJson() string {
	return "{type: \"PING\"}"
}
