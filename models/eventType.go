package models

// EventType is the type of a websocket event.
type EventType = string

const (
	// MessageSent is the event sent when a chat event takes place.
	MessageSent EventType = "CHAT"
	// UserJoined is the event sent when a chat user join action takes place.
	UserJoined EventType = "USER_JOINED"
	// UserParted is the event sent when a chat user parted action takes place.
	UserParted EventType = "USER_PARTED"
	// UserNameChanged is the event sent when a chat username change takes place.
	UserNameChanged EventType = "NAME_CHANGE"
	// FediverseEngagementFollow is the event sent when a user follows the stream.
	FediverseEngagementFollow EventType = "FEDIVERSE_ENGAGEMENT_FOLLOW"
	// FediverseEngagementLike is the internal event sent when a remote user likes a local post.
	FediverseEngagementLike EventType = "FEDIVERSE_ENGAGEMENT_LIKE"
	// FediverseEngagementRepost is the internal event sent when a remote user reposts a local post.
	FediverseEngagementRepost EventType = "FEDIVERSE_ENGAGEMENT_REPOST"
	// FediverseEngagementQuote is the internal event sent when a remote user quotes a local post.
	FediverseEngagementQuote EventType = "FEDIVERSE_ENGAGEMENT_QUOTE"
	// FediverseMention is the internal event sent when a remote post mentions the stream.
	FediverseMention EventType = "FEDIVERSE_MENTION"
	// FediverseReply is the internal event sent when a remote user replies to a local post.
	FediverseReply EventType = "FEDIVERSE_REPLY"
	// FediverseActivity is the internal event carrying a verified inbound ActivityPub activity.
	FediverseActivity EventType = "FEDIVERSE_ACTIVITY"
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
	// StreamTitleUpdated is the event sent when a stream's title changes.
	StreamTitleUpdated EventType = "STREAM_TITLE_UPDATED"
	// SystemMessageSent is the event sent when a system message is sent.
	SystemMessageSent EventType = "SYSTEM"
	// ChatActionSent is a generic chat action that can be used for anything that doesn't need specific handling or formatting.
	ChatActionSent EventType = "CHAT_ACTION"
)
