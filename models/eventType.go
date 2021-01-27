package models

type EventType = string

const (
	MessageSent       EventType = "CHAT"
	UserJoined        EventType = "USER_JOINED"
	UserNameChanged   EventType = "NAME_CHANGE"
	VisibiltyToggled  EventType = "VISIBILITY-UPDATE"
	PING              EventType = "PING"
	PONG              EventType = "PONG"
	StreamStarted     EventType = "STREAM_STARTED"
	StreamStopped     EventType = "STREAM_STOPPED"
	SystemMessageSent EventType = "SYSTEM"
	ChatActionSent    EventType = "CHAT_ACTION" // Generic chat action that can be used for anything that doesn't need specific handling or formatting.
)
