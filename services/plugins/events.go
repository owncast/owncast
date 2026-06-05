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

	// SSE connection lifecycle. Fired to the plugin that owns a
	// Server-Sent-Events channel when a browser opens or closes a connection
	// to it, so the plugin can track who is currently connected. The payload
	// is an SSEConnectionEvent. Requires the http.sse permission (the same
	// gate as serving the stream).
	EventSSEConnect    = "sse.connect"
	EventSSEDisconnect = "sse.disconnect"

	// EventTick fires to every plugin that subscribes (defines onTick) once a
	// second, for periodic work. Payload is a TickEvent.
	EventTick = "tick"

	// EventTimerFire is delivered to a plugin when one of its host-scheduled
	// timers elapses (see TimerHub). Payload is a TimerFireEvent. Internal:
	// the guest SDK maps the id back to the author's callback, so plugins
	// don't subscribe to or handle this directly.
	EventTimerFire = "timer.fire"

	// Fediverse events. Engagement (follow/like/repost) carries only
	// actor + target metadata; mention/reply also carry the post content
	// so plugins can act on what was actually said.
	EventFediverseFollow  = "fediverse.follow"
	EventFediverseLike    = "fediverse.like"
	EventFediverseRepost  = "fediverse.repost"
	EventFediverseMention = "fediverse.mention"
	EventFediverseReply   = "fediverse.reply"
)
