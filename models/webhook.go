package models

import "time"

type Webhook struct {
	ID        int         `json:"id"`
	Url       string      `json:"url"`
	Events    []EventType `json:"events"`
	Timestamp time.Time   `json:"timestamp"`
	LastUsed  time.Time   `json:"lastUsed"`
}

// For an event to be seen as "valid" it must live in this slice.
var validEvents = []EventType{
	MessageSent,
	UserJoined,
	UserNameChanged,
	VisibiltyToggled,
	StreamStarted,
	StreamStopped,
}

// HasValidEvents will verify that all the events provided are valid.
// This is not a efficient method.
func HasValidEvents(events []EventType) bool {
	for _, event := range events {
		if !findItemInSlice(validEvents, event) {
			return false
		}
	}
	return true
}

func findItemInSlice(slice []EventType, value EventType) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
