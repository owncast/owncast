package models

// EventType is the type of a websocket event.
type EventType = string

const (
	// MessageSent is the event sent when a chat event takes place.
	MessageSent EventType = "CHAT"
	// UserJoined is the event sent when a chat user join action takes place.
	UserJoined EventType = "USER_JOINED"
	// UserNameChanged is the event sent when a chat username change takes place.
	UserNameChanged EventType = "NAME_CHANGE"
	// VisibiltyToggled is the event sent when a chat message's visibility changes.
	VisibiltyToggled EventType = "VISIBILITY-UPDATE"
	// PING is a ping message.
	PING EventType = "PING"
	// PONG is a pong message.
	PONG EventType = "PONG"
	// StreamStarted represents a stream started event.
	StreamStarted EventType = "STREAM_STARTED"
	// StreamStopped represents a stream stopped event.
	StreamStopped EventType = "STREAM_STOPPED"
	// SystemMessageSent is the event sent when a system message is sent.
	SystemMessageSent EventType = "SYSTEM"
	// ChatActionSent is a generic chat action that can be used for anything that doesn't need specific handling or formatting.
	ChatActionSent EventType = "CHAT_ACTION"
)
