package models

import "time"

// Scheduled event occurrence statuses.
const (
	// ScheduledEventStatusScheduled is an upcoming (or past, uncancelled) occurrence.
	ScheduledEventStatusScheduled = "scheduled"
	// ScheduledEventStatusCancelled is an occurrence cancelled but kept, both to
	// announce the cancellation to federated followers and to stop the
	// materializer from re-creating the slot.
	ScheduledEventStatusCancelled = "cancelled"
)

// ScheduledEventWebhookAction identifies a scheduled-event lifecycle notification.
type ScheduledEventWebhookAction string

const (
	ScheduledEventWarning ScheduledEventWebhookAction = "10-minute-warning"
	ScheduledEventStarted ScheduledEventWebhookAction = "started"
	ScheduledEventEnded   ScheduledEventWebhookAction = "ended"
)

// ScheduledEventSeries is a recurring schedule rule ("every Monday at 18:00")
// that the scheduler expands into concrete ScheduledEvent occurrences on a
// rolling horizon. The recurrence value is an iCalendar RFC 5545 string:
// a DTSTART line carrying an IANA TZID plus an RRULE line.
type ScheduledEventSeries struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ReminderMessage string `json:"reminderMessage,omitempty"`
	Recurrence      string `json:"recurrence"`
	DurationMinutes int    `json:"durationMinutes"`
	Active          bool   `json:"active"`
}

// ScheduledEvent is one concrete scheduled stream occurrence: either a
// one-off event or a row materialized from a series.
type ScheduledEvent struct {
	ID       string  `json:"id"`
	SeriesID *string `json:"seriesId,omitempty"`
	// OriginalStart is the instant the recurrence rule produced, immutable
	// even if the occurrence is later moved. Nil for one-off events.
	OriginalStart   *time.Time `json:"-"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	ReminderMessage string     `json:"reminderMessage,omitempty"`
	StartTime       time.Time  `json:"startTime"`
	DurationMinutes int        `json:"durationMinutes"`
	Timezone        string     `json:"timezone"`
	Status          string     `json:"status"`
	FederatedAt     *time.Time `json:"-"`
	CreatedAt       *time.Time `json:"-"`
	UpdatedAt       *time.Time `json:"-"`
	Reminder1SentAt *time.Time `json:"-"`
	Reminder2SentAt *time.Time `json:"-"`
}
