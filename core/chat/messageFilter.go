package chat

import (
	goaway "github.com/TwiN/go-away"
)

// ChatMessageFilter is a allow/deny chat message filter.
type ChatMessageFilter struct{}

// NewMessageFilter will return an instance of the chat message filter.
func NewMessageFilter() *ChatMessageFilter {
	return &ChatMessageFilter{}
}

// Allow will test if this message should be allowed to be sent.
func (*ChatMessageFilter) Allow(message string) bool {
	return !goaway.IsProfane(message)
}
