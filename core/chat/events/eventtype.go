package events

// EventType is the type of a websocket event.
type EventType = string

const (
	// MessageSent is the event sent when a chat event takes place.
	MessageSent EventType = "CHAT"
	// UserJoined is the event sent when a chat user join action takes place.
	UserJoined EventType = "USER_JOINED"
	// UserNameChanged is the event sent when a chat username change takes place.
	UserNameChanged EventType = "NAME_CHANGE"
	// VisibiltyUpdate is the event sent when a chat message's visibility changes.
	VisibiltyUpdate EventType = "VISIBILITY-UPDATE"
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
	// ChatDisabled is when a user is explicitly disabled and blocked from using chat.
	ChatDisabled EventType = "CHAT_DISABLED"
	// ConnectedUserInfo is a private event to a user letting them know their user details.
	ConnectedUserInfo EventType = "CONNECTED_USER_INFO"
	// ChatActionSent is a generic chat action that can be used for anything that doesn't need specific handling or formatting.
	ChatActionSent EventType = "CHAT_ACTION"
	// ErrorNeedsRegistration is an error returned when the client needs to perform registration.
	ErrorNeedsRegistration EventType = "ERROR_NEEDS_REGISTRATION"
	// ErrorMaxConnectionsExceeded is an error returned when the server determined it should not handle more connections.
	ErrorMaxConnectionsExceeded EventType = "ERROR_MAX_CONNECTIONS_EXCEEDED"
	// ErrorUserDisabled is an error returned when the connecting user has been previously banned/disabled.
	ErrorUserDisabled EventType = "ERROR_USER_DISABLED"
)
