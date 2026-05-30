package plugins

// Built-in Owncast event types. These mirror the SDK's `Events` const so the
// host and plugins agree on names without typo risk. Plugin-emitted custom
// events are arbitrary strings; define those at the call site.
const (
	// Chat events.
	EventChatMessageReceived  = "chat.message.received"
	EventChatUserJoined       = "chat.user.joined"
	EventChatUserParted       = "chat.user.parted"
	EventChatUserRenamed      = "chat.user.renamed"
	EventChatMessageModerated = "chat.message.moderated"

	// Stream lifecycle events.
	EventStreamStarted      = "stream.started"
	EventStreamStopped      = "stream.stopped"
	EventStreamTitleChanged = "stream.title.changed"

	// Fediverse events. Engagement (follow/like/repost) carries only
	// actor + target metadata; mention/reply also carry the post content
	// so plugins can act on what was actually said.
	EventFediverseFollow  = "fediverse.follow"
	EventFediverseLike    = "fediverse.like"
	EventFediverseRepost  = "fediverse.repost"
	EventFediverseMention = "fediverse.mention"
	EventFediverseReply   = "fediverse.reply"
)
