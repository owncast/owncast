package plugins

import (
	"strings"
	"sync"
)

// DefaultMaxSSEConnectionsPerPlugin caps how many simultaneous
// Server-Sent-Events connections a single plugin may hold open. The limit
// is per-plugin so one plugin can't exhaust the host's file descriptors or
// memory at the expense of others.
const DefaultMaxSSEConnectionsPerPlugin = 64

// sseClient is one connected browser. The host owns the connection; the
// stream goroutine reads pre-framed bytes off send and flushes them to the
// client. The buffer absorbs short bursts; a client that can't keep up has
// messages dropped (see SSEHub.Publish) rather than stalling the publisher.
type sseClient struct {
	send chan []byte
}

// SSEHub fans plugin-emitted events out to connected browser clients.
//
// It is the host-owned half of the SSE capability: plugins never hold the
// HTTP connection (their wasm calls are short and mutex-serialized), they
// only Publish. The hub holds the long-lived connections and delivers to
// them. Conceptually it's the event dispatcher with browser clients as the
// sink instead of other plugins.
//
// Subscribers are keyed by (plugin, channel) so a plugin can run several
// independent streams (e.g. "overlay" and "admin-stats") and so one
// plugin's clients never receive another's events.
type SSEHub struct {
	mu               sync.Mutex
	subscribers      map[string]map[*sseClient]struct{}
	connectionCounts map[string]int
	maxPerPlugin     int
}

// NewSSEHub constructs an empty hub with the default per-plugin connection
// cap.
func NewSSEHub() *SSEHub {
	return &SSEHub{
		subscribers:      make(map[string]map[*sseClient]struct{}),
		connectionCounts: make(map[string]int),
		maxPerPlugin:     DefaultMaxSSEConnectionsPerPlugin,
	}
}

func sseKey(pluginName, channel string) string {
	return pluginName + "\x00" + channel
}

// Subscribe registers a browser client for (pluginName, channel). It
// returns a receive-only channel of pre-framed SSE bytes and an unsubscribe
// func the caller MUST invoke when the connection ends. ok is false when
// the plugin is already at its connection cap, in which case the caller
// should reject the request.
func (h *SSEHub) Subscribe(pluginName, channel string) (<-chan []byte, func(), bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connectionCounts[pluginName] >= h.maxPerPlugin {
		return nil, nil, false
	}

	client := &sseClient{send: make(chan []byte, 16)}
	key := sseKey(pluginName, channel)
	if h.subscribers[key] == nil {
		h.subscribers[key] = make(map[*sseClient]struct{})
	}
	h.subscribers[key][client] = struct{}{}
	h.connectionCounts[pluginName]++

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			h.mu.Lock()
			defer h.mu.Unlock()
			if subs, ok := h.subscribers[key]; ok {
				delete(subs, client)
				if len(subs) == 0 {
					delete(h.subscribers, key)
				}
			}
			h.connectionCounts[pluginName]--
		})
	}

	return client.send, unsubscribe, true
}

// Publish frames event+data once and delivers it to every client subscribed
// to (pluginName, channel). Delivery is non-blocking per client: if a
// client's buffer is full (a slow consumer), the message is dropped for that
// client. This keeps a single slow browser from blocking the publishing
// plugin call — correctness of the live stream is best-effort by design.
// Returns the number of clients the frame was delivered to.
//
// The send happens under the hub lock; sends are non-blocking, so the lock
// is held only briefly. Because clients are removed from the map under the
// same lock by unsubscribe, Publish never sends to a torn-down client.
func (h *SSEHub) Publish(pluginName, channel, event string, data []byte) int {
	frame := frameSSE(event, data)

	h.mu.Lock()
	defer h.mu.Unlock()

	delivered := 0
	for client := range h.subscribers[sseKey(pluginName, channel)] {
		select {
		case client.send <- frame:
			delivered++
		default:
			// Slow client; drop this frame rather than block the publisher.
		}
	}
	return delivered
}

// frameSSE formats one Server-Sent-Events message. A named event becomes an
// "event:" line; the data is split across "data:" lines (SSE requires one
// per newline) and terminated by a blank line.
func frameSSE(event string, data []byte) []byte {
	var b strings.Builder
	if event != "" {
		b.WriteString("event: ")
		b.WriteString(event)
		b.WriteByte('\n')
	}
	for _, line := range strings.Split(string(data), "\n") {
		b.WriteString("data: ")
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	return []byte(b.String())
}
