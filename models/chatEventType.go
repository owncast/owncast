package models

type ChatMessageType string

const (
	MessageSent      ChatMessageType = "CHAT"
	UserJoined       ChatMessageType = "USER_JOINED"
	UserNameChanged  ChatMessageType = "NAME_CHANGE"
	VisibiltyToggled ChatMessageType = "VISIBILITY-UPDATE"
	PING             ChatMessageType = "PING"
	PONG             ChatMessageType = "PONG"
)
