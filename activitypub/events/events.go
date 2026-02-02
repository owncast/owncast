package events

import (
	"time"

	"github.com/teris-io/shortid"
)

// Event is any kind of event.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	ID        string    `json:"id"`
}

// SetDefaults will set default properties of all inbound events.
func (e *Event) SetDefaults() {
	e.ID = shortid.MustGenerate()
	e.Timestamp = time.Now()
}
