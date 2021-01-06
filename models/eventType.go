package models

type EventType = string

const (
	MessageSent      EventType = "CHAT"
	UserJoined       EventType = "USER_JOINED"
	UserNameChanged  EventType = "NAME_CHANGE"
	VisibiltyToggled EventType = "VISIBILITY-UPDATE"
	PING             EventType = "PING"
	PONG             EventType = "PONG"
	StreamStarted    EventType = "STREAM_STARTED"
	StreamStopped    EventType = "STREAM_STOPPED"
)
