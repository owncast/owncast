package models

import (
	"time"

	"github.com/teris-io/shortid"
)

// Event is any kind of event.  A type is required to be specified.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type,omitempty"`
	ID        string    `json:"id"`
}

// SetDefaults will set default properties of all inbound events.
func (e *Event) SetDefaults() {
	e.ID = shortid.MustGenerate()
	e.Timestamp = time.Now()
}
