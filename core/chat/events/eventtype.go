package events

// EventType is the type of a websocket event.
type EventType = string

const (
	// MessageSent is the event sent when a chat event takes place.
	Event_MessageSent EventType = "CHAT"
	// UserJoined is the event sent when a chat user join action takes place.
	Event_TypeUserJoined EventType = "USER_JOINED"
	// UserNameChanged is the event sent when a chat username change takes place.
	Event_UserNameChanged EventType = "NAME_CHANGE"
	// VisibiltyToggled is the event sent when a chat message's visibility changes.
	Event_VisibiltyToggled EventType = "VISIBILITY-UPDATE"
	// PING is a ping message.
	Event_PING EventType = "PING"
	// PONG is a pong message.
	Event_PONG EventType = "PONG"
	// StreamStarted represents a stream started event.
	StreamStarted EventType = "STREAM_STARTED"
	// StreamStopped represents a stream stopped event.
	StreamStopped EventType = "STREAM_STOPPED"
	// SystemMessageSent is the event sent when a system message is sent.
	Event_SystemMessageSent EventType = "SYSTEM"
	Event_Chat_Disabled     EventType = "CHAT_DISABLED"

	// ChatActionSent is a generic chat action that can be used for anything that doesn't need specific handling or formatting.
	Event_ChatActionSent                 EventType = "CHAT_ACTION"
	Event_UserJoined                     EventType = "USER_JOINED"
	Event_Error_Needs_Registration       EventType = "ERROR_NEEDS_REGISTRATION"
	Event_Error_Max_Connections_Exceeded EventType = "ERROR_MAX_CONNECTIONS_EXCEEDED"
	Event_Error_User_Disabled            EventType = "ERROR_USER_DISABLED"
)
