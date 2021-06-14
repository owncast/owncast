package events

import (
	"time"

	"github.com/teris-io/shortid"

	"github.com/owncast/owncast/core/user"
)

// EventPayload is a generic key/value map for sending out to chat clients.
type EventPayload map[string]interface{}

type OutboundEvent interface {
	GetBroadcastPayload() EventPayload
	GetMessageType() EventType
}

// Event is any kind of event.  A type is required to be specified.
type Event struct {
	Type      string    `json:"type"`
	Id        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserEvent struct {
	User     *user.User `json:"user"`
	HiddenAt *time.Time `json:"hiddenAt,omitempty"`
}

// MessageEvent is an event that has a message body.
type MessageEvent struct {
	OutboundEvent `json:"-"`
	Body          string `json:"body"`
}

type SystemActionEvent struct {
	Event
	MessageEvent
}

// render will convert any supported markdown into the message body.
func (e *MessageEvent) render() string {
	return e.Body
}

func (e *MessageEvent) sanitize() string {
	return e.Body
}

// SetDefaults will set default properties of all inbound events.
func (e *Event) SetDefaults() {
	e.Id = shortid.MustGenerate()
	e.Timestamp = time.Now()
}
